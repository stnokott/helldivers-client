// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/stnokott/helldivers-client/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appName = "HELLDIVERS_2_CLIENT"

// CollectionName is the name of a collection in the MongoDB database.
type CollectionName string

const (
	// CollPlanets is the collection name for Planets
	CollPlanets CollectionName = "planets"
	// CollCampaigns is the collection name for Campaigns
	CollCampaigns CollectionName = "campaigns"
	// CollDispatches is the collection name for Dispatches
	CollDispatches CollectionName = "dispatches"
	// CollEvents is the collection name for Events
	CollEvents CollectionName = "events"
	// CollAssignments is the collection name for Assignments
	CollAssignments CollectionName = "assignments"
	// CollWars is the collection name for Wars
	CollWars CollectionName = "wars"
	// CollSnapshots is the collection name for Snapshots
	CollSnapshots CollectionName = "snapshots"
)

// Client is the abstraction layer for the MongoDB connector
type Client struct {
	mongo *mongo.Client
	db    *mongo.Database
	log   *log.Logger
}

// New creates a new client and connects it to the DB
func New(cfg *config.Config, database string, logger *log.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
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

// DocsProvider wraps collection name and document slice for further processing.
//
// T should be the Document type.
type DocsProvider[T any] struct {
	CollectionName CollectionName
	Docs           []DocWrapper[T]
}

// DocWrapper holds the document to be processed plus its document ID.
//
// T should be the Document type.
type DocWrapper[T any] struct {
	DocID    any
	Document T
}

// UpsertDocs inserts or updates a list of documents based on their IDs.
//
// Documents are matched by ID ("_id"). If a document with that ID already exists, it is updated.
// If it doesn't exist, it is inserted.
//
// If an error occurs during upsert of one of the documents, processing continues.
// Thus, no error is returned.
func UpsertDocs[T any](ctx context.Context, c *Client, provider *DocsProvider[T]) {
	if provider.Docs == nil || len(provider.Docs) == 0 {
		c.log.Printf("upsert into '%s' aborted, no documents to process", provider.CollectionName)
		return
	}

	coll := c.db.Collection(string(provider.CollectionName))

	models := make([]mongo.WriteModel, len(provider.Docs))
	for i, doc := range provider.Docs {
		models[i] = mongo.NewUpdateOneModel().
			SetFilter(bson.D{{Key: "_id", Value: doc.DocID}}).
			SetUpdate(bson.D{{Key: "$set", Value: doc.Document}}).
			SetUpsert(true)
	}
	result, err := coll.BulkWrite(
		ctx,
		models,
		options.BulkWrite().SetOrdered(false), // prevents from stopping on error
	)
	if result != nil {
		c.log.Printf(
			"upsert into '%s' finished, %d inserted, %d matched, %d updated",
			provider.CollectionName,
			result.UpsertedCount, result.MatchedCount, result.ModifiedCount,
		)
	}
	if err != nil {
		c.log.Printf("error(s) occured during upsert into '%s': %v", coll.Name(), err)
	}
}
