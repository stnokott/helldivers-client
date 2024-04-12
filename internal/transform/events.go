package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Events implements worker.docTransformer
type Events struct{}

func (_ Events) Transform(data APIData) (*db.DocsProvider[structs.Event], error) {
	if data.Events == nil {
		return nil, errors.New("got nil events slice")
	}

	events := *data.Events
	eventDocs := make([]db.DocWrapper[structs.Event], len(events))

	for i, event := range events {
		if event.Id == nil ||
			event.EventType == nil ||
			event.Faction == nil ||
			event.MaxHealth == nil ||
			event.StartTime == nil ||
			event.EndTime == nil {
			return nil, errFromNils(&event)
		}

		eventDocs[i] = db.DocWrapper[structs.Event]{
			DocID: *event.Id,
			Document: structs.Event{
				ID:        *event.Id,
				Type:      *event.EventType,
				Faction:   *event.Faction,
				MaxHealth: *event.MaxHealth,
				StartTime: db.PrimitiveTime(*event.StartTime),
				EndTime:   db.PrimitiveTime(*event.EndTime),
			},
		}
	}
	return &db.DocsProvider[structs.Event]{
		CollectionName: db.CollEvents,
		Docs:           eventDocs,
	}, nil
}
