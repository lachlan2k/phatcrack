package db

import (
	"fmt"

	"gorm.io/gorm/clause"
)

const ConfigKeyIsSetupComplete = "IsSetupComplete"

const ConfigValueTrue = "true"
const ConfigValueFalse = "false"

type ConfigItem struct {
	SimpleBaseModel
	Key   string `gorm:"uniqueIndex"`
	Value string
}

func GetConfigItem(configKey string) (*ConfigItem, error) {
	var item ConfigItem
	err := GetInstance().First(&item, "key = ?", configKey).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func GetConfigItemWithDefault(configKey string, defaultValue string) (*ConfigItem, error) {
	var item ConfigItem
	err := GetInstance().First(&item, "key = ?", configKey).Error
	if err == ErrNotFound {
		return SetConfigItem(configKey, defaultValue)
	}

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func SetConfigItem(configKey string, value string) (*ConfigItem, error) {
	item := &ConfigItem{
		Key:   configKey,
		Value: value,
	}

	err := GetInstance().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func upsertConfig() error {
	var userCount int64
	err := GetInstance().Model(&User{}).Count(&userCount).Error
	if err != nil {
		return err
	}

	if userCount == 0 {
		_, err := RegisterUser("admin", "changeme", "admin")
		if err != nil {
			return err
		}

		fmt.Println("Created default admin user with password \"changeme\"")
		SetConfigItem(ConfigKeyIsSetupComplete, ConfigValueFalse)
	}

	return nil
}
