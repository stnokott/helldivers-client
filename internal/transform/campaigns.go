package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
)

// Campaigns converts API data into mergable DB entities.
func Campaigns(c Converter, data APIData) ([]db.EntityMerger, error) {
	if data.Campaigns == nil {
		return nil, errors.New("got nil campaigns slice")
	}

	src := *data.Campaigns
	mergers := make([]db.EntityMerger, len(src))
	for i, campaign := range src {
		merger, err := c.ConvertCampaign(campaign)
		if err != nil {
			return nil, err
		}
		mergers[i] = merger
	}
	return mergers, nil
}
