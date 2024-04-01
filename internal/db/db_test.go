// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"log"
	"testing"
)

const (
	mongoURI              = "mongodb://root:test@db:27017/"  //NOSONAR
	mongoURIInvalidScheme = "http://root:test@db:27017"      //NOSONAR
	mongoURIInvalidAuth   = "mongodb://root:pass@db:27017"   //NOSONAR
	mongoURIInvalidHost   = "mongodb://root:test@host:27017" //NOSONAR
	mongoURIInvalidPort   = "mongodb://root:test@db:55555"   //NOSONAR
)

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
		{name: "valid", args: args{uri: mongoURI}, wantErr: false},
		{name: "invalid scheme", args: args{uri: mongoURIInvalidScheme}, wantErr: true},
		{name: "invalid auth", args: args{uri: mongoURIInvalidAuth}, wantErr: true},
		{name: "invalid host", args: args{uri: mongoURIInvalidHost}, wantErr: true},
		{name: "invalid port", args: args{uri: mongoURIInvalidPort}, wantErr: true},
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
	client, err := New(mongoURI, "test_client_disconnect", log.Default())
	if err != nil {
		t.Skipf("could not initialize DB connection: %v", err)
	}
	if err := client.Disconnect(); err != nil {
		t.Fatalf("Disconnect() error = %v, want nil", err)
	}
	if err := client.Disconnect(); err == nil {
		t.Fatalf("Disconnect() (while not connected) error = nil, want err")
	}
}
