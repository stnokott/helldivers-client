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
	id, err := tx.GetDispatch(ctx, d.ID)
	exists, err := entityExistsByPK(id, err, d.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing dispatch: %v", err)
	}
	if exists {
		// perform UPDATE
		if _, err = tx.UpdateDispatch(ctx, gen.UpdateDispatchParams(*d)); err != nil {
			return fmt.Errorf("failed to update dispatch (ID=%d): %v", d.ID, err)
		}
		stats.IncrUpdate("Dispatches")
	} else {
		// perform INSERT
		if _, err = tx.InsertDispatch(ctx, gen.InsertDispatchParams(*d)); err != nil {
			return fmt.Errorf("failed to insert dispatch (ID=%d): %v", d.ID, err)
		}
		stats.IncrInsert("Dispatches")
	}
	return nil
}
