package apitypes

type HashlistCreateRequestDTO struct {
	ProjectID    string   `json:"project_id"`
	Name         string   `json:"name"`
	HashType     uint     `json:"hash_type"`
	InputHashes  []string `json:"input_hashes"`
	HasUsernames bool     `json:"has_usernames"`
}

type HashlistCreateResponseDTO struct {
	ID string `json:"id"`
}

type HashlistHashDTO struct {
	InputHash      string `json:"input_hash"`
	NormalizedHash string `json:"normalized_hash"`
	IsCracked      bool   `json:"is_cracked"`
	PlaintextHex   string `json:"plaintext_hex"`
}

type HashlistDTO struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Name        string            `json:"name"`
	TimeCreated int64             `json:"time_created"`
	HashType    uint              `json:"hash_type"`
	Hashes      []HashlistHashDTO `json:"hashes"`
	Version     uint              `json:"version"`
}

type HashlistResponseMultipleDTO struct {
	Hashlists []HashlistDTO `json:"hashlists"`
}
