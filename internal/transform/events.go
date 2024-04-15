package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Events implements worker.DocTransformer
type Events struct{}

// Transform implements the worker.DocTransformer interface
func (Events) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.Event] {
	provider := &db.DocsProvider[structs.Event]{
		CollectionName: db.CollEvents,
		Docs:           []db.DocWrapper[structs.Event]{},
	}

	if data.Planets == nil {
		errFunc(errors.New("got nil planets slice (required for events)"))
		return provider
	}

	planets := *data.Planets

	for _, planet := range planets {
		if planet.Event == nil {
			continue
		}
		event, err := parsePlanetEvent(planet.Event)
		if err != nil {
			errFunc(err)
			continue
		}

		provider.Docs = append(provider.Docs, db.DocWrapper[structs.Event]{
			DocID: *event.Id,
			Document: structs.Event{
				ID:        *event.Id,
				Type:      *event.EventType,
				Faction:   *event.Faction,
				MaxHealth: *event.MaxHealth,
				StartTime: primitive.NewDateTimeFromTime(*event.StartTime),
				EndTime:   primitive.NewDateTimeFromTime(*event.EndTime),
			},
		})
	}
	return provider
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
