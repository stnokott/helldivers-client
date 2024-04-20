package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func Planets(data APIData) ([]db.EntityMerger, error) {
	if data.Planets == nil {
		return nil, errors.New("got nil planets slice")
	}

	src := *data.Planets
	planets := make([]db.EntityMerger, len(src))
	for i, planet := range src {
		if planet.Index == nil ||
			planet.Name == nil ||
			planet.Sector == nil ||
			planet.Position == nil ||
			planet.Waypoints == nil ||
			planet.Disabled == nil ||
			planet.Biome == nil ||
			planet.Hazards == nil ||
			planet.MaxHealth == nil ||
			planet.InitialOwner == nil {
			return nil, errFromNils(&planet)
		}

		pos, err := parsePlanetPosition(planet.Position)
		if err != nil {
			return nil, err
		}

		biome, err := parsePlanetBiome(planet.Biome)
		if err != nil {
			return nil, err
		}

		hazards, err := convertPlanetHazards(planet.Hazards)
		if err != nil {
			return nil, err
		}
		hazardNames := make([]string, len(hazards))
		for i, hazard := range hazards {
			hazardNames[i] = hazard.Name
		}

		planets[i] = &db.Planet{
			Planet: gen.Planet{
				ID:           *planet.Index,
				Name:         *planet.Name,
				Sector:       *planet.Sector,
				Position:     []float64{*pos.X, *pos.Y},
				WaypointIds:  *planet.Waypoints,
				Disabled:     *planet.Disabled,
				BiomeName:    *biome.Name,
				HazardNames:  hazardNames,
				MaxHealth:    *planet.MaxHealth,
				InitialOwner: *planet.InitialOwner,
			},
			Biome: db.Biome{
				Name:        *biome.Name,
				Description: *biome.Description,
			},
			Hazards: hazards,
		}
	}
	return planets, nil
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

func convertPlanetHazards(in *[]api.Hazard) ([]db.Hazard, error) {
	hazards := make([]db.Hazard, len(*in))
	for i, hazard := range *in {
		if hazard.Name == nil || hazard.Description == nil {
			return nil, errFromNils(&hazard)
		}
		hazards[i] = db.Hazard{
			Name:        *hazard.Name,
			Description: *hazard.Description,
		}
	}
	return hazards, nil
}
