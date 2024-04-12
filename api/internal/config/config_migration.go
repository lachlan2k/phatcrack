package config

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type v1runtimeConfig struct {
	IsSetupComplete                   bool  `json:"is_setup_complete"`
	IsMFARequired                     bool  `json:"is_mfa_required"`
	IsMaintenanceMode                 bool  `json:"is_maintenance_mode"`
	AutomaticallySyncListfiles        bool  `json:"auto_sync_listfiles"`
	SplitJobsPerAgent                 int   `json:"split_jobs_per_agent"`
	RequirePasswordChangeOnFirstLogin bool  `json:"require_password_change_on_first_login"`
	MaximumUploadedFileSize           int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize   int64 `json:"maximum_uploaded_file_line_scan_size"`
}

func migrateV1ToV2(oldConfig v1runtimeConfig) RuntimeConfig {
	return RuntimeConfig{
		ConfigVersion:   2,
		IsSetupComplete: oldConfig.IsSetupComplete,

		Auth: AuthConfig{
			EnabledMethods: []string{AuthMethodCredentials},

			IsMFARequired:                     oldConfig.IsMFARequired,
			RequirePasswordChangeOnFirstLogin: oldConfig.RequirePasswordChangeOnFirstLogin,

			// all oidc fields are new, and can happily default to an empty string
			OIDCAdditionalScopes: []string{},
		},

		Agent: AgentConfig{
			AutomaticallySyncListfiles: oldConfig.AutomaticallySyncListfiles,
			SplitJobsPerAgent:          oldConfig.SplitJobsPerAgent,
		},

		General: GeneralConfig{
			IsMaintenanceMode:               oldConfig.IsMaintenanceMode,
			MaximumUploadedFileSize:         oldConfig.MaximumUploadedFileSize,
			MaximumUploadedFileLineScanSize: oldConfig.MaximumUploadedFileLineScanSize,
		},
	}
}

// Currently only migrates from 1 -> 2
func runMigrations(configJsonStr string) (*RuntimeConfig, error) {
	var justVersion struct {
		ConfigVersion int `json:"version"`
	}

	err := json.Unmarshal([]byte(configJsonStr), &justVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	if justVersion.ConfigVersion >= latestConfigVersion {
		return nil, fmt.Errorf("no migrations to perform, version is %d", justVersion.ConfigVersion)
	}

	log.Warnf("Migrating from config version %d to %d", justVersion.ConfigVersion, latestConfigVersion)

	var oldConf v1runtimeConfig

	err = json.Unmarshal([]byte(configJsonStr), &oldConf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal as old conf: %v", err)
	}

	newConf := migrateV1ToV2(oldConf)
	return &newConf, nil
}
