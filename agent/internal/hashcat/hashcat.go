package hashcat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"

	"github.com/hpcloud/tail"
)

type HashcatParams hashcattypes.HashcatParams

type HashcatStatusGuess hashcattypes.HashcatStatusGuess

type HashcatStatusDevice hashcattypes.HashcatStatusDevice

type HashcatStatus = hashcattypes.HashcatStatus

type HashcatResult hashcattypes.HashcatResult

const (
	AttackModeDictionary = 0
	AttackModeCombinator = 1
	AttackModeMask       = 3
	AttackModeHybridDM   = 6
	AttackModeHybridMD   = 7
)

func (params HashcatParams) Validate() error {
	switch params.AttackMode {
	case AttackModeDictionary:
		if len(params.WordlistFilenames) != 1 {
			return fmt.Errorf("expected 1 wordlist for dictionary attack (%d), but %d given", AttackModeDictionary, len(params.WordlistFilenames))
		}

	case AttackModeCombinator:
		if len(params.WordlistFilenames) != 2 {
			return fmt.Errorf("expected 2 wordlists for combinator attack (%d), but %d given", AttackModeCombinator, len(params.WordlistFilenames))
		}

	case AttackModeMask:
		if params.Mask == "" {
			return fmt.Errorf("using mask attack (%d), but no mask was given", AttackModeMask)
		}

	case AttackModeHybridDM, AttackModeHybridMD:
		if params.Mask == "" {
			return fmt.Errorf("using hybrid attack (%d), but no mask was given", params.AttackMode)
		}
		if len(params.WordlistFilenames) != 1 {
			return fmt.Errorf("using hybrid attack (%d), but %d wordlist were given", params.AttackMode, len(params.WordlistFilenames))
		}

	default:
		return fmt.Errorf("unsupported attack mode %d", params.AttackMode)
	}

	return nil
}

func (params HashcatParams) maskArgs() ([]string, error) {
	if len(params.MaskCustomCharsets) > 4 {
		return nil, fmt.Errorf("too many custom charsets supplied (%d), the max is 4", len(params.MaskCustomCharsets))
	}

	args := []string{}

	for i, charset := range params.MaskCustomCharsets {
		// Hashcat accepts paramters --custom-charset1 to --custom-charset4
		args = append(args, fmt.Sprintf("--custom-charset%d", i+1), charset)
	}

	if params.MaskIncrement {
		args = append(args, "--increment")

		if params.MaskIncrementMin > 0 {
			args = append(args, "--increment-min", strconv.Itoa(int(params.MaskIncrementMin)))
		}

		if params.MaskIncrementMax > 0 {
			args = append(args, "--increment-max", strconv.Itoa(int(params.MaskIncrementMax)))
		}
	}

	return args, nil
}

func (params HashcatParams) ToCmdArgs(conf *config.Config, session, tempHashFile string, outFile string) (args []string, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	args = append(
		args,
		"--quiet",
		"--session", "sess-"+session+"_"+uuid.New().String(),
		"--outfile-format", "1,3,5",
		"--outfile", outFile,
		"--status",
		"--status-json",
		"--status-timer", "3",
		"--potfile-disable",
		"-a", strconv.Itoa(int(params.AttackMode)),
		"-m", strconv.Itoa(int(params.HashType)),
	)

	args = append(args, params.AdditionalArgs...)

	if params.OptimizedKernels {
		args = append(args, "-O")
	}

	if params.SlowCandidates {
		args = append(args, "-S")
	}

	wordlists := make([]string, len(params.WordlistFilenames))
	for i, list := range params.WordlistFilenames {
		wordlists[i] = path.Join(conf.ListfileDirectory, path.Clean(list))
		if _, err = os.Stat(wordlists[i]); err != nil {
			err = fmt.Errorf("provided wordlist %s couldn't be opened on filesystem", wordlists[i])
			return
		}
	}

	rules := make([]string, len(params.RulesFilenames))
	for i, rule := range params.RulesFilenames {
		rules[i] = path.Join(conf.ListfileDirectory, path.Clean(rule))
		if _, err = os.Stat(rules[i]); err != nil {
			err = fmt.Errorf("provided rules file %s couldn't be opened on filesystem", wordlists[i])
			return
		}
	}

	args = append(args, tempHashFile)

	switch params.AttackMode {
	case AttackModeDictionary:
		for _, rule := range rules {
			args = append(args, "-r", rule)
		}
		args = append(args, wordlists[0])

	case AttackModeCombinator:
		args = append(args, wordlists[0], wordlists[1])

	case AttackModeMask:
		args = append(args, params.Mask)

	case AttackModeHybridDM:
		args = append(args, wordlists[0], params.Mask)

	case AttackModeHybridMD:
		args = append(args, params.Mask, wordlists[0])
	}

	switch params.AttackMode {
	case AttackModeMask, AttackModeHybridDM, AttackModeHybridMD:
		maskArgs, err := params.maskArgs()
		if err != nil {
			return nil, err
		}
		args = append(args, maskArgs...)
	}

	return
}

