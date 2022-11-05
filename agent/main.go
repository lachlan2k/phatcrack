package main

import (
	"flag"
	"log"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/handler"
	_ "github.com/lachlan2k/phatcrack/agent/pkg/apitypes"
)

func main() {
	configPath := flag.String("config", "/etc/phatcrack-agent/config.json", "Location of config file")
	flag.Parse()

	conf := config.LoadConfig(*configPath)

	log.Printf("Starting agent")
	for {
		err := handler.Run(&conf)
		if err != nil {
			log.Printf("Error when running agent, reconnecting: %v", err)
		}
		time.Sleep(time.Second)
	}
}
