package apitypes

type ListfileDTO struct {
	ID              string `json:"id"`
	FileType        string `json:"file_type"`
	Name            string `json:"name"`
	SizeInBytes     uint64 `json:"size_in_bytes"`
	Lines           uint64 `json:"lines"`
	IsLocked        bool   `json:"is_locked"`
	AvailableForUse bool   `json:"available_for_use"`
}

type GetAllWordlistsDTO struct {
	Wordlists []ListfileDTO `json:"wordlists"`
}

type GetAllRuleFilesDTO struct {
	RuleFiles []ListfileDTO `json:"rulefiles"`
}

type ListfileUploadResponseDTO struct {
	Listfile ListfileDTO `json:"listfile"`
}
