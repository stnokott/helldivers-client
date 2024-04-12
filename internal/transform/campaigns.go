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

func (_ Campaigns) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.Campaign] {
	provider := &db.DocsProvider[structs.Campaign]{
		CollectionName: db.CollCampaigns,
		Docs:           []db.DocWrapper[structs.Campaign]{},
	}

	if data.Campaigns == nil {
		errFunc(errors.New("got nil campaigns slice"))
		return provider
	}

	campaigns := *data.Campaigns

	for _, campaign := range campaigns {
		if campaign.Id == nil ||
			campaign.Planet == nil ||
			campaign.Type == nil ||
			campaign.Count == nil {
			errFunc(errFromNils(&campaign))
			continue
		}

		planetRef, err := parseCampaignPlanet(campaign.Planet)
		if err != nil {
			errFunc(err)
			continue
		}

		provider.Docs = append(provider.Docs, db.DocWrapper[structs.Campaign]{
			DocID: *campaign.Id,
			Document: structs.Campaign{
				ID:       *campaign.Id,
				PlanetID: *planetRef.Index,
				Type:     *campaign.Type,
				Count:    *campaign.Count,
			},
		})
	}
	return provider
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
