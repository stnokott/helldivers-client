package db

import (
	"context"
	"fmt"
	"log"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Campaign)(nil)

// Campaign implements EntityMerger
type Campaign gen.Campaign

func (c *Campaign) Merge(ctx context.Context, tx *gen.Queries, stats *MergeStats, logger *log.Logger) error {
	logger.Printf("** merging campaign ID=%d", c.ID)

	id, err := tx.GetCampaign(ctx, c.ID)
	exists, err := entityExistsByPK(id, err, c.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing campaign: %v", err)
	}
	if exists {
		// perform UPDATE
		if _, err = tx.UpdateCampaign(ctx, gen.UpdateCampaignParams(*c)); err != nil {
			return fmt.Errorf("failed to update campaign (ID=%d): %v", c.ID, err)
		}
		stats.Updates++
	} else {
		// perform INSERT
		if _, err = tx.InsertCampaign(ctx, gen.InsertCampaignParams(*c)); err != nil {
			return fmt.Errorf("failed to insert assignment (ID=%d): %v", c.ID, err)
		}
		stats.Inserts++
	}
	return nil
}
