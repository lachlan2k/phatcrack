package apitypes

type ListsWordlistCreateDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
	Filename    string `json:"filename" validate:"required,min=5"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type ListsRuleFileCreateDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
	Filename    string `json:"filename" validate:"required,min=5"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type ListsWordlistResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type ListsRuleFileResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type ListsGetAllWordlistsDTO struct {
	Wordlists []ListsWordlistResponseDTO `json:"wordlists"`
}

type ListsGetAllRuleFilesDTO struct {
	RuleFiles []ListsRuleFileResponseDTO `json:"rulefiles"`
}
