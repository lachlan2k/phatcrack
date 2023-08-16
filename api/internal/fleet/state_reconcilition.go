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

	// Create a convienient map so we can look up agents by ID later
	agentMap := make(map[string]*db.Agent, 0)

	// Jobs that we deem to have "failed" because the agent is dead
	jobsFailed := make([]string, 0)
	// Golang doesn't have a Set type, so using this like a set to check if elements are present
	jobsOk := make(map[string]interface{}, 0)

	for _, agent := range allAgents {
		info := agent.AgentInfo.Data
		activeJobs := agent.AgentInfo.Data.ActiveJobIDs
		needsSave := false

		agentMap[agent.ID.String()] = &agent

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

		if newInfo.Status == db.AgentStatusDead && len(activeJobs) > 0 {
			jobsFailed = append(jobsFailed, activeJobs...)
			newInfo.ActiveJobIDs = []string{}
			needsSave = true
		} else {
			// Agent is healthy, put all of its jobs into our set
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

	incompleteJobs, err := db.GetAllIncompleteJobs(true)
	if err != nil {
		return err
	}

	for _, job := range incompleteJobs {

		switch job.RuntimeData.Status {
		case db.JobStatusAwaitingStart:
			{
				if time.Since(job.RuntimeData.StartRequestTime) > acceptableJobStartTime {
					// The job didn't start in time

					_, jobOk := jobsOk[job.ID.String()]
					if jobOk {
						// Actually, it looks like it did start? One of our agents says they are running it!
						// This condition *shouldn't* be reached, but lets handle it anyway
						log.Printf("Warn: Job was marked as awaiaitng-start in database, but an agent was found to be running it. This probably indicates a race condition, and we'll let it slide for now.")
						db.SetJobStarted(job.ID.String(), "Unknown", time.Now())
					} else {
						// As expected, no agent is running the job.
						db.SetJobExited(job.ID.String(), db.JobStopReasonTimeout, "The job did not start in time", time.Now())

						// Tell agent to kill this job, incase it *is* running but it just didn't make it through, or its in a broken state.
						// This is an unlikely error condition, but let's handle it just in case.
						agentConnection, ok := fleet[job.AssignedAgent.ID.String()]
						if ok {
							agentConnection.sendMessage(wstypes.JobKillType, wstypes.JobKillDTO{
								JobID: job.ID.String(),
							})
						}
					}
				}
			}

		// Make sure the job is running somewhere
		case db.JobStatusCreated, db.JobStatusStarted:
			{
				jobId := job.ID.String()
				_, isJobOk := jobsOk[jobId]
				if isJobOk {
					// We observed that this job is running, no problem here.
					continue
				}

				if job.AssignedAgentID == nil {
					// Job hasn't been assigned an agent yet
					// I don't think this state is reacahble, but always worthwhile before we try to nil deference
					continue
				}

				agentThatShouldBeRunningJob, isAgentOk := agentMap[job.AssignedAgentID.String()]
				if !isAgentOk {
					// Hmm, the agent doesn't exist at all?
					log.Printf("Warn: Job %s was found assigned to an agent that doesn't exist (%s), this shouldn't happen. Ignoring this job", jobId, job.AssignedAgentID.String())
					continue
				}

				isJobInList := false
				for _, jobRunningOnAgent := range agentThatShouldBeRunningJob.AgentInfo.Data.ActiveJobIDs {
					if jobRunningOnAgent == jobId {
						isJobInList = true
						break
					}
				}

				// We expected the agent to be running the job, but it wasn't
				if !isJobInList {
					db.SetJobExited(jobId, db.JobStopReasonFailed, "The job unexpectedly disappeared from the agent's list of running jobs", time.Now())
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
		case <-time.After(30 * time.Second):
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
	select {
	case stateReconciliationQueue <- nil:
	default: // Channel already full, already been signalled, no need to block
	}
}

// Requests that a state reconciliation should happen in the next 3 seconds
// This is useful when we don't have an urgent need for one to happen now
// And to avoid a storm if multiple agents heartbeat in rapid succession
// Effectively a throttle/de-bounce
var lazyTimer *time.Timer = nil

func LazyQueueStateReconciliation() {
	if lazyTimer != nil {
		return
	}

	lazyTimer = time.AfterFunc(3*time.Second, func() {
		QueueStateReconciliation()
		lazyTimer = nil
	})
}

func Setup() error {
	agents, err := db.GetAllAgents()
	if err != nil {
		return nil
	}

	// this is a bit manual and could be achieved in one UPDATE query, but I think this is fine for now
	for _, agent := range agents {
		if agent.AgentInfo.Data.Status == db.AgentStatusDead {
			continue
		}
		err = db.UpdateAgentStatus(agent.ID.String(), db.AgentStatusUnhealthyAndDisconnected)
		if err != nil {
			return err
		}
	}

	// This state re-conciliation we manually invoke will go out and mark any agents as dead and jobs as failed, as necessar1y
	err = stateReconciliation()
	if err != nil {
		return fmt.Errorf("failed to perform initial state reconciliation: %v", err)
	}

	go stateReconciliationTask()
	return nil
}
