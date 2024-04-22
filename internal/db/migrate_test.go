package db

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // use pgx as driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // load migrations from file
	"github.com/stnokott/helldivers-client/internal/config"
)

func withClient(t *testing.T, do func(client *Client, migration *migrate.Migrate)) {
	cfg := config.Get()

	client, err := New(cfg, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatalf("could not initialize DB connection: %v", err)
	}
	defer func() {
		if err = client.Disconnect(); err != nil {
			t.Logf("failed to disconnect: %v", err)
		}
	}()
	migration, err := client.newMigration("../../scripts/migrations")
	if err != nil {
		t.Fatalf("client.newMigration() error = %v, want nil", err)
	}
	defer func() {
		if err = migration.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			t.Fatalf("failed to migrate down: %v", err)
		}
	}()
	do(client, migration)
}

func TestMigrateUp(t *testing.T) {
	withClient(t, func(client *Client, _ *migrate.Migrate) {
		if err := client.MigrateUp("../../scripts/migrations"); err != nil {
			t.Errorf("client.MigrateUp() error = %v, expected nil", err)
		}
	})
}

var tableNames = []string{
	"planets",
	"biomes",
	"hazards",
	"campaigns",
	"dispatches",
	"events",
	"assignments",
	"assignment_tasks",
	"wars",
	"snapshots",
	"war_snapshots",
	"event_snapshots",
	"planet_snapshots",
	"snapshot_statistics",
}

func TestTablesExist(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		if err := migration.Up(); err != nil {
			t.Errorf("failed to migrate up: %v", err)
			return
		}

		fnTables := func() []string {
			rows, err := client.conn.Query(
				context.Background(),
				`SELECT table_name FROM information_schema.tables WHERE table_name = any($1);`,
				tableNames,
			)
			if err != nil {
				t.Errorf("could not list tables: %v", err)
				return []string{}
			}
			defer rows.Close()

			names := []string{}
			for rows.Next() {
				var tableName string
				rows.Scan(&tableName)
				names = append(names, tableName)
			}
			return names
		}
		if colls := fnTables(); len(colls) != len(tableNames) {
			t.Errorf("expected %d tables, got %d (%v)", len(tableNames), len(colls), colls)
			return
		}
		if err := migration.Down(); err != nil {
			t.Errorf("failed to migrate down: %v", err)
			return
		}
		if colls := fnTables(); len(colls) > 0 {
			t.Errorf("expected no tables, got %d (%v)", len(colls), colls)
			return
		}
	})
}
