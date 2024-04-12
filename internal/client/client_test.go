// Package client wraps the API specs into a client
package client

import (
	"context"
	"log"
	"testing"

	"github.com/stnokott/helldivers-client/internal/config"
)

var logger = log.Default()

func mustClient() *Client {
	config := config.Get()

	client, err := New(config, logger)
	if err != nil {
		panic(err)
	}
	return client
}

// global client used for rate-limiting across tests.
//
// FIXME: not required anymore once we can disable rate limiting in API container
var globalClient = mustClient()

func TestClientHosts(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO

	host := config.Get().APIRootURL
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "trailing slash",
			cfg:     &config.Config{APIRootURL: host + "/"},
			wantErr: false,
		},
		{
			name:    "no trailing slash",
			cfg:     &config.Config{APIRootURL: host},
			wantErr: false,
		},
		{
			name:    "wrong endpoint",
			cfg:     &config.Config{APIRootURL: host + "/api"},
			wantErr: true,
		},
		{
			name:    "wrong host",
			cfg:     &config.Config{APIRootURL: "127.0.0.1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := New(tt.cfg, logger)
			_, err := client.War(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.War() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClientWarId(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	got, err := globalClient.WarID(context.Background())
	if err != nil {
		t.Errorf("Client.WarID() error = %v, want nil", err)
		return
	}
	if got == nil {
		t.Error("Client.WarID() returned nil, want non-nil")
		return
	}
	if got.Id == nil || *got.Id == 0 {
		t.Error("got.ClientVersion is empty, expected non-empty")
		return
	}
}

func TestClientWar(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	got, err := globalClient.War(context.Background())
	if err != nil {
		t.Errorf("Client.War() error = %v, want nil", err)
		return
	}
	if got == nil {
		t.Error("Client.War() returned nil, want non-nil")
		return
	}
	if got.ClientVersion == nil || *got.ClientVersion == "" {
		t.Error("got.ClientVersion is empty, expected non-empty")
		return
	}
}

func TestClientAssignments(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	got, err := globalClient.Assignments(context.Background())
	if err != nil {
		t.Errorf("Client.Assignments() error = %v, want nil", err)
		return
	}
	if got == nil {
		t.Error("Client.Assignments() returned nil, want non-nil")
		return
	}
	firstItem := (*got)[0]
	if firstItem.Title == nil || *firstItem.Title == "" {
		t.Error("got[0].Title is empty, expected non-empty")
		return
	}
}

func TestClientCampaigns(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	got, err := globalClient.Campaigns(context.Background())
	if err != nil {
		t.Errorf("Client.Campaigns() error = %v, want nil", err)
		return
	}
	if got == nil {
		t.Error("Client.Campaigns() returned nil, want non-nil")
		return
	}
	firstItem := (*got)[0]
	if firstItem.Id == nil || *firstItem.Id == 0 {
		t.Error("got[0].Id is empty, expected non-empty")
		return
	}
}

func TestClientDispatches(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	got, err := globalClient.Dispatches(context.Background())
	if err != nil {
		t.Errorf("Client.Dispatches() error = %v, want nil", err)
		return
	}
	if got == nil {
		t.Error("Client.Dispatches() returned nil, want non-nil")
		return
	}
	firstItem := (*got)[0]
	if firstItem.Message == nil || *firstItem.Message == "" {
		t.Error("got[0].Message is empty, expected non-empty")
		return
	}
}

func TestClientPlanets(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	got, err := globalClient.Planets(context.Background())
	if err != nil {
		t.Errorf("Client.Planets() error = %v, want nil", err)
		return
	}
	if got == nil {
		t.Error("Client.Planets() returned nil, want non-nil")
		return
	}
	firstItem := (*got)[0]
	if firstItem.CurrentOwner == nil || *firstItem.CurrentOwner == "" {
		t.Error("got[0].CurrentOwner is empty, expected non-empty")
		return
	}
}
