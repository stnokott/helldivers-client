// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const appName = "HELLDIVERS_2_CLIENT"

// Client is the abstraction layer for the MongoDB connector
type Client struct {
	mongo       *mongo.Client
	log         *log.Logger
	collections struct {
		Seasons *mongo.Collection
		Planets *mongo.Collection
	}
}

// New creates a new client and connects it to the DB
func New(uri string, logger *log.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri).SetAppName(appName)
	logger.Printf("connecting to MongoDB instance at %v", clientOptions.Hosts)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB instance: %w", err)
	}
	logger.Println("connected")
	return &Client{
		mongo: client,
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

// PrepareDB ensure required collections exist.
//
// This should be run before any other operations.
func (c *Client) PrepareDB() (err error) {
	db := c.mongo.Database("helldivers2")
	c.log.Printf("using database '%s'", db.Name())
	c.log.Println("ensuring required collections exist")
	if c.collections.Seasons, err = c.ensureCollection(db, "seasons"); err != nil {
		return
	}
	if c.collections.Planets, err = c.ensureCollection(db, "planets"); err != nil {
		return
	}
	return
}

func (c *Client) ensureCollection(db *mongo.Database, name string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.CreateCollection(ctx, name); err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			c.log.Printf("collection '%s' already exists", name)
		} else {
			return nil, fmt.Errorf("could not create collection '%s': %w", name, err)
		}
	} else {
		c.log.Printf("collection '%s' created", name)
	}
	return db.Collection(name, options.Collection().SetWriteConcern(writeconcern.Majority())), nil
}
