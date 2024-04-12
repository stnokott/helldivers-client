package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Campaigns implements worker.docTransformer
type Campaigns struct{}

func (_ Campaigns) Transform(data APIData) (*db.DocsProvider[structs.Campaign], error) {
	if data.Campaigns == nil {
		return nil, errors.New("got nil campaigns slice")
	}

	campaigns := *data.Campaigns
	campaignDocs := make([]db.DocWrapper[structs.Campaign], len(campaigns))

	for i, campaign := range campaigns {
		if campaign.Id == nil ||
			campaign.Planet == nil ||
			campaign.Type == nil ||
			campaign.Count == nil {
			return nil, errFromNils(&campaign)
		}

		planetRef, err := parseCampaignPlanet(campaign.Planet)
		if err != nil {
			return nil, err
		}

		campaignDocs[i] = db.DocWrapper[structs.Campaign]{
			DocID: *campaign.Id,
			Document: structs.Campaign{
				ID:       *campaign.Id,
				PlanetID: *planetRef.Index,
				Type:     *campaign.Type,
				Count:    *campaign.Count,
			},
		}
	}
	return &db.DocsProvider[structs.Campaign]{
		CollectionName: db.CollCampaigns,
		Docs:           campaignDocs,
	}, nil
}

func parseCampaignPlanet(in *api.Campaign2_Planet) (api.Planet, error) {
	planet, err := in.AsPlanet()
	if err != nil {
		return api.Planet{}, fmt.Errorf("cannot parse campaign planet: %w", err)
	}
	if planet.Index == nil {
		return api.Planet{}, errFromNils(&planet)
	}
	return planet, nil
}
