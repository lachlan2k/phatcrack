package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wordlist struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Filename    string             `bson:"filename"`
	Size        uint64             `bson:"size"`
	Lines       uint64             `bson:"lines"`
}

type RuleFile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Filename    string             `bson:"filename"`
	Size        uint64             `bson:"size"`
	Lines       uint64             `bson:"lines"`
}

func AddWordlist(w Wordlist) error {
	_, err := GetWordlistColl().InsertOne(context.Background(), w)
	if err != nil {
		return fmt.Errorf("couldn't insert wordlist: %v", err)
	}
	return nil
}

func GetWordlist(id string) (*Wordlist, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := GetWordlistColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	list := new(Wordlist)
	err = res.Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func GetAllWordlists() ([]Wordlist, error) {
	cursor, err := GetWordlistColl().Find(
		context.Background(),
		bson.M{},
	)
	if err != nil {
		return nil, err
	}

	var wordlists []Wordlist
	err = cursor.All(context.Background(), &wordlists)

	if err != nil {
		return nil, err
	}
	return wordlists, nil
}

func AddRuleFile(r RuleFile) error {
	_, err := GetRulesColl().InsertOne(context.Background(), r)
	if err != nil {
		return fmt.Errorf("couldn't insert wordlist: %v", err)
	}
	return nil
}

func GetRuleFile(id string) (*RuleFile, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := GetRulesColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	rulefile := new(RuleFile)
	err = res.Decode(&rulefile)
	if err != nil {
		return nil, err
	}

	return rulefile, nil
}

func GetAllRuleFiles() ([]RuleFile, error) {
	cursor, err := GetRulesColl().Find(
		context.Background(),
		bson.M{},
	)
	if err != nil {
		return nil, err
	}

	var rules []RuleFile
	err = cursor.All(context.Background(), &rules)

	if err != nil {
		return nil, err
	}
	return rules, nil
}
