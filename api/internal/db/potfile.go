package db

import (
	"gorm.io/gorm/clause"
)

type PotfileEntry struct {
	UUIDBaseModel
	Hash         string `gorm:"index:idx_uniq,unique"`
	PlaintextHex string `gorm:"index:idx_uniq,unique"`
	HashType     uint   `gorm:"index:idx_uniq,unique"`
}

func AddPotfileEntry(newEntry *PotfileEntry) (*PotfileEntry, error) {
	return newEntry, GetInstance().
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(newEntry).Error
}
