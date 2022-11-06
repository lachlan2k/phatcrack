package apitypes

type AgentCreateRequestDTO struct {
	Name string `json:"name" validate:"required,min=5,max=30"`
}

type AgentCreateResponseDTO struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}

type UserCreateRequestDTO struct {
	Username string `json:"username" validate:"required,min=4,max=64"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Role     string `json:"role" validate:"required,userrole"`
}

type UserCreateResponseDTO struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	Role     string `json:"role"`
}
