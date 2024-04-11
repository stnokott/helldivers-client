package transform

import (
	"context"
	"errors"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// War implements worker.docTransformer[warRequestData]
type War struct{}

type warRequestData struct {
	WarID *api.WarId
	War   *api.War
}

func (_ War) Request(api *client.Client, ctx context.Context) (data warRequestData, err error) {
	data.WarID, err = apiWithTimeout(api.WarID, 1*time.Second)
	if err != nil {
		return
	}
	data.War, err = apiWithTimeout(api.War, 5*time.Second)
	return
}

func (_ War) Transform(data warRequestData) (*db.DocsProvider, error) {
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
