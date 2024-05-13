package apitypes

type PublicOIDCConfigDTO struct {
	Prompt string `json:"prompt"`
}

type PublicAuthConfigDTO struct {
	EnabledMethods []string            `json:"enabled_methods"`
	OIDC           PublicOIDCConfigDTO `json:"oidc"`
}

type PublicGeneralConfigDTO struct {
	IsMaintenanceMode               bool  `json:"is_maintenance_mode"`
	MaximumUploadedFileSize         int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize int64 `json:"maximum_uploaded_file_line_scan_size"`
}

type PublicConfigDTO struct {
	Auth    PublicAuthConfigDTO    `json:"auth"`
	General PublicGeneralConfigDTO `json:"general"`
}
