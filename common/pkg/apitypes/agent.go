package apitypes

import "github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"

type AgentDTO struct {
	ID           string                             `json:"id"`
	Name         string                             `json:"name"`
	AgentInfo    AgentInfoDTO                       `json:"agent_info"`
	AgentDevices []hashcattypes.HashcatStatusDevice `json:"agent_devices"`
}

type AgentFileDTO struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type AgentInfoDTO struct {
	Status             string         `json:"status"`
	LastCheckInTime    int64          `json:"last_checkin,omitempty"`
	AvailableListfiles []AgentFileDTO `json:"available_listfiles,omitempty"`
	ActiveJobIDs       []string       `json:"active_job_ids,omitempty"`
}

type AgentGetAllResponseDTO struct {
	Agents []AgentDTO `json:"agents"`
}
