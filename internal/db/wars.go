package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*War)(nil)

// War implements EntityMerger
type War gen.War

// Merge implements EntityMerger.
func (w *War) Merge(ctx context.Context, tx *gen.Queries, onMerge onMergeFunc) error {
	exists, err := tx.WarExists(ctx, w.ID)
	if err != nil {
		return fmt.Errorf("failed to check if war ID=%d exists: %v", w.ID, err)
	}

	rows, err := tx.MergeWar(ctx, gen.MergeWarParams(*w))
	if err != nil {
		return fmt.Errorf("failed to merge war (ID=%d): %v", w.ID, err)
	}
	onMerge(gen.TableWars, exists, rows)
	return nil
}
