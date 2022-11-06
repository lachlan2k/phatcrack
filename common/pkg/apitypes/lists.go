package apitypes

type WordlistCreateDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
	Filename    string `json:"filename" validate:"required,min=5"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type RuleFileCreateDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
	Filename    string `json:"filename" validate:"required,min=5"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type WordlistResponseDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
	Filename    string `json:"filename" validate:"required,min=5"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type RuleFileResponseDTO struct {
	Name        string `json:"name" validate:"required,min=4,max=64"`
	Description string `json:"description" validate:"required,max=1000"`
	Filename    string `json:"filename" validate:"required,min=5"`
	Size        uint64 `json:"size"`
	Lines       uint64 `json:"lines"`
}

type GetAllWordlistsDTO struct {
	Wordlists []WordlistResponseDTO `json:"wordlists"`
}

type GetAllRuleFilesDTO struct {
	RuleFiles []RuleFileResponseDTO `json:"rulefiles"`
}
