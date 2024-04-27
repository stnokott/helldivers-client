// Package main provides the very simplest main function
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/worker"
)

var (
	projectName string
	version     string
	commit      string
	buildDate   string
)

const databaseName = "helldivers2"

func main() {
	fmt.Printf("%s v%s %s built %s\n\n", projectName, version, commit, buildDate)

	cfg := config.MustGet()
	logger := loggerFor("main")

	// TODO: wait until DB available

	dbClient, err := db.New(cfg, loggerFor("postgresql"))
	if err != nil {
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

	// TODO: wait until API available

	apiClient, err := client.New(cfg, loggerFor("api"))
	if err != nil {
		logger.Fatal(err)
	}

	worker := worker.New(apiClient, dbClient, loggerFor("worker"))
	stopWorkerChan := make(chan struct{}, 1)

	stopSignal(stopWorkerChan, logger)
	worker.Run(cfg.WorkerInterval, stopWorkerChan)
}

func stopSignal(stopChan chan<- struct{}, logger *log.Logger) {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt)
	go func() {
		s := <-osSignalChan
		logger.Printf("received %s signal, stopping once current process finishes", s.String())
		stopChan <- struct{}{}
	}()
}

func loggerFor(name string) *log.Logger {
	return log.New(os.Stdout, name+" | ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
