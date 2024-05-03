package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
)

// Events converts API data into mergable DB entities.
func Events(c Converter, data APIData) ([]db.EntityMerger, error) {
	if data.Planets == nil {
		return nil, errors.New(("got nil planets slice (required for events)"))
	}
	if data.Campaigns == nil {
		return nil, errors.New("got nil campaigns slice (required for events)")
	}

	src := *data.Planets
	// we could also use the planet-events endpoint directly which returns only planets with active events.
	// this would introduce an additional API query though which causes more overhead than a simple client-side check in this function.
	events := []db.EntityMerger{}
	for _, planet := range src {
		if planet.Event == nil {
			// event is optional
			continue
		}
		event, err := planet.Event.AsEvent()
		if err != nil {
			return nil, err
		}
		merger, err := c.ConvertEvent(event)
		if err != nil {
			return nil, err
		}
		events = append(events, merger)
	}
	return events, nil
}
