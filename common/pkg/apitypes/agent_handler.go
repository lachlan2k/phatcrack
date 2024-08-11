package apitypes

type AgentRegisterRequestDTO struct {
	Name string `json:"name" validate:"max=64"`
}

type AgentRegisterResponseDTO struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}
