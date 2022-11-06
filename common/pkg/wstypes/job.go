package wstypes

import (
	"time"

	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

const (
	// server -> agent types
	JobStartType = "JobStart"

	// agent -> server types
	JobStartedType       = "JobStarted"
	JobFailedToStartType = "JobFailedToStart"
	JobCrackedHashType   = "JobCrackedHash"
	JobStdLineType       = "JobStdLine"
	JobExitedType        = "JobExited"
	JobStatusUpdateType  = "JobStatusUpdate"
)

// JobStart
type JobStartDTO struct {
	ID            string                     `json:"id"`
	HashcatParams hashcattypes.HashcatParams `json:"hashcat_parms"`
	Hashes        []string                   `json:"hashes"`
}

// JobFailedToStart
type JobFailedToStartDTO struct {
	JobID string
	Time  time.Time `json:"time"`
	Error error
}

// JobStarted
type JobStartedDTO struct {
	JobID string
	Time  time.Time `json:"time"`
}

// JobCrackedHash
type JobCrackedHashDTO struct {
	JobID  string                     `json:"job_id"`
	Result hashcattypes.HashcatResult `json:"result"`
}

// JobStdLine
const (
	JobStdLineStreamStdout = "stdout"
	JobStdLineStreamStderr = "stderr"
)

type JobStdLineDTO struct {
	JobID  string `json:"job_id"`
	Line   string `json:"line"`
	Stream string `json:"stream"`
}

// JobExited
type JobExitedDTO struct {
	JobID string    `json:"job_id"`
	Time  time.Time `json:"time"`
	Error error     `json:"error"`
}

// JobStatusUpdate
type JobStatusUpdateDTO struct {
	JobID  string                     `json:"job_id"`
	Status hashcattypes.HashcatStatus `json:"status"`
}
