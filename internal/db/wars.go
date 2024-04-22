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

func (w *War) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	if _, err := tx.MergeWar(ctx, gen.MergeWarParams(*w)); err != nil {
		return fmt.Errorf("failed to merge war (ID=%d): %v", w.ID, err)
	}
	stats.IncrInsert("Wars")
	return nil
}
