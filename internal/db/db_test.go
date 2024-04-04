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
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("could not load %s: %v", envFile, err)
	}
	code := m.Run()
	os.Exit(code)
}

// getMongoURI reads the config from ENV and returns the mongo URI inside
func getMongoURI() string {
	return config.Get().MongoURI
}

func TestNew(t *testing.T) {
	logger := log.Default()
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid", args: args{uri: getMongoURI()}, wantErr: false},
		{name: "invalid", args: args{uri: "http://localhost"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.args.uri, tt.name+" db", logger)
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
	mongoURI := getMongoURI()
	client, err := New(mongoURI, "test_client_disconnect", log.Default())
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
