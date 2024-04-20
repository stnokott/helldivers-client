package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // use pgx as driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // load migrations from file
)

func (c *Client) newMigration(scriptFolder string) (*migrate.Migrate, error) {
	cfg := c.conn.Config()
	uri := fmt.Sprintf("pgx5://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	migration, err := migrate.New(
		"file://"+scriptFolder,
		uri,
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
func (c *Client) MigrateUp(migrationsFolder string) error {
	migration, err := c.newMigration(migrationsFolder)
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
