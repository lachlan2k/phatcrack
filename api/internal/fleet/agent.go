package fleet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	conn            *websocket.Conn
	agentId         string
	ready           bool
	latestAgentInfo *db.AgentInfo
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

	case wstypes.JobFailedToStartType:
		return a.handleJobFailedToStart(msg)

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

	activeJobObjIDs := make([]primitive.ObjectID, len(payload.ActiveJobIDs))
	for i, id := range payload.ActiveJobIDs {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("received invalid object ID hex from agent %s: (%s) %v", a.agentId, id, err)
		}
		activeJobObjIDs[i] = objId
	}

	availableWordlists := make([]db.AgentFile, len(payload.Wordlists))
	availableRuleFiles := make([]db.AgentFile, len(payload.RuleFiles))

	for i, list := range payload.Wordlists {
		availableWordlists[i] = db.AgentFile{
			Name: list.Name,
			Size: list.Size,
		}
	}

	for i, file := range payload.RuleFiles {
		availableRuleFiles[i] = db.AgentFile{
			Name: file.Name,
			Size: file.Size,
		}
	}

	info := db.AgentInfo{
		Status: db.AgentStatusAlive,
		LastCheckIn: db.AgentLastCheckIn{
			Time: util.MongoNow(),
		},
		AvailableWordlists: availableWordlists,
		AvailableRuleFiles: availableRuleFiles,
		ActiveJobIDs:       activeJobObjIDs,
	}

	err = db.UpdateAgentInfo(a.agentId, info)

	if err != nil {
		return err
	}

	a.latestAgentInfo = &info
	a.ready = true

	return nil
}

func (a *Agent) IsHealthy() bool {
	if !a.ready {
		return false
	}

	if a.latestAgentInfo == nil {
		return false
	}

	lastHeard := time.Unix(int64(a.latestAgentInfo.LastCheckIn.Time.T), 0)

	// If we've not heard from it for 60 seconds, consider it unhealthy
	return lastHeard.Add(time.Minute).After(time.Now())
}

func (a *Agent) ActiveJobCount() int {
	if !a.ready || a.latestAgentInfo == nil || a.latestAgentInfo.ActiveJobIDs == nil {
		return -1
	}
	return len(a.latestAgentInfo.ActiveJobIDs)
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
