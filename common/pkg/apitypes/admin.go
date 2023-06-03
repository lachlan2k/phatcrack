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
	Roles    []string `json:"roles" validate:"required,userrole"`
}

type AdminUserCreateResponseDTO struct {
	Username string   `json:"username"`
	ID       string   `json:"id"`
	Roles    []string `json:"roles"`
}

type AdminIsSetupCompleteResponseDTO struct {
	IsComplete bool `json:"is_complete"`
}
