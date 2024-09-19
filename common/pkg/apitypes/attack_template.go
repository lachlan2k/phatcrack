package apitypes

import (
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type AttackTemplateDTO struct {
	ID string `json:"id"`

	Type          string
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	HashcatParams *hashcattypes.HashcatParams `json:"hashcat_params,omitempty"`

	CreatedByUserID string `json:"created_by_user_id"`

	AttackTemplateIDs []string `json:"attack_template_ids,omitempty"`
}

type AttackTemplateGetAllResponseDTO struct {
	AttackTemplates []AttackTemplateDTO `json:"attack_templates"`
}

type AttackTemplateCreateRequestDTO struct {
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	HashcatParams hashcattypes.HashcatParams `json:"hashcat_params" validate:"required"`
}

type AttaackTemplateCreateSetRequestDTO struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	AttackTemplateIDs []string `json:"attack_template_ids"`
}
