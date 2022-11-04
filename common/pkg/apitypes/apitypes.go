package apitypes

type AgentCreateRequestDTO struct {
	Name string `json:"name" validate:"required,min=5,max=30"`
}

type AgentCreateResponseDTO struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}
