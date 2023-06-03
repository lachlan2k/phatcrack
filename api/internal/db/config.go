package db

import (
	"encoding/json"

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
