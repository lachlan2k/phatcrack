package fleet

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

var ErrJobDoesntExist = errors.New("job doesn't exist")
var ErrJobAlreadyScheduled = errors.New("job already scheduled to start")

var ErrNoAgentsOnline = errors.New("no agents online")

var fleetLock sync.Mutex
var fleet = make(map[string]*Agent)

var observerMapLock sync.Mutex
var observerMap = make(map[string][]chan wstypes.Message)

func Observe(jobId string) chan wstypes.Message {
	observerMapLock.Lock()
	defer observerMapLock.Unlock()

	_, sliceExists := observerMap[jobId]
	if !sliceExists {
		observerMap[jobId] = make([]chan wstypes.Message, 1)
	}

	newObserverChan := make(chan wstypes.Message)
	observerMap[jobId] = append(observerMap[jobId], newObserverChan)

	return newObserverChan
}

func RemoveObserver(observer chan wstypes.Message, jobId string) bool {
	observerMapLock.Lock()
	defer observerMapLock.Unlock()
	defer close(observer)

	jobObservers, exists := observerMap[jobId]
	if !exists {
		return false
	}

	for i, loopObserver := range jobObservers {
		if loopObserver == observer {
			observerMap[jobId] = append(observerMap[jobId][:i], observerMap[jobId][i+1:]...)
			return true
		}
	}

	return false
}

func notifyObservers(jobId string, msg wstypes.Message) {
	observerMapLock.Lock()
	defer observerMapLock.Unlock()

	jobObservers, exists := observerMap[jobId]
	if !exists {
		return
	}

	for _, observer := range jobObservers {
		observerToNotify := observer
		go func() {
			observerToNotify <- msg
		}()
	}
}

func closeObservers(jobId string) {
	observerMapLock.Lock()
	defer observerMapLock.Unlock()

	jobObservers, exists := observerMap[jobId]
	if !exists {
		return
	}

	for _, observer := range jobObservers {
		close(observer)
	}

	delete(observerMap, jobId)
}

func RegisterAgentFromWebsocket(conn *websocket.Conn, agentId string) (*Agent, error) {
	agent := &Agent{
		conn:            conn,
		agentId:         agentId,
		ready:           false,
		latestAgentInfo: nil,
	}

	fleetLock.Lock()
	defer fleetLock.Unlock()

	if _, agentExists := fleet[agent.agentId]; agentExists {
		return nil, fmt.Errorf("Agent %s was already active", agent.agentId)
	}

	fleet[agent.agentId] = agent
	return agent, nil
}

func RemoveAgentByID(projectId string) {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	delete(fleet, projectId)
}

func ScheduleJob(jobId string) (string, error) {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	if len(fleet) == 0 {
		return "", ErrNoAgentsOnline
	}

	jobDb, err := db.GetJob(jobId)
	if err == db.ErrNotFound {
		return "", ErrJobDoesntExist
	}
	job := jobDb.ToDTO()

	// TOOD
	// if job.RuntimeData.Status != db.JobStatusCreated {
	// return "", ErrJobAlreadyScheduled
	// }

	var leastBusyAgent *Agent = nil
	for _, agent := range fleet {
		if !agent.IsHealthy() {
			continue
		}

		if leastBusyAgent == nil {
			leastBusyAgent = agent
			continue
		}

		aJobs := agent.ActiveJobCount()
		lJobs := leastBusyAgent.ActiveJobCount()

		if aJobs == lJobs {
			// Biased semi-random assignment as tie-braker
			if rand.Intn(2) == 1 {
				leastBusyAgent = agent
			}
		}

		if aJobs < lJobs {
			leastBusyAgent = agent
		}
	}

	if leastBusyAgent == nil {
		return "", ErrNoAgentsOnline
	}

	err = db.SetJobScheduled(job.ID, leastBusyAgent.agentId)
	if err != nil {
		return "", fmt.Errorf("failed to set job as scheduled in db: %v", err)
	}

	leastBusyAgent.sendMessage(wstypes.JobStartType, wstypes.JobStartDTO{
		ID:            job.ID,
		HashcatParams: job.HashcatParams,
		TargetHashes:  job.TargetHashes,
	})

	return leastBusyAgent.agentId, nil
}
