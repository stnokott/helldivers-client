package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Dispatch)(nil)

// Dispatch implements EntityMerger
type Dispatch gen.Dispatch

func (d *Dispatch) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	exists, err := tx.DispatchExists(ctx, d.ID)
	if err != nil {
		return fmt.Errorf("failed to check if dispatch ID=%d exists: %v", d.ID, err)
	}

	rows, err := tx.MergeDispatch(ctx, gen.MergeDispatchParams(*d))
	if err != nil {
		return fmt.Errorf("failed to merge dispatch (ID=%d): %v", d.ID, err)
	}
	stats.Incr("Dispatches", exists, rows)
	return nil
}
