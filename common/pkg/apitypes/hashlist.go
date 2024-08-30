package apitypes

type HashlistCreateRequestDTO struct {
	ProjectID    string   `json:"project_id" validate:"required,uuid"`
	Name         string   `json:"name" validate:"required,standardname,min=5,max=30"`
	HashType     int      `json:"hash_type" validate:"hashtype"`
	InputHashes  []string `json:"input_hashes" validate:"required,min=1,dive,required,min=4"`
	HasUsernames bool     `json:"has_usernames"`
}

type HashlistCreateResponseDTO struct {
	ID                      string `json:"id"`
	NumPopulatedFromPotfile int64  `json:"num_populated_from_potfile"`
}

type HashlistHashDTO struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	InputHash      string `json:"input_hash"`
	NormalizedHash string `json:"normalized_hash"`
	IsCracked      bool   `json:"is_cracked"`
	IsUnexpected   bool   `json:"is_unexpected"`
	PlaintextHex   string `json:"plaintext_hex"`
}

type HashlistDTO struct {
	ID           string            `json:"id"`
	ProjectID    string            `json:"project_id"`
	Name         string            `json:"name"`
	TimeCreated  int64             `json:"time_created"`
	HashType     int               `json:"hash_type"`
	Hashes       []HashlistHashDTO `json:"hashes"`
	Version      uint              `json:"version"`
	HasUsernames bool              `json:"has_usernames"`
}

type HashlistResponseMultipleDTO struct {
	Hashlists []HashlistDTO `json:"hashlists"`
}
