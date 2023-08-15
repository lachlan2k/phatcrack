package fleet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

var ErrJobDoesntExist = errors.New("job doesn't exist")
var ErrJobAlreadyScheduled = errors.New("job already scheduled to start")

var ErrNoAgentsOnline = errors.New("no agents online")

var fleetLock sync.Mutex
var fleet = make(map[string]*AgentConnection)

func RegisterAgentFromWebsocket(conn *websocket.Conn, agentId string) *AgentConnection {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	existingAgent, ok := fleet[agentId]
	if ok {
		// Something went wrong here, we already have the agent in our map.
		// The most likely scenario I imagine this happening is a bug in our code..
		// ..that caused it to not be cleaned up properly on disconnect.
		// OR, the agent has been started twice, OR the same API key re-used on multiple machines

		// So, let's try and close the pprevious websocket connection just to be safe, then delete it, and start again.
		if existingAgent.conn != nil {
			// Gracefully ignore error, it might already be dead
			existingAgent.conn.Close()
		}

		delete(fleet, agentId)
	}

	newAgent := &AgentConnection{
		conn:    conn,
		agentId: agentId,
	}

	fleet[agentId] = newAgent
	return newAgent
}

func RemoveAgentByID(projectId string) {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	delete(fleet, projectId)
}

func ScheduleJobs(jobIds []string) ([]string, error) {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	return scheduleJobUnsafe(jobIds)
}

func NumActiveAgents() int {
	agents, err := db.GetAllAgents()
	if err != nil {
		return 0
	}

	count := 0
	for _, a := range agents {
		if a.AgentInfo.Data.Status == db.AgentStatusHealthy {
			count++
		}
	}

	return count
}

// Jobs will be evenly spread across agents
func scheduleJobUnsafe(jobIds []string) ([]string, error) {
	if len(fleet) == 0 {
		return nil, ErrNoAgentsOnline
	}

	var jobs []apitypes.JobDTO

	for _, jobId := range jobIds {
		jobDb, err := db.GetJob(jobId, false)
		if err == db.ErrNotFound {
			return nil, ErrJobDoesntExist
		}
		jobs = append(jobs, jobDb.ToDTO())
	}

	healthyAgents, err := db.GetAllHealthyAgents()
	if err != nil {
		return nil, err
	}

	for _, agent := range healthyAgents {
		_, ok := fleet[agent.ID.String()]
		if !ok {
			return nil, fmt.Errorf("agent %s was supposed to be healthy, but couldn't be found in the fleet", agent.ID.String())
		}
	}

	if len(healthyAgents) == 0 {
		return nil, ErrNoAgentsOnline
	}

	agentsJobsScheduledTo := []string{}

	for len(jobs) > 0 {
		for _, agent := range healthyAgents {
			if len(jobs) == 0 {
				break
			}
			job := jobs[0]
			jobs = jobs[1:]

			agentConnection := fleet[agent.ID.String()]
			agentsJobsScheduledTo = append(agentsJobsScheduledTo, agent.ID.String())

			db.SetJobScheduled(job.ID, agent.ID.String())
			agentConnection.sendMessage(wstypes.JobStartType, wstypes.JobStartDTO{
				ID:            job.ID,
				HashcatParams: job.HashcatParams,
				TargetHashes:  job.TargetHashes,
			})
		}
	}

	return agentsJobsScheduledTo, nil
}

func RequestFileDownload(fileIDs ...uuid.UUID) {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	for _, agent := range fleet {
		agent.RequestFileDownload(fileIDs...)
	}
}
