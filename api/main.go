package main

import (
	"log"
	"os"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/webserver"
)

func main() {
	if os.Getenv("MONGO_URI") == "" {
		log.Fatal("MONGO_URI was not specified")
	}

	err := db.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	err = webserver.Listen(port)
	if err != nil {
		log.Fatalf("couldn't run server: %v", err)
	}
}
