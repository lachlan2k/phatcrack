package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lachlan2k/phatcrack/agent/pkg/hashcattypes"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// General failure
	JobStopReasonFailed = "JobStopReason-Failed"
	// Agent timed out and we lost contact
	JobStopReasonTimeout = "JobStopReason-Timeout"
)

const (
	JobStdLineStreamStdout = "stdout"
	JobStdLineStreamStderr = "stderr"
)

type JobOutputLine struct {
	Line   string `bson:"line"`
	Stream string `bson:"stream"`
}

type RuntimeData struct {
	// when we ask the job to start
	StartRequestTime primitive.Timestamp `bson:"start_request_time,omitempty"`
	// when it actually starts on the agent
	StartedTime   primitive.Timestamp          `bson:"started_time,omitempty"`
	StoppedTime   primitive.Timestamp          `bson:"stopped_time,omitempty"`
	Status        string                       `bson:"status,omitempty"`
	StopReason    string                       `bson:"stop_reason,omitempty"`
	ErrorString   string                       `bson:"error_string,omitempty"`
	OutputLines   []JobOutputLine              `bson:"output_line"`
	StatusUpdates []hashcattypes.HashcatStatus `bson:"status_updates"`
}

type JobCrackedHash struct {
	Hash         string `bson:"hash,omitempty"`
	PlaintextHex string `bson:"plaintext_hex,omitempty"`
}

type Job struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"`
	CreatedTime     primitive.Timestamp `bson:"created_time,omitempty"`
	HashcatParams   HashcatParams       `bson:"hashcat_params"`
	Hashes          []string            `bson:"hashes"`
	HashType        int                 `bson:"hash_type"`
	Name            string              `bson:"name"`
	Description     string              `bson:"description"`
	AssignedAgentID primitive.ObjectID  `bson:"assigned_agent_id,omitempty"`
	RuntimeData     RuntimeData         `bson:"runtime_data"`
	CrackedHashes   []JobCrackedHash    `bson:"cracked_hashes"`
}

func SetJobStarted(jobId string, timestamp time.Time) error {
	_, err := GetJobsColl().UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: jobId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{Key: "runtime_data.status", Value: JobStatusStarted},
				{Key: "runtime_data.started_time", Value: primitive.Timestamp{T: uint32(timestamp.Unix())}},
			},
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to set job as started in db: %v", err)
	}
	return nil
}

func SetJobExited(jobId string, reason string, timestamp time.Time) error {
	_, err := GetJobsColl().UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: jobId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{Key: "runtime_data.status", Value: JobStatusExited},
				{Key: "runtime_data.stop_reason", Value: reason},
				{Key: "runtime_data.stopped_time", Value: primitive.Timestamp{T: uint32(timestamp.Unix())}},
			},
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to set job as started in db: %v", err)
	}
	return nil
}

func AddJobCrackedHash(jobId, hash, plaintextHex string) error {
	_, err := GetJobsColl().UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: jobId}},
		bson.D{{
			Key: "$push",
			Value: bson.D{{Key: "cracked_hashes", Value: JobCrackedHash{
				Hash:         hash,
				PlaintextHex: plaintextHex,
			}}},
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to add new cracked hash to db: %v", err)
	}
	return nil
}

func AddJobStdline(jobId, line, stream string) error {
	if stream != JobStdLineStreamStderr && stream != JobStdLineStreamStdout {
		return fmt.Errorf("unrecognized job line stream %s", stream)
	}
	_, err := GetJobsColl().UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: jobId}},
		bson.D{{
			Key: "$push",
			Value: bson.D{{Key: "runtime_data.output_lines", Value: JobOutputLine{
				Line:   line,
				Stream: stream,
			}}},
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to add new job output line to db: %v", err)
	}
	return nil
}

func AddJobStatusUpdate(jobId string, status hashcattypes.HashcatStatus) error {
	_, err := GetJobsColl().UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: jobId}},
		bson.D{{
			Key:   "$push",
			Value: bson.D{{Key: "runtime_data.status_updates", Value: status}},
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to add new job status update to db: %v", err)
	}
	return nil
}

func GetJobHashtype(jobId string) (int, error) {
	res := GetJobsColl().FindOne(
		context.Background(),
		bson.D{{Key: "_id", Value: jobId}},
		&options.FindOneOptions{
			Projection: bson.D{{Key: "hash_type", Value: 1}},
		},
	)

	err := res.Err()
	if err != nil {
		return 0, fmt.Errorf("couldn't get job hashtype: %v", err)
	}

	var out struct {
		HashType int `bson:"hash_type"`
	}

	err = res.Decode(&out)
	if err != nil {
		return 0, fmt.Errorf("couldn't decode hash type: %v", err)
	}

	return out.HashType, nil
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
