package db

import (
	"context"
	"io"
	"log"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func withClient(t *testing.T, do func(client *Client, migration *migrate.Migrate)) {
	cfg := config.Get()
	dbName := strings.ReplaceAll(t.Name(), "/", "_")

	client, err := New(cfg, dbName, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatalf("could not initialize DB connection: %v", err)
	}
	defer func() {
		if err = client.db.Drop(context.Background()); err != nil {
			t.Logf("could not drop database: %v", err)
		}
		if err = client.Disconnect(); err != nil {
			t.Logf("could not disconnect: %v", err)
		}
	}()
	migration, err := client.newMigration("../../migrations")
	if err != nil {
		t.Fatalf("client.newMigration() error = %v, want nil", err)
	}
	do(client, migration)
}

func TestMigrateUp(t *testing.T) {
	withClient(t, func(client *Client, _ *migrate.Migrate) {
		if err := client.MigrateUp("../../migrations"); err != nil {
			t.Errorf("client.MigrateUp() error = %v, expected nil", err)
		}
	})
}

var collNames = []string{
	"planets",
	"campaigns",
	"dispatches",
	"events",
	"assignments",
	"wars",
	"snapshots",
}

func TestCollectionsExist(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		if err := migration.Up(); err != nil {
			t.Errorf("failed to migrate up: %v", err)
			return
		}

		fnPlanetCollections := func() []string {
			colls, errList := client.mongo.Database(t.Name()).ListCollectionNames(
				context.Background(),
				bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: collNames}}}},
			)
			if errList != nil {
				t.Errorf("could not list collections: %v", errList)
				return []string{}
			}
			return colls
		}
		if colls := fnPlanetCollections(); len(colls) != len(collNames) {
			t.Errorf("expected %d collections, got %d (%v)", len(collNames), len(colls), colls)
			return
		}
		if err := migration.Down(); err != nil {
			t.Errorf("failed to migrate down: %v", err)
			return
		}
		if colls := fnPlanetCollections(); len(colls) > 0 {
			t.Errorf("expected no collections, got %d (%v)", len(colls), colls)
			return
		}
	})
}

func TestIndexesExist(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		if err := migration.Up(); err != nil {
			t.Errorf("failed to migrate up: %v", err)
			return
		}

		for _, collection := range collNames {
			coll := client.mongo.Database(t.Name()).Collection(collection)
			indexes, err := coll.Indexes().List(context.Background())
			if err != nil {
				t.Errorf("failed to retrieve indexes: %v", err)
				return
			}
			var results []any
			if err = indexes.All(context.Background(), &results); err != nil {
				t.Errorf("failed to decode indexes response: %v", err)
				return
			}
			if len(results) == 0 {
				t.Error("expected len(indexes) > 0, got 0")
				return
			}
		}
	})
}

func TestPlanetsSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.Planet{
				ID:        1,
				Name:      "Foo",
				Sector:    "Bar",
				Position:  structs.PlanetPosition{X: 1, Y: 2},
				Waypoints: []int32{1, 2, 3},
				Disabled:  false,
				Biome: structs.Biome{
					Name:        "Forest",
					Description: "Lush forest",
				},
				Hazards: []structs.Hazard{
					{
						Name:        "Moist",
						Description: "Very very moist",
					},
				},
				MaxHealth:      1000,
				InitialOwner:   "Super Humans",
				RegenPerSecond: 50.0,
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Planet{
				ID:        1,
				Name:      "Foo",
				Sector:    "Bar",
				Position:  structs.PlanetPosition{X: 1, Y: 2},
				Waypoints: []int32{1, 2, 3},
				Disabled:  false,
				Biome: structs.Biome{
					Name:        "Forest",
					Description: "Lush forest",
				},
				Hazards: []structs.Hazard{
					{
						Name:        "Moist",
						Description: "Very very moist",
					},
				},
				MaxHealth:      1000,
				RegenPerSecond: 50.0,
			},
			wantErr: true,
		},
		{
			name: "valid struct missing slice",
			doc: structs.Planet{
				ID:        1,
				Name:      "Foo",
				Sector:    "Bar",
				Position:  structs.PlanetPosition{X: 1, Y: 2},
				Waypoints: nil,
				Disabled:  false,
				Biome: structs.Biome{
					Name:        "Forest",
					Description: "Lush forest",
				},
				Hazards: []structs.Hazard{
					{
						Name:        "Moist",
						Description: "Very very moist",
					},
				},
				MaxHealth:      1000,
				InitialOwner:   "Humans",
				RegenPerSecond: 50.0,
			},
			wantErr: true,
		},
		{
			name: "negative health",
			doc: structs.Planet{
				ID:        1,
				Name:      "Foo",
				Sector:    "Bar",
				Position:  structs.PlanetPosition{X: 1, Y: 2},
				Waypoints: []int32{1, 2, 3},
				Disabled:  false,
				Biome: structs.Biome{
					Name:        "Forest",
					Description: "Lush forest",
				},
				Hazards: []structs.Hazard{
					{
						Name:        "Moist",
						Description: "Very very moist",
					},
				},
				MaxHealth:      -1,
				InitialOwner:   "Super Humans",
				RegenPerSecond: 50.0,
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now()),
				Factions:  []string{"Humans", "Automatons"},
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("planets")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Planet
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
					return
				}
			})
		})
	}
}

func TestCampaignsSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.Campaign{
				ID:       1,
				PlanetID: 3,
				Type:     5,
				Count:    10,
			},
			wantErr: false,
		},
		{
			name: "negative count",
			doc: structs.Campaign{
				ID:       1,
				PlanetID: 3,
				Type:     5,
				Count:    -5,
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now()),
				Factions:  []string{"Humans", "Automatons"},
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("campaigns")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Campaign
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
					return
				}
			})
		})
	}
}

func TestDispatchesSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.Dispatch{
				ID:         1,
				CreateTime: primitive.NewDateTimeFromTime(time.Now()),
				Type:       3,
				Message:    "Foobar",
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Dispatch{
				ID:         1,
				CreateTime: primitive.NewDateTimeFromTime(time.Now()),
				Type:       3,
				Message:    "",
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now()),
				Factions:  []string{"Humans", "Automatons"},
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("dispatches")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Dispatch
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
					return
				}
			})
		})
	}
}

func TestEventsSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.Event{
				ID:        1,
				Type:      3,
				Faction:   "Foobar",
				MaxHealth: 100,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(10 * 24 * time.Hour)),
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Event{
				ID:        1,
				Type:      3,
				Faction:   "Foobar",
				MaxHealth: 100,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				// EndTime: toPrimitiveTs(time.Now().Add(10 * 24 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "endtime gt starttime",
			doc: structs.Event{
				ID:        1,
				Type:      3,
				Faction:   "Foobar",
				MaxHealth: 100,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(-1 * 10 * 24 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "negative health",
			doc: structs.Event{
				ID:        1,
				Type:      3,
				Faction:   "Foobar",
				MaxHealth: -1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(10 * 24 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now()),
				Factions:  []string{"Humans", "Automatons"},
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("events")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Event
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
					return
				}
			})
		})
	}
}

func TestAssignmentsSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.Assignment{
				ID:          1,
				Title:       "Foobar",
				Briefing:    "Briefing text",
				Description: "Description text, but a bit longer",
				Expiration:  primitive.NewDateTimeFromTime(time.Now().Add(5 * 24 * time.Hour)),
				Progress:    []int32{2, 3, 4},
				Tasks: []structs.AssignmentTask{
					{
						Type:       2,
						Values:     []int32{1, 2, 3},
						ValueTypes: []int32{5, 6, 7},
					},
				},
				Reward: structs.AssignmentReward{
					Type:   4,
					Amount: 8,
				},
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Assignment{
				ID:          1,
				Title:       "Foobar",
				Briefing:    "Briefing text",
				Description: "Description text, but a bit longer",
				Reward: structs.AssignmentReward{
					Type:   4,
					Amount: 8,
				},
			},
			wantErr: true,
		},
		{
			name: "negative reward amount",
			doc: structs.Assignment{
				ID:          1,
				Title:       "Foobar",
				Briefing:    "Briefing text",
				Description: "Description text, but a bit longer",
				Tasks: []structs.AssignmentTask{
					{
						Type:       2,
						Values:     []int32{1, 2, 3},
						ValueTypes: []int32{5, 6, 7},
					},
				},
				Reward: structs.AssignmentReward{
					Type:   4,
					Amount: -1,
				},
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now()),
				Factions:  []string{"Humans", "Automatons"},
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("assignments")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Assignment
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
					return
				}
			})
		})
	}
}

func TestWarsSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(5 * 24 * time.Hour)),
				Factions: []string{
					"Humans", "Automatons",
				},
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				Factions: []string{
					"Humans", "Automatons",
				},
			},
			wantErr: true,
		},
		{
			name: "endtime gt starttime",
			doc: structs.War{
				ID:        1,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(-1 * 5 * 24 * time.Hour)),
				Factions: []string{
					"Humans", "Automatons",
				},
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.Event{
				ID:        1,
				Type:      3,
				Faction:   "Foobar",
				MaxHealth: 100,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(5 * 24 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("wars")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.War
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
					return
				}
			})
		})
	}
}

func TestSnapshotsSchema(t *testing.T) {
	type document any
	tests := []struct {
		name    string
		doc     document
		wantErr bool
	}{
		{
			name: "valid struct complete",
			doc: structs.Snapshot{
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
				WarSnapshot: structs.WarSnapshot{
					WarID:            6,
					ImpactMultiplier: 50.0,
				},
				AssignmentIDs: []int64{2, 3, 4},
				CampaignIDs:   []int32{6, 7, 8},
				DispatchIDs:   []int32{10, 11, 12},
				Planets: []structs.PlanetSnapshot{
					{
						ID:           3,
						Health:       100,
						CurrentOwner: "Humans",
						Event: &structs.EventSnapshot{
							EventID: 5,
							Health:  700,
						},
						Attacking: []int32{2, 3, 4},
						Statistics: structs.PlanetStatistics{
							MissionsWon:  44323,
							MissionsLost: 53555,
							MissionTime:  445566,
							Kills: structs.StatisticsKills{
								Terminid:   432432443244,
								Automaton:  34333312212222,
								Illuminate: 2333333333,
							},
							BulletsFired: 888999399393222,
							BulletsHit:   49324924499449222,
							TimePlayed:   structs.BSONLong(365 * 24 * time.Hour),
							Deaths:       55223535,
							Revives:      44442,
							Friendlies:   2221111,
							PlayerCount:  12345678,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid struct high number",
			doc: structs.Snapshot{
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
				WarSnapshot: structs.WarSnapshot{
					WarID:            6,
					ImpactMultiplier: 50.0,
				},
				AssignmentIDs: []int64{2, 3, 4},
				CampaignIDs:   []int32{6, 7, 8},
				DispatchIDs:   []int32{10, 11, 12},
				Planets: []structs.PlanetSnapshot{
					{
						ID:           3,
						Health:       100,
						CurrentOwner: "Humans",
						Event: &structs.EventSnapshot{
							EventID: 5,
							Health:  700,
						},
						Attacking: []int32{},
						Statistics: structs.PlanetStatistics{
							MissionsWon:  math.MaxUint64,
							MissionsLost: math.MaxUint64,
							MissionTime:  math.MaxUint64,
							Kills: structs.StatisticsKills{
								Terminid:   math.MaxUint64,
								Automaton:  math.MaxUint64,
								Illuminate: math.MaxUint64,
							},
							BulletsFired: math.MaxUint64,
							BulletsHit:   math.MaxUint64,
							TimePlayed:   math.MaxUint64,
							Deaths:       math.MaxUint64,
							Revives:      math.MaxUint64,
							Friendlies:   math.MaxUint64,
							PlayerCount:  math.MaxUint64,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Snapshot{
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
				WarSnapshot: structs.WarSnapshot{
					WarID:            6,
					ImpactMultiplier: 50.0,
				},
				AssignmentIDs: []int64{2, 3, 4},
				CampaignIDs:   []int32{6, 7, 8},
				DispatchIDs:   []int32{10, 11, 12},
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.Event{
				ID:        1,
				Type:      3,
				Faction:   "Foobar",
				MaxHealth: 100,
				StartTime: primitive.NewDateTimeFromTime(time.Now()),
				EndTime:   primitive.NewDateTimeFromTime(time.Now().Add(5 * 24 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "invalid struct",
			doc: struct {
				Foo string
			}{
				Foo: "bar",
			},
			wantErr: true,
		},
		{
			name:    "nil struct",
			doc:     nil,
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

				coll := client.db.Collection("snapshots")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					// exit prematurely, since all following assertions depend on the transaction result
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Error("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Snapshot
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Errorf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Errorf("fetched result = %v, want %v", decoded, tt.doc)
				}
			})
		})
	}
}
