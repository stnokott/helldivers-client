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
func (c *Campaign) Merge(ctx context.Context, tx *gen.Queries, onMerge onMergeFunc) error {
	exists, err := tx.CampaignExists(ctx, c.ID)
	if err != nil {
		return fmt.Errorf("failed to check if campaign ID=%d exists: %v", c.ID, err)
	}

	rows, err := tx.MergeCampaign(ctx, gen.MergeCampaignParams(*c))
	if err != nil {
		return fmt.Errorf("failed to merge campaign (ID=%d): %v", c.ID, err)
	}

	onMerge(gen.TableCampaigns, exists, rows)
	return nil
}
