package apitypes

type AdminAgentCreateRequestDTO struct {
	Name      string `json:"name" validate:"required,min=4,max=64,username"`
	Ephemeral bool   `json:"ephemeral"`
}

type AdminAgentCreateResponseDTO struct {
	Ephemeral bool   `json:"ephemeral"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Key       string `json:"key"`
}

type AdminAgentRegistrationKeyCreateRequestDTO struct {
	Name      string `json:"name" validate:"required,min=4,max=64,username"`
	Ephemeral bool   `json:"ephemeral"`
}

type AdminAgentRegistrationKeyCreateResponseDTO struct {
	Ephemeral bool   `json:"ephemeral"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Key       string `json:"key"`
}

type AdminUserCreateRequestDTO struct {
	Username     string   `json:"username" validate:"required,min=4,max=64,username"`
	Password     string   `json:"password" validate:"max=128"`
	GenPassword  bool     `json:"gen_password"`
	LockPassword bool     `json:"lock_password"`
	Roles        []string `json:"roles" validate:"required,userroles"`
}

type AdminUserCreateResponseDTO struct {
	Username          string   `json:"username"`
	ID                string   `json:"id"`
	Roles             []string `json:"roles"`
	GeneratedPassword string   `json:"generated_password"`
}

type AdminUserUpdateRequestDTO struct {
	Username string   `json:"username" validate:"required,min=4,max=64,username"`
	Roles    []string `json:"roles" validate:"required,userroles"`
}

type AdminServiceAccountCreateRequestDTO struct {
	Username string   `json:"username" validate:"required,min=4,max=64,username"`
	Roles    []string `json:"roles" validate:"required,userroles"`
}
type AdminServiceAccountCreateResponseDTO struct {
	Username string   `json:"username"`
	ID       string   `json:"id"`
	Roles    []string `json:"roles"`
	APIKey   string   `json:"api_key"`
}

type AdminGetAllUsersResponseDTO struct {
	Users []UserDTO `json:"users"`
}

type AdminAgentSetMaintanceRequestDTO struct {
	IsMaintenanceMode bool `json:"is_maintenance_mode"`
}
