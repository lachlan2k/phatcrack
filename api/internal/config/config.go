package config

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

var lock sync.Mutex

const AuthMethodCredentials = "method_credentials"
const AuthMethodOIDC = "method_oidc"

type AuthOIDCConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	IssuerURL   string `json:"issuer_url"`
	RedirectURL string `json:"redirect_url"`

	AutomaticUserCreation bool   `json:"automatic_user_creation"`
	UsernameClaim         string `json:"username_claim"`

	RolesClaim   string `json:"role_field"`
	RequiredRole string `json:"required_role"`
	Prompt       string `json:"prompt"`

	AdditionalScopes []string `json:"scopes"`
}

type AuthGeneralConfig struct {
	EnabledMethods                    []string `json:"enabled_methods"`
	IsMFARequired                     bool     `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool     `json:"require_password_change_on_first_login"`
}

type AuthConfig struct {
	General AuthGeneralConfig `json:"general"`
	OIDC    AuthOIDCConfig    `json:"oidc"`
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
	oidcClientSecret := ""
	if len(conf.Auth.OIDC.ClientSecret) > 0 {
		oidcClientSecret = "redacted"
	}

	return apitypes.AdminConfigResponseDTO{
		Auth: apitypes.AuthConfigDTO{
			General: &apitypes.GeneralAuthConfigDTO{
				EnabledMethods:                    conf.Auth.General.EnabledMethods,
				IsMFARequired:                     conf.Auth.General.IsMFARequired,
				RequirePasswordChangeOnFirstLogin: conf.Auth.General.RequirePasswordChangeOnFirstLogin,
			},

			OIDC: &apitypes.AuthOIDCConfigDTO{
				ClientID:     conf.Auth.OIDC.ClientID,
				ClientSecret: oidcClientSecret,

				IssuerURL:   conf.Auth.OIDC.IssuerURL,
				RedirectURL: conf.Auth.OIDC.RedirectURL,

				AutomaticUserCreation: conf.Auth.OIDC.AutomaticUserCreation,
				UsernameClaim:         conf.Auth.OIDC.UsernameClaim,
				Prompt:                conf.Auth.OIDC.Prompt,

				RolesClaim:       conf.Auth.OIDC.RolesClaim,
				RequiredRole:     conf.Auth.OIDC.RequiredRole,
				AdditionalScopes: conf.Auth.OIDC.AdditionalScopes,
			},
		},

		Agent: apitypes.AgentConfigDTO{
			AutomaticallySyncListfiles: conf.Agent.AutomaticallySyncListfiles,
			SplitJobsPerAgent:          conf.Agent.SplitJobsPerAgent,
		},

		General: apitypes.GeneralConfigDTO{
			IsMaintenanceMode:               conf.General.IsMaintenanceMode,
			MaximumUploadedFileSize:         conf.General.MaximumUploadedFileSize,
			MaximumUploadedFileLineScanSize: conf.General.MaximumUploadedFileLineScanSize,
		},
	}
}

func (conf RuntimeConfig) ToPublicDTO() apitypes.PublicConfigDTO {
	return apitypes.PublicConfigDTO{
		Auth: apitypes.PublicAuthConfigDTO{
			EnabledMethods: conf.Auth.General.EnabledMethods,
			OIDC: apitypes.PublicOIDCConfigDTO{
				Prompt: conf.Auth.OIDC.Prompt,
			},
		},

		General: apitypes.PublicGeneralConfigDTO{
			IsMaintenanceMode:               conf.General.IsMaintenanceMode,
			MaximumUploadedFileSize:         conf.General.MaximumUploadedFileSize,
			MaximumUploadedFileLineScanSize: conf.General.MaximumUploadedFileLineScanSize,
		},
	}
}

var runningConf RuntimeConfig

func MakeDefaultConfig() RuntimeConfig {
	return RuntimeConfig{
		IsSetupComplete: true,
		ConfigVersion:   latestConfigVersion,

		Auth: AuthConfig{
			General: AuthGeneralConfig{
				EnabledMethods: []string{AuthMethodCredentials},

				IsMFARequired:                     false,
				RequirePasswordChangeOnFirstLogin: true,
			},

			OIDC: AuthOIDCConfig{
				AdditionalScopes:      []string{},
				AutomaticUserCreation: true,
				RolesClaim:            "groups",
				UsernameClaim:         "email",
			},
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
