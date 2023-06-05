package db

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
	Hashlists   []Hashlist `gorm:"constraint:OnDelete:CASCADE;"`

	OwnerUser   User      `gorm:"constraint:OnDelete:SET NULL;"`
	OwnerUserID uuid.UUID `gorm:"type:uuid"`

	ProjectShare []ProjectShare
}

func (p *Project) ToDTO() apitypes.ProjectDTO {
	return apitypes.ProjectDTO{
		ID:          p.ID.String(),
		TimeCreated: p.CreatedAt.Unix(),
		Name:        p.Name,
		Description: p.Description,
		OwnerUserID: p.OwnerUserID.String(),
	}
}

func CreateProject(proj *Project) (*Project, error) {
	return proj, GetInstance().Create(proj).Error
}

type ProjectShare struct {
	SimpleBaseModel

	ProjectID uuid.UUID `gorm:"type:uuid"`
	Project   *Project  `gorm:"constraint:OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid"`
	User   *User     `gorm:"constraint:OnDelete:CASCADE;"`
}

type Hashlist struct {
	UUIDBaseModel
	ProjectID uuid.UUID `gorm:"type:uuid"`

	Name    string
	Version uint

	HashType uint
	Hashes   []HashlistHash `gorm:"constraint:OnDelete:CASCADE;"`

	Attacks []Attack `gorm:"constraint:OnDelete:CASCADE;"`
}

func (h *Hashlist) ToDTO(withHashes bool) apitypes.HashlistDTO {
	var hashes []apitypes.HashlistHashDTO = nil
	if withHashes {
		hashes = make([]apitypes.HashlistHashDTO, len(h.Hashes))
		for i, hash := range h.Hashes {
			hashes[i] = hash.ToDTO()
		}
	}

	return apitypes.HashlistDTO{
		ID:          h.ID.String(),
		ProjectID:   h.ProjectID.String(),
		Name:        h.Name,
		TimeCreated: h.CreatedAt.Unix(),
		HashType:    h.HashType,
		Hashes:      hashes,
		Version:     h.Version,
	}
}

func CreateHashlist(hashlist *Hashlist) (*Hashlist, error) {
	return hashlist, GetInstance().Create(hashlist).Error
}

type HashlistHash struct {
	SimpleBaseModel
	HashlistID     uuid.UUID `gorm:"type:uuid"`
	NormalizedHash string
	InputHash      string
	PlaintextHex   string
	IsCracked      bool
}

func (h *HashlistHash) ToDTO() apitypes.HashlistHashDTO {
	return apitypes.HashlistHashDTO{
		InputHash:      h.InputHash,
		NormalizedHash: h.InputHash,
		PlaintextHex:   h.PlaintextHex,
		IsCracked:      h.IsCracked,
	}
}

type Attack struct {
	UUIDBaseModel
	HashcatParams datatypes.JSONType[hashcattypes.HashcatParams]

	Jobs       []Job     `gorm:"constraint:OnDelete:CASCADE;"`
	HashlistID uuid.UUID `gorm:"type:uuid"`
}

func CreateAttack(attack *Attack) (*Attack, error) {
	return attack, GetInstance().Create(attack).Error
}

func (a *Attack) ToDTO() apitypes.AttackDTO {
	return apitypes.AttackDTO{
		ID:            a.ID.String(),
		HashlistID:    a.HashlistID.String(),
		HashcatParams: a.HashcatParams.Data,
	}
}

func GetProjectForUser(projId, userId string) (*Project, error) {
	proj := new(Project)

	accessControlQuery := GetInstance().
		Where("owner_user_id = ?", userId).
		Or("project_shares.user_id = ?", userId)

	err := GetInstance().
		Preload("ProjectShare").
		Select("distinct on (projects.id) projects.*").
		Joins("left join project_shares on project_shares.project_id = projects.id").
		Where("projects.id = ?", projId).
		Where(accessControlQuery).First(proj).Error

	if err != nil {
		return nil, err
	}
	return proj, err
}

func GetAllProjectsForUser(userId string) ([]Project, error) {
	projs := []Project{}

	subquery := GetInstance().
		Table("projects").
		Preload("ProjectShare").
		Select("distinct on (projects.id) projects.*").
		Joins("left join project_shares on project_shares.project_id = projects.id").
		Where("owner_user_id = ?", userId).
		Or("project_shares.user_id = ?", userId)

	err := GetInstance().
		Table("(?) as p", subquery).
		Order("p.created_at DESC").
		Find(&projs).Error

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

func GetHashlistWithHashes(hashlistId string) (*Hashlist, error) {
	var hashlist Hashlist
	err := GetInstance().Preload("Hashes").First(&hashlist, "id = ?", hashlistId).Error
	if err != nil {
		return nil, err
	}
	return &hashlist, nil
}

func GetHashlistProjID(hashlistId string) (string, error) {
	var result struct {
		ProjectID uuid.UUID
	}

	err := GetInstance().Model(&Hashlist{}).First(&result, "id = ?", hashlistId).Error

	if err != nil {
		return "", err
	}
	return result.ProjectID.String(), nil
}

func GetAllHashlistsForProject(projId string) ([]Hashlist, error) {
	hashlists := []Hashlist{}
	err := GetInstance().Find(&hashlists, "project_id = ?", projId).Error
	if err != nil {
		return nil, err
	}
	return hashlists, err
}

func GetAttack(attackId string) (*Attack, error) {
	var attack Attack
	err := GetInstance().First(&attack, "id = ?", attackId).Error
	if err != nil {
		return nil, err
	}
	return &attack, nil
}

func GetAllAttacksForHashlist(hashlistId string) ([]Attack, error) {
	attacks := []Attack{}
	err := GetInstance().Find(&attacks, "hashlist_id = ?", hashlistId).Error
	if err != nil {
		return nil, err
	}
	return attacks, nil
}

func GetAttackProjID(attackId string) (string, error) {
	var result struct {
		ProjectID uuid.UUID
	}

	err := GetInstance().
		Table("attacks").
		Select("hashlists.project_id as project_id").
		Joins("join hashlists on hashlists.id = attacks.hashlist_id").
		Where("attacks.id = ?", attackId).
		Scan(&result).Error

	if err != nil {
		return "", err
	}
	return result.ProjectID.String(), nil
}
