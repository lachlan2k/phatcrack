package db

import (
	"strconv"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/lachlan2k/phatcrack/api/internal/roles"
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

	ProjectShare []ProjectShare `gorm:"constraint:OnDelete:CASCADE;"`
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

func CreateProjectShare(share *ProjectShare) (*ProjectShare, error) {
	return share, GetInstance().Create(share).Error
}

func DeleteProjectShare(projId string, userId string) error {
	return GetInstance().Delete(&ProjectShare{}, "project_id = ? and user_id = ?", projId, userId).Error
}

type ProjectShares []ProjectShare

func (shares ProjectShares) ToDTO() apitypes.ProjectSharesDTO {
	userIDs := []string{}
	for _, share := range shares {
		userIDs = append(userIDs, share.UserID.String())
	}
	return apitypes.ProjectSharesDTO{
		UserIDs: userIDs,
	}
}

func GetProjectShares(projId string) (ProjectShares, error) {
	shares := []ProjectShare{}
	err := GetInstance().Where("project_id = ?", projId).Find(&shares).Error
	if err != nil {
		return nil, err
	}
	return shares, nil
}

type Hashlist struct {
	UUIDBaseModel
	ProjectID uuid.UUID `gorm:"type:uuid"`

	Name    string
	Version uint

	HashType     int
	Hashes       []HashlistHash `gorm:"constraint:OnDelete:CASCADE;"`
	HasUsernames bool

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
		ID:           h.ID.String(),
		ProjectID:    h.ProjectID.String(),
		Name:         h.Name,
		TimeCreated:  h.CreatedAt.Unix(),
		HashType:     h.HashType,
		Hashes:       hashes,
		Version:      h.Version,
		HasUsernames: h.HasUsernames,
	}
}

func CreateHashlist(hashlist *Hashlist) (*Hashlist, error) {
	err := GetInstance().Transaction(func(tx *gorm.DB) error {
		hashes := hashlist.Hashes
		hashlist.Hashes = []HashlistHash{}

		err := tx.Create(hashlist).Error
		if err != nil {
			return err
		}

		for i := range hashes {
			hashes[i].HashlistID = hashlist.ID
		}

		return tx.CreateInBatches(hashes, 1000).Error
	})

	return hashlist, err
}

// Caller must ensure HashlistHash.HashlistID is set correctly
func AppendToHashlist(newHashes []HashlistHash) (int64, error) {
	count := int64(0)

	err := GetInstance().Transaction(func(tx *gorm.DB) error {
		for _, newHash := range newHashes {
			existingCount := int64(0)
			err := tx.Table("hashlist_hashes").Where(
				"hashlist_id = ? and input_hash = ? and username = ?", newHash.HashlistID.String(), newHash.InputHash, newHash.Username,
			).Count(&existingCount).Error

			if err != nil && err != ErrNotFound {
				return err
			}

			if existingCount == 0 {
				err := tx.Create(&newHash).Error
				if err != nil {
					return err
				}
				count++
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func PopulateHashlistFromPotfile(hashlistId string) (int64, error) {
	res := GetInstance().Exec(
		`UPDATE hashlist_hashes
		SET plaintext_hex = potfile_entries.plaintext_hex, is_cracked = true
		FROM potfile_entries, hashlists 
		WHERE hashlist_hashes.hashlist_id = ?
		AND potfile_entries.hash = hashlist_hashes.normalized_hash 
		AND potfile_entries.hash_type = hashlists.hash_type
		AND hashlist_hashes.is_cracked = false`,
		hashlistId,
	)
	return res.RowsAffected, res.Error
}

type HashlistHash struct {
	SimpleBaseModel
	HashlistID     uuid.UUID `gorm:"type:uuid"`
	NormalizedHash string
	InputHash      string
	PlaintextHex   string
	Username       string
	IsCracked      bool
	IsUnexpected   bool
}

func (h *HashlistHash) ToDTO() apitypes.HashlistHashDTO {
	return apitypes.HashlistHashDTO{
		ID:             strconv.FormatInt(int64(h.ID), 10),
		InputHash:      h.InputHash,
		Username:       h.Username,
		NormalizedHash: h.InputHash,
		PlaintextHex:   h.PlaintextHex,
		IsCracked:      h.IsCracked,
		IsUnexpected:   h.IsUnexpected,
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
		HashcatParams:  a.HashcatParams.Data(),
		IsDistributed:  a.IsDistributed,
		ProgressString: a.ProgressString,
	}
}

func GetProject(projId string) (*Project, error) {
	proj := &Project{}
	err := GetInstance().First(&proj, "id = ?", projId).Error
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func GetProjectForUser(projId string, user *User) (*Project, error) {
	proj := &Project{}
	if user.HasRole(roles.UserRoleAdmin) {
		return GetProject(projId)
	}

	err := GetInstance().
		Preload("ProjectShare").
		Select("distinct on (projects.id) projects.*").
		Joins("left join project_shares on project_shares.project_id = projects.id").
		Where("projects.id = ?", projId).
		Where("owner_user_id = ? or project_shares.user_id = ?", user.ID, user.ID).First(proj).Error

	if err != nil {
		return nil, err
	}
	return proj, err
}

func GetAllProjects() ([]Project, error) {
	projs := []Project{}
	err := GetInstance().Order("created_at DESC").Find(&projs).Error
	if err != nil {
		return nil, err
	}
	return projs, nil
}

func GetAllProjectsForUser(user *User) ([]Project, error) {
	if user.HasRole(roles.UserRoleAdmin) {
		return GetAllProjects()
	}
	projs := []Project{}

	subquery := GetInstance().
		Table("projects").
		Select("distinct on (projects.id) projects.*").
		Joins("left join project_shares on project_shares.project_id = projects.id").
		Where("owner_user_id = ? or project_shares.user_id = ?", user.ID, user.ID)

	err := GetInstance().
		Table("(?) as p", subquery).
		Preload("ProjectShare").
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
	err := GetInstance().Order("created_at DESC").Find(&hashlists, "project_id = ?", projId).Error
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

func GetAllAttacksForHashlist(hashlistId string) ([]Attack, error) {
	attacks := []Attack{}
	err := GetInstance().Order("created_at DESC").Find(&attacks, "hashlist_id = ?", hashlistId).Error
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

type AttackIDTree struct {
	ProjectID  string
	HashlistID string
	AttackID   string
}

func (a *AttackIDTree) ToDTO() apitypes.AttackIDTreeDTO {
	return apitypes.AttackIDTreeDTO{
		ProjectID:  a.ProjectID,
		HashlistID: a.HashlistID,
		AttackID:   a.AttackID,
	}
}

func GetAllAttacksWithProgressStringsForUser(user *User) ([]AttackIDTree, error) {
	attacks := []AttackIDTree{}

	query := GetInstance().
		Select("projects.id as project_id, hashlists.id as hashlist_id, attacks.id as attack_id").
		Table("attacks").
		Joins("join hashlists on hashlists.id = attacks.hashlist_id").
		Joins("join projects on projects.id = hashlists.project_id").
		Joins("left join project_shares on project_shares.project_id = projects.id").
		Where("starts_with(progress_string, 'Processing')")

	if !user.HasRole(roles.UserRoleAdmin) {
		query = query.Where("owner_user_id = ? or project_shares.user_id = ?", user.ID, user.ID)
	}

	err := query.Find(&attacks).Error

	if err != nil {
		return nil, err
	}
	return attacks, err
}

func SetAttackProgressString(attackId string, progressString string) error {
	return GetInstance().
		Table("attacks").
		Where("id = ?", attackId).
		Update("progress_string", progressString).Error
}
