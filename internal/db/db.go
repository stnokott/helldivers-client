// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appName = "HELLDIVERS_2_CLIENT"

type CollectionName string

const (
	CollPlanets     CollectionName = "planets"
	CollCampaigns   CollectionName = "campaigns"
	CollDispatches  CollectionName = "dispatches"
	CollEvents      CollectionName = "events"
	CollAssignments CollectionName = "assignments"
	CollWars        CollectionName = "wars"
	CollSnapshots   CollectionName = "snapshots"
)

// Client is the abstraction layer for the MongoDB connector
type Client struct {
	mongo *mongo.Client
	db    *mongo.Database
	log   *log.Logger
}

// New creates a new client and connects it to the DB
func New(uri string, database string, logger *log.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientOptions := options.Client().
		ApplyURI(uri).
		SetAppName(appName).
		SetDirect(true)

	logger.Printf("connecting to MongoDB instance at %v", clientOptions.Hosts)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	// ensure connection is stable
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB instance: %w", err)
	}
	logger.Println("connected")
	db := client.Database(database)
	return &Client{
		mongo: client,
		db:    db,
		log:   logger,
	}, nil
}

// Disconnect disconnects from the MongoDB instance
func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.mongo.Disconnect(ctx); err != nil {
		return fmt.Errorf("could not disconnect from MongoDB: %w", err)
	}
	c.log.Println("disconnected from MongoDB")
	return nil
}

type DocProvider interface {
	DocID() any
	Document() any
	CollectionName() CollectionName
}

func (c *Client) UpsertDoc(provider DocProvider, ctx context.Context) (inserted bool, err error) {
	coll := c.db.Collection(string(provider.CollectionName()))
	result, err := coll.UpdateByID(
		ctx,
		provider.DocID(),
		bson.D{{Key: "$set", Value: provider.Document()}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return
	}
	inserted = result.MatchedCount == 0
	return
}
