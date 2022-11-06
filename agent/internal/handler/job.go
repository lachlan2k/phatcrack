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
	h.sendMessage(wstypes.JobExitedType, wstypes.JobExitedDTO{
		JobID: jobId,
		Error: err,
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

func (h *Handler) runJob(job wstypes.JobStartDTO) error {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()

	_, alredyExists := h.activeJobs[job.ID]
	if alredyExists {
		return fmt.Errorf("job %s already exists", job.ID)
	}

	params := hashcat.HashcatParams{
		AttackMode:        job.HashcatParams.AttackMode,
		HashType:          job.HashcatParams.HashType,
		Mask:              job.HashcatParams.Mask,
		WordlistFilenames: job.HashcatParams.WordlistFilenames,
		RulesFilenames:    job.HashcatParams.RulesFilenames,
		AdditionalArgs:    job.HashcatParams.AdditionalArgs,
		OptimizedKernels:  job.HashcatParams.OptimizedKernels,
	}

	sess, err := hashcat.NewHashcatSession(job.ID, job.Hashes, params, h.conf)
	if err != nil {
		h.sendJobFailedToStart(job.ID, err)
		return err
	}

	log.Printf("Job is %v", job)
	log.Printf("Starting job %s", job.ID)

	sess.Start()

	go func() {

		h.sendJobStarted(job.ID)

	procLoop:
		for {
			log.Printf("proc loop")
			select {
			case stdoutLine := <-sess.StdoutLines:
				h.sendJobStdoutLine(job.ID, stdoutLine)

			case stderrLine := <-sess.StdoutLines:
				h.sendJobStderrLine(job.ID, stderrLine)

			case result := <-sess.CrackedHashes:
				h.sendJobCrackedHash(job.ID, hashcattypes.HashcatResult{
					Timestamp:    result.Timestamp,
					Hash:         result.Hash,
					PlaintextHex: result.PlaintextHex,
				})

			case status := <-sess.StatusUpdates:
				h.sendJobStatusUpdate(job.ID, status)

			case err := <-sess.DoneChan:
				h.sendJobExited(job.ID, err)
				break procLoop
			}
		}

		sess.Kill()

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
