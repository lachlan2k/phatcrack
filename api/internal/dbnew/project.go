package dbnew

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type Project struct {
	UUIDBaseModel
	Name        string
	Description string
	Hashlists   []Hashlist

	OwnerUser   User
	OwnerUserID uuid.UUID `gorm:"type:uuid"`

	ProjectShare []ProjectShare
}

func (p *Project) ToDTO() apitypes.ProjectDTO {
	return apitypes.ProjectDTO{
		ID:          p.ID.String(),
		TimeCreated: p.CreatedAt.UnixMilli(),
		Name:        p.Name,
		Description: p.Description,
	}
}

func CreateProject(proj *Project) (*Project, error) {
	return proj, GetInstance().Create(proj).Error
}

type ProjectShare struct {
	SimpleBaseModel

	ProjectID uuid.UUID `gorm:"type:uuid"`
	Project   *Project

	UserID uuid.UUID `gorm:"type:uuid"`
	User   *User
}

type Hashlist struct {
	UUIDBaseModel
	ProjectID uuid.UUID `gorm:"type:uuid"`

	Name    string
	Version uint

	HashType uint
	Hashes   []HashlistHash

	Attacks []Attack
}

func CreateHashlist(hashlist *Hashlist) (*Hashlist, error) {
	return hashlist, GetInstance().Create(hashlist).Error
}

type HashlistHash struct {
	SimpleBaseModel
	HashlistID     uuid.UUID `gorm:"type:uuid"`
	NormalizedHash string
	InputHash      string
}

type Attack struct {
	UUIDBaseModel
	HashcatParams datatypes.JSONType[hashcattypes.HashcatParams]

	Jobs       []Job
	HashlistID uuid.UUID `gorm:"type:uuid"`
}

func (a *Attack) ToDTO() apitypes.AttackDTO {
	return apitypes.AttackDTO{
		ID:            a.ID.String(),
		HashcatParams: a.HashcatParams.Data,
	}
}

func GetProjectForUser(projId, userId string) (*Project, error) {
	proj := new(Project)

	accessControlQuery := GetInstance().Where(
		"owner_user_id = ?", userId,
	).Or(
		"project_shares.user_id = ?", userId,
	)

	err := GetInstance().Preload("ProjectShare").Select(
		"distinct on (projects.id) projects.*, project_shares.*",
	).Joins(
		"join project_shares on project_shares.project_id = projects.id",
	).Where(
		"projects.id = ?", projId,
	).Where(accessControlQuery).First(proj).Error

	if err != nil {
		return nil, err
	}
	return proj, err
}

func GetAllProjectsForUser(userId string) ([]Project, error) {
	projs := []Project{}

	err := GetInstance().Preload("ProjectShare").Select(
		"distinct on (projects.id) projects.*",
	).Joins(
		"join project_shares on project_shares.project_id = projects.id",
	).Where(
		"owner_user_id = ?", userId,
	).Or(
		"project_shares.user_id = ?", userId,
	).Find(&projs).Error

	if err != nil {
		return nil, err
	}

	return projs, err
}

func GetHashlist(hashlistId string) (*Hashlist, error) {
	var hashlist Hashlist
	err := GetInstance().First(&hashlist, "id = ?", hashlistId).Error
	if err != nil {
		return nil, err
	}
	return &hashlist, nil
}

func GetAttack(attackId string) (*Attack, error) {
	var attack Attack
	err := GetInstance().First(&attack, "id = ?", attackId).Error
	if err != nil {
		return nil, err
	}
	return &attack, nil
}
