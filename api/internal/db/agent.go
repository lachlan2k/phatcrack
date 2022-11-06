package db

import (
	"context"
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AgentStatusAlive        = "AgentStatusAlive"
	AgentStatusDisconnected = "AgentStatusDisconnected"
	AgentStatusNeverSeen    = "AgentStatusNeverSeen"
)

type AgentLastCheckIn struct {
	Time primitive.Timestamp `bson:"time,omitempty"`
}

type AgentInfo struct {
	Status      string           `bson:"status,omitempty"`
	LastCheckIn AgentLastCheckIn `bson:"last_checkin,omitempty"`
}

type Agent struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name,omitempty"`
	KeyHash string             `bson:"agent_key_hash,omitempty"`
	Info    AgentInfo          `bson:"agent_info,omitempty"`
}

func UpdateAgentCheckin(agentId string) error {
	objId, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return err
	}
	_, err = GetAgentsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},
		bson.D{{Key: "$set", Value: bson.D{{Key: "agent_info.last_checkin", Value: util.MongoNow()}}}},
	)
	if err != nil {
		return fmt.Errorf("failed to save new agent status in database: %v", err)
	}
	return nil
}

func UpdateAgentStatus(newStatus, agentId string) error {
	objId, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return err
	}
	_, err = GetAgentsColl().UpdateOne(
		context.Background(),
		bson.M{"_id": objId},
		bson.D{{Key: "$set", Value: bson.D{{Key: "agent_info.status", Value: newStatus}}}},
	)
	if err != nil {
		return fmt.Errorf("failed to save new agent status in database: %v", err)
	}
	return nil
}

func CreateAgent(name string) (agentId, plainKey string, err error) {
	plainKey, keyHash, err := util.GenAgentKeyAndHash()
	if err != nil {
		return
	}

	result, err := GetAgentsColl().InsertOne(context.Background(), Agent{
		Name:    name,
		KeyHash: keyHash,
		Info: AgentInfo{
			Status: AgentStatusNeverSeen,
		},
	})

	if err != nil {
		return "", "", fmt.Errorf("couldn't insert agent to database: %v", err)
	}

	if objectId, ok := result.InsertedID.(primitive.ObjectID); ok {
		agentId = objectId.Hex()
	} else {
		return "", "", fmt.Errorf("couldn't cast new object id: %v", result.InsertedID)
	}

	return
}

func FindAgentByAuthKey(authKey string) (*Agent, error) {
	keyHash := util.HashAgentKey(authKey)

	filter := bson.D{{
		Key:   "agent_key_hash",
		Value: keyHash,
	}}

	result := GetAgentsColl().FindOne(context.Background(), filter)

	err := result.Err()
	if err != nil {
		return nil, err
	}

	var agentData Agent
	result.Decode(&agentData)

	return &agentData, nil
}
