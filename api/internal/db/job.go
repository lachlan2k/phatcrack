package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	StartedTime     primitive.Timestamp          `bson:"started_time,omitempty"`
	StoppedTime     primitive.Timestamp          `bson:"stopped_time,omitempty"`
	Status          string                       `bson:"status,omitempty"`
	StopReason      string                       `bson:"stop_reason,omitempty"`
	ErrorString     string                       `bson:"error_string,omitempty"`
	AssignedAgentID primitive.ObjectID           `bson:"assigned_agent_id,omitempty"`
	OutputLines     []JobOutputLine              `bson:"output_lines,omitempty"`
	StatusUpdates   []hashcattypes.HashcatStatus `bson:"status_updates,omitempty"`
}

type JobCrackedHash struct {
	Hash         string `bson:"hash,omitempty"`
	PlaintextHex string `bson:"plaintext_hex,omitempty"`
}

type Job struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID       primitive.ObjectID `bson:"project_id"`
	HashlistID      primitive.ObjectID `bson:"hashlist_id"`
	HashlistVersion uint               `bson:"hashlist_version"`
	AttackID        primitive.ObjectID `bson:"attack_id"`

	HashcatParams HashcatParams    `bson:"hashcat_params"`
	Hashes        []string         `bson:"hashes"`
	HashType      uint             `bson:"hash_type"`
	RuntimeData   RuntimeData      `bson:"runtime_data"`
	CrackedHashes []JobCrackedHash `bson:"cracked_hashes,omitempty"`
}

func SetJobScheduled(jobId, agentId string) error {
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}

	agentObjId, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return err
	}

	_, err = GetJobsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},

		bson.D{{
			Key: "$set",
			Value: bson.D{
				{Key: "runtime_data.status", Value: JobStatusAwaitingStart},
				{Key: "runtime_data.assigned_agent_id", Value: agentObjId},
				{Key: "runtime_data.start_request_time", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
			},
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to set job as started in db: %v", err)
	}
	return nil
}

func SetJobStarted(jobId string, timestamp time.Time) error {
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}

	_, err = GetJobsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},

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
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}
	_, err = GetJobsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},

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
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}
	_, err = GetJobsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},
		bson.D{{
			Key: "$push",
			Value: bson.D{{Key: "cracked_hashes", Value: JobCrackedHash{
				Hash:         hash,
				PlaintextHex: plaintextHex,
			}}},
		}},
		options.Update().SetUpsert(true),
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
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}
	_, err = GetJobsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},

		bson.D{{
			Key: "$push",
			Value: bson.D{{Key: "runtime_data.output_lines", Value: JobOutputLine{
				Line:   line,
				Stream: stream,
			}}},
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to add new job output line to db: %v", err)
	}
	return nil
}

func AddJobStatusUpdate(jobId string, status hashcattypes.HashcatStatus) error {
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}
	_, err = GetJobsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},

		bson.D{{
			Key:   "$push",
			Value: bson.D{{Key: "runtime_data.status_updates", Value: status}},
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to add new job status update to db: %v", err)
	}
	return nil
}

func GetJobHashtype(jobId string) (int, error) {
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return 0, err
	}
	res := GetJobsColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},

		&options.FindOneOptions{
			Projection: bson.D{{Key: "hash_type", Value: 1}},
		},
	)

	err = res.Err()
	if err != nil {
		return 0, err
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

func GetJob(jobId string) (*Job, error) {
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, err
	}

	res := GetJobsColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	job := new(Job)
	err = res.Decode(job)
	if err != nil {
		return nil, fmt.Errorf("failed to decode job result (%v): %v", res, err)
	}

	return job, nil
}

func GetJobProjID(jobId string) (string, error) {
	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return "", err
	}

	res := GetJobsColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},
		options.FindOne().SetProjection(bson.M{"project_id": 1}),
	)

	err = res.Err()
	if err != nil {
		return "", err
	}

	job := new(Job)
	err = res.Decode(job)
	if err != nil {
		return "", fmt.Errorf("failed to decode job result (%v): %v", res, err)
	}

	return job.ID.Hex(), nil
}

func GetJobsForAttack(projId string, hashlistId string, attackId string) ([]Job, error) {
	projObjId, err := primitive.ObjectIDFromHex(projId)
	if err != nil {
		return nil, err
	}

	hashlistObjId, err := primitive.ObjectIDFromHex(hashlistId)
	if err != nil {
		return nil, err
	}

	attackObjId, err := primitive.ObjectIDFromHex(attackId)
	if err != nil {
		return nil, err
	}

	cursor, err := GetJobsColl().Find(
		context.Background(),
		bson.M{
			"$and": []bson.M{
				{"project_id": projObjId},
				{"hashlist_id": hashlistObjId},
				{"attack_id": attackObjId},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0)
	err = cursor.All(context.Background(), &jobs)
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		cursor.Decode(&job)
	}

	return jobs, nil
}

func GetJobForAttack(projId string, hashlistId string, attackId string, jobId string) (*Job, error) {
	projObjId, err := primitive.ObjectIDFromHex(projId)
	if err != nil {
		return nil, err
	}

	hashlistObjId, err := primitive.ObjectIDFromHex(hashlistId)
	if err != nil {
		return nil, err
	}

	attackObjId, err := primitive.ObjectIDFromHex(attackId)
	if err != nil {
		return nil, err
	}

	objId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, err
	}

	res := GetJobsColl().FindOne(
		context.Background(),
		bson.M{
			"$and": []bson.M{
				{"_id": objId},
				{"project_id": projObjId},
				{"hashlist_id": hashlistObjId},
				{"attack_id": attackObjId},
			},
		},
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	job := new(Job)
	err = res.Decode(job)
	if err != nil {
		return nil, fmt.Errorf("failed to decode job result (%v): %v", res, err)
	}

	return job, nil
}

func CreateJob(job Job) (newJobId string, err error) {
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
