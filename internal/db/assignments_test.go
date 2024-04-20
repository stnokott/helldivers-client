package db

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

var validAssignment = Assignment{
	Assignment: gen.Assignment{
		ID:           3,
		Title:        "Footitle",
		Briefing:     "Foobriefing",
		Description:  "Bardescription",
		Expiration:   PGTimestamp(time.Date(2024, 1, 2, 3, 4, 5, 6, time.Local)),
		Progress:     []int32{1, 2, 3},
		TaskIds:      []int64{5, 6, 7},
		RewardType:   8,
		RewardAmount: 100,
	},
	Tasks: []AssignmentTask{
		{
			Type:       9,
			Values:     []int32{7, 8, 9},
			ValueTypes: []int32{42, 44, 46},
		},
	},
}

func TestAssignmentsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Assignment)
	tests := []struct {
		name          string
		modifier      modifier
		wantTaskErr   bool
		wantAssignErr bool
	}{
		{
			name:          "valid",
			modifier:      func(*Assignment) {},
			wantTaskErr:   false,
			wantAssignErr: false,
		},
		{
			name: "empty required title",
			modifier: func(a *Assignment) {
				a.Title = ""
			},
			wantTaskErr:   false,
			wantAssignErr: true,
		},
		{
			name: "empty required expiration time",
			modifier: func(a *Assignment) {
				a.Expiration = pgtype.Timestamp{}
			},
			wantTaskErr:   false,
			wantAssignErr: true,
		},
		{
			name: "mismatched assignment task array lengths",
			modifier: func(a *Assignment) {
				a.Tasks[0].Values = []int32{2, 3, 4}
				a.Tasks[0].ValueTypes = []int32{5, 6}
			},
			wantTaskErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withClient(t, func(client *Client, migration *migrate.Migrate) {
				if err := migration.Up(); err != nil {
					t.Errorf("failed to migrate up: %v", err)
					return
				}
				assignment := validAssignment
				tt.modifier(&assignment)

				for _, task := range assignment.Tasks {
					_, err := client.queries.InsertAssignmentTask(context.Background(), gen.InsertAssignmentTaskParams{
						Type:       task.Type,
						Values:     task.Values,
						ValueTypes: task.ValueTypes,
					})
					if (err != nil) != tt.wantTaskErr {
						t.Errorf("InsertAssignmentTask() error = %v, wantTaskErr = %v", err, tt.wantTaskErr)
						return
					}
					if err != nil {
						return
					}
				}

				_, err := client.queries.InsertAssignment(context.Background(), gen.InsertAssignmentParams(assignment.Assignment))
				if (err != nil) != tt.wantAssignErr {
					t.Errorf("InsertAssignment() error = %v, wantErr = %v", err, tt.wantAssignErr)
					return
				}
				if err != nil {
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
