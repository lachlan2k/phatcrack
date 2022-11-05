package apitypes

import (
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type JobCreateRequestDTO struct {
	HashcatParams    hashcattypes.HashcatParams `json:"hashcat_params"`
	Hashes           []string                   `json:"hashes"`
	StartImmediately bool                       `json:"start_immediately"`
	Name             string                     `json:"name"`
	Description      string                     `json:"description"`
}

type JobCreateResponseDTO struct {
	ID string `json:"id"`
}

type JobStartResponseDTO struct {
	AgentID string `json:"agent_id"`
}
