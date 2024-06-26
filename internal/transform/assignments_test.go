//go:build !goverter

package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func mustAssignmentTitle(from api.Assignment2Title0) *api.Assignment2_Title {
	assignmentTitle := new(api.Assignment2_Title)
	if err := assignmentTitle.FromAssignment2Title0(from); err != nil {
		panic(err)
	}
	return assignmentTitle
}

func mustAssignmentBriefing(from api.Assignment2Briefing0) *api.Assignment2_Briefing {
	assignmentBriefing := new(api.Assignment2_Briefing)
	if err := assignmentBriefing.FromAssignment2Briefing0(from); err != nil {
		panic(err)
	}
	return assignmentBriefing
}

func mustAssignmentDescription(from api.Assignment2Description0) *api.Assignment2_Description {
	assignmentDescription := new(api.Assignment2_Description)
	if err := assignmentDescription.FromAssignment2Description0(from); err != nil {
		panic(err)
	}
	return assignmentDescription
}

func mustAssignment2Reward(from api.Reward2) *api.Assignment2_Reward {
	assignmentReward := new(api.Assignment2_Reward)
	if err := assignmentReward.FromReward2(from); err != nil {
		panic(err)
	}
	return assignmentReward
}

var validAssignment = api.Assignment2{
	Id:          ptr(int64(7)),
	Title:       mustAssignmentTitle("Foo"),
	Briefing:    mustAssignmentBriefing("Foo briefing"),
	Description: mustAssignmentDescription("Foo description"),
	Expiration:  ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
	Progress:    &[]uint64{1, 2, 3},
	Reward: mustAssignment2Reward(api.Reward2{
		Amount: ptr(uint64(100)),
		Type:   ptr(int32(3)),
	}),
	Tasks: &[]api.Task2{
		{
			Type:       ptr(int32(4)),
			ValueTypes: &[]uint64{2, 3, 4},
			Values:     &[]uint64{5, 6, 7},
		},
	},
}

func TestAssignments(t *testing.T) {
	// modifier changes the valid assignment to one that is suited for the test
	type modifier func(*api.Assignment2)
	tests := []struct {
		name     string
		modifier modifier
		want     []db.EntityMerger
		wantErr  bool
	}{
		{
			name: "valid",
			modifier: func(a *api.Assignment2) {
				// keep valid
			},
			want: []db.EntityMerger{
				&db.Assignment{
					Assignment: gen.Assignment{
						ID:           7,
						Title:        "Foo",
						Briefing:     "Foo briefing",
						Description:  "Foo description",
						Expiration:   db.PGTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
						RewardType:   3,
						RewardAmount: db.PGUint64(100),
					},
					Tasks: []gen.AssignmentTask{
						{
							TaskType:   4,
							ValueTypes: []pgtype.Numeric{db.PGUint64(2), db.PGUint64(3), db.PGUint64(4)},
							Values:     []pgtype.Numeric{db.PGUint64(5), db.PGUint64(6), db.PGUint64(7)},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing required title",
			modifier: func(a *api.Assignment2) {
				a.Title = nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing required reward amount",
			modifier: func(a *api.Assignment2) {
				a.Reward = mustAssignment2Reward(api.Reward2{
					Amount: nil,
					Type:   ptr(int32(3)),
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing required task values",
			modifier: func(a *api.Assignment2) {
				(*a.Tasks)[0].Values = nil
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var assignment api.Assignment2
			if err := copytest.DeepCopy(&assignment, &validAssignment); err != nil {
				t.Errorf("failed to create assignment struct copy: %v", err)
				return
			}
			// call modifier on valid assignment copy
			tt.modifier(&assignment)
			data := APIData{
				Assignments: &[]api.Assignment2{
					assignment,
				},
			}
			converter := &ConverterImpl{}
			got, err := Assignments(converter, data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assignments() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Assignments() = %v, want %v", got, tt.want)
			}
		})
	}
}
