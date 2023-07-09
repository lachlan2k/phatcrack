package fleet

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func (a *Agent) handleJobStarted(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobStartedDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job started dto: %v", msg.Payload, err)
	}

	return db.SetJobStarted(payload.JobID, payload.Time)
}

func (a *Agent) handleJobCrackedHash(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobCrackedHashDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to cracked hash dto: %v", msg.Payload, err)
	}

	err = db.AddJobCrackedHash(payload.JobID, payload.Result.Hash, payload.Result.PlaintextHex)
	if err != nil {
		return err
	}

	hashType, err := db.GetJobHashtype(payload.JobID)
	if err != nil {
		return err
	}

	_, err = db.AddPotfileEntry(&db.PotfileEntry{
		Hash:         payload.Result.Hash,
		PlaintextHex: payload.Result.PlaintextHex,
		HashType:     uint(hashType),
	})
	return err
}

func (a *Agent) handleJobStdLine(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobStdLineDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job stdline dto: %v", msg.Payload, err)
	}

	notifyObservers(payload.JobID, *msg)
	return db.AddJobStdline(payload.JobID, payload.Line, payload.Stream)
}

func (a *Agent) handleJobStatusUpdate(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobStatusUpdateDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job status update dto: %v", msg.Payload, err)
	}

	if len(payload.Status.Devices) > 0 {
		db.UpdateAgentDevices(a.agentId, payload.Status.Devices)
	}

	notifyObservers(payload.JobID, *msg)
	return db.AddJobStatusUpdate(payload.JobID, payload.Status)
}

func (a *Agent) handleJobExited(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobExitedDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job exited dto: %v", msg.Payload, err)
	}

	reason := db.JobStopReasonFinished
	if payload.Error != "" {
		reason = db.JobStopReasonFailed
	}

	notifyObservers(payload.JobID, *msg)
	closeObservers(payload.JobID)

	return db.SetJobExited(payload.JobID, reason, payload.Error, payload.Time)
}

func (a *Agent) handleJobFailedToStart(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobFailedToStartDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job failed to start dto: %v", msg.Payload, err)
	}

	notifyObservers(payload.JobID, *msg)
	closeObservers(payload.JobID)

	return db.SetJobExited(payload.JobID, db.JobStopReasonFailedToStart, "", payload.Time)
}
