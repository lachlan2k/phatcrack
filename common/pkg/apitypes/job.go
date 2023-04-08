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

type JobRuntimeDataDTO struct{}

type JobCrackedHashDTO struct {
	Hash         string
	PlaintextHex string
}

type JobDTO struct {
	ID              string                     `json:"id"`
	HashlistVersion uint                       `json:"hashlist_version"`
	AttackID        string                     `json:"attack_id"`
	HashcatParams   hashcattypes.HashcatParams `json:"hashcat_params"`
	TargetHashes    []string                   `json:"target_hashes"`
	HashType        uint                       `json:"hash_type"`
	RuntimeData     JobRuntimeDataDTO          `json:"runtime_data"`
	AssignedAgentID string                     `json:"assigned_agent_id"`
	CrackedHashes   []JobCrackedHashDTO        `json:"cracked_hashes"`
}

type JobSimpleDTO struct {
	ID              string `json:"id"`
	HashlistVersion uint   `json:"hashlist_version"`
	AttackID        string `json:"attack_id"`
	HashType        uint   `json:"hash_type"`
	AssignedAgentID string `json:"assigned_agent_id"`
}

type JobMultipleDTO struct {
	Jobs []JobSimpleDTO `json:"jobs"`
}
