package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func mustWarStatistics(from api.Statistics) *api.War_Statistics {
	warStats := new(api.War_Statistics)
	if err := warStats.FromStatistics(from); err != nil {
		panic(err)
	}
	return warStats
}

var validWarIDSnapshot = api.WarId{
	Id: ptr(int32(999)),
}

var validWarSnapshot = api.War{
	ImpactMultiplier: ptr(float64(0.0004)),
	Statistics: mustWarStatistics(api.Statistics{
		MissionsWon:     ptr(uint64(564643344)),
		MissionsLost:    ptr(uint64(4324332)),
		MissionTime:     ptr(uint64(432432532552)),
		TerminidKills:   ptr(uint64(66878822)),
		AutomatonKills:  ptr(uint64(73737274)),
		IlluminateKills: ptr(uint64(112212441)),
		BulletsFired:    ptr(uint64(424444421112)),
		BulletsHit:      ptr(uint64(33444465767)),
		TimePlayed:      ptr(uint64(365 * 24 * 60 * time.Second)),
		Deaths:          ptr(uint64(885545432)),
		Revives:         ptr(uint64(8765333)),
		Friendlies:      ptr(uint64(444432232)),
		PlayerCount:     ptr(uint64(44899)),
	}),
}

var validPlanetSnapshot = api.Planet{
	Index:        ptr(int32(456)),
	Health:       ptr(int64(2441141122)),
	CurrentOwner: ptr("Humans"),
	Event: mustPlanetEvent(api.Event{
		Id:     ptr(int32(667)),
		Health: ptr(int64(889999)),
	}),
	Attacking:      &[]int32{4, 5, 6},
	RegenPerSecond: ptr(float64(0.003)),
	Statistics: mustPlanetStatistics(api.Statistics{
		MissionsWon:     ptr(uint64(654345432324)),
		MissionsLost:    ptr(uint64(43234242355)),
		MissionTime:     ptr(uint64(6558685)),
		TerminidKills:   ptr(uint64(527838425)),
		AutomatonKills:  ptr(uint64(2854382888)),
		IlluminateKills: ptr(uint64(32845248882)),
		BulletsFired:    ptr(uint64(823344447)),
		BulletsHit:      ptr(uint64(885454545465645466)),
		TimePlayed:      ptr(uint64(2 * 24 * 60 * time.Second)),
		Deaths:          ptr(uint64(7574557545454)),
		Revives:         ptr(uint64(2232342344223)),
		Friendlies:      ptr(uint64(99976547755)),
		PlayerCount:     ptr(uint64(4242442443)),
	}),
}

var validAssignmentSnapshot = api.Assignment2{
	Id: ptr(int64(7)),
}

var validCampaignSnapshot = api.Campaign2{
	Id: ptr(int32(987)),
}

var validDispatchSnapshot = api.Dispatch{
	Id: ptr(int32(678)),
}

