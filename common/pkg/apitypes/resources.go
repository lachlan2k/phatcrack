package apitypes

import "github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"

type HashTypesDTO struct {
	HashTypes hashcattypes.HashTypeMap `json:"hashtypes"`
}

type DetectHashTypeRequestDTO struct {
	TestHash string `json:"test_hash"`
}

type DetectHashTypeResponseDTO struct {
	PossibleTypes []int `json:"possible_types"`
}
