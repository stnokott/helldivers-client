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
	if _, err := tx.MergeDispatch(ctx, gen.MergeDispatchParams(*d)); err != nil {
		return fmt.Errorf("failed to merge dispatch (ID=%d): %v", d.ID, err)
	}
	stats.IncrInsert("Dispatches")
	return nil
}
