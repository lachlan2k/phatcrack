package installer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func RegisterWithKey(conf *InstallConfig) (*apitypes.AgentRegisterResponseDTO, error) {
	u, err := url.Parse(conf.APIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API endpoint to register agent: %v", err)
	}

	if conf.Name == "" {
		name, err := os.Hostname()
		if err != nil || name == "" {
			log.Printf("Warn: couldn't get hostname for agent registration: %v\n", err)
			name = "unknown"
		}
		conf.Name = name
	}

	reqBody := apitypes.AgentRegisterRequestDTO{
		Name: conf.Name,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent registration request: %v", err)
	}

	u.Path = path.Join(u.Path, "/agent-handler/register")

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create agent registration request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", conf.RegistrationKey)

	resp, err := makeHttpClient(conf.DisableTLSVerification).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to register agent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to register agent: got status code %d", resp.StatusCode)
	}

	var respBody apitypes.AgentRegisterResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode agent registration response: %v", err)
	}

	log.Print("Registered agent with server as " + respBody.Name + " with ID " + respBody.ID)
	return &respBody, nil
}

func registerIfRequired(installConf *InstallConfig) {
	if installConf.RegistrationKey != "" {
		_, err := RegisterWithKey(installConf)
		if err != nil {
			log.Fatalf("failed to register agent: %v", err)
		}
	}
}