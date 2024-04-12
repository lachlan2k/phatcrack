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

type WrappedServerError struct {
	internal error
	id       string
}

func (e WrappedServerError) Error() string {
	return e.internal.Error()
}

func (e WrappedServerError) Unwrap() error {
	return e.internal
}

func (e WrappedServerError) ID() string {
	return e.id
}

func ServerError(message string, internal error) *echo.HTTPError {
	wrapped := WrappedServerError{
		internal: internal,
		id:       uuid.NewString(),
	}

	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("%s (error id %s)", message, wrapped.id)).SetInternal(wrapped)
}

func GenericServerError(internal error) *echo.HTTPError {
	return ServerError("Something went wrong", internal)
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
		err = fmt.Errorf("couldn't generate random agent key: %w", err)
		return
	}

	keyStr = strings.ToLower(hex.EncodeToString(key))
	hashStr = HashAgentKey(keyStr)

	return
}

// All just the same as the Agent Key stuff, but duplicating the code to keep it uncoupled
func HashAPIKey(keyStr string) string {
	hash := sha256.Sum256([]byte(keyStr))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

func CheckAPIKey(keyStr string, hashStr string) bool {
	return HashAPIKey(keyStr) == hashStr
}

const apiKeyLen = 32

func GenAPIKeyAndHash() (keyStr string, hashStr string, err error) {
	key := make([]byte, apiKeyLen)
	_, err = rand.Read(key)
	if err != nil {
		err = fmt.Errorf("couldn't generate random api key: %w", err)
		return
	}

	keyStr = strings.ToLower(hex.EncodeToString(key))
	hashStr = HashAPIKey(keyStr)

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
