package apitypes

type AuthOIDCConfigDTO struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	IssuerURL   string `json:"issuer_url"`
	RedirectURL string `json:"redirect_url"`

	Prompt string `json:"prompt"`

	AutomaticUserCreation bool   `json:"automatic_user_creation"`
	UsernameClaim         string `json:"username_claim"`

	RolesClaim   string `json:"role_field"`
	RequiredRole string `json:"required_role"`

	AdditionalScopes []string `json:"scopes"`
}

type GeneralAuthConfigDTO struct {
	EnabledMethods []string `json:"enabled_methods"`

	IsMFARequired                     bool `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool `json:"require_password_change_on_first_login"`
}

type AuthConfigDTO struct {
	General *GeneralAuthConfigDTO `json:"general"`
	OIDC    *AuthOIDCConfigDTO    `json:"oidc"`
}

type AgentConfigDTO struct {
	AutomaticallySyncListfiles bool `json:"auto_sync_listfiles"`
	SplitJobsPerAgent          int  `json:"split_jobs_per_agent"`
}

type GeneralConfigDTO struct {
	IsMaintenanceMode               bool  `json:"is_maintenance_mode"`
	MaximumUploadedFileSize         int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize int64 `json:"maximum_uploaded_file_line_scan_size"`
}

type AdminConfigRequestDTO struct {
	Auth    *AuthConfigDTO    `json:"auth"`
	Agent   *AgentConfigDTO   `json:"agent"`
	General *GeneralConfigDTO `json:"general"`
}

type AdminConfigResponseDTO struct {
	Auth    AuthConfigDTO    `json:"auth"`
	Agent   AgentConfigDTO   `json:"agent"`
	General GeneralConfigDTO `json:"general"`
}
