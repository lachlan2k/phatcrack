package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectHashlistHash struct {
	InputHash      string `bson:"input_hash"`
	NormalizedHash string `bson:"normalized_hash"`
}

type ProjectHashlist struct {
	ID       primitive.ObjectID    `bson:"_id,omitempty"`
	Name     string                `bson:"name"`
	HashType uint                  `bson:"hash_type"`
	Hashes   []ProjectHashlistHash `bson:"hashes"`
	Version  uint                  `bson:"version"`
}

type Project struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty"`
	Name              string               `bson:"name"`
	Description       string               `bson:"description"`
	Hashlists         []ProjectHashlist    `bson:"hashlists"`
	OwnerUserID       primitive.ObjectID   `bson:"owner_user_id,omitempty"`
	SharedWithUserIDs []primitive.ObjectID `bson:"shared_with_user_ids,omitempty"`
}

func GetProject(id string) (*Project, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := GetProjectColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	proj := new(Project)
	err = res.Decode(&proj)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func GetProjectForUser(id, userId string) (*Project, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{
		Key: "$and",
		Value: bson.A{
			bson.M{"_id": objId},

			bson.D{{Key: "$or", Value: bson.A{
				bson.M{"owner_user_id": userObjId},
				bson.M{"shared_with_user_ids": bson.M{"$in": bson.A{userObjId}}},
			}}},
		},
	}}

	res := GetProjectColl().FindOne(
		context.Background(),
		filter,
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	proj := new(Project)
	err = res.Decode(&proj)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func GetAllProjectForUser(userId string) ([]Project, error) {
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter :=
		bson.D{{Key: "$or", Value: bson.A{
			bson.M{"owner_user_id": userObjId},
			bson.M{"shared_with_user_ids": bson.M{"$in": bson.A{userObjId}}},
		}}}

	cursor, err := GetProjectColl().Find(
		context.Background(),
		filter,
	)

	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	err = cursor.All(context.Background(), &projects)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		cursor.Decode(&project)
	}

	return projects, nil
}

func CreateProject(proj Project, ownerId string) (newProjectId string, err error) {
	ownerObjId, err := primitive.ObjectIDFromHex(ownerId)
	if err != nil {
		return "", err
	}

	proj.OwnerUserID = ownerObjId
	result, err := GetProjectColl().InsertOne(context.Background(), proj)

	if err != nil {
		return "", fmt.Errorf("couldn't insert project to database: %v", err)
	}

	if objectId, ok := result.InsertedID.(primitive.ObjectID); ok {
		newProjectId = objectId.Hex()
	} else {
		return "", fmt.Errorf("couldn't cast new object id: %v", result.InsertedID)
	}

	return
}

func AddHashlistToProject(projectId string, hashlist ProjectHashlist) (newHashlistId string, err error) {
	projectObjId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return
	}

	hashlist.ID = primitive.NewObjectID()

	output, err := GetProjectColl().UpdateOne(
		context.Background(),
		bson.M{"_id": projectObjId},

		bson.D{{
			Key:   "$push",
			Value: bson.D{{Key: "hashlists", Value: hashlist}},
		}},
	)
	if err != nil {
		err = fmt.Errorf("failed to add new hashlist to project %s: %v", projectId, err)
		return
	}

	// TODO
	fmt.Printf("Added new hashlist: %v\n", output)
	return hashlist.ID.Hex(), nil
}
