package db

import (
	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"gorm.io/datatypes"
)

type AttackTemplate struct {
	UUIDBaseModel

	Name          string
	HashcatParams datatypes.JSONType[hashcattypes.HashcatParams]

	CreatedByUser   User      `gorm:"constraint:OnDelete:SET NULL;"`
	CreatedByUserID uuid.UUID `gorm:"type:uuid"`
}

type AttackTemplateSet struct {
	UUIDBaseModel

	Name              string
	AttackTemplateIDs datatypes.JSONSlice[string]

	CreatedByUser   User      `gorm:"constraint:OnDelete:SET NULL;"`
	CreatedByUserID uuid.UUID `gorm:"type:uuid"`
}

const AttackTemplateType = "attack-template"
const AttackTemplateSetType = "attack-template-set"

func (at AttackTemplate) ToDTO() apitypes.AttackTemplateDTO {
	params := at.HashcatParams.Data()

	return apitypes.AttackTemplateDTO{
		ID: at.ID.String(),

		Type:              AttackTemplateType,
		Name:              at.Name,
		HashcatParams:     &params,
		AttackTemplateIDs: nil,

		CreatedByUserID: at.CreatedByUserID.String(),
	}
}

func (at AttackTemplateSet) ToDTO() apitypes.AttackTemplateDTO {
	var ids []string = at.AttackTemplateIDs

	return apitypes.AttackTemplateDTO{
		ID: at.ID.String(),

		Type:              AttackTemplateSetType,
		Name:              at.Name,
		HashcatParams:     nil,
		AttackTemplateIDs: ids,

		CreatedByUserID: at.CreatedByUserID.String(),
	}
}

func CreateAttackTemplate(attackTemplate *AttackTemplate) (*AttackTemplate, error) {
	return attackTemplate, GetInstance().Create(attackTemplate).Error
}

func CreateAttackTemplateSet(templateSet *AttackTemplateSet) (*AttackTemplateSet, error) {
	return templateSet, GetInstance().Create(templateSet).Error
}

func GetAllAttackTemplates() ([]AttackTemplate, error) {
	return GetAll[AttackTemplate]()
}

func GetAllAttackTemplateSets() ([]AttackTemplateSet, error) {
	return GetAll[AttackTemplateSet]()
}

func GetAttackTemplate(id string) (*AttackTemplate, error) {
	return GetByID[AttackTemplate](id)
}

func GetAttackTemplateSet(id string) (*AttackTemplateSet, error) {
	return GetByID[AttackTemplateSet](id)
}

func DeleteAttackTemplate(id string) error {
	return GetInstance().Unscoped().Delete(&AttackTemplate{}, "id = ?",id).Error
}

func DeleteAttackTemplateSet(id string) error {
	return GetInstance().Unscoped().Delete(&AttackTemplateSet{}, "id = ?", id).Error
}
