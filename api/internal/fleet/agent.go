package fleet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
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
	log.Printf("received: %v", msg)

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

func (a *Agent) sendMessage(msgType string, payload interface{}) error {
	if a.conn == nil {
		return errors.New("connection closed")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Sending %v %v\n", msgType, string(payloadBytes))

	return a.conn.WriteJSON(wstypes.Message{
		Type:    msgType,
		Payload: string(payloadBytes),
	})
}

func (a *Agent) handleHeartbeat(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.HeartbeatDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to hearbeat dto: %v", msg.Payload, err)
	}

	a.activeJobIDs = payload.ActiveJobIDs

	err = db.UpdateAgentCheckin(a.agentId)
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
			return fmt.Errorf("error when handling %s message: %v", msg.Type, err)
		}

		log.Printf("Received: %v", msg)
	}
}
