package config

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lachlan2k/phatcrack/api/internal/db"
)

var lock sync.Mutex

type RuntimeConfig struct {
	IsSetupComplete                   bool  `json:"is_setup_complete"`
	IsMFARequired                     bool  `json:"is_mfa_required"`
	AutomaticallySyncListfiles        bool  `json:"auto_sync_listfiles"`
	SplitJobsPerAgent                 int   `json:"split_jobs_per_agent"`
	RequirePasswordChangeOnFirstLogin bool  `json:"require_password_change_on_first_login"`
	MaximumUploadedFileSize           int64 `json:"maximum_uploaded_file_size"`
	MaximumUploadedFileLineScanSize   int64 `json:"maximum_uploaded_file_line_scan_size"`
}

var runningConf RuntimeConfig

func MakeDefaultConfig() RuntimeConfig {
	return RuntimeConfig{
		IsSetupComplete:                   true,
		IsMFARequired:                     false,
		AutomaticallySyncListfiles:        true,
		RequirePasswordChangeOnFirstLogin: true,
		SplitJobsPerAgent:                 1,
		MaximumUploadedFileSize:           10 * 1000 * 1000 * 1000, // 10GB
		MaximumUploadedFileLineScanSize:   500 * 1000 * 1000,       // 100MB
	}
}

func Reload() error {
	lock.Lock()
	defer lock.Unlock()
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
