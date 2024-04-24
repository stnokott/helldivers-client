package db

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // use pgx as driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // load migrations from file
)

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
	"wars",
	"snapshots",
	"war_snapshots",
	"event_snapshots",
	"planet_snapshots",
	"snapshot_statistics",
}

func TestTablesExist(t *testing.T) {
	withClientMigrated(t, func(client *Client) {
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
	})
}