func TestSnapshots(t *testing.T) {
	// modifier changes the valid assignment to one that is suited for the test
	type modifiers struct {
		WarID      func(*api.WarId)
		War        func(*api.War)
		Assignment func(*api.Assignment2)
		Campaign   func(*api.Campaign2)
		Dispatch   func(*api.Dispatch)
		Planet     func(*api.Planet)
	}
	tests := []struct {
		name      string
		modifiers modifiers
		want      []db.EntityMerger
		wantErr   bool
	}{
		{
			name: "valid",
			modifiers: modifiers{
				WarID:      func(wi *api.WarId) {},
				War:        func(w *api.War) {},
				Assignment: func(a *api.Assignment2) {},
				Campaign:   func(c *api.Campaign2) {},
				Dispatch:   func(d *api.Dispatch) {},
				Planet:     func(p *api.Planet) {},
			},
			want: []db.EntityMerger{
				&db.Snapshot{
					Snapshot: gen.Snapshot{
						AssignmentIds:     []int64{7},
						CampaignIds:       []int32{987},
						DispatchIds:       []int32{678},
						StatisticsID:      -1,
						WarSnapshotID:     -1,
						PlanetSnapshotIds: nil,
					},
					WarSnapshot: gen.WarSnapshot{
						WarID:            999,
						ImpactMultiplier: 0.0004,
					},
					PlanetSnapshots: []db.PlanetSnapshot{
						{
							PlanetSnapshot: gen.PlanetSnapshot{
								PlanetID:           456,
								Health:             2441141122,
								CurrentOwner:       "Humans",
								AttackingPlanetIds: []int32{4, 5, 6},
								RegenPerSecond:     0.003,
								StatisticsID:       -1,
							},
							Event: &gen.EventSnapshot{
								EventID: 667,
								Health:  889999,
							},
							Statistics: gen.SnapshotStatistic{
								MissionsWon:     db.PGUint64(654345432324),
								MissionsLost:    db.PGUint64(43234242355),
								MissionTime:     db.PGUint64(6558685),
								TerminidKills:   db.PGUint64(527838425),
								AutomatonKills:  db.PGUint64(2854382888),
								IlluminateKills: db.PGUint64(32845248882),
								BulletsFired:    db.PGUint64(823344447),
								BulletsHit:      db.PGUint64(885454545465645466),
								TimePlayed:      db.PGUint64(uint64(2 * 24 * 60 * time.Second)),
								Deaths:          db.PGUint64(7574557545454),
								Revives:         db.PGUint64(2232342344223),
								Friendlies:      db.PGUint64(99976547755),
								PlayerCount:     db.PGUint64(4242442443),
							},
						},
					},
					Statistics: gen.SnapshotStatistic{
						MissionsWon:     db.PGUint64(564643344),
						MissionsLost:    db.PGUint64(4324332),
						MissionTime:     db.PGUint64(432432532552),
						TerminidKills:   db.PGUint64(66878822),
						AutomatonKills:  db.PGUint64(73737274),
						IlluminateKills: db.PGUint64(112212441),
						BulletsFired:    db.PGUint64(424444421112),
						BulletsHit:      db.PGUint64(33444465767),
						TimePlayed:      db.PGUint64(uint64(365 * 24 * 60 * time.Second)),
						Deaths:          db.PGUint64(885545432),
						Revives:         db.PGUint64(8765333),
						Friendlies:      db.PGUint64(444432232),
						PlayerCount:     db.PGUint64(44899),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty war impact multiplier",
			modifiers: modifiers{
				WarID: func(wi *api.WarId) {},
				War: func(w *api.War) {
					w.ImpactMultiplier = nil
				},
				Assignment: func(a *api.Assignment2) {},
				Campaign:   func(c *api.Campaign2) {},
				Dispatch:   func(d *api.Dispatch) {},
				Planet:     func(p *api.Planet) {},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty war statistic",
			modifiers: modifiers{
				WarID: func(wi *api.WarId) {},
				War: func(w *api.War) {
					w.Statistics = mustWarStatistics(api.Statistics{})
				},
				Assignment: func(a *api.Assignment2) {},
				Campaign:   func(c *api.Campaign2) {},
				Dispatch:   func(d *api.Dispatch) {},
				Planet:     func(p *api.Planet) {},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty regen",
			modifiers: modifiers{
				WarID:      func(wi *api.WarId) {},
				War:        func(w *api.War) {},
				Assignment: func(a *api.Assignment2) {},
				Campaign:   func(c *api.Campaign2) {},
				Dispatch:   func(d *api.Dispatch) {},
				Planet: func(p *api.Planet) {
					p.RegenPerSecond = nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warID := validWarIDSnapshot
			war := validWarSnapshot
			planet := validPlanetSnapshot
			assignment := validAssignmentSnapshot
			campaign := validCampaignSnapshot
			dispatch := validDispatchSnapshot

			// call modifiers on valid assignment copies
			tt.modifiers.WarID(&warID)
			tt.modifiers.War(&war)
			tt.modifiers.Planet(&planet)
			tt.modifiers.Assignment(&assignment)
			tt.modifiers.Campaign(&campaign)
			tt.modifiers.Dispatch(&dispatch)
			data := APIData{
				WarID:       &warID,
				War:         &war,
				Planets:     &[]api.Planet{planet},
				Assignments: &[]api.Assignment2{assignment},
				Campaigns:   &[]api.Campaign2{campaign},
				Dispatches:  &[]api.Dispatch{dispatch},
			}
			got, err := Snapshot(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Snapshot() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Snapshot() = %v, want %v", got, tt.want)
			}
		})
	}
}
