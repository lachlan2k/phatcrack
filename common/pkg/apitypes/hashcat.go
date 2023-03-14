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

type VerifyHashesRequestDTO struct {
	Hashes   []string `json:"hashes"`
	HashType uint     `json:"hash_type"`
}

type VerifyHashesResponseDTO struct {
	Valid bool `json:"valid"`
}

type NormalizeHashesRequestDTO = VerifyHashesRequestDTO

type NormalizeHashesResponseDTO struct {
	Valid            bool     `json:"valid"`
	NormalizedHashes []string `json:"normalized_hashes"`
}
