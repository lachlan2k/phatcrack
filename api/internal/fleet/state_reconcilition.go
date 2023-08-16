package fleet

import (
	"fmt"
	"log"
	"time"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

const deadtimeToUnhealthy = 60 * time.Second
const deadtimetoDead = 120 * time.Second
const disconnectTimeToDead = 60 * time.Second

// We expect jobs will start within 5 seconds, else we'll consider them to be failed
const acceptableJobStartTime = 5 * time.Second

func stateReconciliation() error {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	allAgents, err := db.GetAllAgents()
	if err != nil {
		return err
	}

	// Jobs that we deem to have "failed" because the agent is dead
	jobsFailed := make([]string, 0)
	jobsOk := make(map[string]interface{}, 0) // Using like a set

	for _, agent := range allAgents {
		info := agent.AgentInfo.Data
		activeJobs := agent.AgentInfo.Data.ActiveJobIDs
		needsSave := false

		newInfo := info

		// Time-based state transitions
		switch info.Status {
		case db.AgentStatusHealthy:
			if time.Since(info.TimeOfLastHeartbeat) > deadtimeToUnhealthy {
				newInfo.Status = db.AgentStatusUnhealthyButConnected
				needsSave = true
			}

		case db.AgentStatusUnhealthyButConnected:
			if time.Since(info.TimeOfLastHeartbeat) > deadtimetoDead {
				newInfo.Status = db.AgentStatusDead
				needsSave = true
			}

		case db.AgentStatusUnhealthyAndDisconnected:
			if time.Since(info.TimeOfLastDisconnect) > disconnectTimeToDead {
				newInfo.Status = db.AgentStatusDead
				needsSave = true
			}
		}

		if newInfo.Status == db.AgentStatusDead {
			jobsFailed = append(jobsFailed, activeJobs...)
			newInfo.ActiveJobIDs = []string{}
		} else {
			for _, job := range activeJobs {
				jobsOk[job] = nil
			}
		}

		if needsSave {
			err := db.UpdateAgentInfo(agent.ID.String(), newInfo)
			if err != nil {
				log.Printf("Warn: Failed to update status of agent %s: %v", agent.ID.String(), err)
			}
		}
	}

	pendingJobs, err := db.GetAllPendingJobs(true)
	if err != nil {
		return err
	}

	for _, job := range pendingJobs {
		if job.RuntimeData.Status == db.JobStatusAwaitingStart {
			if time.Since(job.RuntimeData.StartRequestTime) > acceptableJobStartTime {
				// todo set it failed
				_, jobOk := jobsOk[job.ID.String()]
				if jobOk {
					log.Printf("Warn: Job was marked as awaiaitng-start in database, but an agent was found to be running it. This probably indicates a race condition, and we'll let it slide for now.")
					db.SetJobStarted(job.ID.String(), "Unknown", time.Now())
				} else {
					db.SetJobExited(job.ID.String(), db.JobStopReasonTimeout, "The job did not start in time", time.Now())

					// Tell agent to kill this job, incase it *is* running but it just didn't make it through.
					// It's an unlikely error condition, but just probably tidy to do
					agentConnection, ok := fleet[job.AssignedAgent.ID.String()]
					if ok {
						agentConnection.sendMessage(wstypes.JobKillType, wstypes.JobKillDTO{
							JobID: job.ID.String(),
						})
					}
				}
			}
		}
	}

	for _, jobId := range jobsFailed {
		db.SetJobExited(jobId, db.JobStopReasonFailed, "The agent running this job died", time.Now())
	}

	return nil
}

var stateReconciliationQueue = make(chan interface{}, 1)

func stateReconciliationTask() {
	for {
		var err error

		select {
		case <-time.After(10 * time.Second):
			err = stateReconciliation()
		case <-stateReconciliationQueue:
			err = stateReconciliation()
		}

		if err != nil {
			log.Printf("Error during state reconciliation: %v", err)
		}
	}
}

func QueueStateReconciliation() {
	if len(stateReconciliationQueue) == 0 {
		stateReconciliationQueue <- nil
	}
}

func Setup() error {
	err := stateReconciliation()
	if err != nil {
		return fmt.Errorf("failed to perform initial state reconciliation: %v", err)
	}

	go stateReconciliationTask()
	return nil
}
