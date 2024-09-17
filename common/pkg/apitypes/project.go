package apitypes

type ProjectCreateRequestDTO struct {
	Name        string `json:"name" validate:"required,standardname,min=3,max=64"`
	Description string `json:"description" validate:"printascii,max=1000"`
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

type ProjectAddShareRequestDTO struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type ProjectSharesDTO struct {
	UserIDs []string `json:"user_ids"`
}
