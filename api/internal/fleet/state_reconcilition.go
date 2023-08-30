package fleet

import (
	"fmt"
	"path/filepath"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

const deadtimeToUnhealthy = 60 * time.Second
const deadtimetoDead = 120 * time.Second
const disconnectTimeToDead = 60 * time.Second

// We expect jobs will start within 5 seconds, else we'll consider them to be failed
const acceptableJobStartTime = 30 * time.Second

// This will soft fail if the agent isn't connected. In which case, we're probably fine
func tellAgentToKillJob(agentId *uuid.UUID, jobId *uuid.UUID, reason string) {
	if agentId == nil || jobId == nil {
		return
	}

	agentConnection, ok := fleet[agentId.String()]
	if ok {
		agentConnection.sendMessage(wstypes.JobKillType, wstypes.JobKillDTO{
			JobID:      jobId.String(),
			StopReason: reason,
		})
	}
}

// This function aims to look for inconsistencies between what we expect agents to be doing, and what they're actually doing
// It's also responsible for evaluating whether agents should be considered unhealthy or dead, etc and marking jobs as failed if the agent has died
func stateReconciliation() error {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	allAgents, err := db.GetAllAgents()
	if err != nil {
		return err
	}

	allListfiles, err := db.GetAllListfiles()
	if err != nil {
		return err
	}

	incompleteJobs, err := db.GetAllIncompleteJobs(true)
	if err != nil {
		return err
	}

	// Create a convienient map so we can look up agents by ID later
	agentMap := make(map[string]*db.Agent, 0)

	// This is a map from running job ID -> agent ID. This lets us quickly check which agent claims to be running a job
	jobsOk := make(map[string]uuid.UUID, 0)

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
				log.
					WithField("agent_id", agent.ID.String()).
					WithField("time_since_last_heartbeat", time.Since(info.TimeOfLastHeartbeat).Seconds()).
					Warn("Setting agent to dead, because it has been too long since last heartbeat")

				newInfo.Status = db.AgentStatusDead
				needsSave = true
			}

		case db.AgentStatusUnhealthyAndDisconnected:
			if info.TimeOfLastHeartbeat.After(info.TimeOfLastDisconnect) {
				// In this case, it means the reason it disconnected was because the API died.
				// This is because we only update TimeOfLastDisconnect when the API is healthy and it's the agent that dies.
				// Therefore, in this case, we will evaluate it as if it was UnhealthyButConnected

				if time.Since(info.TimeOfLastHeartbeat) > deadtimetoDead {
					log.
						WithField("agent_id", agent.ID.String()).
						WithField("time_since_last_heartbeat", time.Since(info.TimeOfLastHeartbeat).Seconds()).
						Warn("Setting disconnected agent to dead, because it has been too long since last heartbeat")

					newInfo.Status = db.AgentStatusDead
					needsSave = true
				}
			} else {
				// Otherwise, the reason its UnhealthyAndDisconnected is due to the agent disconnecting.
				// As such, we can trust that TimeOfLastDisconnect is accurate.
				if time.Since(info.TimeOfLastDisconnect) > disconnectTimeToDead {
					log.
						WithField("agent_id", agent.ID.String()).
						WithField("time_since_last_disconnect", time.Since(info.TimeOfLastDisconnect).Seconds()).
						Warn("Setting disconnect agent to dead, because it has been too long since it last disconnected")

					newInfo.Status = db.AgentStatusDead
					needsSave = true
				}
			}
		}

		if newInfo.Status == db.AgentStatusDead && len(activeJobs) > 0 {
			for _, jobId := range activeJobs {
				err = db.SetJobExited(jobId, db.JobStopReasonFailed, "The agent running this job died", time.Now())
				if err != nil {
					log.
						WithField("agent_id", agent.ID.String()).
						WithField("job_id", jobId).
						WithError(err).
						Error("Failed to update job status in database")
				}
			}

			newInfo.ActiveJobIDs = []string{}
			needsSave = true
		} else {
			// Agent is healthy, put all of its jobs into our map
			for _, jobId := range activeJobs {
				jobsOk[jobId] = agent.ID
			}
		}

		if needsSave {
			err := db.UpdateAgentInfo(agent.ID.String(), newInfo)
			if err != nil {
				log.
					WithField("agent_id", agent.ID.String()).
					WithError(err).
					Error("Failed to update status of agent")
			}
		}

		agent.AgentInfo.Data = newInfo
	}

	// Make sure listfiles are available on all healthy agents
	for _, listfile := range allListfiles {
		if listfile.PendingDelete {
			// Check to see if any incompleteJobs is using it as either a wordlist or rulefile
			isWordlistInUse := slices.ContainsFunc(incompleteJobs, func(job db.Job) bool {
				predicate := func(fname string) bool {
					return filepath.Base(fname) == listfile.ID.String()
				}

				return slices.ContainsFunc(
					job.HashcatParams.Data.RulesFilenames,
					predicate,
				) || slices.ContainsFunc(
					job.HashcatParams.Data.WordlistFilenames,
					predicate,
				)
			})

			if !isWordlistInUse {
				// Wordlist can now be safely deleted
				db.Delete(&listfile)

				if config.Get().AutomaticallySyncListfiles {
					// Tell all agents to delete it
					broadcastMessageUnsafe(wstypes.DeleteFileRequestType, wstypes.DeleteFileRequestDTO{
						FileID: listfile.ID.String(),
					})
				}

				filerepo.Delete(listfile.ID)
			}
		}

		availableOnAll := true
		numPresent := 0

		for _, agent := range agentMap {
			if agent.AgentInfo.Data.Status != db.AgentStatusHealthy {
				continue
			}

			// Triple nested for loops lets gooo
			// (in all seriousness yes this is really gross but i dont quite care enough)
			// (should refactor listfile to a map but meh)
			found := false
			for _, availableListfile := range agent.AgentInfo.Data.AvailableListfiles {
				if availableListfile.Name == listfile.ID.String() && availableListfile.Size == int64(listfile.SizeInBytes) {
					found = true
					break
				}
			}
			if !found {
				availableOnAll = false
				break
			} else {
				numPresent++
			}
		}

		if numPresent > 0 && listfile.AvailableForUse != availableOnAll {
			if availableOnAll {
				log.WithField("listfile_id", listfile.ID.String()).Warn("Marking listfile as available, as it is now present on all agents")
			} else {
				log.WithField("listfile_id", listfile.ID.String()).Warn("Marking listfile as unavailable, since it was missing on at least one agent")
			}

			listfile.AvailableForUse = availableOnAll
			err := listfile.Save()
			if err != nil {
				log.WithError(err).WithField("listfile_id", listfile.ID.String()).Error("Failed to save listfile's new availability")
			}
		}
	}

	for _, job := range incompleteJobs {
		switch job.RuntimeData.Status {

		case db.JobStatusAwaitingStart:
			if time.Since(job.RuntimeData.StartRequestTime) < acceptableJobStartTime {
				// The job still has time to start before we get grumpy
				continue
			}

			// The job didn't start in time
			runningAgent, jobOk := jobsOk[job.ID.String()]
			if jobOk {
				// Actually, it looks like it did start? One of our agents says they are running it!
				// This condition *shouldn't* be reached, but lets handle it anyway
				log.
					WithField("job_id", job.ID.String()).
					WithField("agent_id", runningAgent.String()).
					Warn("Job was marked as awaiting-start in database, but an agent was found to be running it. This probably indactes a race condition? We'll let it slide.")

				err = db.SetJobStarted(job.ID.String(), "Unknown", time.Now())
				if err != nil {
					log.
						WithField("job_id", job.ID.String()).
						WithError(err).
						Error("Failed to update job status in database")
				}
			} else {
				// As expected, no agent is running the job.
				log.
					WithField("job_id", job.ID.String()).
					WithField("agent_id", runningAgent.String()).
					WithField("start_request_time", job.RuntimeData.StartRequestTime).
					Warn("Job did not start in time")

				err = db.SetJobExited(job.ID.String(), db.JobStopReasonTimeout, "The job did not start in time", time.Now())
				if err != nil {
					log.
						WithField("job_id", job.ID.String()).
						WithError(err).
						Error("Failed to update job status in database")
				}

				// Tell agent to kill this job, incase it *is* running but it just didn't make it through, or its in a broken state.
				// This is an unlikely error condition, but let's handle it just in case.
				tellAgentToKillJob(job.AssignedAgentID, &job.ID, db.JobStopReasonFailed)
			}

		// The job is supposed to be running somewhere, so lets make sure of it
		case db.JobStatusStarted:
			if job.AssignedAgentID == nil {
				// Job hasn't been assigned an agent yet
				// I don't think this state is reachable, but always worthwhile to prevent a nil dereference
				continue
			}

			jobId := job.ID.String()

			agentThatShouldBeRunningJob, isAgentOk := agentMap[job.AssignedAgentID.String()]
			if !isAgentOk {
				// Hmm, the agent doesn't exist at all?
				log.
					WithField("job_id", job.ID.String()).
					WithField("assigned_agent_id", job.AssignedAgentID.String()).
					Warn("Job was assigned to an agent that doesn't exist, this shouldn't hapen. Failing this job.")

				err = db.SetJobExited(jobId, db.JobStopReasonFailed, "The job was assigned to an invalid agent.", time.Now())
				if err != nil {
					log.
						WithField("job_id", job.ID.String()).
						WithError(err).
						Error("Failed to update job status in database")
				}

				continue
			}

			if job.RuntimeData.StartRequestTime.After(agentThatShouldBeRunningJob.AgentInfo.Data.TimeOfLastHeartbeat) {
				// The job was started after the last heartbeat
				// So we can't know whether or not its in the list of running jobs
				continue
			}

			agentRunningJob, jobOk := jobsOk[jobId]
			if !jobOk {
				err = db.SetJobExited(jobId, db.JobStopReasonFailed, "The job disappeared from the agent's list of running jobs. The agent probably died or was restarted.", time.Now())
				if err != nil {
					log.
						WithField("job_id", job.ID.String()).
						WithError(err).
						Error("Failed to update job status in database")
				}

				tellAgentToKillJob(&agentRunningJob, &job.ID, db.JobStopReasonFailed)
				continue
			}

			if agentRunningJob.String() != agentThatShouldBeRunningJob.ID.String() {
				// Absolutely paranoid sanity check, there is no way on earth the wrong agent should end up running the wrong job
				// But we can check it, so we might as well check it
				err = db.SetJobExited(jobId, db.JobStopReasonFailed, "The job was unexpectedly found running on the wrong agent.", time.Now())
				if err != nil {
					log.
						WithField("job_id", job.ID.String()).
						WithError(err).
						Error("Failed to update job status in database")
				}

				tellAgentToKillJob(&agentRunningJob, &job.ID, db.JobStopReasonFailed)
				tellAgentToKillJob(&agentThatShouldBeRunningJob.ID, &job.ID, db.JobStopReasonFailed)
				continue
			}

			// If we've reached here then the job is indeed running, and its running on the correct agent, nothing to do :)
		}
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
			log.WithError(err).Error("State reconciliation failed")
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
		return fmt.Errorf("failed to perform initial state reconciliation: %w", err)
	}

	go stateReconciliationTask()
	return nil
}
