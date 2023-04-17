package db

import "github.com/lachlan2k/phatcrack/common/pkg/apitypes"

type Wordlist struct {
	UUIDBaseModel
	Name           string
	Description    string
	FilenameOnDisk string
	SizeInBytes    uint64
	Lines          uint64
}

func (w *Wordlist) ToDTO() apitypes.WordlistDTO {
	return apitypes.WordlistDTO{
		ID:             w.ID.String(),
		Name:           w.Name,
		Description:    w.Description,
		FilenameOnDisk: w.FilenameOnDisk,
		SizeInBytes:    w.SizeInBytes,
		Lines:          w.Lines,
	}
}

type RuleFile struct {
	UUIDBaseModel
	Name           string
	Description    string
	FilenameOnDisk string
	SizeInBytes    uint64
	Lines          uint64
}

func (r *RuleFile) ToDTO() apitypes.RuleFileDTO {
	return apitypes.RuleFileDTO{
		ID:             r.ID.String(),
		Name:           r.Name,
		Description:    r.Description,
		FilenameOnDisk: r.FilenameOnDisk,
		SizeInBytes:    r.SizeInBytes,
		Lines:          r.Lines,
	}
}

func GetWordlist(id string) (*Wordlist, error) {
	var wordlist Wordlist
	err := GetInstance().First(&wordlist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wordlist, nil
}

func GetAllWordlists() ([]Wordlist, error) {
	wordlists := []Wordlist{}
	err := GetInstance().Find(&wordlists).Error
	if err != nil {
		return nil, err
	}
	return wordlists, nil
}

func GetRuleFile(id string) (*RuleFile, error) {
	var rulefile RuleFile
	err := GetInstance().First(&rulefile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rulefile, nil
}

func GetAllRuleFiles() ([]RuleFile, error) {
	rulefiles := []RuleFile{}
	err := GetInstance().Find(&rulefiles).Error
	if err != nil {
		return nil, err
	}
	return rulefiles, nil
}
