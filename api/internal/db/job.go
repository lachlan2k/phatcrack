package db

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

	AssignedAgent   Agent      `gorm:"constraint:OnDelete:SET NULL;"`
	AssignedAgentID *uuid.UUID `gorm:"type:uuid"`
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

	OutputLines   pgJSONBArray[JobRuntimeOutputLine]
	StatusUpdates pgJSONBArray[hashcattypes.HashcatStatus]
	CrackedHashes pgJSONBArray[JobCrackedHash]
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
	dto := apitypes.JobDTO{
		ID:              j.ID.String(),
		HashlistVersion: j.HashlistVersion,
		AttackID:        j.AttackID.String(),
		HashcatParams:   j.HashcatParams.Data,
		TargetHashes:    j.TargetHashes,
		HashType:        j.HashType,
		RuntimeData:     j.RuntimeData.ToDTO(),
		RuntimeSummary:  j.RuntimeData.ToSummaryDTO(),
	}

	if j.AssignedAgentID == nil {
		dto.AssignedAgentID = ""
	} else {
		dto.AssignedAgentID = j.AssignedAgentID.String()
	}

	return dto
}

func (r *JobRuntimeData) ToDTO() apitypes.JobRuntimeDataDTO {
	// TODO
	outlines := make([]apitypes.JobRuntimeOutputLineDTO, len(r.OutputLines.Data))
	cracked := make([]apitypes.JobCrackedHashDTO, len(r.CrackedHashes.Data))

	for i, line := range r.OutputLines.Data {
		outlines[i] = apitypes.JobRuntimeOutputLineDTO{
			Stream: line.Data.Stream,
			Line:   line.Data.Line,
		}
	}

	for i, hash := range r.CrackedHashes.Data {
		cracked[i] = hash.Data.ToDTO()
	}

	return apitypes.JobRuntimeDataDTO{
		StartRequestTime: r.StartRequestTime,

		StartedTime: r.StartedTime,
		StoppedTime: r.StoppedTime,
		Status:      r.Status,
		StopReason:  r.StopReason,
		ErrorString: r.ErrorString,

		OutputLines:   outlines,
		StatusUpdates: r.StatusUpdates.Unwrap(),
		CrackedHashes: cracked,
	}
}

func (r *JobRuntimeData) ToSummaryDTO() apitypes.JobRuntimeSummaryDTO {
	dto := apitypes.JobRuntimeSummaryDTO{
		StartedTime:            r.StartedTime.Unix(),
		StoppedTime:            r.StoppedTime.Unix(),
		Hashrate:               0,
		PercentComplete:        -1,
		EstimatedTimeRemaining: -1,
	}

	statusUpdates := r.StatusUpdates.Unwrap()
	if len(statusUpdates) > 0 {
		status := statusUpdates[len(statusUpdates)-1]

		for _, dev := range status.Devices {
			dto.Hashrate += dev.Speed
		}

		if r.Status == JobStatusStarted {
			if status.EstimatedStop > time.Now().Unix() {
				dto.EstimatedTimeRemaining = status.EstimatedStop - time.Now().Unix()
			}

			// This is a tuple with (keyspace covered, total keyspace)
			if len(status.Progress) == 2 {
				dto.PercentComplete = float32(status.Progress[0]) / float32(status.Progress[1])
			}
		}

		if r.Status == JobStatusExited && r.StopReason == JobStopReasonFinished {
			dto.PercentComplete = 100.0
			dto.EstimatedTimeRemaining = 0
		}
	}

	return dto
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

func GetJob(jobId string, includeRuntimeData bool) (*Job, error) {
	var job Job
	inst := GetInstance()
	if includeRuntimeData {
		inst = inst.Preload("RuntimeData")
	}

	err := inst.First(&job, "id = ?", jobId).Error
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func GetAllPendingJobs(includeRuntimeData bool) ([]Job, error) {
	jobs := []Job{}
	inst := GetInstance()
	if includeRuntimeData {
		inst = inst.Preload("RuntimeData")
	}

	err := inst.
		Joins("join job_runtime_data on job_runtime_data.job_id = jobs.id").
		Where("job_runtime_data.status = ?", JobStatusAwaitingStart).
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}
	return jobs, err
}

func GetAllIncompleteJobs(includeRuntimeData bool) ([]Job, error) {
	jobs := []Job{}
	inst := GetInstance()
	if includeRuntimeData {
		inst = inst.Preload("RuntimeData")
	}

	err := inst.
		Joins("join job_runtime_data on job_runtime_data.job_id = jobs.id").
		Where("job_runtime_data.status != ?", JobStatusExited).
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}
	return jobs, err
}

