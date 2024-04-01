package db

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file" // load migrations from file
)

func (c *Client) newMigration() (*migrate.Migrate, error) {
	// create new migration instance from existing connection
	driver, err := mongodb.WithInstance(c.mongo, &mongodb.Config{
		DatabaseName: c.dbName,
	})
	if err != nil {
		return nil, err
	}
	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		c.dbName,
		driver,
	)
	if err != nil {
		return nil, err
	}
	migration.Log = &migrationLogger{c.log}
	return migration, nil
}

// MigrateUp runs required migrations to get to the latest version.
//
// This should be run before any other operations.
func (c *Client) MigrateUp() error {
	migration, err := c.newMigration()
	if err != nil {
		return err
	}
	if err = migration.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// migrationLogger wraps log.Logger for usage with migrate package
type migrationLogger struct {
	*log.Logger
}

func (l *migrationLogger) Verbose() bool {
	return false
}
