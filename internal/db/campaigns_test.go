package db

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/copier"
)

var validCampaign = Campaign{
	ID:       5,
	PlanetID: 6,
	Type:     8,
	Count:    100,
}
var validCampaignPlanet = validPlanet

func TestCampaignsSchema(t *testing.T) {
	// synchronize planet IDs so we have a valid starting point
	validCampaignPlanet.ID = validCampaign.PlanetID

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

				if err := validCampaignPlanet.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to merge campaign planet (check planet tests): %v", err)
					return
				}

				var campaign Campaign
				// deep copy will copy values behind pointers instead of the pointers themselves
				copyOption := copier.Option{DeepCopy: true}
				if err := copier.CopyWithOption(&campaign, &validCampaign, copyOption); err != nil {
					t.Errorf("failed to create campaign struct copy: %v", err)
					return
				}

				tt.modifier(&campaign)

				err := campaign.Merge(context.Background(), client.queries, tableMergeStats{})
				if (err != nil) != tt.wantErr {
					t.Errorf("Campaign.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
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
