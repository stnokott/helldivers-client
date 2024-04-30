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

func (e *Event) Merge(ctx context.Context, tx *gen.Queries, onMerge onMergeFunc) error {
	exists, err := tx.DispatchExists(ctx, e.ID)
	if err != nil {
		return fmt.Errorf("failed to check if event ID=%d exists: %v", e.ID, err)
	}

	rows, err := tx.MergeEvent(ctx, gen.MergeEventParams(*e))
	if err != nil {
		return fmt.Errorf("failed to merge event (ID=%d): %v", e.ID, err)
	}
	onMerge(gen.TableEvents, exists, rows)
	return nil
}
