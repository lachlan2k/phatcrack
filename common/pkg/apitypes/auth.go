package apitypes

type LoginRequestDTO struct {
	Username string `json:"username" validate:"required,min=4,max=64"`
	Password string `json:"password" validate:"required"`
}

type UserMeDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginResponseDTO struct {
	User UserMeDTO `json:"user"`
}

type WhoamiResponseDTO struct {
	IsLoggedIn bool       `json:"is_logged_in"`
	User       *UserMeDTO `json:"user,omitempty"`
}
