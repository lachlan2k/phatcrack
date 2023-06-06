package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"gorm.io/datatypes"
)

const (
	AgentStatusAlive        = "AgentStatusAlive"
	AgentStatusDisconnected = "AgentStatusDisconnected"
	AgentStatusNeverSeen    = "AgentStatusNeverSeen"
)

type Agent struct {
	UUIDBaseModel
	Name      string
	KeyHash   string
	AgentInfo datatypes.JSONType[AgentInfo]
}

type AgentFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type AgentInfo struct {
	Status             string      `json:"status"`
	LastCheckIn        time.Time   `json:"last_checkin,omitempty"`
	AvailableListfiles []AgentFile `json:"available_listfiles,omitempty"`
	ActiveJobIDs       []string    `json:"active_job_ids,omitempty"`
}

func CreateAgent(name string) (newAgent *Agent, plaintextKey string, err error) {
	plaintextKey, keyHash, err := util.GenAgentKeyAndHash()
	if err != nil {
		return
	}

	agent := &Agent{
		Name:    name,
		KeyHash: keyHash,
	}

	err = GetInstance().Create(agent).Error
	if err != nil {
		return
	}

	newAgent = agent
	return
}

func FindAgentByAuthKey(authKey string) (*Agent, error) {
	keyHash := util.HashAgentKey(authKey)
	agent := &Agent{}
	err := GetInstance().Where(&Agent{KeyHash: keyHash}).First(agent).Error
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func FindAgentIDByAuthKey(authKey string) (string, error) {
	var result struct {
		ID uuid.UUID
	}

	keyHash := util.HashAgentKey(authKey)
	err := GetInstance().Model(&Agent{}).Where(&Agent{KeyHash: keyHash}).First(&result).Error

	if err != nil {
		return "", err
	}

	return result.ID.String(), nil
}

func UpdateAgentStatus(agentID string, status string) error {
	// TODO: learn how to deal with JSONB properly
	return nil
}

func UpdateAgentInfo(agentId string, info AgentInfo) error {
	// TODO: learn how to deal with JSONB properly
	return nil
}
