package db

import (
	"strconv"
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
	Name              string
	KeyHash           string
	IsMaintenanceMode bool `gorm:"default:false; not null"`
	Ephemeral         bool
	AgentInfo         datatypes.JSONType[AgentInfo]
	AgentDevices      datatypes.JSONType[AgentDeviceInfo]
}

type AgentRegistrationKey struct {
	SimpleBaseModel

	Name    string
	KeyHint string
	KeyHash string

	ForEphemeralAgent bool
}

func (a AgentRegistrationKey) ToDTO() apitypes.AdminGetAgentRegistrationKeyDTO {
	return apitypes.AdminGetAgentRegistrationKeyDTO{
		ID:                strconv.FormatUint(uint64(a.ID), 10),
		Name:              a.Name,
		KeyHint:           a.KeyHint,
		ForEphemeralAgent: a.ForEphemeralAgent,
	}
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
	Version              string      `json:"version"`
	TimeOfLastHeartbeat  time.Time   `json:"time_of_last_heartbeat,omitempty"`
	TimeOfLastDisconnect time.Time   `json:"time_of_last_disconnect,omitempty"`
	TimeOfLastConnect    time.Time   `json:"time_of_last_connect,omitempty"`
	AvailableListfiles   []AgentFile `json:"available_listfiles,omitempty"`
	ActiveJobIDs         []string    `json:"active_job_ids,omitempty"`
}

func (a AgentFile) ToDTO() apitypes.AgentFileDTO {
	return apitypes.AgentFileDTO{
		Name: a.Name,
		Size: a.Size,
	}
}

func (a AgentInfo) ToDTO() apitypes.AgentInfoDTO {
	listfileDTOs := make([]apitypes.AgentFileDTO, len(a.AvailableListfiles))
	for i, f := range a.AvailableListfiles {
		listfileDTOs[i] = f.ToDTO()
	}

	return apitypes.AgentInfoDTO{
		Status:             a.Status,
		Version:            a.Version,
		LastCheckInTime:    a.TimeOfLastHeartbeat.Unix(),
		AvailableListfiles: listfileDTOs,
		ActiveJobIDs:       a.ActiveJobIDs,
	}
}

func (a Agent) ToDTO() apitypes.AgentDTO {
	return apitypes.AgentDTO{
		ID:                a.ID.String(),
		Name:              a.Name,
		IsMaintenanceMode: a.IsMaintenanceMode,
		AgentInfo:         a.AgentInfo.Data().ToDTO(),
		AgentDevices:      a.AgentDevices.Data().Devices,
	}
}

func CreateAgent(name string, ephemeral bool) (newAgent *Agent, plaintextKey string, err error) {
	plaintextKey, keyHash, err := util.GenAgentKeyAndHash()
	if err != nil {
		return
	}

	agent := &Agent{
		Name:      name,
		KeyHash:   keyHash,
		Ephemeral: ephemeral,
	}

	err = GetInstance().Create(agent).Error
	if err != nil {
		return
	}

	newAgent = agent
	return
}

func CreateAgentRegistrationKey(name string, ephemeral bool) (newKey *AgentRegistrationKey, plaintextKey string, err error) {
	plaintextKey, keyHash, err := util.GenAgentKeyAndHash()
	if err != nil {
		return
	}

	keyHint := plaintextKey[:4] + "..." + plaintextKey[len(plaintextKey)-4:]

	key := &AgentRegistrationKey{
		Name:              name,
		KeyHash:           keyHash,
		KeyHint:           keyHint,
		ForEphemeralAgent: ephemeral,
	}

	err = GetInstance().Create(key).Error
	if err != nil {
		return
	}

	newKey = key
	return
}

func GetAllAgentRegistrationKeys() ([]AgentRegistrationKey, error) {
	agentRegistrationKeys := []AgentRegistrationKey{}
	err := GetInstance().Find(&agentRegistrationKeys).Error
	if err != nil {
		return nil, err
	}
	return agentRegistrationKeys, nil
}

func GetAgentRegistrationKeyByKey(key string) (*AgentRegistrationKey, error) {
	keyHash := util.HashAgentKey(key)
	registrationKey := &AgentRegistrationKey{}
	err := GetInstance().Where(&AgentRegistrationKey{KeyHash: keyHash}).First(registrationKey).Error
	if err != nil {
		return nil, err
	}
	return registrationKey, nil
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
	err := GetInstance().Find(&agents, "agent_info->>'status' = ?", AgentStatusHealthy).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func GetAllSchedulableAgents() ([]Agent, error) {
	agents := []Agent{}
	err := GetInstance().Find(&agents, "agent_info->>'status' = ? and is_maintenance_mode = false", AgentStatusHealthy).Error
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

func UpdateAgentMaintenanceMode(agentId string, isMaintenance bool) error {
	return GetInstance().Table("agents").Where("id", agentId).Update("is_maintenance_mode", isMaintenance).Error
}
