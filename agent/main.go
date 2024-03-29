package main

import (
	"flag"
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
	flag.Parse()

	conf := config.LoadConfig(*configPath)

	log.Printf("Starting phatcrack-agent " + version.Version())
	err := handler.Run(&conf)

	if err != nil {
		log.Fatalf("%v", err)
	}
}
