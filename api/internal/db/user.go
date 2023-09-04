package db

import (
	"errors"
	"strings"

	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type User struct {
	UUIDBaseModel
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string
	Roles        datatypes.JSONSlice[string]
	APIKeyHash   string // Only for service accounts

	MFAType string
	MFAData datatypes.JSON
}

func (u *User) ToDTO() apitypes.UserDTO {
	return apitypes.UserDTO{
		ID:       u.ID.String(),
		Username: u.Username,
		Roles:    u.Roles,
	}
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

func GetAllUsers() ([]User, error) {
	users := []User{}
	err := GetInstance().Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
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

// callee is responsible for ensuring service account role is present
func RegisterServiceAccount(username string, apiKey string, roles []string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	keyHash := util.HashAPIKey(apiKey)

	user := &User{
		Username:     NormalizeUsername(username),
		PasswordHash: string(passwordHash),
		APIKeyHash:   keyHash,
		Roles:        roles,
	}

	err = GetInstance().Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetServiceAccountByAPIKey(key string) (*User, error) {
	if key == "" {
		return nil, errors.New("empty key provided")
	}

	hash := util.HashAPIKey(key)
	user := &User{}
	err := GetInstance().Where(&User{APIKeyHash: hash}).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
