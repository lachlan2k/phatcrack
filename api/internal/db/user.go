package db

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type User struct {
	UUIDBaseModel
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string
	Roles        datatypes.JSONSlice[string]

	MFAType string
	MFAData datatypes.JSON
}

func (u *User) HasRole(roleToCheck string) bool {
	for _, r := range u.Roles {
		if r == roleToCheck {
			return true
		}
	}

	return false
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func GetUserByID(id string) (*User, error) {
	var user User
	err := GetInstance().First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := GetInstance().First(&user, "username = ?", NormalizeUsername(username)).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RegisterUser(username, password string, roles []string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     NormalizeUsername(username),
		PasswordHash: string(passwordHash),
		Roles:        roles,
	}

	err = GetInstance().Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
