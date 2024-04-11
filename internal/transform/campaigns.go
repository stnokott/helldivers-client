package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Campaigns implements worker.docTransformer
type Campaigns struct{}

func (_ Campaigns) Transform(data APIData) (*db.DocsProvider, error) {
	if data.Campaigns == nil {
		return nil, errors.New("got nil campaigns slice")
	}

	campaigns := *data.Campaigns
	campaignDocs := make([]db.DocWrapper, len(campaigns))

	for i, campaign := range campaigns {
		if campaign.Id == nil ||
			campaign.Planet == nil ||
			campaign.Type == nil ||
			campaign.Count == nil {
			return nil, errFromNils(&campaign)
		}

		planetRef, err := campaign.Planet.AsPlanet()
		if err != nil {
			return nil, fmt.Errorf("cannot parse campaign planet: %w", err)
		}
		if planetRef.Index == nil {
			return nil, errors.New("campaign planet ID is nil")
		}

		campaignDocs[i] = db.DocWrapper{
			DocID: *campaign.Id,
			Document: structs.Campaign{
				ID:       *campaign.Id,
				PlanetID: *planetRef.Index,
				Type:     *campaign.Type,
				Count:    *campaign.Count,
			},
		}
	}
	return &db.DocsProvider{
		CollectionName: db.CollCampaigns,
		Docs:           campaignDocs,
	}, nil
}
