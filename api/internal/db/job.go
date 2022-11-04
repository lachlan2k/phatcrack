package db

import (
	"context"
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HashcatParams struct {
	AttackMode        uint8    `bson:"attack_mode"`
	HashType          uint     `bson:"hash_type"`
	Mask              string   `bson:"mask"`
	WordlistFilenames []string `bson:"wordlist_filenames"`
	RulesFilenames    []string `bson:"rules_filenames"`
	AdditionalArgs    []string `bson:"additional_args"`
	OptimizedKernels  bool     `bson:"optimized_kernels"`
}

const (
	// TODO: paused and checkpointed?
	JobStatusStarted = "JobStatus-Started"
	JobStatusStopped = "JobStatus-Stopped"
	JobStatusFailed  = "JobStatus-Failed"
)

const (
	// User stopped it
	JobStopReasonUserStopped = "JobStopReason-UserStopped"
	// General failure
	JobStopReasonFailed = "JobStopReason-Failed"
	// Agent timed out and we lost contact
	JobStopReasonTimeout = "JobStopReason-Timeout"
)

type RuntimeData struct {
	// when we start the job
	StartRequestTime primitive.Timestamp `bson:"start_request_time,omitempty"`
	// when it actually starts on the agent
	StartedTime primitive.Timestamp `bson:"started_time,omitempty"`
	StoppedTime primitive.Timestamp `bson:"stopped_time,omitempty"`
	Status      string              `bson:"status,omitempty"`
	StopReason  string              `bson:"stop_reason,omitempty"`
	StderrLines []string            `bson:"stderr_lines,omitempty"`
	StdoutLines []string            `bson:"stdout_lines,omitempty"`
}

type Job struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty"`
	CreatedTime   primitive.Timestamp `bson:"created_time,omitempty"`
	HashcatParams HashcatParams       `bson:"hashcat_params"`
	Hashes        []string            `bson:"hashes"`
	Name          string              `bson:"name"`
	Description   string              `bson:"description"`
	RuntimeData   []RuntimeData       `bson:"runtime_data"`
}

func CreateJob(job Job) (newJobId string, err error) {
	job.CreatedTime = util.MongoNow()
	result, err := GetJobsColl().InsertOne(context.Background(), job)

	if err != nil {
		return "", fmt.Errorf("couldn't insert job to database: %v", err)
	}

	if objectId, ok := result.InsertedID.(primitive.ObjectID); ok {
		newJobId = objectId.Hex()
	} else {
		return "", fmt.Errorf("couldn't cast new object id: %v", result.InsertedID)
	}

	return
}
