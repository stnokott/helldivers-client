package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
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
			campaign.Planet == nil ||
			campaign.Type == nil ||
			campaign.Count == nil {
			return nil, errFromNils(&campaign)
		}

		planetRef, err := parseCampaignPlanet(campaign.Planet)
		if err != nil {
			return nil, err
		}

		mergers[i] = &db.Campaign{
			ID:       *campaign.Id,
			PlanetID: *planetRef.Index,
			Type:     *campaign.Type,
			Count:    *campaign.Count,
		}
	}
	return mergers, nil
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
