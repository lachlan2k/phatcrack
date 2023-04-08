package dbnew

import (
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"gorm.io/datatypes"
)

const (
	JobStatusCreated       = "JobStatus-Created"
	JobStatusAwaitingStart = "JobStatus-AwaitingStart"
	JobStatusStarted       = "JobStatus-Started"
	JobStatusExited        = "JobStatus-Exited"
)

const (
	// Clean exit
	JobStopReasonFinished = "JobStopReason-Finished"
	// User stopped it
	JobStopReasonUserStopped = "JobStopReason-UserStopped"
	// Never started in the first place
	JobStopReasonFailedToStart = "JobStopReason-FailedToStart"
	// General failure
	JobStopReasonFailed = "JobStopReason-Failed"
	// Agent timed out and we lost contact
	JobStopReasonTimeout = "JobStopReason-Timeout"
)

type Job struct {
	UUIDBaseModel

	HashlistVersion uint

	Attack   Attack
	AttackID uuid.UUID `gorm:"type:uuid"`

	HashcatParams datatypes.JSONType[hashcattypes.HashcatParams]

	TargetHashes []string `gorm:"type:text[]"`
	HashType     uint

	RuntimeData JobRuntimeData

	AssignedAgent   Agent
	AssignedAgentID uuid.UUID `gorm:"type:uuid"`

	CrackedHashes []JobCrackedHash
}

type JobRuntimeData struct {
	JobID uuid.UUID `gorm:"type:uuid"`

	// when we ask the job to start
	StartRequestTime time.Time
	// when it actually starts on the agent
	StartedTime     time.Time
	StoppedTime     time.Time
	Status          string
	StopReason      string
	ErrorString     string
	AssignedAgentID uuid.UUID

	// []{ stream: string, line: string }
	OutputLines datatypes.JSONSlice[JobRuntimeOutputLine]
	// []hashcattypes.HashcatStatus
	StatusUpdates datatypes.JSONSlice[hashcattypes.HashcatStatus]
}

const (
	JobStdLineStreamStdout = "stdout"
	JobStdLineStreamStderr = "stderr"
)

type JobRuntimeOutputLine struct {
	Stream string
	Line   string
}

func (j *Job) ToDTO() apitypes.JobDTO {
	cracked := make([]apitypes.JobCrackedHashDTO, len(j.CrackedHashes))
	for i, h := range j.CrackedHashes {
		cracked[i] = h.ToDTO()
	}

	return apitypes.JobDTO{
		ID:              j.ID.String(),
		HashlistVersion: j.HashlistVersion,
		AttackID:        j.AttackID.String(),
		HashcatParams:   j.HashcatParams.Data,
		TargetHashes:    j.TargetHashes,
		HashType:        j.HashType,
		RuntimeData:     apitypes.JobRuntimeDataDTO{},
		AssignedAgentID: j.AssignedAgentID.String(),
		CrackedHashes:   cracked,
	}
}

func (j *Job) ToSimpleDTO() apitypes.JobSimpleDTO {
	return apitypes.JobSimpleDTO{
		ID:              j.ID.String(),
		HashlistVersion: j.HashlistVersion,
		AttackID:        j.AttackID.String(),
		HashType:        j.HashType,
		AssignedAgentID: j.AssignedAgentID.String(),
	}
}

type JobCrackedHash struct {
	SimpleBaseModel
	Hash         string
	PlaintextHex string
	JobID        string
}

func (h *JobCrackedHash) ToDTO() apitypes.JobCrackedHashDTO {
	return apitypes.JobCrackedHashDTO{
		Hash:         h.Hash,
		PlaintextHex: h.PlaintextHex,
	}
}

func GetJob(jobId string) (*Job, error) {
	var job Job
	err := GetInstance().First(&job, "id = ?", jobId).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func GetJobProjID(jobId string) (*uuid.UUID, error) {
	var result struct {
		ProjectID uuid.UUID
	}

	err := GetInstance().Table("jobs").Select("name").Where("id = ?", jobId).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result.ProjectID, nil
}

func CreateJob(job *Job) (*Job, error) {
	return job, GetInstance().Create(job).Error
}

func SetJobStarted(jobId string, startTime time.Time) error {
	// TODO
	return nil
}

func SetJobExited(jobId string, reason string, startTime time.Time) error {
	// TODO
	return nil
}

func SetJobScheduled(jobId string, agentId string) error {
	// TODO
	return nil
}

func AddJobCrackedHash(jobId string, hash string, plaintextHex string) error {
	// TODO
	return nil
}

func AddJobStdline(jobId string, line string, stream string) error {
	// TODO
	return nil
}

func AddJobStatusUpdate(jobId string, status hashcattypes.HashcatStatus) error {
	// TODO
	return nil
}

func GetJobHashtype(jobId string) (uint, error) {
	var result struct {
		HashType uint
	}
	err := GetInstance().Model(&Job{}).Select("Hashlist.HashType").Joins("Hashlist").Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.HashType, nil
}

// TODO: rejig for access control?
func GetJobsForAttack(attackId string) ([]Job, error) {
	jobs := []Job{}
	err := GetInstance().Where("AttackID = ?", attackId).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
