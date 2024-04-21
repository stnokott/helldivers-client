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
	id, err := tx.GetWar(ctx, w.ID)
	exists, err := entityExistsByPK(id, err, w.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing war: %v", err)
	}
	if exists {
		// perform UPDATE
		if _, err = tx.UpdateWar(ctx, gen.UpdateWarParams(*w)); err != nil {
			return fmt.Errorf("failed to update war (ID=%d): %v", w.ID, err)
		}
		stats.IncrUpdate("Wars")
	} else {
		// perform INSERT
		if _, err = tx.InsertWar(ctx, gen.InsertWarParams(*w)); err != nil {
			return fmt.Errorf("failed to insert war (ID=%d): %v", w.ID, err)
		}
		stats.IncrInsert("Wars")
	}
	return nil
}
