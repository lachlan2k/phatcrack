package main

import (
	"flag"
	"log"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/handler"
)

func main() {
	configPath := flag.String("config", "/etc/phatcrack-agent/config.json", "Location of config file")
	flag.Parse()

	conf := config.LoadConfig(*configPath)

	log.Printf("Starting agent")
	handler.Run(&conf)
}
