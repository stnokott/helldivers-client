// Package db handles interactions with the MongoDB instance and works as an abstraction layer
package db

import (
	"context"
	"errors"
	"fmt"
	"log"
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

type EntityMerger interface {
	Merge(ctx context.Context, tx *gen.Queries, stats *MergeStats, logger *log.Logger) error
}

type MergeStats struct {
	Inserts int
	Updates int
}

func (c *Client) Merge(ctx context.Context, mergers ...[]EntityMerger) (err error) {
	tx, err := c.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	qtx := c.queries.WithTx(tx)
	stats := MergeStats{}
	defer func() {
		if err != nil {
			// roll back on error
			c.log.Println("error occured during merge, rolling back changes")
			if errRb := tx.Rollback(ctx); errRb != nil {
				c.log.Printf("failed to rollback: %v", errRb)
			}
		} else {
			// commit when no error
			c.log.Println("committing changes")
			if errComm := tx.Commit(ctx); errComm != nil {
				c.log.Printf("failed to commit: %v", errComm)
			} else {
				c.log.Printf("committed %d inserts and %d updates", stats.Inserts, stats.Updates)
			}
		}
	}()

	for _, mSlice := range mergers {
		if len(mSlice) == 0 {
			c.log.Println("WARN: got 0 entities to merge")
		}
		for _, merger := range mSlice {
			if err = merger.Merge(ctx, qtx, &stats, c.log); err != nil {
				return
			}
		}
	}
	return
}

func PGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

func entityExistsByPK[PK comparable](pk PK, err error, expected PK) (bool, error) {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return pk == expected, nil
}
