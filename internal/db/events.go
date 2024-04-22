package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Event)(nil)

// Event implements EntityMerger
type Event gen.Event

func (e *Event) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	id, err := tx.GetEvent(ctx, e.ID)
	exists, err := entityExistsByPK(id, err, e.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing event: %v", err)
	}
	if exists {
		// perform UPDATE
		if _, err = tx.UpdateEvent(ctx, gen.UpdateEventParams(*e)); err != nil {
			return fmt.Errorf("failed to update event (ID=%d): %v", e.ID, err)
		}
		stats.IncrUpdate("Events")
	} else {
		// perform INSERT
		if _, err = tx.InsertEvent(ctx, gen.InsertEventParams(*e)); err != nil {
			return fmt.Errorf("failed to insert event (ID=%d): %v", e.ID, err)
		}
		stats.IncrInsert("Events")
	}
	return nil
}
