package apitypes

import (
	"time"

	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type JobCreateRequestDTO struct {
	HashcatParams    hashcattypes.HashcatParams `json:"hashcat_params" validate:"required"`
	Hashes           []string                   `json:"hashes" validate:"required,min=1,dive,min=4,required"`
	StartImmediately bool                       `json:"start_immediately"`
	Name             string                     `json:"name" validate:"required,standardname,min=5,max=30"`
	Description      string                     `json:"description" validate:"printascii,max=1000"`
}

type JobCreateResponseDTO struct {
	ID string `json:"id"`
}

type JobStartResponseDTO struct {
	AgentID string `json:"agent_id"`
}

type JobRuntimeOutputLineDTO struct {
	Stream string `json:"stream"`
	Line   string `json:"line"`
}

type JobRuntimeDataDTO struct {
	StartRequestTime time.Time `json:"start_request_time"`
	StartedTime      time.Time `json:"started_time"`
	StoppedTime      time.Time `json:"stopped_time"`
	Status           string    `json:"status"`
	StopReason       string    `json:"stop_reason"`
	ErrorString      string    `json:"error_string"`
	CmdLine          string    `json:"cmd_line"`

	OutputLines   []JobRuntimeOutputLineDTO    `json:"output_lines"`
	StatusUpdates []hashcattypes.HashcatStatus `json:"status_updates"`
}

type JobRuntimeSummaryDTO struct {
	Hashrate               int     `json:"hashrate"`
	EstimatedTimeRemaining int64   `json:"estimated_time_remaining"`
	PercentComplete        float32 `json:"percent_complete"`
	StartedTime            int64   `json:"started_time"`
	StoppedTime            int64   `json:"stopped_time"`
	CmdLine                string  `json:"cmd_line"`
}

type JobDTO struct {
	ID              string                     `json:"id"`
	HashlistVersion uint                       `json:"hashlist_version"`
	AttackID        string                     `json:"attack_id"`
	HashcatParams   hashcattypes.HashcatParams `json:"hashcat_params"`
	TargetHashes    []string                   `json:"target_hashes"`
	HashType        int                        `json:"hash_type"`
	RuntimeData     JobRuntimeDataDTO          `json:"runtime_data"`
	RuntimeSummary  JobRuntimeSummaryDTO       `json:"runtime_summary"`
	AssignedAgentID string                     `json:"assigned_agent_id"`
}

type JobSimpleDTO struct {
	ID              string `json:"id"`
	HashlistVersion uint   `json:"hashlist_version"`
	AttackID        string `json:"attack_id"`
	HashType        int    `json:"hash_type"`
	AssignedAgentID string `json:"assigned_agent_id"`
}

type JobMultipleDTO struct {
	Jobs []JobSimpleDTO `json:"jobs"`
}

type RunningJobForUserDTO struct {
	ProjectID  string `json:"project_id"`
	HashlistID string `json:"hashlist_id"`
	AttackID   string `json:"attack_id"`
	JobID      string `json:"job_id"`
}

type RunningJobsForUserResponseDTO struct {
	Jobs []RunningJobForUserDTO `json:"jobs"`
}
