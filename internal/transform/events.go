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

	src := *data.Planets
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

		// TODO: change to Planet->Event->Campaign
		// campaign relation (via CampaignId) is performed via Planet->Event / Planet->Campaign
		events = append(events, &db.Event{
			ID:        *event.Id,
			StartTime: db.PGTimestamp(*event.StartTime),
			EndTime:   db.PGTimestamp(*event.EndTime),
			Type:      *event.EventType,
			Faction:   *event.Faction,
			MaxHealth: *event.MaxHealth,
		})
	}
	return events, nil
}

func parsePlanetEvent(in *api.Planet_Event) (api.Event, error) {
	event, err := in.AsEvent()
	if err != nil {
		return api.Event{}, fmt.Errorf("cannot parse planet event: %w", err)
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
