package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Events implements worker.docTransformer
type Events struct{}

func (_ Events) Transform(data APIData) (*db.DocsProvider[structs.Event], error) {
	if data.Planets == nil {
		return nil, errors.New("got nil planets slice (required for events)")
	}

	planets := *data.Planets
	eventDocs := make([]db.DocWrapper[structs.Event], 0, len(planets))

	for _, planet := range planets {
		if planet.Event == nil {
			continue
		}
		event, err := planet.Event.AsEvent()
		if err != nil {
			return nil, err
		}
		if event.Id == nil ||
			event.EventType == nil ||
			event.Faction == nil ||
			event.MaxHealth == nil ||
			event.StartTime == nil ||
			event.EndTime == nil {
			return nil, errFromNils(&event)
		}

		eventDocs = append(eventDocs, db.DocWrapper[structs.Event]{
			DocID: *event.Id,
			Document: structs.Event{
				ID:        *event.Id,
				Type:      *event.EventType,
				Faction:   *event.Faction,
				MaxHealth: *event.MaxHealth,
				StartTime: db.PrimitiveTime(*event.StartTime),
				EndTime:   db.PrimitiveTime(*event.EndTime),
			},
		})
	}
	return &db.DocsProvider[structs.Event]{
		CollectionName: db.CollEvents,
		Docs:           eventDocs,
	}, nil
}
