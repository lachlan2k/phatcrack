package main

import (
	"log"
	"os"

	"github.com/lachlan2k/phatcrack/api/internal/dbnew"
	"github.com/lachlan2k/phatcrack/api/internal/webserver"
)

func main() {
	if os.Getenv("DB_DSN") == "" {
		log.Fatal("DB_DSN was not specified")
	}

	err := dbnew.Connect(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
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
