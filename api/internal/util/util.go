package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ServerError(message string, internal error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, message).SetInternal(internal)
}

func CleanPath(filePath string) string {
	_, file := path.Split(path.Clean(filePath))
	return file
}

func UnmarshalJSON[T interface{}](jsonBlob string) (out T, err error) {
	err = json.Unmarshal([]byte(jsonBlob), &out)
	return
}

const agentKeyLen = 32 // 256-bit

func HashAgentKey(keyStr string) string {
	hash := sha256.Sum256([]byte(keyStr))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

func CheckAgentKey(keyStr string, hashStr string) bool {
	return HashAgentKey(keyStr) == hashStr
}

// Note: because we're using a large, random key, there's no need to salt.
// Plus, we get more efficient lookups as we can deterministically hash the incoming key...
// ...when we do an auth lookup.
func GenAgentKeyAndHash() (keyStr string, hashStr string, err error) {
	key := make([]byte, agentKeyLen)
	_, err = rand.Read(key)
	if err != nil {
		err = fmt.Errorf("couldn't generate random agent key: %v", err)
		return
	}

	keyStr = strings.ToLower(hex.EncodeToString(key))
	hashStr = HashAgentKey(keyStr)

	return
}

func isValidUUID(id string) bool {
	if id == "" {
		return false
	}
	_, err := uuid.Parse(id)
	return err == nil
}

func AreValidUUIDs(candidates ...string) bool {
	for _, candidate := range candidates {
		if !isValidUUID(candidate) {
			return false
		}
	}

	return true
}
