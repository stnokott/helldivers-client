// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stnokott/helldivers-client/internal/config"
)

func TestMain(m *testing.M) {
	envFile := "../../.env.test"
	// we only try to load .env.test if it is present.
	// The usecase for this is local development when running through Docker is not available.
	// Env variables will then be supplied through the env file instead of the Docker container.
	// This is required because VSCode tasks.json doesn't allow loading from a .env file.
	if _, err := os.Stat(envFile); err == nil {
		log.Printf("using env file for tests: %s", envFile)
		if err = godotenv.Load(envFile); err != nil {
			log.Fatalf("could not load env file for tests: %v", err)
		}
	}
	code := m.Run()
	os.Exit(code)
}

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
		{name: "invalid", args: args{&config.Config{MongoURI: "http://localhost"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.args.cfg, tt.name+" db", logger)
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
	client, err := New(cfg, "test_client_disconnect", log.Default())
	if err != nil {
		t.Fatalf("could not initialize DB connection: %v", err)
	}
	if err := client.Disconnect(); err != nil {
		t.Fatalf("Disconnect() error = %v, want nil", err)
	}
	if err := client.Disconnect(); err == nil {
		t.Fatalf("Disconnect() (while not connected) error = nil, want err")
	}
}
