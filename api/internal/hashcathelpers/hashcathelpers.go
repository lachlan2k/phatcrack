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
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
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

func CalculateKeyspace(params hashcattypes.HashcatParams) (int64, error) {
	wordlistPaths := []string{}
	for _, wordlist := range params.WordlistFilenames {
		wordlistId, err := uuid.Parse(wordlist)
		if err != nil {
			return 0, fmt.Errorf("invalid wordlist id provided: %q", wordlist)
		}

		filePath, err := filerepo.GetPathToFile(wordlistId)
		if err != nil {
			return 0, err
		}

		wordlistPaths = append(wordlistPaths, filePath)
	}

	rulefilePaths := []string{}
	for _, rulefile := range params.RulesFilenames {
		rulefileId, err := uuid.Parse(rulefile)
		if err != nil {
			return 0, fmt.Errorf("invalid rulefile id provided: %q", rulefile)
		}

		filePath, err := filerepo.GetPathToFile(rulefileId)
		if err != nil {
			return 0, err
		}

		rulefilePaths = append(rulefilePaths, filePath)
	}

	switch params.AttackMode {
	case hashcattypes.AttackModeDictionary:
		if len(rulefilePaths) > 1 || len(wordlistPaths) != 1 {
			return 0, fmt.Errorf("keyspace calculation for dicitonary attack requires exactly 0 or 1 rulefiles and 1 wordlist. found %d and %d", len(rulefilePaths), len(wordlistPaths))
		}

	case hashcattypes.AttackModeCombinator:
		if len(rulefilePaths) != 0 || len(wordlistPaths) != 2 {
			return 0, fmt.Errorf("keyspace calculation for combinator requires exactly 0 rulefiles and 2 wordlists. found %d and %d", len(rulefilePaths), len(wordlistPaths))
		}

	default:
		return 0, fmt.Errorf("keyspace calculation not implemented for attack mode %d", params.AttackMode)
	}

	args := []string{"--keyspace", "-m", strconv.Itoa(int(params.AttackMode))}
	for _, rule := range rulefilePaths {
		args = append(args, "-r", rule)
	}

	args = append(args, wordlistPaths...)

	if params.OptimizedKernels {
		args = append(args, "-O")
	}

	if params.SlowCandidates {
		args = append(args, "-S")
	}

	cmd, err := hashcatCommandWithRandSession(args...)
	if err != nil {
		return 0, err
	}

	out, err := cmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if ok {
			return 0, fmt.Errorf("hashcat gave an exit error: %w, %q, %q", ee, string(out), string(ee.Stderr))
		}
		return 0, fmt.Errorf("couldn't run hashcat: %w", err)
	}

	trimmedOut := strings.TrimSpace(strings.TrimSuffix(string(out), "\n"))

	return strconv.ParseInt(trimmedOut, 10, 64)
}

func IdentifyHashTypes(exampleHash string, hasUsername bool) ([]int, error) {
	tmpFile, err := os.CreateTemp("/tmp", "phatcrack-hash-identify")
	if err != nil {
		return nil, fmt.Errorf("couldn't create temporary file to store hashes: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	tmpFile.Chmod(0600)

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
			return nil, fmt.Errorf("hashcat gave an exit error: %w, %q, %q", ee, string(out), string(ee.Stderr))
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

	tmpFile.Chmod(0600)

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
			return nil, fmt.Errorf("username delim (:) not found in hashline: %q", hashline)
		}

		index, err := strconv.Atoi(usernameField)
		if err != nil {
			return nil, fmt.Errorf("found invalid index in hashline: %q", hashline)
		}
		if index >= len(hashes) || index < 0 {
			return nil, fmt.Errorf("hashcat gave us back an index out of range in hashline (%d): %q", index, hashline)
		}

		normalizedHashes[index] = hash
	}

	return normalizedHashes, nil
}
