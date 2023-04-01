package apitypes

type ProjectCreateDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
}

type ProjectSimpleDetailsDTO struct {
	ID          string `json:"id"`
	TimeCreated int64  `json:"time_created"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Hashlist
type ProjectHashlistCreateDTO struct {
	ProjectID    string   `json:"project_id"`
	Name         string   `json:"name"`
	HashType     uint     `json:"hash_type"`
	InputHashes  []string `json:"input_hashes"`
	HasUsernames bool     `json:"has_usernames"`
}

type ProjectHashlistCreateResponseDTO = struct {
	ID string `json:"id"`
}

type ProjectHashlistHashDTO struct {
	InputHash      string `json:"input_hash"`
	NormalizedHash string `json:"normalized_hash"`
}

type ProjectHashlistDTO struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	TimeCreated int64                    `json:"time_created"`
	HashType    uint                     `json:"hash_type"`
	Hashes      []ProjectHashlistHashDTO `json:"hashes"`
	Version     uint                     `json:"version"`
}

type ProjectsFullDetailsDTO struct {
	ID          string `json:"id"`
	TimeCreated int64  `json:"time_created"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectResponseMultipleDTO struct {
	Projects []ProjectSimpleDetailsDTO `json:"projects"`
}

type ProjectCreateResponseDTO = ProjectSimpleDetailsDTO
