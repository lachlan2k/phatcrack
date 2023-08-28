package apitypes

type AccountChangePasswordRequestDTO struct {
	NewPassword     string `json:"new_password" validate:"required,min=16,max=128"`
	CurrentPassword string `json:"current_password" validate:"required"`
}
