package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
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
		LastCheckInTime:    a.LastCheckIn.Unix(),
		AvailableListfiles: listfileDTOs,
		ActiveJobIDs:       a.ActiveJobIDs,
	}
}

func (a *Agent) ToDTO() apitypes.AgentDTO {
	return apitypes.AgentDTO{
		ID:        a.ID.String(),
		Name:      a.Name,
		AgentInfo: a.AgentInfo.Data.ToDTO(),
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
	// return GetInstance().Exec(
	// "UPDATE agents SET \"agent_info\" = jsonb_set(\"agent_info\"::jsonb, '{status}', ?) WHERE id = ?", status, agentID,
	// ).Error
	return GetInstance().
		Table("agents").
		Where("id", agentID).
		UpdateColumn("agent_info",
			datatypes.JSONSet("agent_info").Set("{status}", status),
		).
		Error
}

func UpdateAgentInfo(agentId string, info AgentInfo) error {
	return GetInstance().Table("agents").Where("id", agentId).Update("agent_info", info).Error
}
