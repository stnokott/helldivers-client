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
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func withClient(t *testing.T, do func(client *Client, migration *migrate.Migrate)) {
	mongoURI := getMongoURI()
	dbName := strings.ReplaceAll(t.Name(), "/", "_")

	client, err := New(mongoURI, dbName, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatalf("could not initialize DB connection: %v", err)
	}
	defer func() {
		if err = client.mongo.Database(dbName).Drop(context.Background()); err != nil {
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

var collections = []string{
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
			t.Fatalf("failed to migrate up: %v", err)
		}

		fnPlanetCollections := func() []string {
			colls, errList := client.mongo.Database(t.Name()).ListCollectionNames(
				context.Background(),
				bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: collections}}}},
			)
			if errList != nil {
				t.Errorf("could not list collections: %v", errList)
				return []string{}
			}
			return colls
		}
		if colls := fnPlanetCollections(); len(colls) != len(collections) {
			t.Fatalf("expected %d collections, got %d (%v)", len(collections), len(colls), colls)
		}
		if err := migration.Down(); err != nil {
			t.Fatalf("failed to migrate down: %v", err)
		}
		if colls := fnPlanetCollections(); len(colls) > 0 {
			t.Fatalf("expected no collections, got %d (%v)", len(colls), colls)
		}
	})
}

