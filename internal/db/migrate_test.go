package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson"
)

func withClient(t *testing.T, do func(client *Client, migration *migrate.Migrate)) {
	mongoURI := getMongoURI()
	client, err := New(mongoURI, t.Name(), log.Default())
	if err != nil {
		t.Fatalf("could not initialize DB connection: %v", err)
	}
	defer func() {
		if err = client.mongo.Database(client.dbName).Drop(context.Background()); err != nil {
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

func TestClientMigration(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		if err := migration.Up(); err != nil {
			t.Fatalf("failed to migrate up: %v", err)
		}
		fnPlanetCollectionExists := func() bool {
			colls, errList := client.database().ListCollectionNames(context.Background(), bson.D{{Key: "name", Value: "war_seasons"}})
			if errList != nil {
				t.Errorf("could not list collections: %v", errList)
				return false
			}
			return len(colls) == 1
		}
		if !fnPlanetCollectionExists() {
			t.Error("expected collection with name 'war_seasons', none found")
		}
		if err := migration.Down(); err != nil {
			t.Fatalf("failed to migrate down: %v", err)
		}
		if fnPlanetCollectionExists() {
			t.Error("expected collection with name 'war_seasons' to not exist")
		}
	})
}

func TestWarSeasonsSchema(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		type document any
		tests := []struct {
			name    string
			doc     document
			wantErr bool
		}{
			{
				name: "valid struct complete",
				doc: structs.WarSeason{
					ID:                     1,
					Capitals:               []any{},
					PlanetPermanentEffects: []any{},
					StartDate:              time.Now(),
					EndDate:                time.Now().Add(24 * time.Hour),
					Planets: []structs.Planet{
						structs.Planet{
							ID:           1,
							Name:         "foo",
							Disabled:     false,
							InitialOwner: "bar",
							MaxHealth:    100.0,
							Position:     structs.Position{X: 1, Y: 2},
							Sector:       "Alpha Centauri",
							Waypoints:    []int{42},
							History: []structs.PlanetHistory{
								structs.PlanetHistory{
									Timestamp:      time.Now(),
									Health:         95.2,
									Liberation:     4.8,
									Owner:          "Humans",
									PlayerCount:    1234567,
									RegenPerSecond: 1.3,
								},
							},
						},
					},
				},
				wantErr: false,
			},
			{
				name: "valid struct incomplete",
				doc: structs.WarSeason{
					ID:       1,
					Capitals: []any{},
				},
				wantErr: true,
			},
			{
				name: "valid struct missing embedded",
				doc: structs.WarSeason{
					ID:                     1,
					Capitals:               []any{},
					PlanetPermanentEffects: []any{},
					StartDate:              time.Now(),
					EndDate:                time.Now().Add(24 * time.Hour),
					Planets:                []structs.Planet{},
				},
				wantErr: true,
			},
			{
				name: "wrong struct",
				doc: structs.Planet{
					ID:           1,
					Name:         "foobar",
					Disabled:     false,
					InitialOwner: "gopher",
					MaxHealth:    100.0,
					Position:     structs.Position{X: 1, Y: 3},
					Sector:       "Alpha Centauri",
					Waypoints:    []int{1, 2, 3},
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
		if err := migration.Up(); err != nil {
			t.Fatalf("failed to migrate up: %v", err)
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				coll := client.database().Collection("war_seasons")
				_, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}
	})
}
