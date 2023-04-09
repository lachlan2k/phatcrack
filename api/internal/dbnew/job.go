package dbnew

import (
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	AttackID *uuid.UUID `gorm:"type:uuid"`

	HashcatParams datatypes.JSONType[hashcattypes.HashcatParams]

	TargetHashes pq.StringArray `gorm:"type:text[]"`
	HashType     uint

	RuntimeData JobRuntimeData

	AssignedAgent   Agent
	AssignedAgentID *uuid.UUID `gorm:"type:uuid"`

	CrackedHashes datatypes.JSONSlice[JobCrackedHash]
}

type JobRuntimeData struct {
	SimpleBaseModel
	JobID uuid.UUID `gorm:"type:uuid"`

	// when we ask the job to start
	StartRequestTime time.Time
	// when it actually starts on the agent
	StartedTime time.Time
	StoppedTime time.Time
	Status      string
	StopReason  string
	ErrorString string

	// TODO: evaluate JSONB([]line) vs JOSNB(line)[]
	OutputLines   datatypes.JSONSlice[JobRuntimeOutputLine]
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
	Hash         string
	PlaintextHex string
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

func GetJobProjID(jobId string) (string, error) {
	var result struct {
		ProjectID uuid.UUID
	}

	err := GetInstance().Table("jobs").Select("name").Where("id = ?", jobId).Scan(&result).Error
	if err != nil {
		return "", err
	}
	return result.ProjectID.String(), nil
}

func CreateJob(job *Job) (*Job, error) {
	return job, GetInstance().Create(job).Error
}

func SetJobStarted(jobId string, startTime time.Time) error {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		return err
	}

	return GetInstance().Where("job_id = ?", jobUuid).Updates(&JobRuntimeData{
		JobID:       jobUuid,
		Status:      JobStatusStarted,
		StartedTime: startTime,
	}).Error
}

func SetJobExited(jobId string, reason string, exitTime time.Time) error {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		return err
	}

	return GetInstance().Where("job_id = ?", jobUuid).Updates(&JobRuntimeData{
		JobID:       jobUuid,
		Status:      JobStatusExited,
		StopReason:  reason,
		StoppedTime: exitTime,
	}).Error
}

func SetJobScheduled(jobId string, agentId string) error {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		return err
	}

	agentUuid, err := uuid.Parse(agentId)
	if err != nil {
		return err
	}

	return GetInstance().Transaction(func(tx *gorm.DB) error {
		err := tx.Where("job_id = ?", jobUuid).Updates(&JobRuntimeData{
			JobID:            jobUuid,
			Status:           JobStatusAwaitingStart,
			StartRequestTime: time.Now(),
		}).Error
		if err != nil {
			return err
		}

		err = tx.Where("id = ?", jobUuid).Updates(&Job{
			AssignedAgentID: &agentUuid,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func AddJobCrackedHash(jobId string, hash string, plaintextHex string) error {
	dbLine := datatypes.NewJSONType(JobCrackedHash{
		Hash:         hash,
		PlaintextHex: plaintextHex,
	})

	err := GetInstance().Exec(
		"update jobs set cracked_hashes = cracked_hashes || ? where id = ?",
		dbLine, jobId,
	).Error

	return err
}

func AddJobStdline(jobId string, line string, stream string) error {
	dbLine := datatypes.NewJSONType(JobRuntimeOutputLine{
		Stream: stream,
		Line:   line,
	})

	err := GetInstance().Exec(
		"update job_runtime_data set output_lines = output_lines || ? where job_id = ?",
		dbLine, jobId,
	).Error

	return err
}

func AddJobStatusUpdate(jobId string, status hashcattypes.HashcatStatus) error {
	err := GetInstance().Exec(
		"update job_runtime_data set status_updates = status_updates || ? where job_id = ?",
		datatypes.NewJSONType(status), jobId,
	).Error

	return err
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

func GetJobsForAttack(attackId string, projectId string) ([]Job, error) {
	jobs := []Job{}

	err := GetInstance().Select(
		"distinct on (jobs.id), jobs.*",
	).Joins(
		"join attacks on jobs.attack_id = attacks.id",
	).Joins(
		"join hashlists on attack.hashlist_id = hashlists.id",
	).Where(
		"hashlists.project_id = ?", projectId,
	).Where(
		"jobs.attack_id = ?", attackId,
	).Find(&jobs).Error

	if err != nil {
		return nil, err
	}
	return jobs, nil
}
