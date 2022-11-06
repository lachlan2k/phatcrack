package db

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRoleStandard = "standard"
	UserRoleAdmin    = "admin"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"password_hash"`
	Role         string             `bson:"role"`
}

func seedUser() error {
	count, err := GetUsersColl().CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	log.Printf("Seeding default admin user (admin:changeme)")

	_, err = RegisterUser("admin", "changeme", "admin")
	return err
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func RegisterUser(username, password, role string) (newUserId string, err error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	res, err := GetUsersColl().InsertOne(
		context.Background(),
		User{
			Username:     NormalizeUsername(username),
			PasswordHash: string(passwordHash),
			Role:         role,
		},
	)
	if err != nil {
		return "", fmt.Errorf("couldn't register user in database: %v", err)
	}
	if objectId, ok := res.InsertedID.(primitive.ObjectID); ok {
		newUserId = objectId.Hex()
	} else {
		return "", fmt.Errorf("couldn't cast new object id: %v", res.InsertedID)
	}
	return
}

func GetUserByID(id string) (*User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := GetUsersColl().FindOne(
		context.Background(),
		bson.M{"_id": objId},
	)

	err = res.Err()
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = res.Decode(user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user result (%v): %v", res, err)
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	res := GetUsersColl().FindOne(
		context.Background(),
		bson.M{"username": username},
	)

	err := res.Err()
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = res.Decode(user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user result (%v): %v", res, err)
	}

	return user, nil
}
