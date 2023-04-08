package dbnew

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUIDBaseModel
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	err := GetInstance().First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(id string) (*User, error) {
	var user User
	err := GetInstance().First(&user, "username = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RegisterUser(username, password, role string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     NormalizeUsername(username),
		PasswordHash: string(passwordHash),
		Role:         role,
	}

	err = GetInstance().Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
