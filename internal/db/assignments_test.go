//go:build integration

package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

var validAssignment = Assignment{
	Assignment: gen.Assignment{
		ID:           3,
		Title:        "Footitle",
		Briefing:     "Foobriefing",
		Description:  "Bardescription",
		Expiration:   PGTimestamp(time.Date(2024, 1, 2, 3, 4, 5, 6, time.UTC)),
		RewardType:   8,
		RewardAmount: PGUint64(100),
	},
	Tasks: []gen.AssignmentTask{
		{
			TaskType:   9,
			Values:     []pgtype.Numeric{PGUint64(7), PGUint64(8), PGUint64(9)},
			ValueTypes: []pgtype.Numeric{PGUint64(42), PGUint64(44), PGUint64(46)},
		},
	},
}

func TestAssignmentsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Assignment)
	tests := []struct {
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(*Assignment) {},
			wantErr:  false,
		},
		{
			name: "empty required title",
			modifier: func(a *Assignment) {
				a.Title = ""
			},
			wantErr: true,
		},
		{
			name: "empty required expiration time",
			modifier: func(a *Assignment) {
				a.Expiration = pgtype.Timestamp{}
			},
			wantErr: true,
		},
		{
			name: "mismatched assignment task array lengths",
			modifier: func(a *Assignment) {
				a.Tasks[0].Values = []pgtype.Numeric{PGUint64(2), PGUint64(3), PGUint64(4)}
				a.Tasks[0].ValueTypes = []pgtype.Numeric{PGUint64(5), PGUint64(6)}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withClientMigrated(t, func(client *Client) {
				var assignment Assignment
				if err := copytest.DeepCopy(&assignment, &validAssignment); err != nil {
					t.Errorf("failed to create assignment struct copy: %v", err)
					return
				}

				tt.modifier(&assignment)

				err := assignment.Merge(context.Background(), client.queries, func(gen.Table, bool, int64) {})
				if (err != nil) != tt.wantErr {
					t.Errorf("Assignment.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
					return
				}
				fetchedResult, err := client.queries.GetAssignment(context.Background(), assignment.ID)
				if err != nil {
					t.Errorf("failed to fetch inserted assignment: %v", err)
					return
				}
				if fetchedResult != assignment.ID {
					t.Errorf("failed to validate INSERT: inserted data has ID %d, DB returned %d", assignment.ID, fetchedResult)
				}
			})
		})
	}
}
