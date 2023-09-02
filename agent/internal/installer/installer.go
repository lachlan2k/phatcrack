package installer

import (
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"text/template"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
)

type InstallConfig struct {
	Defaults bool

	AgentUser    string
	AgentGroup   string
	AgentBinPath string

	AuthKey     string
	AuthKeyFile string
	ConfigPath  string

	HashcatPath            string
	ListfileDirectory      string
	APIEndpoint            string
	DisableTLSVerification bool
}

//go:embed template.service
var serviceFileTemplateString string

const serviceUnitFilePath = "/etc/systemd/system/phatcrack-agent.service"

func getUidAndGid(installConf InstallConfig) (int, int) {
	runningUser, err := user.Lookup(installConf.AgentUser)
	if runningUser == nil || err != nil {
		log.Fatalf("Couldn't look up user %q for installation: %v", installConf.AgentUser, err)
	}

	runningGroup, err := user.LookupGroup(installConf.AgentGroup)
	if runningGroup == nil || err != nil {
		log.Fatalf("Couldn't look up user group %q for installation: %v", installConf.AgentGroup, err)
	}

	uid, _ := strconv.Atoi(runningUser.Uid)
	gid, _ := strconv.Atoi(runningGroup.Gid)
	return uid, gid
}

func writeAuthKeyFile(installConf InstallConfig) {
	f, err := os.OpenFile(installConf.AuthKeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Couldn't open auth key file for writing: %v", err)
	}

	defer f.Close()
	_, err = f.Write([]byte(installConf.AuthKey))
	if err != nil {
		log.Fatalf("Couldn't write key to auth key file: %v", err)
	}

	uid, gid := getUidAndGid(installConf)
	f.Chown(uid, gid)
}

func writeServiceFile(installConf InstallConfig) {
	serviceFileTmpl, err := template.New("Service File").Parse(serviceFileTemplateString)
	if err != nil {
		log.Fatalf("Couldn't compile template: %v", err)
	}

	f, err := os.OpenFile(serviceUnitFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Couldn't open systemd unit file for writing: %v", err)
	}
	defer f.Close()

	err = serviceFileTmpl.Execute(f, installConf)
	if err != nil {
		log.Fatalf("Failed to write service file: %v", err)
	}

	f.Chown(0, 0)
}

func writeConfigFile(installConf InstallConfig) {
	uid, gid := getUidAndGid(installConf)

	marshalled, err := json.MarshalIndent(config.Config{
		AuthKeyFile:            installConf.AuthKeyFile,
		AuthKey:                "",
		HashcatPath:            installConf.HashcatPath,
		ListfileDirectory:      installConf.ListfileDirectory,
		APIEndpoint:            installConf.APIEndpoint,
		DisableTLSVerification: installConf.DisableTLSVerification,
	}, "", "  ")
	if err != nil {
		log.Fatalf("Couldn't marshal config file: %v", err)
	}

	f, err := os.OpenFile(installConf.ConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		log.Fatalf("Couldn't open config path for writing")
	}
	defer f.Close()

	_, err = f.Write(marshalled)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	f.Chown(uid, gid)
}

func setupPaths(installConf InstallConfig) {
	uid, gid := getUidAndGid(installConf)

	if installConf.Defaults {
		os.MkdirAll("/opt/phatcrack-agent/", 0700)
		os.Chown("/opt/phatcrack-agent/", uid, gid)
	}

	os.MkdirAll(installConf.ListfileDirectory, 0700)
	os.Chown(installConf.ListfileDirectory, uid, gid)
}

func Run(installConf InstallConfig) {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Couldn't get current user for installation: %v", err)
	}

	if user.Uid != "0" || user.Gid != "0" {
		log.Fatalf("Agent installer must be run as root")
	}

	log.Println("Setting up paths...")
	setupPaths(installConf)

	log.Println("Writing service file...")
	writeServiceFile(installConf)

	log.Println("Writing config file...")
	writeConfigFile(installConf)

	log.Println("Writing key file...")
	writeAuthKeyFile(installConf)

	log.Println("Done! Run 'systemctl enable --now phatcrack-agent' to start the agent")
}

func applyDefaults(installConf *InstallConfig) {
	if installConf.AgentUser == "" {
		u, err := user.Lookup("phatcrack-agent")
		if err != nil || u == nil {
			installConf.AgentUser = "root"
		} else {
			installConf.AgentUser = "phatcrack-agent"
		}
	}

	if installConf.AgentGroup == "" {
		g, err := user.LookupGroup("phatcrack-agent")
		if err != nil || g == nil {
			installConf.AgentGroup = "root"
		} else {
			installConf.AgentGroup = "phatcrack-agent"
		}
	}

	if installConf.AgentBinPath == "" {
		exe, err := os.Executable()
		if err == nil {
			installConf.AgentBinPath = exe
		}
	}

	if installConf.AuthKeyFile == "" {
		installConf.AuthKeyFile = "/opt/phatcrack-agent/auth.key"
	}

	if installConf.ConfigPath == "" {
		installConf.ConfigPath = "/opt/phatcrack-agent/config.json"
	}

	if installConf.HashcatPath == "" {
		path, err := exec.LookPath("hashcat")
		if err != nil || path == "" {
			path, err = exec.LookPath("hashcat.bin")
			if err != nil || path == "" {
				path = "/opt/phatcrack-agent/hashcat/hashcat.bin"
			}
		}
		installConf.HashcatPath = path
	}

	if installConf.ListfileDirectory == "" {
		installConf.ListfileDirectory = "/opt/phatcrack-agent/listfiles/"
	}
}

