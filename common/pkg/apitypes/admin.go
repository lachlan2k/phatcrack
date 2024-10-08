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
	Name              string `json:"name" validate:"max=64,username"`
	ForEphemeralAgent bool   `json:"for_ephemeral_agent"`
}

type AdminAgentRegistrationKeyCreateResponseDTO struct {
	Name              string `json:"name"`
	ID                string `json:"id"`
	Key               string `json:"key"`
	ForEphemeralAgent bool   `json:"for_ephemeral_agent"`
}

type AdminGetAgentRegistrationKeyDTO struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	KeyHint           string `json:"key_hint"`
	ForEphemeralAgent bool   `json:"for_ephemeral_agent"`
}

type AdminGetAllAgentRegistrationKeysResponseDTO struct {
	AgentRegistrationKeys []AdminGetAgentRegistrationKeyDTO `json:"agent_registration_keys"`
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

type AdminUserUpdatePasswordRequestDTO struct {
	Action string `json:"action"`
}
type AdminUserUpdatePasswordResponseDTO struct {
	GeneratedPassword string `json:"generated_password"`
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
	Users []AdminGetUserDTO `json:"users"`
}

type AdminGetUserDTO struct {
	ID               string   `json:"id"`
	Username         string   `json:"username"`
	Roles            []string `json:"roles"`
	IsPasswordLocked bool     `json:"is_password_locked"`
}

type AdminAgentSetMaintanceRequestDTO struct {
	IsMaintenanceMode bool `json:"is_maintenance_mode"`
}
