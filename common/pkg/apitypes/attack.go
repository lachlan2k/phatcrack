package apitypes

import "github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"

type AttackDTO struct {
	ID            string                     `json:"id"`
	HashType      uint                       `json:"hash_type"`
	Hashes        []HashlistHashDTO          `json:"hashes"`
	HashcatParams hashcattypes.HashcatParams `json:"hashcat_params"`
}

type AttackMultipleDTO struct {
	Attacks []AttackDTO `json:"attacks"`
}

type AttackCreateRequestDTO struct {
	HashcatParams    hashcattypes.HashcatParams `json:"hashcat_params" validate:"required"`
	Hashes           []string                   `json:"hashes" validate:"required"`
	StartImmediately bool                       `json:"start_immediately"`
	Name             string                     `json:"name" validate:"required,min=5,max=30"`
	Description      string                     `json:"description" validate:"max=1000"`
}

type AttackStartResponseDTO struct {
	JobIDs []string `json:"new_job_id"`
}
