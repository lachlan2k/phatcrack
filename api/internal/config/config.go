package config

import (
	"fmt"
	"sync"

	"github.com/lachlan2k/phatcrack/api/internal/db"
)

var lock sync.Mutex

type RuntimeConfig struct {
	IsSetupComplete                   bool `json:"is_setup_complete"`
	IsMFARequired                     bool `json:"is_mfa_required"`
	RequirePasswordChangeOnFirstLogin bool `json:"require_password_change_on_first_login"`
}

var runningConf RuntimeConfig

func Reload() error {
	lock.Lock()
	defer lock.Unlock()
	newConf, err := db.GetConfig[RuntimeConfig]()
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
		return fmt.Errorf("failed to save new config in db: %v", err)
	}

	runningConf = newConf
	return nil
}

func Get() RuntimeConfig {
	return runningConf
}
