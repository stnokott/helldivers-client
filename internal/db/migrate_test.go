package db

import (
	"context"
	"log"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/bson"
)

func withClient(t *testing.T, do func(client *Client)) {
	mongoURI := getMongoURI()
	client, err := New(mongoURI, t.Name(), log.Default())
	if err != nil {
		t.Skipf("could not initialize DB connection: %v", err)
	}
	do(client)
	defer func() {
		if err = client.mongo.Database(client.dbName).Drop(context.Background()); err != nil {
			t.Logf("could not drop database: %v", err)
		}
		if err = client.Disconnect(); err != nil {
			t.Logf("could not disconnect: %v", err)
		}
	}()
}

func TestClientMigration(t *testing.T) {
	withClient(t, func(client *Client) {
		migration, err := client.newMigration("../../migrations")
		if err != nil {
			t.Fatalf("client.newMigration() error = %v, want nil", err)
		}

		if err = migration.Up(); err != nil {
			t.Fatalf("failed to migrate up: %v", err)
		}
		fnPlanetCollectionExists := func() bool {
			colls, errList := client.database().ListCollectionNames(context.Background(), bson.D{{Key: "name", Value: "planets"}})
			if errList != nil {
				t.Errorf("could not list collections: %v", errList)
				return false
			}
			return len(colls) == 1
		}
		if !fnPlanetCollectionExists() {
			t.Error("expected collection with name 'planets', none found")
		}
		if err = migration.Down(); err != nil {
			t.Fatalf("failed to migrate down: %v", err)
		}
		if fnPlanetCollectionExists() {
			t.Error("expected collection with name 'planets' to not exist")
		}
	})
}
