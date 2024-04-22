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
		return nil, fmt.Errorf("failed to parse PostgreSQL config from ENV: %w", err)
	}
	pgxConfig.RuntimeParams["application_name"] = appName

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Printf("connecting to PostgreSQL instance at %s:%d/%s", pgxConfig.Host, pgxConfig.Port, pgxConfig.Database)
	conn, err := pgx.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	queries := gen.New(conn)

	// ensure connection is stable
	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL instance: %w", err)
	}
	logger.Println("connected")
	return &Client{
		conn:    conn,
		queries: queries,
		log:     logger,
	}, nil
}

// Disconnect disconnects from the MongoDB instance
func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.conn.Close(ctx); err != nil {
		return fmt.Errorf("could not disconnect from PostgreSQL: %w", err)
	}
	c.log.Println("disconnected from PostgreSQL")
	return nil
}

// TODO: implement hashing on non-PK to differentiate between simple updates and actual changes
type EntityMerger interface {
	Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error
}

func (c *Client) Merge(ctx context.Context, mergers ...[]EntityMerger) (err error) {
	tx, err := c.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	qtx := c.queries.WithTx(tx)

	// defer commit/rollback depending on error
	defer func() {
		if err != nil {
			// roll back on error
			c.log.Println("error occured during merge, rolling back changes")
			if errRb := tx.Rollback(ctx); errRb != nil {
				c.log.Printf("failed to rollback: %v", errRb)
			}
		} else {
			// commit when no error
			if errComm := tx.Commit(ctx); errComm != nil {
				c.log.Printf("failed to commit: %v", errComm)
			} else {
				c.log.Println("changes committed")
			}
		}
	}()

	// prepare insert/update statistics
	stats := tableMergeStats{}
	// run merges
	for _, mSlice := range mergers {
		if len(mSlice) == 0 {
			c.log.Println("WARN: got 0 entities to merge")
		}
		for _, merger := range mSlice {
			if err = merger.Merge(ctx, qtx, stats); err != nil {
				return
			}
		}
	}
	stats.Print(c.log)
	return
}

func PGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

func PGUint64(x uint64) pgtype.Numeric {
	return pgtype.Numeric{Int: new(big.Int).SetUint64(x), Valid: true}
}
