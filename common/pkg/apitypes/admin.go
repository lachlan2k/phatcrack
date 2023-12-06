package apitypes

type AdminAgentCreateRequestDTO struct {
	Name string `json:"name" validate:"required,min=4,max=64,username"`
}

type AdminAgentCreateResponseDTO struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}

type AdminUserCreateRequestDTO struct {
	Username    string   `json:"username" validate:"required,min=4,max=64,username"`
	Password    string   `json:"password" validate:"max=128"`
	GenPassword bool     `json:"gen_password"`
	Roles       []string `json:"roles" validate:"required,userroles"`
}

type AdminUserCreateResponseDTO struct {
	Username          string   `json:"username"`
	ID                string   `json:"id"`
	Roles             []string `json:"roles"`
	GeneratedPassword string   `json:"generated_password"`
}

type AdminServiceAccountCreateRequestDTO struct {
	Username string   `json:"username" validate:"required,min=4,max=64,username"`
	Roles    []string `json:"roles" validate:"required,userroles"`
}
type AdminServiceAccountCreateResponseDTO struct {
	Username string   `json:"username"`
	ID       string   `json:"id"`
	Roles    []string `json:"roles"`
	APIKey   string   `json:"api_key"`
}

type AdminIsSetupCompleteResponseDTO struct {
	IsComplete bool `json:"is_complete"`
}

type AdminConfigResponseDTO struct {
	SplitJobsPerAgent                 int   `json:"split_jobs_per_agent"`
	IsSetupComplete                   bool  `json:"is_setup_complete"`
	IsMFARequired                     bool  `json:"is_mfa_required"`
	IsMaintenanceMode                 bool  `json:"is_maintenance_mode"`
	AutomaticallySyncListfiles        bool  `json:"auto_sync_listfiles"`
	RequirePasswordChangeOnFirstLogin bool  `json:"require_password_change_on_first_login"`
	MaximumUploadedFileSize           int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize   int64 `json:"maximum_uploaded_file_line_scan_size"`
}

type AdminConfigRequestDTO struct {
	SplitJobsPerAgent                 int   `json:"split_jobs_per_agent" validate:"min=1,max=4"`
	IsMFARequired                     bool  `json:"is_mfa_required"`
	IsMaintenanceMode                 bool  `json:"is_maintenance_mode"`
	AutomaticallySyncListfiles        bool  `json:"auto_sync_listfiles"`
	RequirePasswordChangeOnFirstLogin bool  `json:"require_password_change_on_first_login"`
	MaximumUploadedFileSize           int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize   int64 `json:"maximum_uploaded_file_line_scan_size"`
}

type AdminGetAllUsersResponseDTO struct {
	Users []UserDTO `json:"users"`
}

type AdminAgentSetMaintanceRequestDTO struct {
	IsMaintenanceMode bool `json:"is_maintenance_mode"`
}
