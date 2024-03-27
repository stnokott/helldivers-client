// Package client wraps the API specs into a client
package client

import (
	"context"
	"testing"
	"time"
)

const (
	host string = "http://api:4000"
)

func TestClientSeasons(t *testing.T) {
	client, err := New(host)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Seasons(ctx)
	if err != nil {
		t.Fatalf("client.Seasons() = %v", err)
	}
	if resp.Current == nil {
		t.Errorf(".Current is nil")
	}
	if resp.Seasons == nil {
		t.Errorf(".Seasons is nil")
	}
}

func TestClientHosts(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{
			name:    "trailing slash",
			host:    host + "/",
			wantErr: false,
		},
		{
			name:    "no trailing slash",
			host:    host,
			wantErr: false,
		},
		{
			name:    "wrong endpoint",
			host:    host + "/api",
			wantErr: true,
		},
		{
			name:    "wrong host",
			host:    "127.0.0.1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.host)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err = client.Seasons(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Seasons() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