func GetJobProjID(jobId string) (string, error) {
	var result struct {
		ProjectID uuid.UUID
	}

	err := GetInstance().Table("jobs").
		Select("hashlists.project_id as project_id").
		Joins("join attacks on attacks.id = jobs.attack_id").
		Joins("join hashlists on hashlists.id = attacks.hashlist_id").
		Where("jobs.id = ?", jobId).
		Scan(&result).Error

	if err != nil {
		return "", err
	}
	return result.ProjectID.String(), nil
}

func CreateJob(job *Job) (*Job, error) {
	if job.RuntimeData.Status == "" {
		job.RuntimeData.Status = JobStatusCreated
		job.RuntimeData.OutputLines.Init()
		job.RuntimeData.StatusUpdates.Init()
		job.RuntimeData.CrackedHashes.Init()
	}

	return job, GetInstance().Create(job).Error
}

func SetJobStarted(jobId string, startTime time.Time) error {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		return err
	}

	return GetInstance().
		Where("job_id = ?", jobUuid).
		Updates(&JobRuntimeData{
			JobID:       jobUuid,
			Status:      JobStatusStarted,
			StartedTime: startTime,
		}).Error
}

func SetJobExited(jobId string, reason string, errStr string, exitTime time.Time) error {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		return err
	}

	return GetInstance().
		Where("job_id = ?", jobUuid).
		Updates(&JobRuntimeData{
			JobID:       jobUuid,
			Status:      JobStatusExited,
			StopReason:  reason,
			StoppedTime: exitTime,
			ErrorString: errStr,
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
		err := tx.
			Where("job_id = ?", jobUuid).
			Updates(&JobRuntimeData{
				JobID:            jobUuid,
				Status:           JobStatusAwaitingStart,
				StartRequestTime: time.Now(),
			}).Error
		if err != nil {
			return err
		}

		return tx.
			Where("id = ?", jobUuid).
			Updates(&Job{
				AssignedAgentID: &agentUuid,
			}).Error
	})
}

func AddJobCrackedHash(jobId string, hash string, plaintextHex string) error {
	dbLine := datatypes.NewJSONType(JobCrackedHash{
		Hash:         hash,
		PlaintextHex: plaintextHex,
	})

	err := GetInstance().Exec(
		"update job_runtime_data set cracked_hashes = array_append(cracked_hashes, ?) where job_id = ?",
		dbLine, jobId,
	).Error
	if err != nil {
		return err
	}

	return GetInstance().
		Table("hashlist_hashes").
		Where("normalized_hash = ?", hash).
		Updates(&HashlistHash{
			PlaintextHex: plaintextHex,
			IsCracked:    true,
		}).Error
}

// TODO: actually, on second thought, I want to keep all stderr lines, and only roll-over stdout lines
const MaxJobOutputs = 10

func AddJobStdline(jobId string, stream string, line string) error {
	dbLine := datatypes.NewJSONType(JobRuntimeOutputLine{
		Stream: stream,
		Line:   line,
	})

	return GetInstance().Exec(
		"update job_runtime_data set output_lines = array_append(output_lines[array_upper(output_lines, 1) - ?:], ?) where job_id = ?",
		MaxJobOutputs-2, dbLine, jobId,
	).Error
}

func AddJobStatusUpdate(jobId string, status hashcattypes.HashcatStatus) error {
	return GetInstance().Exec(
		"update job_runtime_data set status_updates = array_append(status_updates[array_upper(status_updates, 1) - ?:], ?) where job_id = ?",
		MaxJobOutputs-2, datatypes.NewJSONType(status), jobId,
	).Error
}

func GetJobHashtype(jobId string) (uint, error) {
	var result struct {
		HashType uint
	}

	err := GetInstance().
		Table("jobs").
		Select("hashlists.hash_type as hash_type").
		Joins("join attacks on attacks.id = jobs.attack_id").
		Joins("join hashlists on attacks.hashlist_id = hashlists.id").
		Where("jobs.id = ?", jobId).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.HashType, nil
}

func GetJobsForAttack(attackId string, includeRuntimeData bool) ([]Job, error) {
	jobs := []Job{}

	inst := GetInstance()
	if includeRuntimeData {
		inst = inst.Preload("RuntimeData")
	}

	err := inst.
		Select("distinct on (jobs.id) jobs.*").
		Where("jobs.attack_id = ?", attackId).
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}
	return jobs, nil
}
