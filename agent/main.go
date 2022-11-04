package main

import (
	"flag"
	"log"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"golang.org/x/net/websocket"
)

func main() {
	configPath := flag.String("config", "/etc/phatcrack-agent/config.json", "Location of config file")
	flag.Parse()

	conf := config.LoadConfig(*configPath)

	p := hashcat.HashcatParams{
		AttackMode:        0,
		HashType:          0,
		WordlistFilenames: []string{"a.txt"},
		OptimizedKernels:  true,
	}

	err := hashcat.RunHashcat([]string{"2ab96390c7dbe3439de74d0c9b0b1767"}, p, conf)
	if err != nil {
		panic(err)
	}

	wsConfig, err := websocket.NewConfig(conf.Endpoint, "http://dummy-origin")
	if err != nil {
		log.Fatalf("couldn't create websocket config: %v", err)
	}

	wsConfig.Header.Add("Authorization", "bearer "+conf.AuthKey)

	time.Sleep(time.Minute)

	log.Printf("Starting agent")
	for {
		err := run(wsConfig, conf)
		if err != nil {
			log.Printf("Error when running agent, reconnecting: %v", err)
		}
		time.Sleep(time.Second)
	}
}
