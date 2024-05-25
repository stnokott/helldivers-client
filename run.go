// Package main provides the very simplest main function
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/worker"
)

var (
	projectName = "helldivers-client"
	version     = "0.0.0"
	commit      = "dev"
	buildDate   = "now"
)

const (
	databaseName   = "helldivers2"
	dbReadyTimeout = 30 * time.Second
)

const apiReadyTimeout = 30 * time.Second

func run(stopChan <-chan struct{}) {
	fmt.Printf("%s v%s %s built %s\n\n", projectName, version, commit, buildDate)

	cfg := config.MustGet()
	logger := loggerFor("main")

	dbClient, err := db.New(cfg, loggerFor("postgresql"))
	if err != nil {
		logger.Fatal(err)
	}
	if err = waitFor(dbClient, dbReadyTimeout, logger); err != nil {
		logger.Fatal(err)
	}

	defer func() {
		if errInner := dbClient.Disconnect(); errInner != nil {
			logger.Println(errInner)
		}
	}()
	if err = dbClient.MigrateUp("./scripts/migrations"); err != nil {
		logger.Fatal(err)
	}

	apiClient, err := client.New(cfg, loggerFor("api"))
	if err != nil {
		logger.Fatal(err)
	}
	if err = waitFor(apiClient, apiReadyTimeout, logger); err != nil {
		logger.Fatal(err)
	}

	worker, err := worker.New(apiClient, dbClient, cfg, loggerFor("worker"))
	if err != nil {
		logger.Fatal(err)
	}

	worker.Run(cfg.WorkerInterval, stopChan)
}

func loggerFor(name string) *log.Logger {
	return log.New(os.Stdout, name+" | ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
