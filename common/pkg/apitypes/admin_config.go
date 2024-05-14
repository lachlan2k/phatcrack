package apitypes

type AuthOIDCConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	IssuerURL   string `json:"issuer_url"`
	RedirectURL string `json:"redirect_url"`

	AutomaticUserCreation bool   `json:"automatic_creation"`
	UsernameClaim         string `json:"username_field"`
	Prompt                string `json:"prompt"`

	RolesClaim   string `json:"role_field"`
	RequiredRole string `json:"required_role"`

	AdditionalScopes []string `json:"scopes"`
}

type GeneralAuthConfig struct {
	EnabledMethods []string `json:"enabled_methods"`

	IsMFARequired                     bool `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool `json:"require_password_change_on_first_login"`
}

type AuthConfig struct {
	General *GeneralAuthConfig `json:"general"`
	OIDC    *AuthOIDCConfig    `json:"oidc"`
}

type AgentConfig struct {
	AutomaticallySyncListfiles bool `json:"auto_sync_listfiles"`
	SplitJobsPerAgent          int  `json:"split_jobs_per_agent"`
}

type GeneralConfig struct {
	IsMaintenanceMode               bool  `json:"is_maintenance_mode"`
	MaximumUploadedFileSize         int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize int64 `json:"maximum_uploaded_file_line_scan_size"`
}

type AdminConfigRequestDTO struct {
	Auth    *AuthConfig    `json:"auth"`
	Agent   *AgentConfig   `json:"agent"`
	General *GeneralConfig `json:"general"`
}

type AdminConfigResponseDTO struct {
	Auth    AuthConfig    `json:"auth"`
	Agent   AgentConfig   `json:"agent"`
	General GeneralConfig `json:"general"`
}
