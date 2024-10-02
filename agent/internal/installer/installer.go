package installer

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

type InstallConfig struct {
	Defaults bool

	AgentUser    string
	AgentGroup   string
	AgentBinPath string

	AuthKey     string
	AuthKeyFile string
	RegistrationKey string
	ConfigPath  string

	HashcatPath            string
	ListfileDirectory      string
	APIEndpoint            string
	DisableTLSVerification bool
	InstallHashcat         bool
}

//go:embed template.service
var serviceFileTemplateString string

const serviceUnitFilePath = "/etc/systemd/system/phatcrack-agent.service"

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

	adjustPerms(f, installConf)
}

func writeConfigFile(installConf InstallConfig) {

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

	adjustPerms(f, installConf)
}

func setupPaths(installConf InstallConfig) {

	if installConf.Defaults {
		os.MkdirAll(DefaultPathJoin(), 0700)
		adjustPermsPath(DefaultPathJoin(), installConf)
	}

	os.MkdirAll(installConf.ListfileDirectory, 0700)
	adjustPermsPath(installConf.ListfileDirectory, installConf)
}

func Run(installConf InstallConfig) {

	if err := isElevated(); err != nil {
		log.Fatal(err)
	}

	log.Println("Setting up paths...")
	setupPaths(installConf)

	log.Println("Writing config file...")
	writeConfigFile(installConf)

	log.Println("Writing key file...")
	writeAuthKeyFile(installConf)

	if installConf.InstallHashcat {
		installHashcat(installConf)
	}

	log.Println("Installing service...")
	installService(installConf)

	if runtime.GOOS != "windows" {
		log.Println("Done! Run 'systemctl enable --now phatcrack-agent' to start the agent")
	} else {
		log.Println("Done! Execute agent binary as administrator to start agent")

	}

}

