package util

import "encoding/json"

func UnmarshalJSON[T interface{}](jsonBlob string) (out T, err error) {
	err = json.Unmarshal([]byte(jsonBlob), &out)
	return
}
