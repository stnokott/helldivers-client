package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Campaign)(nil)

// Campaign implements EntityMerger
type Campaign gen.Campaign

// Merge implements EntityMerger. It is assumed that the currently known planets are already present
// in the database.
func (c *Campaign) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	if _, err := tx.MergeCampaign(ctx, gen.MergeCampaignParams(*c)); err != nil {
		return fmt.Errorf("failed to merge assignment (ID=%d): %v", c.ID, err)
	}
	stats.IncrInsert("Campaigns")
	return nil
}
