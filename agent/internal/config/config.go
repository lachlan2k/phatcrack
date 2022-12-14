package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Config struct {
	AuthKeyFile        string `json:"auth_key_file"`
	AuthKey            string `json:"auth_key"`
	HashcatBinary      string `json:"hashcat_binary"`
	RulesDirectory     string `json:"rules_directory"`
	WordlistsDirectory string `json:"wordlists_directory"`
	WSEndpoint         string `json:"ws_endpoint"`
}

func LoadConfig(configPath string) (config Config) {
	configJSON, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("couldn't read config file (%s): %v", configPath, err)
	}

	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Fatalf("error when unmarshalling config json: %v", err)
	}

	if config.AuthKey == "" {
		if config.AuthKeyFile == "" {
			log.Fatalf("Neither auth_key nor auth_key_file was provided")
		}

		authKeyBytes, err := os.ReadFile(config.AuthKeyFile)
		if err != nil {
			log.Fatalf("couldn't read provided auth key file (%s): %v", config.AuthKeyFile, err)
		}

		config.AuthKey = strings.TrimSpace(string(authKeyBytes))
	}

	return
}
