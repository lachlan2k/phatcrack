package apitypes

import "github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"

type AttackDTO struct {
	ID            string                     `json:"id"`
	HashlistID    string                     `json:"hashlist_id"`
	HashcatParams hashcattypes.HashcatParams `json:"hashcat_params"`
}

type AttackWithJobsDTO struct {
	AttackDTO
	Jobs []JobDTO `json:"jobs"`
}

type AttackWithJobsMultipleDTO struct {
	Attacks []AttackWithJobsDTO `json:"attacks"`
}

type AttackMultipleDTO struct {
	Attacks []AttackDTO `json:"attacks"`
}

type AttackCreateRequestDTO struct {
	HashlistID       string                     `json:"hashlist_id" validate:"required"`
	HashcatParams    hashcattypes.HashcatParams `json:"hashcat_params" validate:"required"`
	StartImmediately bool                       `json:"start_immediately"`
	Name             string                     `json:"name" validate:"required,min=5,max=30"`
	Description      string                     `json:"description" validate:"max=1000"`
}

type AttackStartResponseDTO struct {
	JobIDs []string `json:"new_job_id"`
}
