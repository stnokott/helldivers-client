package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// War implements worker.docTransformer
type War struct{}

func (_ War) Transform(data APIData) (*db.DocsProvider, error) {
	warID := data.WarID
	if warID.Id == nil {
		return nil, errors.New("got nil war ID")
	}

	war := data.War
	if war.Started == nil ||
		war.Ended == nil ||
		war.ImpactMultiplier == nil ||
		war.Factions == nil {
		return nil, errFromNils(war)
	}
	return &db.DocsProvider{
		CollectionName: db.CollWars,
		Docs: []db.DocWrapper{
			{
				DocID: *warID.Id,
				Document: structs.War{
					ID:               *warID.Id,
					StartTime:        db.PrimitiveTime(*war.Started),
					EndTime:          db.PrimitiveTime(*war.Ended),
					ImpactMultiplier: *war.ImpactMultiplier,
					Factions:         *war.Factions,
				},
			},
		},
	}, nil
}
