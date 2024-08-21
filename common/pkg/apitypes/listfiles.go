package apitypes

type ListfileDTO struct {
	ID                string `json:"id"`
	FileType          string `json:"file_type"`
	Name              string `json:"name"`
	SizeInBytes       uint64 `json:"size_in_bytes"`
	Lines             uint64 `json:"lines"`
	AvailableForUse   bool   `json:"available_for_use"`
	PendingDelete     bool   `json:"pending_delete"`
	CreatedByUserID   string `json:"created_by_user_id"`
	AttachedProjectID string `json:"associated_project_id"`
}

type GetAllWordlistsDTO struct {
	Wordlists []ListfileDTO `json:"wordlists"`
}

type GetAllRuleFilesDTO struct {
	RuleFiles []ListfileDTO `json:"rulefiles"`
}

type GetAllListfilesDTO struct {
	Listfiles []ListfileDTO `json:"listfiles"`
}

type ListfileUploadResponseDTO struct {
	Listfile ListfileDTO `json:"listfile"`
}
