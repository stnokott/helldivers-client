package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

func Events(data APIData) ([]db.EntityMerger, error) {
	if data.Planets == nil {
		return nil, errors.New(("got nil planets slice (required for events)"))
	}
	if data.Campaigns == nil {
		return nil, errors.New("got nil campaigns slice (required for events)")
	}

	src := *data.Planets
	campaignMap, err := mapCampaigns(*data.Campaigns)
	if err != nil {
		return nil, err
	}
	// TODO: query and use planet-events instead to avoid confusion here
	events := []db.EntityMerger{}
	for _, planet := range src {
		if planet.Event == nil {
			// event is optional
			continue
		}
		event, err := parsePlanetEvent(planet.Event)
		if err != nil {
			return nil, err
		}
		if planet.Index == nil ||
			planet.Name == nil {
			return nil, errFromNils(&planet)
		}
		campaignID, ok := campaignMap[*planet.Index]
		if !ok {
			return nil, fmt.Errorf("planet '%s' has event, but no campaign", *planet.Name)
		}

		events = append(events, &db.Event{
			ID:         *event.Id,
			CampaignID: campaignID,
			StartTime:  db.PGTimestamp(*event.StartTime),
			EndTime:    db.PGTimestamp(*event.EndTime),
			Type:       *event.EventType,
			Faction:    *event.Faction,
			MaxHealth:  *event.MaxHealth,
		})
	}
	return events, nil
}

// mapCampaigns maps planet IDs to their campaign ID
func mapCampaigns(campaigns []api.Campaign2) (map[int32]int32, error) {
	m := map[int32]int32{}
	for _, c := range campaigns {
		if c.Id == nil {
			return nil, errors.New("got nil campaign ID")
		}
		planet, err := parseCampaignPlanet(c.Planet)
		if err != nil {
			return nil, err
		}
		if planet.Index == nil {
			return nil, errors.New("got nil planet index")
		}
		m[*planet.Index] = *c.Id
	}
	return m, nil
}

func parseCampaignPlanet(in *api.Campaign2_Planet) (api.Planet, error) {
	if in == nil {
		return api.Planet{}, errors.New("got nil campaign planet")
	}
	planet, err := in.AsPlanet()
	if err != nil {
		return api.Planet{}, fmt.Errorf("parse campaign planet: %w", err)
	}
	if planet.Index == nil {
		return api.Planet{}, errFromNils(&planet)
	}
	return planet, nil
}

func parsePlanetEvent(in *api.Planet_Event) (api.Event, error) {
	event, err := in.AsEvent()
	if err != nil {
		return api.Event{}, fmt.Errorf("parse planet event: %w", err)
	}
	if event.Id == nil ||
		event.EventType == nil ||
		event.Faction == nil ||
		event.MaxHealth == nil ||
		event.StartTime == nil ||
		event.EndTime == nil {
		return api.Event{}, errFromNils(&event)
	}
	return event, nil
}
