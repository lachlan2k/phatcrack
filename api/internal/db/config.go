package db

import (
	"encoding/json"
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/util"
	"gorm.io/datatypes"
)

type Config struct {
	SimpleBaseModel
	Config datatypes.JSON
}

func (c Config) TableName() string {
	return "config"
}

func GetConfig[ConfigT interface{}]() (*ConfigT, error) {
	var configRow Config
	err := GetInstance().First(&configRow).Error
	if err != nil {
		return nil, err
	}

	conf, err := util.UnmarshalJSON[ConfigT](configRow.Config.String())
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func SeedConfig[ConfigT interface{}](defaultConfig ConfigT) error {
	confBytes, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	var configRowCount int64
	err = GetInstance().Model(&Config{}).Count(&configRowCount).Error
	if err != nil {
		return err
	}

	if configRowCount == 0 {
		return GetInstance().Create(&Config{
			Config: datatypes.JSON(confBytes),
		}).Error
	}

	if configRowCount > 1 {
		return fmt.Errorf("found %d config entries in database, there should only be 1", configRowCount)
	}

	return nil
}

func SetConfig[ConfigT interface{}](newConf ConfigT) error {
	confBytes, err := json.Marshal(newConf)
	if err != nil {
		return err
	}

	var configRow Config
	err = GetInstance().First(&configRow).Error
	if err != nil {
		return err
	}

	return GetInstance().Model(&configRow).Updates(&Config{
		Config: datatypes.JSON(confBytes),
	}).Error
}
