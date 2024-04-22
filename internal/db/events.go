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
	if _, err := tx.MergeEvent(ctx, gen.MergeEventParams(*e)); err != nil {
		return fmt.Errorf("failed to merge event (ID=%d): %v", e.ID, err)
	}
	stats.IncrInsert("Events")
	return nil
}
