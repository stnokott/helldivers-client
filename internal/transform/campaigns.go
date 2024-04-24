package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
)

func Campaigns(data APIData) ([]db.EntityMerger, error) {
	if data.Campaigns == nil {
		return nil, errors.New("got nil campaigns slice")
	}

	src := *data.Campaigns
	mergers := make([]db.EntityMerger, len(src))
	for i, campaign := range src {
		if campaign.Id == nil ||
			campaign.Type == nil ||
			campaign.Count == nil {
			return nil, errFromNils(&campaign)
		}

		mergers[i] = &db.Campaign{
			ID:    *campaign.Id,
			Type:  *campaign.Type,
			Count: *campaign.Count,
		}
	}
	return mergers, nil
}
