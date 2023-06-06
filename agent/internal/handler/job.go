package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"github.com/lachlan2k/phatcrack/agent/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func (h *Handler) handleJobStart(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobStartDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job start dto: %v", msg.Payload, err)
	}

	return h.runJob(payload)
}

func (h *Handler) sendJobStarted(jobId string) {
	h.sendMessage(wstypes.JobStartedType, wstypes.JobStartedDTO{
		JobID: jobId,
		Time:  time.Now(),
	})
}

func (h *Handler) sendJobStdoutLine(jobId, line string) {
	h.sendMessage(wstypes.JobStdLineType, wstypes.JobStdLineDTO{
		JobID:  jobId,
		Line:   line,
		Stream: wstypes.JobStdLineStreamStdout,
	})
}

func (h *Handler) sendJobStderrLine(jobId, line string) {
	h.sendMessage(wstypes.JobStdLineType, wstypes.JobStdLineDTO{
		JobID:  jobId,
		Line:   line,
		Stream: wstypes.JobStdLineStreamStderr,
	})
}

func (h *Handler) sendJobExited(jobId string, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	// Hashcat gives exit code of 1 when a wordlist is exhausted
	// But not when a mask doesn't crack everything
	// So, we're just gonna silently ignore this?
	// I hate it
	if errStr == "exit status 1" {
		errStr = ""
	}

	h.sendMessage(wstypes.JobExitedType, wstypes.JobExitedDTO{
		JobID: jobId,
		Error: errStr,
		Time:  time.Now(),
	})
}

func (h *Handler) sendJobCrackedHash(jobId string, result hashcattypes.HashcatResult) {
	h.sendMessage(wstypes.JobCrackedHashType, wstypes.JobCrackedHashDTO{
		JobID:  jobId,
		Result: result,
	})
}

func (h *Handler) sendJobStatusUpdate(jobId string, status hashcattypes.HashcatStatus) {
	h.sendMessage(wstypes.JobStatusUpdateType, wstypes.JobStatusUpdateDTO{
		JobID:  jobId,
		Status: status,
	})
}

func (h *Handler) sendJobFailedToStart(jobId string, err error) {
	h.sendMessage(wstypes.JobFailedToStartType, wstypes.JobFailedToStartDTO{
		JobID: jobId,
		Time:  time.Now(),
		Error: err,
	})
}

var backoffTable = []util.BackoffEntry{
	{
		AfterTime: time.Duration(0),
		TimeApart: time.Duration(0),
	},
	{
		AfterTime: time.Minute,
		TimeApart: 10 * time.Second,
	},
	{
		AfterTime: 10 * time.Minute,
		TimeApart: 30 * time.Second,
	},
	{
		AfterTime: time.Hour,
		TimeApart: time.Minute,
	},
	{
		AfterTime: 8 * time.Hour,
		TimeApart: 5 * time.Minute,
	},
}

func (h *Handler) runJob(job wstypes.JobStartDTO) error {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()

	_, alredyExists := h.activeJobs[job.ID]
	if alredyExists {
		return fmt.Errorf("job %s already exists", job.ID)
	}

	sess, err := hashcat.NewHashcatSession(job.ID, job.TargetHashes, hashcat.HashcatParams(job.HashcatParams), h.conf)
	if err != nil {
		h.sendJobFailedToStart(job.ID, err)
		return err
	}

	log.Printf("Starting job %s", job.ID)

	err = sess.Start()
	if err != nil {
		h.sendJobFailedToStart(job.ID, err)
		return err
	}

	go func() {

		h.sendJobStarted(job.ID)

		statusBackoff := util.Backoff{
			Entries: backoffTable,
		}
		stdoutBackoff := util.Backoff{
			Entries: backoffTable,
		}

		statusBackoff.Start()
		stdoutBackoff.Start()

	procLoop:
		for {
			select {
			case stdoutLine := <-sess.StdoutLines:
				// If it's a json status update, use our rate limit
				if len(stdoutLine) == 0 || stdoutLine[0] != '{' || stdoutBackoff.Ready() {
					h.sendJobStdoutLine(job.ID, stdoutLine)
				}

			case stderrLine := <-sess.StderrMessages:
				h.sendJobStderrLine(job.ID, stderrLine)

			case result := <-sess.CrackedHashes:
				h.sendJobCrackedHash(job.ID, hashcattypes.HashcatResult{
					Timestamp:    result.Timestamp,
					Hash:         result.Hash,
					PlaintextHex: result.PlaintextHex,
				})

			case status := <-sess.StatusUpdates:
				if statusBackoff.Ready() {
					h.sendJobStatusUpdate(job.ID, status)
				}

			case err := <-sess.DoneChan:
				h.sendJobExited(job.ID, err)
				break procLoop
			}
		}

		sess.Kill()
		sess.Cleanup()

		h.jobsLock.Lock()
		defer h.jobsLock.Unlock()

		delete(h.activeJobs, job.ID)

	}()

	h.activeJobs[job.ID] = ActiveJob{
		job:  job,
		sess: sess,
	}

	return nil
}
