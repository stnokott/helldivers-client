// Package main provides the very simplest main function
package main

import (
	"log"
	"os"

	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
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
}

func loggerFor(name string) *log.Logger {
	return log.New(os.Stdout, name+" | ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
