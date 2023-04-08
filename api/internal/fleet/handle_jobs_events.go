package fleet

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/dbnew"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func (a *Agent) handleJobStarted(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobStartedDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job started dto: %v", msg.Payload, err)
	}

	return dbnew.SetJobStarted(payload.JobID, payload.Time)
}

func (a *Agent) handleJobCrackedHash(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobCrackedHashDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to cracked hash dto: %v", msg.Payload, err)
	}

	err = dbnew.AddJobCrackedHash(payload.JobID, payload.Result.Hash, payload.Result.PlaintextHex)
	if err != nil {
		return err
	}

	hashType, err := dbnew.GetJobHashtype(payload.JobID)
	if err != nil {
		return err
	}

	_, err = dbnew.AddPotfileEntry(&dbnew.PotfileEntry{
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
	return dbnew.AddJobStdline(payload.JobID, payload.Line, payload.Stream)
}

func (a *Agent) handleJobStatusUpdate(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobStatusUpdateDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job status update dto: %v", msg.Payload, err)
	}

	notifyObservers(payload.JobID, *msg)
	return dbnew.AddJobStatusUpdate(payload.JobID, payload.Status)
}

func (a *Agent) handleJobExited(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobExitedDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job exited dto: %v", msg.Payload, err)
	}

	reason := dbnew.JobStopReasonFinished
	if payload.Error != nil {
		reason = dbnew.JobStopReasonFailed
	}

	notifyObservers(payload.JobID, *msg)
	closeObservers(payload.JobID)

	return dbnew.SetJobExited(payload.JobID, reason, payload.Time)
}

func (a *Agent) handleJobFailedToStart(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.JobFailedToStartDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job failed to start dto: %v", msg.Payload, err)
	}

	notifyObservers(payload.JobID, *msg)
	closeObservers(payload.JobID)

	return dbnew.SetJobExited(payload.JobID, dbnew.JobStopReasonFailedToStart, payload.Time)
}
