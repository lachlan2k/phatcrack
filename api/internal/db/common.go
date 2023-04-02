package db

import (
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

type HashcatParams struct {
	AttackMode        uint8    `bson:"attack_mode"`
	HashType          uint     `bson:"hash_type"`
	Mask              string   `bson:"mask"`
	WordlistFilenames []string `bson:"wordlist_filenames"`
	RulesFilenames    []string `bson:"rules_filenames"`
	AdditionalArgs    []string `bson:"additional_args"`
	OptimizedKernels  bool     `bson:"optimized_kernels"`
	SlowCandidates    bool     `bson:"slow_candidates"`
}

func (h *HashcatParams) ToDTO() hashcattypes.HashcatParams {
	return hashcattypes.HashcatParams{
		AttackMode:        h.AttackMode,
		HashType:          h.HashType,
		Mask:              h.Mask,
		WordlistFilenames: h.WordlistFilenames,
		RulesFilenames:    h.RulesFilenames,
		AdditionalArgs:    h.AdditionalArgs,
		OptimizedKernels:  h.OptimizedKernels,
		SlowCandidates:    h.SlowCandidates,
	}
}

type HashlistHash struct {
	InputHash      string `bson:"input_hash"`
	NormalizedHash string `bson:"normalized_hash"`
}

func (h *HashlistHash) ToDTO() apitypes.HashlistHashDTO {
	return apitypes.HashlistHashDTO{
		InputHash:      h.InputHash,
		NormalizedHash: h.NormalizedHash,
	}
}
