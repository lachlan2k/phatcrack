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
	if len(os.Args) >= 2 && os.Args[1] == "install" {
		installer.RunInteractive()
		return
	}

	configPath := flag.String("config", "/opt/phatcrack-agent/config.json", "Location of config file")
	versionP := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *versionP {
		fmt.Printf("Phatcrack Agent version %s", version.Version())
		return
	}

	conf := config.LoadConfig(*configPath)

	log.Printf("Starting phatcrack-agent " + version.Version())
	err := handler.Run(&conf)

	if err != nil {
		log.Fatalf("%v", err)
	}
}
