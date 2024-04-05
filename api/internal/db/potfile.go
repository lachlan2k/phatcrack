package db

import (
	"errors"

	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"gorm.io/gorm"
)

type PotfileEntry struct {
	UUIDBaseModel
	Hash         string `gorm:"index:idx_potfile_hash,type:hash"`
	PlaintextHex string
	HashType     uint
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
	err := GetInstance().Transaction(func(tx *gorm.DB) error {
		var foundEntry PotfileEntry

		err := tx.First(&foundEntry, "hash = ? and hash_type = ?", newEntry.Hash, newEntry.HashType).Error

		// If it wasn't found, create it
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(newEntry).Error
		}
		if err != nil {
			return err
		}

		// or, if it was found, but is a colission, create it
		if foundEntry.PlaintextHex != newEntry.PlaintextHex {
			return tx.Create(newEntry).Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newEntry, nil
}

func SearchPotfile(hashes []string) ([]PotfileSearchResult, error) {
	results := make([]PotfileSearchResult, 0)

	err := GetInstance().Transaction(func(tx *gorm.DB) error {
		for _, loopHashToSearch := range hashes {
			hashToSearch := loopHashToSearch

			foundEntries := []PotfileEntry{}

			err := tx.Where("hash = ?", hashToSearch).Find(&foundEntries).Error

			if errors.Is(err, gorm.ErrRecordNotFound) || len(foundEntries) == 0 {
				results = append(results, PotfileSearchResult{
					Entry: nil,
					Hash:  hashToSearch,
					Found: false,
				})
				continue
			}

			if err != nil {
				return err
			}

			for _, loopFoundEntry := range foundEntries {
				foundEntry := loopFoundEntry
				results = append(results, PotfileSearchResult{
					Entry: &foundEntry,
					Hash:  hashToSearch,
					Found: true,
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