func applyDefaults(installConf *InstallConfig) {
	if installConf.AgentUser == "" {
		u, err := user.Lookup("phatcrack-agent")
		if err != nil || u == nil {

			installConf.AgentUser = "root"
			if runtime.GOOS == "windows" {
				// On windows we dont have the concept of a root user, so just use the current user. We check that the binary has been elevated before doing any of the actual install stuff
				user, _ := user.Current()
				installConf.AgentUser = user.Uid
			}

		} else {
			installConf.AgentUser = "phatcrack-agent"
		}
	}

	installConf.InstallHashcat = true

	if installConf.AgentGroup == "" {
		g, err := user.LookupGroup("phatcrack-agent")
		if err != nil || g == nil {
			installConf.AgentGroup = "root"
			if runtime.GOOS == "windows" {
				user, _ := user.Current()
				installConf.AgentGroup = user.Gid
			}
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
		installConf.AuthKeyFile = DefaultPathJoin("auth.key")
	}

	if installConf.ConfigPath == "" {
		installConf.ConfigPath = DefaultPathJoin("config.json")
	}

	if installConf.HashcatPath == "" {
		path, err := exec.LookPath("hashcat")
		if err != nil || path == "" {
			path, err = exec.LookPath(hashcat.Hashcat)
			if err != nil || path == "" {
				path = DefaultPathJoin("hashcat/" + hashcat.Hashcat)
			}
		}
		installConf.HashcatPath = path
	}

	if installConf.ListfileDirectory == "" {
		installConf.ListfileDirectory = DefaultPathJoin("listfiles/")
	}

}

func input(p string, a ...any) string {
	var res string
	fmt.Printf(p, a...)
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

	if installConf.AuthKey == "" && installConf.RegistrationKey == "" {
		installConf.RegistrationKey	= input("Registration key from server (leave blank to specify an auth key): ")
		if installConf.RegistrationKey == "" {
			installConf.AuthKey = input("Auth key from server (this is okay to leave blank for now): ")
		}
	}

	if installConf.AuthKeyFile == "" {
		installConf.AuthKeyFile = input("Auth key file to write (default: %s): ", DefaultPathJoin("auth.key"))
	}

	if installConf.ConfigPath == "" {
		installConf.ConfigPath = input("Config file to write (default: %s): ", DefaultPathJoin("config.json"))
	}

	if installConf.HashcatPath == "" {
		installConf.HashcatPath = input("Path to hashcat executable (default: searches PATH): ")
	}

	if installConf.ListfileDirectory == "" {
		installConf.ListfileDirectory = input("Directory to store listfiles (default: %s): ", DefaultPathJoin("listfiles/"))
	}

	if installConf.APIEndpoint == "" {
		installConf.APIEndpoint = input("API Endpoint (format: https://phatcrack.lan/api/v1): ")
	}
}

func registerIfRequired(installConf *InstallConfig) {
	if installConf.RegistrationKey != "" {
		u, err := url.Parse(installConf.APIEndpoint)
		if err != nil {
			log.Fatal("failed to parse API endpoint to register agent: ", err)
			return
		}

		name, err := os.Hostname()
		if err != nil || name == "" {
			log.Printf("Warn: couldn't get hostname for agent registration: %v\n", err)
			name = "unknown"
		}

		reqBody := apitypes.AgentRegisterRequestDTO{
			Name: name,
		}

		reqBodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatalf("Failed to marshal agent registration request: %v", err)
		}

		u.Path = path.Join(u.Path, "/agent-handler/register")

		req, err := http.NewRequest("POST", u.String(), bytes.NewReader(reqBodyBytes))
		if err != nil {
			log.Fatalf("Failed to create agent registration request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", installConf.RegistrationKey)

		resp, err := makeHttpClient(*installConf).Do(req)
		if err != nil {
			log.Fatalf("Failed to register agent: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to register agent: got status code %d", resp.StatusCode)
		}
		
		var respBody apitypes.AgentRegisterResponseDTO
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			log.Fatalf("Failed to decode agent registration response: %v", err)
		}

		log.Print("Registered agent with server as " + respBody.Name + " with ID " + respBody.ID)

		installConf.AuthKey = respBody.Key
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

	// Blank API Endpoint is fine, assuming user is maybe setting up agents before server, but if they're looking for us to install hashcat we have to be up and running
	if installConf.APIEndpoint == "" && installConf.InstallHashcat {
		return errors.New("install hashcat was selected, but no download (api endpoint) server was specified")
	}

	return nil
}

func RunInteractive() {
	flagSet := flag.NewFlagSet("install", flag.ExitOnError)

	useDefaultsP := flagSet.Bool("defaults", false, "Use basic defaults for intallation? (stores everything in "+DefaultPathJoin()+")")

	userP := flagSet.String("user", "", "Which user to run the agent as")
	groupP := flagSet.String("group", "", "Which user group to run the agent as")
	agentBinPathP := flagSet.String("agent-bin", "", "Path to agent (defaults to running executable)")
	registrationKeyP := flagSet.String("registration-key", "", "Registration key for agent")
	authKeyFileP := flagSet.String("auth-keyfile", "", "Path to file containing agent key")
	authKeyP := flagSet.String("auth-key", "", "Path to file containing agent key")
	configPathP := flagSet.String("config-path", "", "Path to config json file to install")
	hashcatPathP := flagSet.String("hashcat-path", "", "Path to hashcat executable")
	listfilePathP := flagSet.String("listfile-directory", "", "Path to directory to hold listfiles")
	apiEndpointP := flagSet.String("api-endpoint", "", "API endpoint (format: https://phatcrack.lan/api/v1)")

	autoInstallHashcatP := flagSet.Bool("download-hashcat", false, "Install hashcat from agent asset server (requires api endpoint to be set)")

	disableTLSVerificationP := flagSet.Bool("disable-tls-verification", false, "Whether to disable TLS Verification")

	flagSet.Parse(os.Args[2:])

	installConf := InstallConfig{
		Defaults: *useDefaultsP,

		AgentUser:    *userP,
		AgentGroup:   *groupP,
		AgentBinPath: *agentBinPathP,

		RegistrationKey: *registrationKeyP,
		AuthKey:     *authKeyP,
		AuthKeyFile: *authKeyFileP,
		ConfigPath:  *configPathP,

		HashcatPath:            *hashcatPathP,
		ListfileDirectory:      *listfilePathP,
		APIEndpoint:            *apiEndpointP,
		DisableTLSVerification: *disableTLSVerificationP,
		InstallHashcat:         *autoInstallHashcatP,
	}

	if installConf.RegistrationKey != "" && installConf.AuthKey != "" {
		log.Fatal("Registration key and auth key cannot be set at the same time")
	}

	if installConf.Defaults {
		applyDefaults(&installConf)
	} else {
		getOptionsInteractive(&installConf)
		applyDefaults(&installConf)
	}

	registerIfRequired(&installConf)

	err := checkConf(installConf)
	if err != nil {
		log.Fatal("config was invalid: ", err)
	}

	Run(installConf)
}

func DefaultPathJoin(parts ...string) string {
	path := "/opt/phatcrack-agent"
	if runtime.GOOS == "windows" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	parts = append([]string{path}, parts...)
	return filepath.Join(parts...)

}

func adjustPermsPath(path string, installConf InstallConfig) {
	// No op on windows
	f, _ := os.OpenFile(path, os.O_WRONLY, 0700)
	adjustPerms(f, installConf)
}
