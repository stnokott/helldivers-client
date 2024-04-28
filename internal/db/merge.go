package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

type EntityMerger interface {
	Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error
}

func (c *Client) Merge(ctx context.Context, mergers ...[]EntityMerger) (err error) {
	tx, err := c.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
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