func findBinary(conf *config.Config) (path string, err error) {
	path = conf.HashcatBinary
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
			err = errors.New("couldn't find hashcat or hashcat.bin in path, and hashcat_binary was not specified")
			return
		}
	}

	return
}

type HashcatSession struct {
	proc           *exec.Cmd
	hashFile       *os.File
	outFile        *os.File
	CrackedHashes  chan HashcatResult
	StatusUpdates  chan HashcatStatus
	StderrMessages chan string
	StdoutLines    chan string
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

	log.Printf("Running hashcat command: %s", sess.proc.String())

	err = sess.proc.Start()
	if err != nil {
		return fmt.Errorf("couldn't start hashcat: %v", err)
	}

	tailer, err := tail.TailFile(sess.outFile.Name(), tail.Config{Follow: true})
	if err != nil {
		sess.Kill()
		return fmt.Errorf("couldn't tail outfile %s: %v", sess.outFile.Name(), err)
	}

	go func() {
		for tLine := range tailer.Lines {
			line := tLine.Text
			values := strings.Split(line, ":")
			if len(values) != 3 {
				log.Printf("unexpected line contents: %s", line)
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
	}()

	go func() {
		scanner := bufio.NewScanner(pStdout)
		for scanner.Scan() {
			line := scanner.Text()

			log.Printf("read line %v", line)
			sess.StdoutLines <- line

			if len(line) == 0 {
				continue
			}

			switch line[0] {
			case '{':
				var status HashcatStatus
				err := json.Unmarshal([]byte(line), &status)
				if err != nil {
					fmt.Printf("WARN: couldn't unmarshal hashcat status: %v", err)
					continue
				}
				sess.StatusUpdates <- status

			default:
				log.Printf("Unexpected stdout line: %v", line)
			}
		}

		sess.DoneChan <- sess.proc.Wait()
		tailer.Kill(nil)
	}()

	go func() {
		scanner := bufio.NewScanner(pStderr)
		for scanner.Scan() {
			log.Printf("read stderr: %s", scanner.Text())
			sess.StderrMessages <- scanner.Text()
		}
	}()

	return nil
}

func (sess *HashcatSession) Kill() error {
	return sess.proc.Process.Kill()
}

func (sess *HashcatSession) Cleanup() {
	if sess.hashFile != nil {
		os.Remove(sess.hashFile.Name())
		sess.hashFile = nil
	}

	if sess.outFile != nil {
		os.Remove(sess.outFile.Name())
		sess.outFile = nil
	}
}

func NewHashcatSession(id string, hashes []string, params HashcatParams, conf *config.Config) (*HashcatSession, error) {
	binaryPath, err := findBinary(conf)
	if err != nil {
		return nil, err
	}

	hashFile, err := os.CreateTemp("/tmp", "phatcrack-hashes")
	if err != nil {
		return nil, fmt.Errorf("couldn't make a temp file to store hashes: %v", err)
	}
	hashFile.Chmod(0600)

	outFile, err := os.CreateTemp("/tmp", "phatcrack-output")
	if err != nil {
		return nil, fmt.Errorf("couldn't make a temp file to store output: %v", err)
	}
	outFile.Chmod(0600)

	args, err := params.ToCmdArgs(conf, id, hashFile.Name(), outFile.Name())
	if err != nil {
		return nil, err
	}

	for _, hash := range hashes {
		_, err = hashFile.WriteString(hash + "\n")
		if err != nil {
			return nil, fmt.Errorf("couldn't write hash to file: %v", err)
		}
	}

	return &HashcatSession{
		proc:           exec.Command(binaryPath, args...),
		hashFile:       hashFile,
		outFile:        outFile,
		CrackedHashes:  make(chan HashcatResult, 5),
		StatusUpdates:  make(chan HashcatStatus, 5),
		StderrMessages: make(chan string, 5),
		StdoutLines:    make(chan string, 5),
		DoneChan:       make(chan error),
	}, nil
}
