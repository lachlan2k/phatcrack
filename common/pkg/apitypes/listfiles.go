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

type WordlistDTO struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	FilenameOnDisk string `json:"filename_on_disk"`
	SizeInBytes    uint64 `json:"size_in_bytes"`
	Lines          uint64 `json:"lines"`
}

type RuleFileDTO struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	FilenameOnDisk string `json:"filename_on_disk"`
	SizeInBytes    uint64 `json:"size_in_bytes"`
	Lines          uint64 `json:"lines"`
}

type GetAllWordlistsDTO struct {
	Wordlists []WordlistDTO `json:"wordlists"`
}

type GetAllRuleFilesDTO struct {
	RuleFiles []RuleFileDTO `json:"rulefiles"`
}
