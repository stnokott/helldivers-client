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

func TestMigrateUp(t *testing.T) {
	withClient(t, func(client *Client, _ *migrate.Migrate) {
		if err := client.MigrateUp("../../migrations"); err != nil {
			t.Errorf("client.MigrateUp() error = %v, expected nil", err)
		}
	})
}

func TestCollectionsExist(t *testing.T) {
	collections := []string{
		"war_seasons",
		"war_news",
	}
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
	collections := []string{
		"war_seasons",
		"war_news",
	}
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
					History: []structs.WarSeasonHistory{
						{
							Timestamp:                   time.Now(),
							ActiveElectionPolicyEffects: []int{2},
							CommunityTargets:            []int{1},
							ImpactMultiplier:            1.5,
							GlobalEvents: []structs.WarSeasonHistoryGlobalEvent{
								{
									Title:     "my event",
									Effects:   []string{"active effect"},
									PlanetIDs: []int{2},
									Race:      "Humans",
									Message: structs.WarNewsMessage{
										DE: "de",
										EN: "en",
										ES: "es",
										FR: "fr",
										IT: "it",
										PL: "pl",
										RU: "ru",
										ZH: "zh",
									},
								},
							},
						},
					},
					Planets: []structs.Planet{
						{
							ID:           1,
							Name:         "foo",
							Disabled:     false,
							InitialOwner: "bar",
							MaxHealth:    100.0,
							Position:     structs.Position{X: 1, Y: 2},
							Sector:       "Alpha Centauri",
							Waypoints:    []int{42},
							History: []structs.PlanetHistory{
								{
									Timestamp:      time.Now(),
									Health:         95.2,
									Liberation:     4.8,
									Owner:          "Humans",
									PlayerCount:    1234567,
									RegenPerSecond: 1.3,
									AttackTargets:  []int{234},
									Campaign: &structs.PlanetCampaign{
										Count: 3,
										Type:  2,
									},
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
					ID:                     1,
					Capitals:               []any{},
					PlanetPermanentEffects: []any{},
					StartDate:              time.Now(),
					EndDate:                time.Now().Add(24 * time.Hour),
					History: []structs.WarSeasonHistory{
						{
							Timestamp:                   time.Now(),
							ActiveElectionPolicyEffects: []int{2},
							CommunityTargets:            []int{1},
							ImpactMultiplier:            1.5,
							GlobalEvents: []structs.WarSeasonHistoryGlobalEvent{
								{
									Title:   "my event",
									Effects: []string{"active effect"},
									Message: structs.WarNewsMessage{
										DE: "de",
										EN: "en",
										ES: "es",
										FR: "fr",
										IT: "it",
										PL: "pl",
										RU: "ru",
										ZH: "zh",
									},
								},
							},
						},
					},
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

func TestWarNewsSchema(t *testing.T) {
	withClient(t, func(client *Client, migration *migrate.Migrate) {
		type document any
		tests := []struct {
			name    string
			doc     document
			wantErr bool
		}{
			{
				name: "valid struct complete",
				doc: structs.WarNews{
					ID: 1,
					Message: structs.WarNewsMessage{
						DE: "Das Rauschen der Meereswellen beruhigt meine Seele",
						EN: "The sound of ocean waves calms my soul",
						ES: "El sonido de las olas del mar calma mi alma",
						FR: "Le bruit des vagues de l'océan calme mon âme",
						IT: "Il suono delle onde dell'oceano calma la mia anima",
						PL: "Dźwięk fal oceanu uspokaja moją duszę",
						RU: "Шум океанских волн успокаивает мою душу",
						ZH: "海浪声让我的心灵平静",
					},
					Published: time.Now(),
					Type:      0,
				},
				wantErr: false,
			},
			{
				name: "valid struct incomplete",
				doc: structs.WarNews{
					ID:        1,
					Published: time.Now(),
					Type:      0,
				},
				wantErr: true,
			},
			{
				name: "valid struct missing embedded",
				doc: structs.WarNews{
					ID:        1,
					Message:   structs.WarNewsMessage{},
					Published: time.Now(),
					Type:      0,
				},
				wantErr: true,
			},
			{
				name: "wrong struct",
				doc: structs.WarNewsMessage{
					DE: "Das Rauschen der Meereswellen beruhigt meine Seele",
					EN: "The sound of ocean waves calms my soul",
					ES: "El sonido de las olas del mar calma mi alma",
					FR: "Le bruit des vagues de l'océan calme mon âme",
					IT: "Il suono delle onde dell'oceano calma la mia anima",
					PL: "Dźwięk fal oceanu uspokaja moją duszę",
					RU: "Шум океанских волн успокаивает мою душу",
					ZH: "海浪声让我的心灵平静",
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
				coll := client.database().Collection("war_news")
				_, err := coll.InsertOne(context.Background(), tt.doc)
				if (err != nil) != tt.wantErr {
					t.Errorf("InsertOne() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}
	})
}
