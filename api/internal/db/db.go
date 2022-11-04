package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DatabaseName = "phatcrack"

	UsersCollectionName   = "users"
	AgentsCollectionName  = "agents"
	PotfileCollectionName = "potfile"
	JobsCollectionName    = "jobs"
)

var dbClient *mongo.Client

func Connect(uri string) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("couldn't connect to mongo: %v", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return fmt.Errorf("couldn't ping mongodb: %v", err)
	}

	log.Printf("Connected to mongo successfully")
	dbClient = client

	return nil
}

func GetClient() *mongo.Client {
	return dbClient
}

func GetDatabase() *mongo.Database {
	return dbClient.Database(DatabaseName)
}

func GetUsersColl() *mongo.Collection {
	return GetDatabase().Collection(UsersCollectionName)
}

func GetAgentsColl() *mongo.Collection {
	return GetDatabase().Collection(AgentsCollectionName)
}

func GetPotfileColl() *mongo.Collection {
	return GetDatabase().Collection(PotfileCollectionName)
}

func GetJobsColl() *mongo.Collection {
	return GetDatabase().Collection(JobsCollectionName)
}
