package apitypes

type LoginRequestDTO struct {
	Username string `json:"username" validate:"required,min=4,max=64"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
