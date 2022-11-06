package apitypes

import (
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type JobCreateRequestDTO struct {
	HashcatParams    hashcattypes.HashcatParams `json:"hashcat_params" validate:"required"`
	Hashes           []string                   `json:"hashes" validate:"required"`
	StartImmediately bool                       `json:"start_immediately"`
	Name             string                     `json:"name" validate:"required,min=5,max=30"`
	Description      string                     `json:"description" validate:"max=1000"`
}

type JobCreateResponseDTO struct {
	ID string `json:"id"`
}

type JobStartResponseDTO struct {
	AgentID string `json:"agent_id"`
}
