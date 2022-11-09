package apitypes

import "github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"

type HashTypesDTO struct {
	HashTypes hashcattypes.HashTypeMap `json:"hashtypes"`
}