func TestIndexesExist(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		if err := migration.Up(); err != nil {
			t.Fatalf("failed to migrate up: %v", err)
		}

		for _, collection := range collections {
			coll := client.mongo.Database(t.Name()).Collection(collection)
			indexes, err := coll.Indexes().List(context.Background())
			if err != nil {
				t.Fatalf("failed to retrieve indexes: %v", err)
			}
			var results []any
			if err = indexes.All(context.Background(), &results); err != nil {
				t.Fatalf("failed to decode indexes response: %v", err)
			}
			if len(results) == 0 {
				t.Error("expected len(indexes) > 0, got 0")
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
				ID:             1,
				Name:           "Foo",
				Sector:         "Bar",
				Position:       structs.PlanetPosition{X: 1, Y: 2},
				Waypoints:      []int{1, 2, 3},
				Disabled:       false,
				MaxHealth:      1000,
				InitialOwner:   "Super Humans",
				RegenPerSecond: 50.0,
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Planet{
				ID:             1,
				Name:           "Foo",
				Sector:         "Bar",
				Position:       structs.PlanetPosition{X: 1, Y: 2},
				Waypoints:      []int{1, 2, 3},
				Disabled:       false,
				MaxHealth:      1000,
				RegenPerSecond: 50.0,
			},
			wantErr: true,
		},
		{
			name: "valid struct missing slice",
			doc: structs.Planet{
				ID:             1,
				Name:           "Foo",
				Sector:         "Bar",
				Position:       structs.PlanetPosition{X: 1, Y: 2},
				Waypoints:      nil,
				Disabled:       false,
				MaxHealth:      1000,
				InitialOwner:   "Humans",
				RegenPerSecond: 50.0,
			},
			wantErr: true,
		},
		{
			name: "negative health",
			doc: structs.Planet{
				ID:             1,
				Name:           "Foo",
				Sector:         "Bar",
				Position:       structs.PlanetPosition{X: 1, Y: 2},
				Waypoints:      []int{1, 2, 3},
				Disabled:       false,
				MaxHealth:      -1,
				InitialOwner:   "Super Humans",
				RegenPerSecond: 50.0,
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now()),
				ImpactMultiplier: 2.0,
				Factions:         []string{"Humans", "Automatons"},
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("planets")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Planet
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
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
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now()),
				ImpactMultiplier: 2.0,
				Factions:         []string{"Humans", "Automatons"},
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("campaigns")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Campaign
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
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
				CreateTime: toPrimitiveTs(time.Now()),
				Type:       3,
				Message:    "Foobar",
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Dispatch{
				ID:         1,
				CreateTime: primitive.Timestamp{},
				Type:       3,
				Message:    "Foobar",
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now()),
				ImpactMultiplier: 2.0,
				Factions:         []string{"Humans", "Automatons"},
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("dispatches")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Dispatch
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
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
				StartTime: toPrimitiveTs(time.Now()),
				EndTime:   toPrimitiveTs(time.Now().Add(10 * 24 * time.Hour)),
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
				StartTime: toPrimitiveTs(time.Now()),
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
				StartTime: toPrimitiveTs(time.Now()),
				EndTime:   toPrimitiveTs(time.Now().Add(-1 * 10 * 24 * time.Hour)),
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
				StartTime: toPrimitiveTs(time.Now()),
				EndTime:   toPrimitiveTs(time.Now().Add(10 * 24 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "wrong struct",
			doc: structs.War{
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now()),
				ImpactMultiplier: 2.0,
				Factions:         []string{"Humans", "Automatons"},
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("events")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Event
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
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
				Tasks: []structs.AssignmentTask{
					{
						Type:       2,
						Values:     []int{1, 2, 3},
						ValueTypes: []int{5, 6, 7},
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
						Values:     []int{1, 2, 3},
						ValueTypes: []int{5, 6, 7},
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
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now()),
				ImpactMultiplier: 2.0,
				Factions:         []string{"Humans", "Automatons"},
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("assignments")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Assignment
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
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
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now().Add(5 * 24 * time.Hour)),
				ImpactMultiplier: 50.0,
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
				StartTime: toPrimitiveTs(time.Now()),
				// Ended:            toPrimitiveTs(time.Now().Add(5 * 24 * time.Hour)),
				ImpactMultiplier: 50.0,
				Factions: []string{
					"Humans", "Automatons",
				},
			},
			wantErr: true,
		},
		{
			name: "endtime gt starttime",
			doc: structs.War{
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now().Add(-1 * 5 * 24 * time.Hour)),
				ImpactMultiplier: 50.0,
				Factions: []string{
					"Humans", "Automatons",
				},
			},
			wantErr: true,
		},
		{
			name: "negative impact multiplier",
			doc: structs.War{
				ID:               1,
				StartTime:        toPrimitiveTs(time.Now()),
				EndTime:          toPrimitiveTs(time.Now().Add(5 * 24 * time.Hour)),
				ImpactMultiplier: -0.5,
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
				StartTime: toPrimitiveTs(time.Now()),
				EndTime:   toPrimitiveTs(time.Now().Add(5 * 24 * time.Hour)),
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("wars")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.War
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
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
				ID:            toPrimitiveTs(time.Now()),
				WarID:         6,
				AssignmentIDs: []int{2, 3, 4},
				CampaignIDs:   []int{6, 7, 8},
				DispatchIDs:   []int{10, 11, 12},
				Planets: []structs.PlanetSnapshot{
					{
						ID:           3,
						Health:       100,
						CurrentOwner: "Humans",
						Event: &structs.EventSnapshot{
							EventID: 5,
							Health:  700,
						},
						Statistics: &structs.PlanetStatistics{
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
							TimePlayed:   int64(365 * 24 * time.Hour),
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
				ID:            toPrimitiveTs(time.Now()),
				WarID:         6,
				AssignmentIDs: []int{2, 3, 4},
				CampaignIDs:   []int{6, 7, 8},
				DispatchIDs:   []int{10, 11, 12},
				Planets: []structs.PlanetSnapshot{
					{
						ID:           3,
						Health:       100,
						CurrentOwner: "Humans",
						Event: &structs.EventSnapshot{
							EventID: 5,
							Health:  700,
						},
						Statistics: &structs.PlanetStatistics{
							MissionsWon:  math.MaxInt64,
							MissionsLost: math.MaxInt64,
							MissionTime:  math.MaxInt64,
							Kills: structs.StatisticsKills{
								Terminid:   math.MaxInt64,
								Automaton:  math.MaxInt64,
								Illuminate: math.MaxInt64,
							},
							BulletsFired: math.MaxInt64,
							BulletsHit:   math.MaxInt64,
							TimePlayed:   math.MaxInt64,
							Deaths:       math.MaxInt64,
							Revives:      math.MaxInt64,
							Friendlies:   math.MaxInt64,
							PlayerCount:  math.MaxInt64,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid struct incomplete",
			doc: structs.Snapshot{
				ID:            toPrimitiveTs(time.Now()),
				WarID:         6,
				AssignmentIDs: []int{2, 3, 4},
				CampaignIDs:   []int{6, 7, 8},
				DispatchIDs:   []int{10, 11, 12},
			},
			wantErr: true,
		},
		{
			name: "negative statistics",
			doc: structs.Snapshot{
				ID:            toPrimitiveTs(time.Now()),
				WarID:         6,
				AssignmentIDs: []int{2, 3, 4},
				CampaignIDs:   []int{6, 7, 8},
				DispatchIDs:   []int{10, 11, 12},
				Planets: []structs.PlanetSnapshot{
					{
						ID:           3,
						Health:       100,
						CurrentOwner: "Humans",
						Event: &structs.EventSnapshot{
							EventID: 5,
							Health:  700,
						},
						Statistics: &structs.PlanetStatistics{
							MissionsWon:  44323,
							MissionsLost: 53555,
							MissionTime:  445566,
							Kills: structs.StatisticsKills{
								Terminid:   -6,
								Automaton:  34333312212222,
								Illuminate: 2333333333,
							},
							BulletsFired: 888999399393222,
							BulletsHit:   49324924499449222,
							TimePlayed:   int64(365 * 24 * time.Hour),
							Deaths:       55223535,
							Revives:      44442,
							Friendlies:   2221111,
							PlayerCount:  12345678,
						},
					},
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
				StartTime: toPrimitiveTs(time.Now()),
				EndTime:   toPrimitiveTs(time.Now().Add(5 * 24 * time.Hour)),
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
					t.Fatalf("failed to migrate up: %v", err)
				}

				coll := client.database().Collection("snapshots")
				insertResult, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Fatalf("InsertOne() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				fetchedResult := coll.FindOne(context.Background(), bson.D{{
					Key: "_id", Value: insertResult.InsertedID,
				}})
				if fetchedResult == nil {
					t.Fatal("fetched result is nil, expected non-nil")
					return
				}
				var decoded structs.Snapshot
				if err = fetchedResult.Decode(&decoded); err != nil {
					t.Fatalf("failed to decode result: %v", err)
					return
				}
				if !reflect.DeepEqual(tt.doc, decoded) {
					t.Fatalf("fetched result = %v, want %v", decoded, tt.doc)
				}
			})
		})
	}
}
