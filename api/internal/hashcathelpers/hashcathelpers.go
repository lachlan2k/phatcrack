package hashcathelpers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func findBinary() (path string, err error) {
	path = os.Getenv("HC_PATH")
	if path != "" {
		_, err = os.Stat(path)
		if err != nil {
			err = fmt.Errorf("failed to stat the specified hashcat binary (%s): %w (check path and permissions?)", path, err)
			path = ""
		}
		return
	}

	path, err = exec.LookPath("hashcat")
	if err != nil {
		path, err = exec.LookPath("hashcat.bin")
		if err != nil {
			err = errors.New("couldn't find hashcat or hashcat.bin in path, and HC_PATH was not specified")
			return
		}
	}

	return
}

func hashcatCommand(args ...string) (*exec.Cmd, error) {
	binPath, err := findBinary()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(binPath, args...)
	return cmd, nil
}

func hashcatCommandWithRandSession(args ...string) (*exec.Cmd, error) {
	fullArgs := []string{"--session", uuid.New().String()}
	fullArgs = append(fullArgs, args...)

	return hashcatCommand(fullArgs...)
}

func IdentifyHashTypes(exampleHash string, hasUsername bool) ([]int, error) {
	tmpFile, err := os.CreateTemp("/tmp", "phatcrack-hash-identify")
	if err != nil {
		return nil, fmt.Errorf("couldn't create temporary file to store hashes: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_ = tmpFile.Chmod(0600)

	if hasUsername {
		_, exampleHash, _ = strings.Cut(exampleHash, ":")
	}

	_, err = tmpFile.WriteString(exampleHash)
	if err != nil {
		return nil, fmt.Errorf("failed to write example hash to file: %w", err)
	}

	cmd, err := hashcatCommandWithRandSession("--identify", tmpFile.Name(), "--machine-readable")
	if err != nil {
		return nil, err
	}

	out, err := cmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if ok {
			if bytes.Contains(ee.Stderr, []byte("No hash-mode matches")) {
				return []int{}, nil
			}
		}

		return nil, fmt.Errorf("couldn't run hashcat: %w", err)
	}

	candidates := make([]int, 0)

	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimPrefix(scanner.Text(), "Autodetecting hash-modes. Please be patient...Autodetected hash-modes")
		candidate, err := strconv.ParseInt(line, 10, 32)
		if err != nil {
			continue
		}
		candidates = append(candidates, int(candidate))
	}

	return candidates, nil
}

func NormalizeHashes(hashes []string, hashType int, hasUsernames bool) ([]string, error) {
	tmpFile, err := os.CreateTemp("/tmp", "phatcrack-hash-normalize")
	if err != nil {
		return nil, fmt.Errorf("couldn't create temporary file to store hashes: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_ = tmpFile.Chmod(0600)

	for index, hash := range hashes {
		if hasUsernames {
			_, hash, _ = strings.Cut(hash, ":")
		}

		// Use the list index as a "username" so hashcat outputs them in a nice way
		_, err = tmpFile.WriteString(strconv.Itoa(index) + ":" + strings.TrimSpace(hash) + "\n")
		if err != nil {
			return nil, fmt.Errorf("failed to write example hash to file: %w", err)
		}
	}

	cmd, err := hashcatCommandWithRandSession("-m", strconv.FormatUint(uint64(hashType), 10), tmpFile.Name(), "--left", "--username", "--potfile-path", "/dev/null")
	if err != nil {
		return nil, err
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("couldn't normalize hashes: %w, out: %v", err, out)
	}

	normalizedHashes := make([]string, len(hashes))

	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		hashline := strings.TrimSpace(scanner.Text())
		usernameField, hash, found := strings.Cut(hashline, ":")
		if !found {
			return nil, fmt.Errorf("username delim (:) not found in hashline: %s", hashline)
		}

		index, err := strconv.Atoi(usernameField)
		if err != nil {
			return nil, fmt.Errorf("found invalid index in hashline: %s", hashline)
		}
		if index >= len(hashes) || index < 0 {
			return nil, fmt.Errorf("hashcat gave us back an index out of range in hashline (%d): %s", index, hashline)
		}

		normalizedHashes[index] = hash
	}

	return normalizedHashes, nil
}
