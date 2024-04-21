package db

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

var validWar = War{
	ID:        999,
	StartTime: PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local)),
	EndTime:   PGTimestamp(time.Date(2025, 1, 1, 1, 1, 1, 1, time.Local)),
	Factions:  []string{"Humans", "Automatons"},
}

func TestWarsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*War)
	tests := []struct {
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(*War) {},
			wantErr:  false,
		},
		{
			name: "start time after end time",
			modifier: func(w *War) {
				w.StartTime = PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 2, time.Local))
				w.EndTime = PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local))
			},
			wantErr: true,
		},
		{
			name: "empty faction list",
			modifier: func(w *War) {
				w.Factions = []string{}
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
				war := validWar
				tt.modifier(&war)

				err := war.Merge(context.Background(), client.queries, tableMergeStats{})
				if (err != nil) != tt.wantErr {
					t.Errorf("War.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
					return
				}

				fetchedResult, err := client.queries.GetWar(context.Background(), war.ID)
				if err != nil {
					t.Errorf("failed to fetch inserted war: %v", err)
					return
				}
				if fetchedResult != war.ID {
					t.Errorf("failed to validate INSERT: inserted data has ID %d, DB returned %d", war.ID, fetchedResult)
				}
			})
		})
	}
}
