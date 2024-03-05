package apitypes

type PotfileSearchRequestDTO struct {
	Hashes []string `json:"hashes"`
}

type PotfileSearchResultDTO struct {
	Hash         string `json:"hash"`
	HashType     uint   `json:"hash_type"`
	PlaintextHex string `json:"plaintext_hex"`
	Found        bool   `json:"found"`
}

type PotfileSearchResponseDTO struct {
	Results []PotfileSearchResultDTO `json:"results"`
}
