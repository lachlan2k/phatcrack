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
)

func findBinary() (path string, err error) {
	path = os.Getenv("HC_PATH")
	if path != "" {
		_, err = os.Stat(path)
		if err != nil {
			err = fmt.Errorf("failed to stat the specified hashcat binary (%s): %v (check path and permissions?)", path, err)
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

func IdentifyHashTypes(exampleHash string) ([]int, error) {
	tmpFile, err := os.CreateTemp("/tmp", "phatcrack-hash-identify")
	if err != nil {
		return nil, fmt.Errorf("couldn't create temporary file to store hashes: %v", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_ = tmpFile.Chmod(0600)

	_, err = tmpFile.WriteString(exampleHash)
	if err != nil {
		return nil, fmt.Errorf("failed to write example hash to file: %v", err)
	}

	cmd, err := hashcatCommand("--identify", tmpFile.Name(), "--machine-readable")
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

		return nil, fmt.Errorf("couldn't run hashcat: %v", err)
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

func NormalizeHashes(hashes []string, hashMode int) ([]string, error) {
	tmpFile, err := os.CreateTemp("/tmp", "phatcrack-hash-normalize")
	if err != nil {
		return nil, fmt.Errorf("couldn't create temporary file to store hashes: %v", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_ = tmpFile.Chmod(0600)

	for _, hash := range hashes {
		_, err = tmpFile.WriteString(strings.TrimSpace(hash) + "\n")
		if err != nil {
			return nil, fmt.Errorf("failed to write example hash to file: %v", err)
		}
	}

	cmd, err := hashcatCommand("-m", strconv.Itoa(hashMode), tmpFile.Name(), "--left", "--potfile-path", "/dev/null")
	if err != nil {
		return nil, err
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("couldn't normalize hashes: %v", err)
	}

	normalizedHashes := make([]string, 0)

	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		hash := strings.TrimSpace(scanner.Text())
		if hash != "" {
			normalizedHashes = append(normalizedHashes, scanner.Text())
		}
	}

	return normalizedHashes, nil
}
