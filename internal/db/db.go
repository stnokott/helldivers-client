// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

const appName = "HELLDIVERS_2_CLIENT"

// Client is the abstraction layer for the MongoDB connector
type Client struct {
	conn    *pgx.Conn
	queries *gen.Queries
	log     *log.Logger
}

// New creates a new client and connects it to the DB
func New(cfg *config.Config, logger *log.Logger) (*Client, error) {
	pgxConfig, err := pgx.ParseConfig(cfg.PostgresURI)
	if err != nil {
		return nil, fmt.Errorf("parse config from ENV: %w", err)
	}
	pgxConfig.RuntimeParams["application_name"] = appName

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Printf("connecting to %s:%d/%s", pgxConfig.Host, pgxConfig.Port, pgxConfig.Database)
	conn, err := pgx.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	queries := gen.New(conn)

	return &Client{
		conn:    conn,
		queries: queries,
		log:     logger,
	}, nil
}

// Connect implements main.ConnectWaiter.
func (c *Client) Connect(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// ensure connection is stable
			if err := c.conn.Ping(ctx); err != nil {
				c.log.Printf("connect: %v", err)
				continue
			}
			return nil
		}
	}
}

// Disconnect disconnects from the MongoDB instance
func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.conn.Close(ctx); err != nil {
		return fmt.Errorf("disconnect: %w", err)
	}
	c.log.Println("disconnected")
	return nil
}

// PGTimestamp converts a `time.Time` to a `pgx`-compatible `pgtype.Timestamp`.
func PGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

// PGUint64 converts a `uint64` to a `pgx`-compatible `pgtype.Numeric`.
func PGUint64(x uint64) pgtype.Numeric {
	return pgtype.Numeric{Int: new(big.Int).SetUint64(x), Valid: true}
}
