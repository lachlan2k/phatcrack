package db

import (
	"errors"

	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PotfileEntry struct {
	UUIDBaseModel
	Hash         string `gorm:"index:idx_uniq,unique"`
	PlaintextHex string `gorm:"index:idx_uniq,unique"`
	HashType     uint   `gorm:"index:idx_uniq,unique"`
}

type PotfileSearchResult struct {
	Entry *PotfileEntry
	Hash  string
	Found bool
}

func (r PotfileSearchResult) ToDTO() apitypes.PotfileSearchResultDTO {
	dto := apitypes.PotfileSearchResultDTO{
		Hash:  r.Hash,
		Found: r.Found,
	}

	if r.Entry != nil {
		dto.PlaintextHex = r.Entry.PlaintextHex
		dto.HashType = r.Entry.HashType
	}

	return dto
}

func AddPotfileEntry(newEntry *PotfileEntry) (*PotfileEntry, error) {
	return newEntry, GetInstance().
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(newEntry).Error
}

func SearchPotfile(hashes []string) ([]PotfileSearchResult, error) {
	results := make([]PotfileSearchResult, len(hashes))

	err := GetInstance().Transaction(func(tx *gorm.DB) error {
		for i, hashToSearch := range hashes {
			var foundEntry PotfileEntry

			err := tx.First(&foundEntry, "hash = ?", hashToSearch).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				results[i] = PotfileSearchResult{
					Entry: nil,
					Hash:  hashToSearch,
					Found: false,
				}
				continue
			}
			if err != nil {
				return err
			}

			results[i] = PotfileSearchResult{
				Entry: &foundEntry,
				Hash:  hashToSearch,
				Found: true,
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
