package transform

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Planets implements worker.docTransformer[warRequestData]
type Planets struct{}

type planetsRequestData *[]api.Planet

func (_ Planets) Request(api *client.Client, ctx context.Context) (data planetsRequestData, err error) {
	data, err = apiWithTimeout(api.Planets, 1*time.Second)
	return
}

func (_ Planets) Transform(data planetsRequestData) (*db.DocsProvider, error) {
	if data == nil {
		return nil, errors.New("got nil planets slice")
	}

	planets := *data
	planetDocs := make([]db.DocWrapper, len(planets))

	for i, planet := range planets {
		if planet.Index == nil ||
			planet.Name == nil ||
			planet.Sector == nil ||
			planet.Position == nil ||
			planet.Waypoints == nil ||
			planet.Disabled == nil ||
			planet.MaxHealth == nil ||
			planet.InitialOwner == nil ||
			planet.RegenPerSecond == nil {
			return nil, errFromNils(&planet)
		}

		pos, err := planet.Position.AsPosition()
		if err != nil {
			return nil, fmt.Errorf("cannot parse planet position: %w", err)
		}
		if pos.X == nil || pos.Y == nil {
			return nil, errFromNils(&pos)
		}
		planetDocs[i] = db.DocWrapper{
			DocID: *planet.Index,
			Document: structs.Planet{
				ID:             *planet.Index,
				Name:           *planet.Name,
				Sector:         *planet.Sector,
				Position:       structs.PlanetPosition{X: *pos.X, Y: *pos.Y},
				Waypoints:      *planet.Waypoints,
				Disabled:       *planet.Disabled,
				MaxHealth:      *planet.MaxHealth,
				InitialOwner:   *planet.InitialOwner,
				RegenPerSecond: *planet.RegenPerSecond,
			},
		}
	}
	return &db.DocsProvider{
		CollectionName: db.CollPlanets,
		Docs:           planetDocs,
	}, nil
}
