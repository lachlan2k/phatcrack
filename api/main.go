package main

import (
	"log"
	"os"

	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/webserver"
)

func main() {
	if os.Getenv("DB_DSN") == "" {
		log.Fatal("DB_DSN was not specified")
	}

	err := db.Connect(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	err = config.Reload()
	if err != nil {
		log.Fatalf("failed to load runtime config: %v", err)
	}

	err = config.Save()
	if err != nil {
		log.Fatalf("failed to write config to db for consistency: %v", err)
	}

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if os.Getenv("HC_PATH") == "" {
		log.Printf("HC_PATH was not specified, some API endpoints may not work if hashcat is not in PATH\n")
	}

	err = webserver.Listen(port)
	if err != nil {
		log.Fatalf("couldn't run server: %v", err)
	}
}
