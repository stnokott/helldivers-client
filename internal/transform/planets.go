package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Planets implements worker.docTransformer
type Planets struct{}

func (_ Planets) Transform(data APIData) (*db.DocsProvider[structs.Planet], error) {
	if data.Planets == nil {
		return nil, errors.New("got nil planets slice")
	}

	planets := *data.Planets
	planetDocs := make([]db.DocWrapper[structs.Planet], len(planets))

	for i, planet := range planets {
		if planet.Index == nil ||
			planet.Name == nil ||
			planet.Sector == nil ||
			planet.Position == nil ||
			planet.Waypoints == nil ||
			planet.Disabled == nil ||
			planet.Biome == nil ||
			planet.Hazards == nil ||
			planet.MaxHealth == nil ||
			planet.InitialOwner == nil ||
			planet.RegenPerSecond == nil {
			return nil, errFromNils(&planet)
		}

		// TODO: move to function
		// TODO: keep processing on error
		pos, err := planet.Position.AsPosition()
		if err != nil {
			return nil, fmt.Errorf("cannot parse planet position: %w", err)
		}
		if pos.X == nil || pos.Y == nil {
			return nil, errFromNils(&pos)
		}

		biome, err := planet.Biome.AsBiome()
		if err != nil {
			return nil, fmt.Errorf("cannot parse planet biome: %w", err)
		}
		if biome.Name == nil || biome.Description == nil {
			return nil, errFromNils(&biome)
		}

		hazardsRaw := *planet.Hazards
		hazards := make([]structs.Hazard, len(hazardsRaw))
		for i, hazard := range hazardsRaw {
			if hazard.Name == nil || hazard.Description == nil {
				return nil, errFromNils(&hazard)
			}
			hazards[i] = structs.Hazard{
				Name:        *hazard.Name,
				Description: *hazard.Description,
			}
		}
		planetDocs[i] = db.DocWrapper[structs.Planet]{
			DocID: *planet.Index,
			Document: structs.Planet{
				ID:        *planet.Index,
				Name:      *planet.Name,
				Sector:    *planet.Sector,
				Position:  structs.PlanetPosition{X: *pos.X, Y: *pos.Y},
				Waypoints: *planet.Waypoints,
				Disabled:  *planet.Disabled,
				Biome: structs.Biome{
					Name:        *biome.Name,
					Description: *biome.Description,
				},
				Hazards:        hazards,
				MaxHealth:      *planet.MaxHealth,
				InitialOwner:   *planet.InitialOwner,
				RegenPerSecond: *planet.RegenPerSecond,
			},
		}
	}
	return &db.DocsProvider[structs.Planet]{
		CollectionName: db.CollPlanets,
		Docs:           planetDocs,
	}, nil
}
