// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"errors"
	"io"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stnokott/helldivers-client/internal/config"
)

func TestNew(t *testing.T) {
	logger := log.Default()
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid", args: args{config.Get()}, wantErr: false},
		{name: "invalid", args: args{&config.Config{PostgresURI: "http://localhost"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.args.cfg, logger)
			defer func() {
				if client != nil {
					client.Disconnect()
				}
			}()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClientDisconnect(t *testing.T) {
	cfg := config.Get()
	client, err := New(cfg, log.Default())
	if err != nil {
		t.Fatalf("could not initialize DB connection: %v", err)
	}
	if err := client.Disconnect(); err != nil {
		t.Fatalf("Disconnect() error = %v, want nil", err)
	}
	if err := client.Disconnect(); err != nil {
		t.Fatalf("Disconnect() (while not connected) error = %v, want nil", err)
	}
}

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
		return
	}
	defer func() {
		if err = migration.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			t.Fatalf("failed to migrate down: %v", err)
			return
		}
	}()
	do(client, migration)
}

func withClientMigrated(t *testing.T, do func(client *Client)) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		if err := migration.Up(); err != nil {
			t.Errorf("failed to migrate up: %v", err)
			return
		}
		do(client)
	})
}
