package apitypes

type AuthLoginRequestDTO struct {
	Username string `json:"username" validate:"required,min=4,max=64"`
	Password string `json:"password" validate:"required"`
}

type AuthCurrentUserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AuthLoginResponseDTO struct {
	User AuthCurrentUserDTO `json:"user"`
}

type AuthWhoamiResponseDTO struct {
	User AuthCurrentUserDTO `json:"user"`
}

type AuthRefreshResponseDTO struct {
	User AuthCurrentUserDTO `json:"user"`
}
