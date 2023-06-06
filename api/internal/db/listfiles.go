package db

import "github.com/lachlan2k/phatcrack/common/pkg/apitypes"

const (
	ListfileTypeWordlist = "Wordlist"
	ListfileTypeRulefile = "Rulefile"
)

type Listfile struct {
	UUIDBaseModel
	Name                 string
	AvailableForDownload bool
	FileType             string
	SizeInBytes          uint64
	Lines                uint64
}

func (w *Listfile) ToDTO() apitypes.ListfileDTO {
	return apitypes.ListfileDTO{
		ID:          w.ID.String(),
		Name:        w.Name,
		FileType:    w.FileType,
		SizeInBytes: w.SizeInBytes,
		Lines:       w.Lines,
	}
}

func GetListfile(id string) (*Listfile, error) {
	var listfile Listfile
	err := GetInstance().First(&listfile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &listfile, nil
}

func CreateListfile(listfile *Listfile) (*Listfile, error) {
	return listfile, GetInstance().Create(listfile).Error
}

func MarkListfileAsAvailable(id string) error {
	return GetInstance().Model(&Listfile{}).Where("id = ?", id).Updates(&Listfile{AvailableForDownload: true}).Error
}

func GetAllRulefiles() ([]Listfile, error) {
	rulefiles := []Listfile{}
	err := GetInstance().Where("file_type = ?", ListfileTypeRulefile).Find(&rulefiles).Error
	if err != nil {
		return nil, err
	}
	return rulefiles, nil
}

func GetAllWordlists() ([]Listfile, error) {
	wordlists := []Listfile{}
	err := GetInstance().Where("file_type = ?", ListfileTypeWordlist).Find(&wordlists).Error
	if err != nil {
		return nil, err
	}
	return wordlists, nil
}

func GetAllListfiles() ([]Listfile, error) {
	listfiles := []Listfile{}
	err := GetInstance().Find(&listfiles).Error
	if err != nil {
		return nil, err
	}
	return listfiles, nil
}
