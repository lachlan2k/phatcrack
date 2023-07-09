package apitypes

type AdminAgentCreateRequestDTO struct {
	Name string `json:"name" validate:"required,min=5,max=30"`
}

type AdminAgentCreateResponseDTO struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}

type AdminUserCreateRequestDTO struct {
	Username string   `json:"username" validate:"required,min=4,max=64"`
	Password string   `json:"password" validate:"required,min=8,max=128"`
	Roles    []string `json:"roles" validate:"required,userroles"`
}

type AdminUserCreateResponseDTO struct {
	Username string   `json:"username"`
	ID       string   `json:"id"`
	Roles    []string `json:"roles"`
}

type AdminIsSetupCompleteResponseDTO struct {
	IsComplete bool `json:"is_complete"`
}

type AdminConfigResponseDTO struct {
	IsSetupComplete                   bool `json:"is_setup_complete"`
	IsMFARequired                     bool `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool `json:"require_password_change_on_first_login"`
}

type AdminConfigRequestDTO struct {
	IsMFARequired                     bool `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool `json:"require_password_change_on_first_login"`
}

type AdminGetAllUsersResponseDTO struct {
	Users []UserDTO `json:"users"`
}
