package db

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stnokott/helldivers-client/internal/copytest"
)

var validDispatch = Dispatch{
	ID:         123,
	CreateTime: PGTimestamp(time.Date(2024, 1, 2, 3, 4, 5, 6, time.UTC)),
	Type:       5,
	Message:    "A valid dispatch",
}

func TestDispatchesSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Dispatch)
	tests := []struct {
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(*Dispatch) {},
			wantErr:  false,
		},
		{
			name: "empty message",
			modifier: func(d *Dispatch) {
				d.Message = ""
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

				var dispatch Dispatch
				if err := copytest.DeepCopy(&dispatch, &validDispatch); err != nil {
					t.Errorf("failed to create dispatch struct copy: %v", err)
					return
				}

				tt.modifier(&dispatch)

				err := dispatch.Merge(context.Background(), client.queries, tableMergeStats{})
				if (err != nil) != tt.wantErr {
					t.Errorf("Dispatch.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
					return
				}

				fetchedResult, err := client.queries.GetDispatch(context.Background(), dispatch.ID)
				if err != nil {
					t.Errorf("failed to fetch inserted dispatch: %v", err)
					return
				}
				if fetchedResult != dispatch.ID {
					t.Errorf("failed to validate INSERT: inserted data has ID %d, DB returned %d", dispatch.ID, fetchedResult)
				}
			})
		})
	}
}
