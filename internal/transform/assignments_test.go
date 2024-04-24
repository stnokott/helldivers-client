package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func mustAssignment2Reward(from api.Reward2) *api.Assignment2_Reward {
	assignmentReward := new(api.Assignment2_Reward)
	if err := assignmentReward.FromReward2(from); err != nil {
		panic(err)
	}
	return assignmentReward
}

var validAssignment = api.Assignment2{
	Id:          ptr(int64(7)),
	Title:       ptr("Foo"),
	Briefing:    ptr("Foo briefing"),
	Description: ptr("Foo description"),
	Expiration:  ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
	Progress:    &[]int32{1, 2, 3},
	Reward: mustAssignment2Reward(api.Reward2{
		Amount: ptr(int32(100)),
		Type:   ptr(int32(3)),
	}),
	Tasks: &[]api.Task2{
		{
			Type:       ptr(int32(4)),
			ValueTypes: &[]int32{2, 3, 4},
			Values:     &[]int32{5, 6, 7},
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
						RewardAmount: 100,
					},
					Tasks: []gen.AssignmentTask{
						{
							TaskType:   4,
							ValueTypes: []int32{2, 3, 4},
							Values:     []int32{5, 6, 7},
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
			got, err := Assignments(data)
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
