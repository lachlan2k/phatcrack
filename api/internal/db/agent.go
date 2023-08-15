package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"gorm.io/datatypes"
)

const (
	AgentStatusHealthy                  = "AgentStatusHealthy"
	AgentStatusUnhealthyButConnected    = "AgentStatusUnhealthyButConnected"
	AgentStatusUnhealthyAndDisconnected = "AgentStatusUnhealthyAndDisconnected"
	AgentStatusDead                     = "AgentStatusDead"
)

type Agent struct {
	UUIDBaseModel
	Name         string
	KeyHash      string
	AgentInfo    datatypes.JSONType[AgentInfo]
	AgentDevices datatypes.JSONType[AgentDeviceInfo]
}

type AgentFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type AgentDeviceInfo struct {
	Devices []hashcattypes.HashcatStatusDevice
}

type AgentInfo struct {
	Status               string      `json:"status"`
	TimeOfLastHeartbeat  time.Time   `json:"time_of_last_heartbeat,omitempty"`
	TimeOfLastDisconnect time.Time   `json:"time_of_last_disconnect,omitempty"`
	TimeOfLastConnect    time.Time   `json:"time_of_last_connect,omitempty"`
	AvailableListfiles   []AgentFile `json:"available_listfiles,omitempty"`
	ActiveJobIDs         []string    `json:"active_job_ids,omitempty"`
}

func (a *AgentFile) ToDTO() apitypes.AgentFileDTO {
	return apitypes.AgentFileDTO{
		Name: a.Name,
		Size: a.Size,
	}
}

func (a *AgentInfo) ToDTO() apitypes.AgentInfoDTO {
	listfileDTOs := make([]apitypes.AgentFileDTO, len(a.AvailableListfiles))
	for i, f := range a.AvailableListfiles {
		listfileDTOs[i] = f.ToDTO()
	}

	return apitypes.AgentInfoDTO{
		Status:             a.Status,
		LastCheckInTime:    a.TimeOfLastHeartbeat.Unix(),
		AvailableListfiles: listfileDTOs,
		ActiveJobIDs:       a.ActiveJobIDs,
	}
}

func (a *Agent) ToDTO() apitypes.AgentDTO {
	return apitypes.AgentDTO{
		ID:           a.ID.String(),
		Name:         a.Name,
		AgentInfo:    a.AgentInfo.Data.ToDTO(),
		AgentDevices: a.AgentDevices.Data.Devices,
	}
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

func GetAllAgents() ([]Agent, error) {
	agents := []Agent{}
	err := GetInstance().Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func GetAllHealthyAgents() ([]Agent, error) {
	agents := []Agent{}
	err := GetInstance().Find(&agents, "agent_info-->'status' = ?", AgentStatusHealthy).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func GetAgent(id string) (*Agent, error) {
	var agent Agent
	err := GetInstance().First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
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

func UpdateAgentDevices(agentId string, devices []hashcattypes.HashcatStatusDevice) error {
	return GetInstance().
		Table("agents").
		Where("id", agentId).
		Update("agent_devices", AgentDeviceInfo{
			Devices: devices,
		}).Error
}

func UpdateAgentStatus(agentId string, status string) error {
	return GetInstance().
		Table("agents").
		Where("id", agentId).
		UpdateColumn("agent_info",
			datatypes.JSONSet("agent_info").Set("{status}", status),
		).
		Error
}

func UpdateAgentInfo(agentId string, info AgentInfo) error {
	return GetInstance().Table("agents").Where("id", agentId).Update("agent_info", info).Error
}
