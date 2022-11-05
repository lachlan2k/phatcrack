package fleet

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/pkg/wstypes"
)

type Agent struct {
	conn         *websocket.Conn
	agentId      string
	activeJobIDs []string
}

func (a *Agent) Kill() {
	a.conn.Close()
	db.UpdateAgentStatus(db.AgentStatusDisconnected, a.agentId)
	RemoveAgentByID(a.agentId)
}

func (a *Agent) handleMessage(msg *wstypes.Message) error {
	switch msg.Type {
	case wstypes.HeartbeatType:
		return a.handleHeartbeat(msg)

	case wstypes.JobStartedType:
		return a.handleJobStarted(msg)

	case wstypes.JobCrackedHashType:
		return a.handleJobCrackedHash(msg)

	case wstypes.JobStdLineType:
		return a.handleJobStdLine(msg)

	case wstypes.JobExitedType:
		return a.handleJobExited(msg)

	case wstypes.JobStatusUpdateType:
		return a.handleJobStatusUpdate(msg)

	default:
		return fmt.Errorf("unrecognized message type: %s", msg.Type)
	}
}

func (a *Agent) handleHeartbeat(msg *wstypes.Message) error {
	payload, ok := msg.Payload.(wstypes.HeartbeatDTO)
	if !ok {
		return fmt.Errorf("couldn't cast %v to heartbeat dto", msg.Payload)
	}

	a.activeJobIDs = payload.ActiveJobIDs

	err := db.UpdateAgentCheckin(a.agentId)
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) Handle() error {
	log.Printf("handling agent")
	defer a.Kill()
	err := db.UpdateAgentStatus(db.AgentStatusAlive, a.agentId)
	if err != nil {
		return err
	}

	for {
		var msg wstypes.Message
		err := a.conn.ReadJSON(&msg)
		if err != nil {
			return fmt.Errorf("error when trying to read websocket JSON: %v", err)
		}

		err = a.handleMessage(&msg)
		if err != nil {
			return fmt.Errorf("error when handling message: %v", err)
		}

		log.Printf("Received: %v", msg)
	}
}
