package db

import (
	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

const (
	ListfileTypeWordlist = "Wordlist"
	ListfileTypeRulefile = "Rulefile"
	ListfileTypeHashlist = "Hashlist"
	ListfileTypeCharset  = "Charset"
)

type Listfile struct {
	UUIDBaseModel
	Name                 string
	AvailableForDownload bool
	AvailableForUse      bool
	FileType             string
	SizeInBytes          uint64
	Lines                uint64
	PendingDelete        bool
	CreatedByUser        User      `gorm:"constraint:OnDelete:SET NULL;"`
	CreatedByUserID      uuid.UUID `gorm:"type:uuid"`

	AttachedProjectID *uuid.UUID `gorm:"type:uuid"`
	AttachedProject   *Project   `gorm:"constraint:OnDelete:SET NULL;"`
}

func (l *Listfile) Save() error {
	return GetInstance().Save(l).Error
}

func (w *Listfile) ToDTO() apitypes.ListfileDTO {
	projId := ""
	if w.AttachedProjectID != nil {
		projId = w.AttachedProjectID.String()
	}

	return apitypes.ListfileDTO{
		ID:                w.ID.String(),
		Name:              w.Name,
		FileType:          w.FileType,
		SizeInBytes:       w.SizeInBytes,
		PendingDelete:     w.PendingDelete,
		Lines:             w.Lines,
		AvailableForUse:   w.AvailableForUse,
		CreatedByUserID:   w.CreatedByUserID.String(),
		AttachedProjectID: projId,
	}
}

func GetListfile(id string) (*Listfile, error) {
	var listfile Listfile
	err := GetInstance().First(&listfile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &listfile, nil
}

func CreateListfile(listfile *Listfile) (*Listfile, error) {
	return listfile, GetInstance().Create(listfile).Error
}

func MarkListfileAsAvailable(id string) error {
	return GetInstance().Model(&Listfile{}).Where("id = ?", id).Updates(&Listfile{AvailableForDownload: true}).Error
}

func MarkListfileForDeletion(id string) error {
	return GetInstance().Model(&Listfile{}).Where("id = ?", id).Updates(&Listfile{PendingDelete: true}).Error
}

// A "Public" listfile is a listfile which is not attached to a specific project
func GetAllPublicRulefiles() ([]Listfile, error) {
	rulefiles := []Listfile{}
	err := GetInstance().Where("file_type = ? and attached_project_id is NULL", ListfileTypeRulefile).Find(&rulefiles).Error
	if err != nil {
		return nil, err
	}
	return rulefiles, nil
}

func GetAllPublicWordlists() ([]Listfile, error) {
	wordlists := []Listfile{}
	err := GetInstance().Where("file_type = ? and attached_project_id is NULL", ListfileTypeWordlist).Find(&wordlists).Error
	if err != nil {
		return nil, err
	}
	return wordlists, nil
}

func GetAllListfilesAvailableToProject(projectID string) ([]Listfile, error) {
	listfiles := []Listfile{}
	err := GetInstance().Where("attached_project_id = ? or attached_project_id is NULL", projectID).Find(&listfiles).Error
	if err != nil {
		return nil, err
	}
	return listfiles, nil
}

func GetAllListfiles() ([]Listfile, error) {
	listfiles := []Listfile{}
	err := GetInstance().Find(&listfiles).Error
	if err != nil {
		return nil, err
	}
	return listfiles, nil
}

func GetAllProjectSpecificListfiles(projectID string) ([]Listfile, error) {
	listfiles := []Listfile{}
	err := GetInstance().Where("attached_project_id = ?", projectID).Find(&listfiles).Error
	if err != nil {
		return nil, err
	}
	return listfiles, nil
}
