package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Planets implements worker.docTransformer
type Planets struct{}

func (_ Planets) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.Planet] {
	provider := &db.DocsProvider[structs.Planet]{
		CollectionName: db.CollPlanets,
		Docs:           []db.DocWrapper[structs.Planet]{},
	}

	if data.Planets == nil {
		errFunc(errors.New("got nil planets slice"))
		return provider
	}

	planets := *data.Planets

	for _, planet := range planets {
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
			errFunc(errFromNils(&planet))
			continue
		}

		pos, err := parsePlanetPosition(planet.Position)
		if err != nil {
			errFunc(err)
			continue
		}

		biome, err := parsePlanetBiome(planet.Biome)
		if err != nil {
			errFunc(err)
			continue
		}

		hazards, err := convertPlanetHazards(planet.Hazards)
		if err != nil {
			errFunc(err)
			continue
		}
		provider.Docs = append(provider.Docs, db.DocWrapper[structs.Planet]{
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
		}) 
	}
	return provider
}

func parsePlanetPosition(in *api.Planet_Position) (api.Position, error) {
	pos, err := in.AsPosition()
	if err != nil {
		return api.Position{}, fmt.Errorf("cannot parse planet position: %w", err)
	}
	if pos.X == nil || pos.Y == nil {
		return api.Position{}, errFromNils(&pos)
	}
	return pos, nil
}

func parsePlanetBiome(in *api.Planet_Biome) (api.Biome, error) {
	biome, err := in.AsBiome()
	if err != nil {
		return api.Biome{}, fmt.Errorf("cannot parse planet position: %w", err)
	}
	if biome.Name == nil || biome.Description == nil {
		return api.Biome{}, errFromNils(&biome)
	}
	return biome, nil
}

func convertPlanetHazards(in *[]api.Hazard) ([]structs.Hazard, error) {
	hazards := make([]structs.Hazard, len(*in))
	for i, hazard := range *in {
		if hazard.Name == nil || hazard.Description == nil {
			return nil, errFromNils(&hazard)
		}
		hazards[i] = structs.Hazard{
			Name:        *hazard.Name,
			Description: *hazard.Description,
		}
	}
	return hazards, nil
}
