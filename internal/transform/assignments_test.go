package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mustAssignment2Reward(from api.Reward2) *api.Assignment2_Reward {
	assignmentReward := new(api.Assignment2_Reward)
	if err := assignmentReward.FromReward2(from); err != nil {
		panic(err)
	}
	return assignmentReward
}

func TestAssignmentsTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		a       Assignments
		args    args
		want    *db.DocsProvider[structs.Assignment]
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					Assignments: &[]api.Assignment2{
						{
							Id:          ptr(int64(7)),
							Title:       ptr("Foo"),
							Briefing:    ptr("Foo briefing"),
							Description: ptr("Foo description"),
							Expiration:  ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)),
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
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Assignment]{
				CollectionName: "assignments",
				Docs: []db.DocWrapper[structs.Assignment]{
					{
						DocID: int64(7),
						Document: structs.Assignment{
							ID:          7,
							Title:       "Foo",
							Briefing:    "Foo briefing",
							Description: "Foo description",
							Expiration:  primitive.NewDateTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)),
							Progress:    []int32{1, 2, 3},
							Reward: structs.AssignmentReward{
								Amount: 100,
								Type:   3,
							},
							Tasks: []structs.AssignmentTask{
								{
									Type:       4,
									ValueTypes: []int32{2, 3, 4},
									Values:     []int32{5, 6, 7},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil briefing",
			args: args{
				data: APIData{
					Assignments: &[]api.Assignment2{
						{
							Id:          ptr(int64(7)),
							Title:       ptr("Foo"),
							Briefing:    nil,
							Description: ptr("Foo description"),
							Expiration:  ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)),
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
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Assignment]{
				CollectionName: db.CollAssignments,
				Docs:           []db.DocWrapper[structs.Assignment]{},
			},
			wantErr: true,
		},
		{
			name: "nil embedded",
			args: args{
				data: APIData{
					Assignments: &[]api.Assignment2{
						{
							Id:          ptr(int64(7)),
							Title:       ptr("Foo"),
							Briefing:    nil,
							Description: ptr("Foo description"),
							Expiration:  ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)),
							Progress:    &[]int32{1, 2, 3},
							Reward: mustAssignment2Reward(api.Reward2{
								Amount: ptr(int32(100)),
								Type:   ptr(int32(3)),
							}),
							Tasks: nil,
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Assignment]{
				CollectionName: db.CollAssignments,
				Docs:           []db.DocWrapper[structs.Assignment]{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := false
			errFunc := func(err error) {
				if !tt.wantErr {
					t.Logf("Assignments.Transform() error: %v", err)
				}
				gotErr = true
			}
			got := tt.a.Transform(tt.args.data, errFunc)
			if gotErr != tt.wantErr {
				t.Errorf("Assignments.Transform() returned errors, wantErr %v", tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Assignments.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
