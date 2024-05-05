package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stnokott/helldivers-client/internal/db/gen"
	"github.com/stnokott/helldivers-client/internal/db/stats"
)

// EntityMerger provide a means of merging an entity to the database.
type EntityMerger interface {
	// Merge merges the implementing entity to the database. It should call `onMerge` when finished.
	Merge(ctx context.Context, tx *gen.Queries, onMerge onMergeFunc) error
}

type onMergeFunc func(table gen.Table, exists bool, affectedRows int64)

// Merge attempts to merge each `EntityMerger` to the database.
//
// It will print statistics once finished.
func (c *Client) Merge(ctx context.Context, mergers ...[]EntityMerger) error {
	mergeFunc := func(qtx *gen.Queries, onMerge onMergeFunc) error {
		// run merges
		for _, mSlice := range mergers {
			if len(mSlice) == 0 {
				c.log.Println("WARN: got 0 entities to merge")
			}
			for _, merger := range mSlice {
				if err := merger.Merge(ctx, qtx, onMerge); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return c.withTx(ctx, mergeFunc)
}

func (c *Client) withTx(ctx context.Context, txFunc func(*gen.Queries, onMergeFunc) error) error {
	tx, err := c.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// prepare insert/update statistics
	stats := stats.NewCollector()
	onMerge := func(table gen.Table, exists bool, affectedRows int64) {
		collectAfterMerge(stats, table, exists, affectedRows)
	}

	qtx := c.queries.WithTx(tx)

	err = txFunc(qtx, onMerge)
	if err != nil {
		// roll back on error
		c.rollback(ctx, tx)
	} else {
		// commit when no error
		if errComm := tx.Commit(ctx); errComm != nil {
			return fmt.Errorf("failed to commit: %w", errComm)
		}
		c.log.Println("changes committed")
		stats.Print(c.log)
	}
	return nil
}

func (c *Client) rollback(ctx context.Context, tx pgx.Tx) {
	c.log.Println("error occured during merge, rolling back changes")
	if errRb := tx.Rollback(ctx); errRb != nil {
		c.log.Printf("failed to rollback: %v", errRb)
	}
}

func collectAfterMerge(s stats.Collector, table gen.Table, exists bool, affectedRows int64) {
	if affectedRows == 0 {
		s.Noop(table, 1)
		return
	}
	if exists {
		s.Updated(table, affectedRows)
	} else {
		s.Inserted(table, affectedRows)
	}
}
