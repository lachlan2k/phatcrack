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
