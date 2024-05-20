//go:build integration

package db

import (
	"context"
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

var validCampaign = Campaign{
	ID:    5,
	Type:  8,
	Count: PGUint64(100),
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
			name: "negative count",
			modifier: func(c *Campaign) {
				c.Count = pgtype.Numeric{Int: big.NewInt(-1), Valid: true}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withClientMigrated(t, func(client *Client) {
				var campaign Campaign
				if err := copytest.DeepCopy(&campaign, &validCampaign); err != nil {
					t.Errorf("failed to create campaign struct copy: %v", err)
					return
				}

				tt.modifier(&campaign)

				err := campaign.Merge(context.Background(), client.queries, func(gen.Table, bool, int64) {})
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
