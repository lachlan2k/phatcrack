package resources

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
)

//go:embed hash_info.json
var hashInfoJsonStr string

var infomap hashcattypes.HashTypeMap

func init() {
	err := json.Unmarshal([]byte(hashInfoJsonStr), &infomap)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal hash_info.json from disk: %w", err))
	}

	for key, val := range infomap {
		val.ID = key
	}
}

func GetHashTypeMap() hashcattypes.HashTypeMap {
	return infomap
}
