package db

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

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

	HashType int
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
	HashcatParams  datatypes.JSONType[hashcattypes.HashcatParams]
	IsDistributed  bool
	ProgressString string

	Jobs       []Job     `gorm:"constraint:OnDelete:CASCADE;"`
	HashlistID uuid.UUID `gorm:"type:uuid"`
}

func CreateAttack(attack *Attack) (*Attack, error) {
	return attack, GetInstance().Create(attack).Error
}

func (a *Attack) ToDTO() apitypes.AttackDTO {
	return apitypes.AttackDTO{
		ID:             a.ID.String(),
		HashlistID:     a.HashlistID.String(),
		HashcatParams:  a.HashcatParams.Data,
		IsDistributed:  a.IsDistributed,
		ProgressString: a.ProgressString,
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

// Also deletes hashlists, attacks, jobs
func DeleteProject(projectId string) error {
	return GetInstance().Transaction(func(tx *gorm.DB) error {
		// Jobs
		err := tx.
			Joins("join attacks on attacks.id = jobs.attack_id ").
			Joins("join hashlists on hashlists.id = attacks.hashlist_id ").
			Where("hashlists.project_id = ?", projectId).
			Delete(&Job{}).Error
		if err != nil {
			return err
		}

		// Attacks
		err = tx.
			Joins("join hashlists on hashlists.id = attacks.hashlist_id").
			Where("hashlists.project_id = ?", projectId).
			Delete(&Attack{}).Error
		if err != nil {
			return err
		}

		// Hashlists
		err = tx.Where("project_id = ?", projectId).Delete(&Hashlist{}).Error
		if err != nil {
			return err
		}

		// Project
		return tx.Delete(&Project{}, projectId).Error
	})
}

// Also deletes attacks, jobs
func DeleteHashlist(hashlistId string) error {
	return GetInstance().Transaction(func(tx *gorm.DB) error {
		// Jobs
		err := tx.
			Joins("join attacks on attacks.id = jobs.attack_id").
			Where("attacks.hashlist_id = ?", hashlistId).
			Delete(&Job{}).Error
		if err != nil {
			return err
		}

		// Attacks
		err = tx.Where("hashlist_id = ?", hashlistId).Delete(&Attack{}).Error
		if err != nil {
			return err
		}

		// Hashlist
		return tx.Delete(&Hashlist{}, hashlistId).Error
	})
}

// // Also deletes jobs
func DeleteAttack(attackId string) error {
	return GetInstance().Transaction(func(tx *gorm.DB) error {
		// Jobs
		err := tx.Where("attack_id = ?", attackId).Delete(&Job{}).Error
		if err != nil {
			return err
		}

		// Attack
		return tx.Delete(&Attack{}, attackId).Error
	})
}

func (p *Project) BeforeDelete(tx *gorm.DB) error {
	log.Printf("Handling delete of project %q", p.ID.String())
	return tx.Where("project_id = ?", p.ID).Delete(&Hashlist{}).Error
}

func (h *Hashlist) BeforeDelete(tx *gorm.DB) error {
	log.Printf("Handling delete of hashlist %q", h.ID.String())
	return tx.Where("hashlist_id = ?", h.ID).Delete(&Attack{}).Error
}

func (a *Attack) BeforeDelete(tx *gorm.DB) error {
	log.Printf("Handling delete of attack %q, %q", a.ID.String(), a.HashlistID.String())
	return tx.Where("attack_id = ?", a.ID).Delete(&Job{}).Error
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

func SetAttackProgressString(attackId string, progressString string) error {
	return GetInstance().
		Table("attacks").
		Update("progress_string", progressString).
		Where("id = ?", attackId).Error
}
