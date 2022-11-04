package fleet

import (
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

var mapLock sync.Mutex
var fleet = make(map[string]*Agent)

func RegisterAgentFromWebsocket(conn *websocket.Conn, agentId string) (*Agent, error) {
	agent, err := agentFromWebsocket(conn, agentId)
	if err != nil {
		return nil, err
	}

	mapLock.Lock()
	defer mapLock.Unlock()

	if _, agentExists := fleet[agent.agentId]; agentExists {
		return nil, fmt.Errorf("Agent %s was already active", agent.agentId)
	}

	fleet[agent.agentId] = agent
	return agent, nil
}

func RemoveAgentByID(projectId string) {
	mapLock.Lock()
	defer mapLock.Unlock()

	delete(fleet, projectId)
}

func ScheduleJob(jobId string) {
	return
}
