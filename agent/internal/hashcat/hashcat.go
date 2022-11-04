package hashcat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/pkg/apitypes"
)

type HashcatParams struct {
	apitypes.HashcatParamsDTO
	tempHashFile string
}

type HashcatStatusGuess struct {
	GuessBase        string  `json:"guess_base"`
	GuessBaseCount   uint64  `json:"guess_base_count"`
	GuessBaseOffset  uint64  `json:"guess_base_offset"`
	GuessBasePercent float32 `json:"guess_base_percent"`

	GuessMod        string  `json:"guess_mod"`
	GuessModCount   uint64  `json:"guess_mod_count"`
	GuessModOffset  uint64  `json:"guess_mod_offset"`
	GuessModPercent float32 `json:"guess_mod_percent"`

	GuessMode int `json:"guess_mode"`
}

type HashcatStatusDevice struct {
	DeviceID   int    `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	Speed      int    `json:"speed"`
	Util       int    `json:"util"`
	Temp       int    `json:"temp"`
}

type HashcatStatus struct {
	OriginalLine string

	Session         string                `json:"session"`
	Guess           HashcatStatusGuess    `json:"guess"`
	Status          int                   `json:"status"`
	Target          string                `json:"target"`
	Progress        []int                 `json:"progress"`
	RestorePoint    int                   `json:"restore_point"`
	RecoveredHashes []int                 `json:"recovered_hashes"`
	RecoveredSalts  []int                 `json:"recovered_salts"`
	Rejected        int                   `json:"rejected"`
	Devices         []HashcatStatusDevice `json:"devices"`
}

type HashcatResult struct {
	Timestamp    time.Time
	Hash         string
	PlaintextHex string
}

const (
	AttackModeDictionary = 0
	AttackModeMask       = 1
	AttackModeHybridDM   = 6
	AttackModeHybridMD   = 7
)

func (params HashcatParams) Validate() error {
	switch params.AttackMode {
	case AttackModeDictionary:
		if len(params.WordlistFilenames) != 1 {
			return fmt.Errorf("expected 1 wordlist for dictionary attack (0), but %d given", len(params.WordlistFilenames))
		}
	case AttackModeMask:
		if params.Mask == "" {
			return errors.New("using mask attack (1), but no mask was given")
		}
	case AttackModeHybridDM, AttackModeHybridMD:
		if params.Mask == "" {
			return fmt.Errorf("using hybrid attack (%d), but no mask was given", params.AttackMode)
		}
		if len(params.WordlistFilenames) != 1 {
			return fmt.Errorf("using hybrid attack (%d), but no wordlist was given", params.AttackMode)
		}
	default:
		return fmt.Errorf("unsupported attack mode %d", params.AttackMode)
	}

	return nil
}

func (params HashcatParams) ToCmdArgs(conf config.Config) (args []string, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	args = append(
		args,
		"--outfile-format", "1,3,5",
		"--quiet",
		"--status-json",
		"--status-timer", "1",
		"--potfile-disable",
		"-a", strconv.Itoa(int(params.AttackMode)),
		"-m", strconv.Itoa(int(params.HashType)),
	)

	args = append(args, params.AdditionalArgs...)

	if params.OptimizedKernels {
		args = append(args, "-O")
	}

	wordlists := make([]string, len(params.WordlistFilenames))
	for i, list := range params.WordlistFilenames {
		wordlists[i] = path.Join(conf.WordlistsDirectory, path.Clean(list))
		if _, err = os.Stat(wordlists[i]); err != nil {
			err = fmt.Errorf("provided wordlist %s couldn't be opened on filesystem", wordlists[i])
			return
		}
	}

	rules := make([]string, len(params.RulesFilenames))
	for i, rule := range params.RulesFilenames {
		rules[i] = path.Join(conf.RulesDirectory, path.Clean(rule))
		if _, err = os.Stat(rules[i]); err != nil {
			err = fmt.Errorf("provided rules file %s couldn't be opened on filesystem", wordlists[i])
			return
		}
	}

	args = append(args, params.tempHashFile)

	switch params.AttackMode {
	case AttackModeDictionary:
		for _, rule := range rules {
			args = append(args, "-r", rule)
		}
		args = append(args, wordlists[0])
	case AttackModeMask:
		args = append(args, params.Mask)
	case AttackModeHybridDM:
		args = append(args, wordlists[0], params.Mask)
	case AttackModeHybridMD:
		args = append(args, params.Mask, wordlists[0])
	}

	return
}

func findBinary(conf config.Config) (path string, err error) {
	if conf.HashcatBinary != "" {
		_, err = os.Stat(conf.HashcatBinary)
		if err != nil {
			err = fmt.Errorf("failed to stat the specified hashcat binary (%s): %v (check path and permissions?)", path, err)
			return
		}
	}

	path, err = exec.LookPath("hashcat")
	if err != nil {
		path, err = exec.LookPath("hashcat.bin")
		if err != nil {
			err = errors.New("couldn't find hashcat or hashcat.bin in path, and hashcat_binary was not specified")
			return
		}
	}

	return
}

type HashcatSession struct {
	proc           *exec.Cmd
	hashFile       *os.File
	CrackedHashes  chan HashcatResult
	StatusUpdates  chan HashcatStatus
	StderrMessages chan string
	DoneChan       chan error
}

func (sess *HashcatSession) Start() error {
	pStdout, err := sess.proc.StdoutPipe()
	if err != nil {
		return fmt.Errorf("couldn't attach stdout to hashcat: %v", err)
	}

	pStderr, err := sess.proc.StderrPipe()
	if err != nil {
		return fmt.Errorf("couldn't attach stderr to hashcat: %v", err)
	}

	err = sess.proc.Start()
	if err != nil {
		return fmt.Errorf("couldn't start hashcat: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(pStdout)
		for scanner.Scan() {
			line := scanner.Text()

			switch line[0] {
			case '{':
				var status HashcatStatus
				json.Unmarshal([]byte(line), &status)
				sess.StatusUpdates <- status
			default:
				values := strings.Split(line, ":")
				if len(values) != 3 {
					log.Printf("unexpected lien contents: %s", line)
					continue
				}
				timestamp := values[0]
				hash := values[1]
				plainHex := values[2]

				timestampI, err := strconv.ParseInt(timestamp, 10, 64)
				if err != nil {
					log.Printf("couldn't parse hashcat timestamp %s: %v", timestamp, err)
					continue
				}

				sess.CrackedHashes <- HashcatResult{
					Timestamp:    time.Unix(timestampI, 0),
					Hash:         hash,
					PlaintextHex: plainHex,
				}
			}
		}

		sess.DoneChan <- sess.proc.Wait()
	}()

	go func() {
		scanner := bufio.NewScanner(pStderr)
		for scanner.Scan() {
			sess.StderrMessages <- scanner.Text()
		}
	}()

	return nil
}

func (sess *HashcatSession) Kill() error {
	return sess.proc.Process.Kill()
}

func NewHashcatSession(hashes []string, params HashcatParams, conf config.Config) (*HashcatSession, error) {
	binaryPath, err := findBinary(conf)
	if err != nil {
		return nil, err
	}

	hashFile, err := ioutil.TempFile("/tmp", "phatcrack-hashes")
	if err != nil {
		return nil, fmt.Errorf("couldn't make a temp file to store hashes: %v", err)
	}

	params.tempHashFile = hashFile.Name()

	args, err := params.ToCmdArgs(conf)
	if err != nil {
		return nil, err
	}

	hashFile.Chmod(0600)
	for _, hash := range hashes {
		_, err = hashFile.WriteString(hash + "\n")
		if err != nil {
			return nil, fmt.Errorf("couldn't write hash to file: %v", err)
		}
	}

	return &HashcatSession{
		proc:           exec.Command(binaryPath, args...),
		hashFile:       hashFile,
		CrackedHashes:  make(chan HashcatResult),
		StatusUpdates:  make(chan HashcatStatus),
		StderrMessages: make(chan string),
	}, nil
}
