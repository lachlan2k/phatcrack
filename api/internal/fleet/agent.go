package fleet

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

type AgentConnection struct {
	conn    *websocket.Conn
	agentId string
}

func (a *AgentConnection) MarkDisconnected() {
	fleetLock.Lock()
	defer fleetLock.Unlock()

	agent, err := db.GetAgent(a.agentId)
	if err != nil {
		db.UpdateAgentStatus(a.agentId, db.AgentStatusUnhealthyAndDisconnected)
		return
	}

	newInfo := agent.AgentInfo.Data
	newInfo.Status = db.AgentStatusUnhealthyAndDisconnected
	newInfo.TimeOfLastDisconnect = time.Now()

	db.UpdateAgentInfo(a.agentId, newInfo)
	a.conn = nil
}

func (a *AgentConnection) handleMessage(msg *wstypes.Message) error {
	fleetLock.Lock()
	defer fleetLock.Unlock()

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
		return fmt.Errorf("unrecognized message type: %q", msg.Type)
	}
}

func (a *AgentConnection) sendMessage(msgType string, payload interface{}) error {
	if a.conn == nil {
		return errors.New("connection closed")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return a.conn.WriteJSON(wstypes.Message{
		Type:    msgType,
		Payload: string(payloadBytes),
	})
}

func (a *AgentConnection) handleHeartbeat(msg *wstypes.Message) error {
	payload, err := util.UnmarshalJSON[wstypes.HeartbeatDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to hearbeat dto: %w", msg.Payload, err)
	}

	availableListfiles := make([]db.AgentFile, len(payload.Listfiles))
	listfilesToCheckMap := make(map[string]*db.AgentFile)

	for i, list := range payload.Listfiles {
		availableListfiles[i] = db.AgentFile{
			Name: list.Name,
			Size: list.Size,
		}
		listfilesToCheckMap[list.Name] = &availableListfiles[i]
	}

	info := db.AgentInfo{
		Status:              db.AgentStatusHealthy,
		TimeOfLastHeartbeat: time.Now(),
		AvailableListfiles:  availableListfiles,
		ActiveJobIDs:        payload.ActiveJobIDs,
	}

	err = db.UpdateAgentInfo(a.agentId, info)

	if err != nil {
		return err
	}

	if !config.Get().AutomaticallySyncListfiles {
		return nil
	}

	// If it's in the middle of downloading a file, let's be polite and not bother it
	if payload.IsDownloadingFile {
		return nil
	}

	expectedListfiles, err := db.GetAllListfiles()
	if err != nil {
		return err
	}

	filesToRequestDownload := []uuid.UUID{}

	for _, expectedFile := range expectedListfiles {
		if !expectedFile.AvailableForDownload {
			continue
		}

		fileOnAgent, ok := listfilesToCheckMap[expectedFile.ID.String()]
		if !ok || fileOnAgent.Size != int64(expectedFile.SizeInBytes) {
			// File is absent, or the length is wrong on disk (probably due to a failed download)
			filesToRequestDownload = append(filesToRequestDownload, expectedFile.ID)
		}
		if ok {
			delete(listfilesToCheckMap, expectedFile.ID.String())
		}
	}

	// Okay, so, any files that are left listfilesToCheckMap are unexpected files.
	// Let's tell the agent to delete those files
	for _, file := range listfilesToCheckMap {
		log.WithField("agent_id", a.agentId).WithField("filename", file.Name).Warn("Found an unexpected file on an agent. Telling it to delete it.")
		a.RequestFileDelete(file.Name)
	}

	if len(filesToRequestDownload) > 0 {
		log.WithField("agent_id", a.agentId).WithField("files_to_download", filesToRequestDownload).Warn("Agent was missing files. Telling it to download them.")
		a.RequestFileDownload(filesToRequestDownload...)
	}

	LazyQueueStateReconciliation()

	return nil
}

func (a *AgentConnection) IsHealthy() bool {
	agent, err := db.GetAgent(a.agentId)
	if err != nil {
		return false
	}
	return agent.AgentInfo.Data.Status == db.AgentStatusHealthy
}

func (a *AgentConnection) RequestFileDownload(fileIDs ...uuid.UUID) error {
	fileIDStrs := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		idStr := id.String()
		if idStr == "" {
			return fmt.Errorf("couldn't parse ID for file download %v", id)
		}
		fileIDStrs[i] = idStr
	}

	return a.sendMessage(wstypes.DownloadFileRequestType, wstypes.DownloadFileRequestDTO{
		FileIDs: fileIDStrs,
	})
}

func (a *AgentConnection) RequestFileDelete(fileID string) error {
	return a.sendMessage(wstypes.DeleteFileRequestType, wstypes.DeleteFileRequestDTO{
		FileID: fileID,
	})
}

func (a *AgentConnection) Handle() error {
	defer a.MarkDisconnected()

	agentInfo, err := db.GetAgent(a.agentId)
	if err != nil {
		return err
	}

	newInfo := agentInfo.AgentInfo.Data
	newInfo.Status = db.AgentStatusUnhealthyButConnected
	newInfo.TimeOfLastConnect = time.Now()

	err = db.UpdateAgentInfo(a.agentId, newInfo)
	if err != nil {
		return err
	}

	for {
		var msg wstypes.Message
		err := a.conn.ReadJSON(&msg)
		if err != nil {
			return fmt.Errorf("error when trying to read websocket JSON: %w", err)
		}

		err = a.handleMessage(&msg)
		if err != nil {
			return fmt.Errorf("error when handling %q message: %w", msg.Type, err)
		}
	}
}
