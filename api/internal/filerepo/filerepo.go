package filerepo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var basePath string

func SetPath(pathToSet string) error {
	inf, err := os.Stat(pathToSet)
	if err != nil {
		return err
	}

	if !inf.IsDir() {
		return fmt.Errorf("provided filerepo path: %v is not a directory", pathToSet)
	}

	basePath = pathToSet
	return nil
}

func GetPathToFile(id uuid.UUID) (string, error) {
	filename := id.String()
	if filename == "" {
		return "", fmt.Errorf("invalid uuid for filename: %v", id)
	}

	return filepath.Join(basePath, filename), nil
}

func Create(id uuid.UUID) (io.WriteCloser, error) {
	filename := id.String()
	if filename == "" {
		return nil, fmt.Errorf("invalid uuid for filename: %v", id)
	}

	return os.OpenFile(filepath.Join(basePath, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
}

func Delete(id uuid.UUID) error {
	filename, err := GetPathToFile(id)
	if err != nil {
		return err
	}

	return os.Remove(filename)
}

func MakeTmp() (*os.File, string, error) {
	tmpFilename := filepath.Join(basePath, fmt.Sprintf("tmp-%s", uuid.New().String()))

	f, err := os.Create(tmpFilename)
	if err != nil {
		return nil, "", err
	}

	return f, tmpFilename, nil
}

func CreateFromTmp(id uuid.UUID, tmpFilename string) error {
	filename, err := GetPathToFile(id)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFilename, filename)
	if err != nil {
		return err
	}
	return nil
}
