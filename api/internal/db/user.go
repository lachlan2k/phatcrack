package db

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func RegisterUser(username, password, role string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = GetUsersColl().InsertOne(
		context.Background(),
		User{
			Username:     NormalizeUsername(username),
			PasswordHash: string(passwordHash),
			Role:         role,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't register user in database: %v", err)
	}
	return nil
}

func LookupUserByUsername(username string) (*User, error) {
	res := GetUsersColl().FindOne(
		context.Background(),
		bson.D{{Key: "username", Value: NormalizeUsername(username)}},
	)

	err := res.Err()
	if err == mongo.ErrNoDocuments {
		return nil, err
	}

	user := new(User)
	err = res.Decode(user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user result (%v): %v", res, err)
	}

	return user, nil
}
