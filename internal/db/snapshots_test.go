package db

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// TODO: fuzzing tests

var validWarSnapshot = War{
	ID:        999,
	StartTime: PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local)),
	EndTime:   PGTimestamp(time.Date(2025, 1, 1, 1, 1, 1, 1, time.Local)),
	Factions:  []string{"Humans", "Automatons"},
}

var validAssignmentSnapshot = Assignment{
	Assignment: gen.Assignment{
		ID:           3,
		Title:        "Footitle",
		Briefing:     "Foobriefing",
		Description:  "Bardescription",
		Expiration:   PGTimestamp(time.Date(2024, 1, 2, 3, 4, 5, 6, time.Local)),
		RewardType:   8,
		RewardAmount: 100,
	},
	Tasks: []gen.AssignmentTask{
		{
			TaskType:   9,
			Values:     []int32{7, 8, 9},
			ValueTypes: []int32{42, 44, 46},
		},
	},
}

var validEventSnapshot = Event{
	ID:        555,
	Type:      7,
	Faction:   "Automatons",
	MaxHealth: 55667788,
	StartTime: PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local)),
	EndTime:   PGTimestamp(time.Date(2025, 1, 1, 1, 1, 1, 1, time.Local)),
}

var validPlanetSnapshot = Planet{
	Planet: gen.Planet{
		ID:           1,
		Name:         "Foo",
		Sector:       "Bar",
		Position:     []float64{1, 2},
		WaypointIds:  []int32{1, 2, 3},
		Disabled:     false,
		BiomeName:    "FooBiome",
		HazardNames:  []string{"BarHazard"},
		MaxHealth:    1000,
		InitialOwner: "Super Humans",
	},
	Biome: gen.Biome{
		Name:        "FooBiome",
		Description: "This biome contains a lot of spaghetti",
	},
	Hazards: []gen.Hazard{
		{
			Name:        "BarHazard",
			Description: "This hazard contains a lot of bugs",
		},
	},
}

var validCampaignSnapshot = Campaign{
	ID:       5,
	PlanetID: 1,
	Type:     8,
	Count:    100,
}

var validDispatchSnapshot = Dispatch{
	ID:         123,
	CreateTime: PGTimestamp(time.Date(2024, 1, 2, 3, 4, 5, 6, time.Local)),
	Type:       5,
	Message:    "A valid dispatch",
}

var validSnapshot = Snapshot{
	Snapshot: gen.Snapshot{
		AssignmentIds:     []int64{3},
		WarSnapshotID:     -1, // will be filled from Merge
		CampaignIds:       []int32{5},
		DispatchIds:       []int32{123},
		StatisticsID:      -1,  // will be filled from Merge
		PlanetSnapshotIds: nil, // will be filled from Merge
	},
	WarSnapshot: gen.WarSnapshot{
		WarID:            999,
		ImpactMultiplier: 0.005,
	},
	PlanetSnapshots: []PlanetSnapshot{
		{
			PlanetSnapshot: gen.PlanetSnapshot{
				PlanetID:           456,
				Health:             556677,
				CurrentOwner:       "Automatons",
				RegenPerSecond:     0.06,
				AttackingPlanetIds: []int32{1},
			},
			Event: &gen.EventSnapshot{
				EventID: 789,
				Health:  999999,
			},
			Statistics: gen.SnapshotStatistic{
				MissionsWon:     PGUint64(6444323),
				MissionsLost:    PGUint64(53555),
				MissionTime:     PGUint64(445566),
				TerminidKills:   PGUint64(6666677),
				AutomatonKills:  PGUint64(7565465454),
				IlluminateKills: PGUint64(5345433455),
				BulletsFired:    PGUint64(888999399393222),
				BulletsHit:      PGUint64(49324924499449222),
				TimePlayed:      PGUint64(uint64(1 * 24 * 60 * time.Second)),
				Deaths:          PGUint64(55223535),
				Revives:         PGUint64(44442),
				Friendlies:      PGUint64(2221111),
				PlayerCount:     PGUint64(12345678),
			},
		},
	},
	Statistics: gen.SnapshotStatistic{
		MissionsWon:     PGUint64(564643344),
		MissionsLost:    PGUint64(4324332),
		MissionTime:     PGUint64(432432532552),
		TerminidKills:   PGUint64(66878822),
		AutomatonKills:  PGUint64(73737274),
		IlluminateKills: PGUint64(112212441),
		BulletsFired:    PGUint64(424444421112),
		BulletsHit:      PGUint64(33444465767),
		TimePlayed:      PGUint64(uint64(365 * 24 * 60 * time.Second)),
		Deaths:          PGUint64(885545432),
		Revives:         PGUint64(8765333),
		Friendlies:      PGUint64(444432232),
		PlayerCount:     PGUint64(44899),
	},
}

func TestSnapshotsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Snapshot)
	tests := []struct {
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(p *Snapshot) {},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withClient(t, func(client *Client, migration *migrate.Migrate) {
				if err := migration.Up(); err != nil {
					t.Errorf("failed to migrate up: %v", err)
					return
				}
				war := validWarSnapshot
				assignment := validAssignmentSnapshot
				event := validEventSnapshot
				planet := validPlanetSnapshot
				campaign := validCampaignSnapshot
				dispatch := validDispatchSnapshot
				snapshot := validSnapshot

				// equalize IDs for FK constraints
				snapshot.WarSnapshot.WarID = war.ID
				snapshot.PlanetSnapshots[0].Event.EventID = event.ID
				snapshot.PlanetSnapshots[0].PlanetID = planet.ID

				tt.modifier(&snapshot)

				if err := event.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to insert event (required for snapshot): %v", err)
					return
				}
				if err := assignment.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to insert assignment (required for snapshot): %v", err)
					return
				}
				if err := planet.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to insert planet (required for snapshot): %v", err)
					return
				}
				if err := campaign.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to insert campaign (required for snapshot): %v", err)
					return
				}
				if err := dispatch.Merge(context.Background(), client.queries, tableMergeStats{}); err != nil {
					t.Errorf("failed to insert dispatch (required for snapshot): %v", err)
					return
				}

				err := snapshot.Merge(context.Background(), client.queries, tableMergeStats{})
				if (err != nil) != tt.wantErr {
					t.Errorf("Snapshot.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
					return
				}

				_, err = client.queries.GetLatestSnapshot(context.Background())
				if err != nil {
					t.Errorf("failed to fetch inserted snapshot: %v", err)
					return
				}
			})
		})
	}
}
