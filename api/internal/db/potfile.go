package db

type PotfileEntry struct {
	UUIDBaseModel
	Hash         string
	PlaintextHex string
	HashType     uint
}

func AddPotfileEntry(entry *PotfileEntry) (*PotfileEntry, error) {
	return entry, GetInstance().Create(entry).Error
}
