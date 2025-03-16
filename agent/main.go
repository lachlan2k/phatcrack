package main

import (
	"flag"
	"fmt"
	"os"

	"log"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/handler"
	"github.com/lachlan2k/phatcrack/agent/internal/installer"
	"github.com/lachlan2k/phatcrack/agent/internal/version"
)

func main() {

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "install":
			{
				installer.RunInteractive()
				return
			}

		case "register":
			{
				flagSet := flag.NewFlagSet("register", flag.ExitOnError)
				keyP := flagSet.String("key", "", "Registration key (to exchange for auth key)")
				apiEndpointP := flagSet.String("api-endpoint", "", "API endpoint (format: https://phatcrack.lan/api/v1)")
				disableTLSVerificationP := flagSet.Bool("disable-tls-verification", false, "Whether to disable TLS Verification")
				nameP := flagSet.String("name", "", "Name of the agent (defaults to hostname)")

				flagSet.Parse(os.Args[2:])

				if *keyP == "" {
					log.Fatal("No key provided")
				}
				if *apiEndpointP == "" {
					log.Fatal("No API endpoint provided")
				}

				conf := installer.InstallConfig{
					RegistrationKey: *keyP,
					APIEndpoint: *apiEndpointP,
					DisableTLSVerification: *disableTLSVerificationP,
					Name: *nameP,
				}

				resp, err := installer.RegisterWithKey(&conf)
				if err != nil {
					log.Fatalf("failed to register agent: %v", err)
				}

				fmt.Printf("\nAgent registered!\n\nID: %s\nName: %s\nAuth key: %s\n\n", resp.ID, resp.Name, resp.Key)

				return
			}

		default:
			{
				run()
			}
		}
	}

}

func run() {
	configPath := flag.String("config", installer.DefaultPathJoin("config.json"), "Location of config file")
	versionP := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *versionP {
		fmt.Printf("Phatcrack Agent version %s", version.Version())
		return
	}

	conf := config.LoadConfig(*configPath)

	log.Printf("Starting phatcrack-agent %s", version.Version())
	err := handler.Run(&conf)

	if err != nil {
		log.Fatalf("%v", err)
	}
}