func input(prompt string) string {
	var res string
	fmt.Print(prompt)
	fmt.Scanln(&res)
	return strings.TrimSuffix(strings.TrimSpace(res), "\n")
}

func getOptionsInteractive(installConf *InstallConfig) {
	fmt.Println("You will be prompted to enter anything that hasn't been configured. Press enter for default.")

	if installConf.AgentUser == "" {
		installConf.AgentUser = input("Which user do you want to run Phatcrack as? (default: phatcrack-agent if present, else root): ")
	}

	if installConf.AgentGroup == "" {
		installConf.AgentGroup = input("Which user group do you want to run Phatcrack as? (default: phatcrack-agent if present, else root): ")
	}

	if installConf.AgentBinPath == "" {
		installConf.AgentBinPath = input("Where is the phatcrack agent binary? (default: current binary): ")
	}

	if installConf.AuthKey == "" {
		installConf.AuthKey = input("Auth key from server (this is okay to leave blank for now): ")
	}

	if installConf.AuthKeyFile == "" {
		installConf.AuthKeyFile = input("Auth key file to write (default: /opt/phatcrack-agent/auth.key): ")
	}

	if installConf.ConfigPath == "" {
		installConf.ConfigPath = input("Config file to write (default: /opt/phatcrack-agent/config.json): ")
	}

	if installConf.HashcatPath == "" {
		installConf.HashcatPath = input("Path to hashcat executable (default: searches PATH): ")
	}

	if installConf.ListfileDirectory == "" {
		installConf.ListfileDirectory = input("Directory to store listfiles (default: /opt/phatcrack-agent/listfiles/): ")
	}

	if installConf.APIEndpoint == "" {
		installConf.APIEndpoint = input("API Endpoint (format: https://phatcrack.lan/api/v1): ")
	}
}

func checkConf(installConf InstallConfig) error {
	if installConf.AgentUser == "" {
		return errors.New("agent user is not set")
	}

	if installConf.AgentGroup == "" {
		return errors.New("agent group is not set")
	}

	if installConf.AgentBinPath == "" {
		return errors.New("agent bin path could not be determined")
	}

	// blank Authkey is fine

	if installConf.AuthKeyFile == "" {
		return errors.New("auth key file path is not set")
	}

	if installConf.ConfigPath == "" {
		return errors.New("config file path is not set")
	}

	if installConf.HashcatPath == "" {
		return errors.New("hashcat path could not be determined: either add to PATH, or specify")
	}

	if installConf.ListfileDirectory == "" {
		return errors.New("listfile directory is not set")
	}

	// Blank API Endpoint is fine, assuming user is maybe setting up agents before server

	return nil
}

func RunInteractive() {
	flagSet := flag.NewFlagSet("install", flag.ExitOnError)

	useDefaultsP := flagSet.Bool("defaults", false, "Use basic defaults for intallation? (stores everything in /opt/phatcrack-agent)")
	userP := flagSet.String("user", "", "Which user to run the agent as")
	groupP := flagSet.String("group", "", "Which user group to run the agent as")
	agentBinPathP := flagSet.String("agent-bin", "", "Path to agent (defaults to running executable)")
	authKeyFileP := flagSet.String("auth-keyfile", "", "Path to file containing agent key")
	authKeyP := flagSet.String("auth-key", "", "Path to file containing agent key")
	configPathP := flagSet.String("config-path", "", "Path to config json file to install")
	hashcatPathP := flagSet.String("hashcat-path", "", "Path to hashcat executable")
	listfilePathP := flagSet.String("listfile-directory", "", "Path to directory to hold listfiles")
	apiEndpointP := flagSet.String("api-endpoint", "", "API endpoint (format: https://phatcrack.lan/api/v1)")
	disableTLSVerificationP := flagSet.Bool("disable-tls-verification", false, "Whether to disable TLS Verification")

	flagSet.Parse(os.Args[2:])

	installConf := InstallConfig{
		Defaults: *useDefaultsP,

		AgentUser:    *userP,
		AgentGroup:   *groupP,
		AgentBinPath: *agentBinPathP,

		AuthKey:     *authKeyP,
		AuthKeyFile: *authKeyFileP,
		ConfigPath:  *configPathP,

		HashcatPath:            *hashcatPathP,
		ListfileDirectory:      *listfilePathP,
		APIEndpoint:            *apiEndpointP,
		DisableTLSVerification: *disableTLSVerificationP,
	}

	if installConf.Defaults {
		applyDefaults(&installConf)
	} else {
		getOptionsInteractive(&installConf)
		applyDefaults(&installConf)
	}

	err := checkConf(installConf)
	if err != nil {
		panic(err)
	}

	Run(installConf)
}
