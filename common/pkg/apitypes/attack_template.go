package apitypes

import (
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type AttackTemplateDTO struct {
	ID string `json:"id"`

	Type          string                      `json:"type"`
	Name          string                      `json:"name"`
	HashcatParams *hashcattypes.HashcatParams `json:"hashcat_params,omitempty"`

	CreatedByUserID string `json:"created_by_user_id"`

	AttackTemplateIDs []string `json:"attack_template_ids,omitempty"`
}

type AttackTemplateGetAllResponseDTO struct {
	AttackTemplates []AttackTemplateDTO `json:"attack_templates"`
}

type AttackTemplateCreateRequestDTO struct {
	Name          string                     `json:"name" validate:"required,standardname,min=3,max=64"`
	HashcatParams hashcattypes.HashcatParams `json:"hashcat_params" validate:"required"`
}

type AttackTemplateCreateSetRequestDTO struct {
	Name              string   `json:"name" validate:"required,standardname,min=3,max=64"`
	AttackTemplateIDs []string `json:"attack_template_ids" validate:"min=1,max=32,dive,uuid4"`
}

type AttackTemplateUpdateRequestDTO struct {
	Type string `json:"type"`
	Name string `json:"name"`

	HashcatParams     *hashcattypes.HashcatParams `json:"hashcat_params,omitempty"`
	AttackTemplateIDs []string                    `json:"attack_template_ids,omitempty"`
}
