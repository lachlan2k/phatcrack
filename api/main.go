package main

import (
	"net/url"
	"os"

	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/webserver"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

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
		log.Warnf("HC_PATH was not specified, some API endpoints may not work if hashcat is not in PATH\n")
	}

	if os.Getenv("FILEREPO_PATH") == "" {
		log.Fatalf("FILEREPO_PATH was not specified")
	}

	if os.Getenv("BASE_URL") == "" {
		log.Fatalf("BASE_URL was not specified")
	}

	parsedBaseURL, err := url.Parse(os.Getenv("BASE_URL"))
	if err != nil {
		log.Fatalf("Provided BASE_URL could not be parsed: %v", err)
	}

	auth.InitWebAuthn(*parsedBaseURL)

	err = filerepo.SetPath(os.Getenv("FILEREPO_PATH"))
	if err != nil {
		log.Fatalf("failed to use specified FILEREPO_PATH: %v", err)
	}

	err = fleet.Setup()
	if err != nil {
		log.Fatal(err)
	}

	err = webserver.Listen(port)
	if err != nil {
		log.Fatalf("couldn't run server: %v", err)
	}
}
