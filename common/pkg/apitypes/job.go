package apitypes

type HashcatParamsDTO struct {
	AttackMode        uint8    `json:"attack_mode"`
	HashType          uint     `json:"hash_type"`
	Mask              string   `json:"mask"`
	WordlistFilenames []string `json:"wordlist_filenames"`
	RulesFilenames    []string `json:"rules_filenames"`
	AdditionalArgs    []string `json:"additional_args"`
	OptimizedKernels  bool     `json:"optimized_kernels"`
}

type JobStartRequestDTO struct {
	HashcatParams    HashcatParamsDTO `json:"hashcat_params"`
	Hashes           []string         `json:"hashes"`
	StartImmediately bool             `json:"start_immediately"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
}
