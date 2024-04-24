package db

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stnokott/helldivers-client/internal/copytest"
)

var validEventCampaign = Campaign{
	ID:    123,
	Type:  55,
	Count: 678,
}
var validEvent = Event{
	CampaignID: 123,
	ID:         555,
	Type:       7,
	Faction:    "Automatons",
	MaxHealth:  55667788,
	StartTime:  PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
	EndTime:    PGTimestamp(time.Date(2025, 1, 1, 1, 1, 1, 1, time.UTC)),
}

func TestEventsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Event)
	tests := []struct {
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(*Event) {},
			wantErr:  false,
		},
		{
			name: "max health high number",
			modifier: func(e *Event) {
				e.MaxHealth = math.MaxInt64
			},
			wantErr: false,
		},
		{
			name: "negative max health",
			modifier: func(e *Event) {
				e.MaxHealth = -1
			},
			wantErr: true,
		},
		{
			name: "empty faction",
			modifier: func(e *Event) {
				e.Faction = ""
			},
			wantErr: true,
		},
		{
			name: "start time after end time",
			modifier: func(e *Event) {
				e.StartTime = PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 2, time.UTC))
				e.EndTime = PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC))
			},
			wantErr: true,
		},
		{
			name: "campaign FK violation",
			modifier: func(e *Event) {
				e.CampaignID++
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

				var (
					campaign Campaign
					event    Event
				)
				if err := copytest.DeepCopy(
					&campaign, &validEventCampaign,
					&event, &validEvent,
				); err != nil {
					t.Errorf("failed to create event struct copy: %v", err)
					return
				}

				tt.modifier(&event)

				if err := campaign.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to merge campaign (required for event): %v", err)
					return
				}

				err := event.Merge(context.Background(), client.queries, tableMergeStats{})
				if (err != nil) != tt.wantErr {
					t.Errorf("Event.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
					return
				}

				fetchedResult, err := client.queries.GetEvent(context.Background(), event.ID)
				if err != nil {
					t.Errorf("failed to fetch inserted event: %v", err)
					return
				}
				if fetchedResult != event.ID {
					t.Errorf("failed to validate INSERT: inserted data has ID %d, DB returned %d", event.ID, fetchedResult)
				}
			})
		})
	}
}
