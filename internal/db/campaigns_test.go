package db

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

var validCampaign = Campaign{
	ID:       5,
	PlanetID: 6,
	Type:     8,
	Count:    100,
}

func TestCampaignsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Campaign)
	tests := []struct {
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(p *Campaign) {},
			wantErr:  false,
		},
		{
			name: "invalid planet reference",
			modifier: func(c *Campaign) {
				c.PlanetID = 99999
			},
			wantErr: true,
		},
		{
			name: "negative count",
			modifier: func(c *Campaign) {
				c.Count = -1
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withClient(t, func(client *Client, migration *migrate.Migrate) {
				if err := migration.Up(); err != nil {
					t.Errorf("failed to migrate up: %v", err)
					return
				}

				campaign := validCampaign
				tt.modifier(&campaign)
				_, err := client.queries.InsertCampaign(context.Background(), gen.InsertCampaignParams(campaign))
				if err != nil {
					return
				}
				fetchedResult, err := client.queries.GetCampaign(context.Background(), campaign.ID)
				if err != nil {
					t.Errorf("failed to fetch inserted campaign: %v", err)
					return
				}
				if fetchedResult != campaign.ID {
					t.Errorf("failed to validate INSERT: inserted data has ID %d, DB returned %d", campaign.ID, fetchedResult)
				}
			})
		})
	}
}
