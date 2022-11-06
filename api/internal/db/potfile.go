package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PotfileEntry struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Hash         string             `bson:"hash,omitempty"`
	PlaintextHex string             `bson:"plaintext_hex,omitempty"`
	HashType     uint               `bson:"hash_type,omitempty"`
}

func FindInPotfile(hash string, hashType uint) (*PotfileEntry, error) {
	filter := bson.D{{
		Key: "$and",
		Value: bson.A{
			bson.M{"hash": hash},
			bson.M{"hash_type": hashType},
		},
	}}

	result := GetPotfileColl().FindOne(context.Background(), filter)

	err := result.Err()
	if err != nil {
		return nil, err
	}

	var potfileEntry PotfileEntry
	result.Decode(&potfileEntry)

	return &potfileEntry, nil
}

func AddPotfileEntry(entry PotfileEntry) error {
	_, err := GetPotfileColl().InsertOne(context.Background(), entry)
	if err != nil {
		return fmt.Errorf("couldn't insert potfile entry: %v", err)
	}
	return nil
}
