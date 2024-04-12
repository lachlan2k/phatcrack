package config

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

var lock sync.Mutex

const AuthMethodCredentials = "auth_method_credentials"
const AuthMethodOIDC = "auth_method_oidc"

type AuthConfig struct {
	EnabledMethods []string `json:"auth_enabled_methods"`

	IsMFARequired                     bool `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool `json:"require_password_change_on_first_login"`

	OIDCClientID     string `json:"auth_oidc_client_id"`
	OIDCClientSecret string `json:"auth_oidc_client_secret"`
	OIDCEndpoint     string `json:"auth_oidc_endpoint"`
	OIDCRedirectURL  string `json:"auth_oidc_redirect_url"`
	OIDCScopes       string `json:"auth_oidc_scopes"`
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

type RuntimeConfig struct {
	ConfigVersion int `json:"version"`

	IsSetupComplete bool `json:"is_setup_complete"`

	Auth    AuthConfig    `json:"auth"`
	Agent   AgentConfig   `json:"agent"`
	General GeneralConfig `json:"general"`
}

const latestConfigVersion = 2

func (conf RuntimeConfig) ToAdminDTO() apitypes.AdminConfigResponseDTO {
	return apitypes.AdminConfigResponseDTO{
		IsSetupComplete:   conf.IsSetupComplete,
		IsMFARequired:     conf.Auth.IsMFARequired,
		IsMaintenanceMode: conf.General.IsMaintenanceMode,

		AutomaticallySyncListfiles:        conf.Agent.AutomaticallySyncListfiles,
		SplitJobsPerAgent:                 conf.Agent.SplitJobsPerAgent,
		RequirePasswordChangeOnFirstLogin: conf.Auth.RequirePasswordChangeOnFirstLogin,

		MaximumUploadedFileSize:         conf.General.MaximumUploadedFileSize,
		MaximumUploadedFileLineScanSize: conf.General.MaximumUploadedFileLineScanSize,
	}
}

func (conf RuntimeConfig) ToPublicDTO() apitypes.ConfigDTO {
	return apitypes.ConfigDTO{
		IsMaintenanceMode:               conf.General.IsMaintenanceMode,
		MaximumUploadedFileSize:         conf.General.MaximumUploadedFileSize,
		MaximumUploadedFileLineScanSize: conf.General.MaximumUploadedFileLineScanSize,
	}
}

var runningConf RuntimeConfig

func MakeDefaultConfig() RuntimeConfig {
	return RuntimeConfig{
		IsSetupComplete: true,
		ConfigVersion:   latestConfigVersion,

		Auth: AuthConfig{
			EnabledMethods: []string{AuthMethodCredentials},

			IsMFARequired:                     false,
			RequirePasswordChangeOnFirstLogin: true,
		},

		Agent: AgentConfig{
			AutomaticallySyncListfiles: true,
			SplitJobsPerAgent:          1,
		},

		General: GeneralConfig{
			IsMaintenanceMode:               false,
			MaximumUploadedFileSize:         10 * 1000 * 1000 * 1000, // 10GB,
			MaximumUploadedFileLineScanSize: 500 * 1000 * 1000,       // 100MB
		},
	}
}

func Reload() error {
	lock.Lock()
	defer lock.Unlock()

	var newConf *RuntimeConfig

	newConf, err := db.GetConfig[RuntimeConfig]()
	if err == db.ErrNotFound {
		err = db.SeedConfig[RuntimeConfig](MakeDefaultConfig())
		if err != nil {
			return err
		}

		newConf, err = db.GetConfig[RuntimeConfig]()
		if err != nil {
			return errors.New("failed to fetch configuration even after seeding: " + err.Error())
		}
	}

	if err != nil {
		return err
	}

	if newConf.ConfigVersion < latestConfigVersion {
		configJsonStr, err := db.GetConfigJSONString()
		if err != nil {
			return fmt.Errorf("failed to get config string for migrations: %v", err)
		}

		migratedConfig, err := runMigrations(configJsonStr)
		if err != nil {
			return fmt.Errorf("failed to perform config migrations: %v", err)
		}

		newConf = migratedConfig
	}

	runningConf = *newConf
	return nil
}

func save(conf RuntimeConfig) error {
	return db.SetConfig[RuntimeConfig](conf)
}

func Save() error {
	lock.Lock()
	defer lock.Unlock()

	return save(runningConf)
}

func Update(updateFunc func(*RuntimeConfig) error) error {
	lock.Lock()
	defer lock.Unlock()

	var newConf RuntimeConfig = runningConf
	err := updateFunc(&newConf)
	if err != nil {
		return err
	}

	err = save(newConf)
	if err != nil {
		return fmt.Errorf("failed to save new config in db: %w", err)
	}

	runningConf = newConf
	return nil
}

func Get() RuntimeConfig {
	return runningConf
}
