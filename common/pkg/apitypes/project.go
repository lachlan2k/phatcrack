package apitypes

type ProjectCreateDTO struct {
	Name        string `json:"name" validate:"required,printascii,min=4,max=64"`
	Description string `json:"description" validate:"required,printascii,max=1000"`
}

type ProjectDTO struct {
	ID          string `json:"id"`
	TimeCreated int64  `json:"time_created"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerUserID string `json:"owner_user_id"`
}

type ProjectResponseMultipleDTO struct {
	Projects []ProjectDTO `json:"projects"`
}

type ProjectCreateResponseDTO = ProjectDTO
