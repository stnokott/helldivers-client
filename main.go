// Package main provides the very simplest main function
package main

import (
	"log"
	"os"

	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
)

func main() {
	cfg := config.Get()

	dbClient, err := db.New(cfg.MongoURI, loggerFor("mongo"))
	if err != nil {
		log.Fatalf("MongoDB client could not be initialized: %v", err)
	}
	defer dbClient.Disconnect()
	if err = dbClient.PrepareDB(); err != nil {
		log.Fatalf("db preparation failed: %v", err)
	}
}

func loggerFor(name string) *log.Logger {
	return log.New(os.Stdout, name+" | ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
