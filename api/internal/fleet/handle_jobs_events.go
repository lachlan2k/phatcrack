package fleet

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func (a *Agent) handleJobStarted(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.JobStartedDTO)
	if !ok {
		return fmt.Errorf("couldn't cast %v to job started dto", msg.Payload)
	}

	return db.SetJobStarted(payload.JobID, payload.Time)
}

func (a *Agent) handleJobCrackedHash(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.JobCrackedHashDTO)
	if !ok {
		return fmt.Errorf("couldn't cast %v to heartbeat dto", msg.Payload)
	}

	err := db.AddJobCrackedHash(payload.JobID, payload.Result.Hash, payload.Result.PlaintextHex)
	if err != nil {
		return err
	}

	hashType, err := db.GetJobHashtype(payload.JobID)
	if err != nil {
		return err
	}

	return db.AddPotfileEntry(db.PotfileEntry{
		Hash:         payload.Result.Hash,
		PlaintextHex: payload.Result.PlaintextHex,
		HashType:     uint(hashType),
	})
}

func (a *Agent) handleJobStdLine(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.JobStdLineDTO)
	if !ok {
		return fmt.Errorf("couldn't cast %v to job stdline dto", msg.Payload)
	}

	notifyObservers(payload.JobID, *msg)
	return db.AddJobStdline(payload.JobID, payload.Line, payload.Stream)
}

func (a *Agent) handleJobStatusUpdate(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.JobStatusUpdateDTO)
	if !ok {
		return fmt.Errorf("couldn't cast %v to job status update dto", msg.Payload)
	}

	notifyObservers(payload.JobID, *msg)
	return db.AddJobStatusUpdate(payload.JobID, payload.Status)
}

func (a *Agent) handleJobExited(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.JobExitedDTO)
	if !ok {
		return fmt.Errorf("couldn't cast %v to job exited dto", msg.Payload)
	}

	reason := db.JobStopReasonFinished
	if payload.Error != nil {
		reason = db.JobStopReasonFailed
	}

	notifyObservers(payload.JobID, *msg)
	closeObservers(payload.JobID)

	return db.SetJobExited(payload.JobID, reason, payload.Time)
}
