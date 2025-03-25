package db

type KeyspaceCache struct {
	SimpleBaseModel

	SettingsHash string `gorm:"index"`
	Keyspace int64
}

func GetKeyspaceCacheEntry(settingsHash string) (int64, error) {
	cache := &KeyspaceCache{}
	err := GetInstance().First(&cache, "settings_hash = ?", settingsHash).Error
	if err != nil {
		return 0, err
	}
	return cache.Keyspace, nil
}

func InsertKeyspaceCacheEntry(settingsHash string, keyspace int64) error {
	err := GetInstance().Create(&KeyspaceCache{
		SettingsHash: settingsHash,
		Keyspace: keyspace,
	}).Error
	return err
}