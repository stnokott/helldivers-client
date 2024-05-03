package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stnokott/helldivers-client/internal/db/gen"
	"github.com/stnokott/helldivers-client/internal/db/stats"
)

type EntityMerger interface {
	Merge(ctx context.Context, tx *gen.Queries, onMerge onMergeFunc) error
}

type onMergeFunc func(table gen.Table, exists bool, affectedRows int64)

func (c *Client) Merge(ctx context.Context, mergers ...[]EntityMerger) (err error) {
	tx, err := c.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	qtx := c.queries.WithTx(tx)

	// prepare insert/update statistics
	stats := stats.NewCollector()
	onMerge := func(table gen.Table, exists bool, affectedRows int64) {
		collectAfterMerge(stats, table, exists, affectedRows)
	}

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
				stats.Print(c.log)
			}
		}
	}()

	// run merges
	for _, mSlice := range mergers {
		if len(mSlice) == 0 {
			c.log.Println("WARN: got 0 entities to merge")
		}
		for _, merger := range mSlice {
			if err = merger.Merge(ctx, qtx, onMerge); err != nil {
				return
			}
		}
	}
	return
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
