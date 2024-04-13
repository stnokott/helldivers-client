// Package transform converts API structs to DB structs
package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

func mustStatistics(from api.Statistics) *api.Planet_Statistics {
	planetStats := new(api.Planet_Statistics)
	if err := planetStats.FromStatistics(from); err != nil {
		panic(err)
	}
	return planetStats
}

func TestSnapshotsTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		s       Snapshots
		args    args
		want    *db.DocsProvider[structs.Snapshot]
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					WarID: &api.WarId{Id: ptr(int32(6))},
					War: &api.War{
						Now: ptr(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
					},
					Assignments: &[]api.Assignment2{
						{Id: ptr(int64(2))},
						{Id: ptr(int64(3))},
						{Id: ptr(int64(4))},
					},
					Campaigns: &[]api.Campaign2{
						{Id: ptr(int32(7))},
						{Id: ptr(int32(8))},
						{Id: ptr(int32(9))},
					},
					Dispatches: &[]api.Dispatch{
						{Id: ptr(int32(10))},
					},
					Planets: &[]api.Planet{
						{
							Index:        ptr(int32(6)),
							Health:       ptr(int64(1234567)),
							CurrentOwner: ptr("Humans"),
							Event: mustPlanetEvent(api.Event{
								Id:     ptr(int32(6)),
								Health: ptr(int64(999)),
							}),
							Statistics: mustStatistics(api.Statistics{
								MissionsWon:     ptr(uint64(100)),
								MissionsLost:    ptr(uint64(55)),
								MissionTime:     ptr(uint64(12345)),
								TerminidKills:   ptr(uint64(10000)),
								AutomatonKills:  ptr(uint64(99999)),
								IlluminateKills: ptr(uint64(333333)),
								BulletsFired:    ptr(uint64(11111)),
								BulletsHit:      ptr(uint64(1111)),
								TimePlayed:      ptr(uint64(123456)),
								Deaths:          ptr(uint64(32134)),
								Revives:         ptr(uint64(94284)),
								Friendlies:      ptr(uint64(12940)),
								PlayerCount:     ptr(uint64(444442)),
							}),
							Attacking: &[]int32{8, 9, 10},
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Snapshot]{
				CollectionName: "snapshots",
				Docs: []db.DocWrapper[structs.Snapshot]{
					{
						DocID: db.PrimitiveTime(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
						Document: structs.Snapshot{
							Timestamp:     db.PrimitiveTime(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
							WarID:         6,
							AssignmentIDs: []int64{2, 3, 4},
							CampaignIDs:   []int32{7, 8, 9},
							DispatchIDs:   []int32{10},
							Planets: []structs.PlanetSnapshot{
								{
									ID:           6,
									Health:       1234567,
									CurrentOwner: "Humans",
									Event: &structs.EventSnapshot{
										EventID: 6,
										Health:  999,
									},
									Statistics: structs.PlanetStatistics{
										MissionsWon:  100,
										MissionsLost: 55,
										MissionTime:  12345,
										Kills: structs.StatisticsKills{
											Terminid:   10000,
											Automaton:  99999,
											Illuminate: 333333,
										},
										BulletsFired: 11111,
										BulletsHit:   1111,
										TimePlayed:   123456,
										Deaths:       32134,
										Revives:      94284,
										Friendlies:   12940,
										PlayerCount:  444442,
									},
									Attacking: []int32{8, 9, 10},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "optional event nil",
			args: args{
				data: APIData{
					WarID: &api.WarId{Id: ptr(int32(6))},
					War: &api.War{
						Now: ptr(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
					},
					Assignments: &[]api.Assignment2{
						{Id: ptr(int64(2))},
						{Id: ptr(int64(3))},
						{Id: ptr(int64(4))},
					},
					Campaigns: &[]api.Campaign2{
						{Id: ptr(int32(7))},
						{Id: ptr(int32(8))},
						{Id: ptr(int32(9))},
					},
					Dispatches: &[]api.Dispatch{
						{Id: ptr(int32(10))},
					},
					Planets: &[]api.Planet{
						{
							Index:        ptr(int32(6)),
							Health:       ptr(int64(1234567)),
							CurrentOwner: ptr("Humans"),
							Event:        nil,
							Statistics: mustStatistics(api.Statistics{
								MissionsWon:     ptr(uint64(100)),
								MissionsLost:    ptr(uint64(55)),
								MissionTime:     ptr(uint64(12345)),
								TerminidKills:   ptr(uint64(10000)),
								AutomatonKills:  ptr(uint64(99999)),
								IlluminateKills: ptr(uint64(333333)),
								BulletsFired:    ptr(uint64(11111)),
								BulletsHit:      ptr(uint64(1111)),
								TimePlayed:      ptr(uint64(123456)),
								Deaths:          ptr(uint64(32134)),
								Revives:         ptr(uint64(94284)),
								Friendlies:      ptr(uint64(12940)),
								PlayerCount:     ptr(uint64(444442)),
							}),
							Attacking: &[]int32{8, 9, 10},
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Snapshot]{
				CollectionName: "snapshots",
				Docs: []db.DocWrapper[structs.Snapshot]{
					{
						DocID: db.PrimitiveTime(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
						Document: structs.Snapshot{
							Timestamp:     db.PrimitiveTime(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
							WarID:         6,
							AssignmentIDs: []int64{2, 3, 4},
							CampaignIDs:   []int32{7, 8, 9},
							DispatchIDs:   []int32{10},
							Planets: []structs.PlanetSnapshot{
								{
									ID:           6,
									Health:       1234567,
									CurrentOwner: "Humans",
									Event:        nil,
									Statistics: structs.PlanetStatistics{
										MissionsWon:  100,
										MissionsLost: 55,
										MissionTime:  12345,
										Kills: structs.StatisticsKills{
											Terminid:   10000,
											Automaton:  99999,
											Illuminate: 333333,
										},
										BulletsFired: 11111,
										BulletsHit:   1111,
										TimePlayed:   123456,
										Deaths:       32134,
										Revives:      94284,
										Friendlies:   12940,
										PlayerCount:  444442,
									},
									Attacking: []int32{8, 9, 10},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil required campaign ID",
			args: args{
				data: APIData{
					WarID: &api.WarId{Id: ptr(int32(6))},
					War: &api.War{
						Now: ptr(time.Date(2024, 2, 3, 4, 5, 6, 7, time.Local)),
					},
					Assignments: &[]api.Assignment2{
						{Id: ptr(int64(2))},
						{Id: ptr(int64(3))},
						{Id: ptr(int64(4))},
					},
					Campaigns: &[]api.Campaign2{
						{Id: nil},
						{Id: ptr(int32(8))},
						{Id: ptr(int32(9))},
					},
					Dispatches: &[]api.Dispatch{
						{Id: ptr(int32(10))},
					},
					Planets: &[]api.Planet{
						{
							Index:        ptr(int32(6)),
							Health:       ptr(int64(1234567)),
							CurrentOwner: ptr("Humans"),
							Event: mustPlanetEvent(api.Event{
								Id:     ptr(int32(6)),
								Health: ptr(int64(999)),
							}),
							Statistics: mustStatistics(api.Statistics{
								MissionsWon:     ptr(uint64(100)),
								MissionsLost:    ptr(uint64(55)),
								MissionTime:     ptr(uint64(12345)),
								TerminidKills:   ptr(uint64(10000)),
								AutomatonKills:  ptr(uint64(99999)),
								IlluminateKills: ptr(uint64(333333)),
								BulletsFired:    ptr(uint64(11111)),
								BulletsHit:      ptr(uint64(1111)),
								TimePlayed:      ptr(uint64(123456)),
								Deaths:          ptr(uint64(32134)),
								Revives:         ptr(uint64(94284)),
								Friendlies:      ptr(uint64(12940)),
								PlayerCount:     ptr(uint64(444442)),
							}),
							Attacking: &[]int32{8, 9, 10},
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Snapshot]{
				CollectionName: "snapshots",
				Docs:           []db.DocWrapper[structs.Snapshot]{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := false
			errFunc := func(err error) {
				if !tt.wantErr {
					t.Logf("War.Transform() error: %v", err)
				}
				gotErr = true
			}
			got := tt.s.Transform(tt.args.data, errFunc)
			if gotErr != tt.wantErr {
				t.Errorf("Snapshots.Transform() returned error, wantErr %v", tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Snapshots.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
