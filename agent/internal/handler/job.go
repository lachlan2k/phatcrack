package handler

import (
	"fmt"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"github.com/lachlan2k/phatcrack/agent/pkg/hashcattypes"
	"github.com/lachlan2k/phatcrack/agent/pkg/wstypes"
)

func (h *Handler) handleJobStart(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.JobStartDTO)
	if !ok {
		return fmt.Errorf("couldn't cast body to JobStartDTO: %v", msg.Payload)
	}

	return h.runJob(payload)
}

func (h *Handler) sendJobStarted(jobId string) {
	h.sendMessage(wstypes.JobStartedType, wstypes.JobStartedDTO{
		Time: time.Now(),
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

	sess, _ := hashcat.NewHashcatSession(job.ID, job.Hashes, params, h.conf)
	sess.Start()

	go func() {

		h.sendJobStarted(job.ID)

	procLoop:
		for {
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
