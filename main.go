// Package main provides the very simplest main function
package main

import (
	"log"
	"os"

	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/worker"
)

const databaseName = "helldivers2"

func main() {
	cfg := config.Get()
	logger := loggerFor("main")

	dbClient, err := db.New(cfg.MongoURI, databaseName, loggerFor("mongo"))
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		logger.Println(dbClient.Disconnect())
	}()
	if err = dbClient.MigrateUp("./migrations"); err != nil {
		logger.Fatal(err)
	}

	apiClient, err := client.New(cfg.APIRootURL, loggerFor("api"))
	if err != nil {
		logger.Fatal(err)
	}

	worker := worker.New(apiClient, dbClient, loggerFor("worker"))
	// TODO: catch interrupt
	worker.Run(cfg.WorkerInterval)
}

func loggerFor(name string) *log.Logger {
	return log.New(os.Stdout, name+" | ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
