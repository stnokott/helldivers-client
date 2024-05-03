package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	gen "github.com/stnokott/helldivers-client/internal/db/gen"
)

// Planets converts API data into mergable DB entities.
func Planets(c Converter, data APIData) ([]db.EntityMerger, error) {
	if data.Planets == nil {
		return nil, errors.New("got nil planets slice")
	}

	src := *data.Planets
	planets := make([]db.EntityMerger, len(src))
	for i, planet := range src {
		converted, err := c.ConvertPlanet(planet)
		if err != nil {
			return nil, err
		}
		planets[i] = converted
	}
	return planets, nil
}

// MustPlanetName returns the default locale for a localized planet name.
func MustPlanetName(source *api.Planet_Name) (string, error) {
	if source == nil {
		return "", errors.New("Planet name is nil")
	}
	return source.AsPlanetName0()
}

// MustPlanetHazards implements a converter for planet hazards.
func MustPlanetHazards(c Converter, source *[]api.Hazard) ([]gen.Hazard, error) {
	if source == nil {
		return nil, errors.New("Planet hazards is nil")
	}
	src := *source
	hazards := make([]gen.Hazard, len(src))
	for i, hazard := range src {
		converted, err := c.ConvertPlanetHazard(hazard)
		if err != nil {
			return nil, err
		}
		hazards[i] = converted
	}
	return hazards, nil
}

// MustPlanetHazardNames implements a converter for planet hazard names.
func MustPlanetHazardNames(source *[]api.Hazard) ([]string, error) {
	if source == nil {
		return nil, errors.New("Planet hazards is nil")
	}
	src := *source
	names := make([]string, len(src))
	for i, hazard := range src {
		if hazard.Name == nil {
			return nil, fmt.Errorf("Hazard name at index %d is nil", i)
		}
		names[i] = *hazard.Name
	}
	return names, nil
}

// MustPlanetPosition implements a converter for a planet position (i.e. coordinate).
func MustPlanetPosition(source *api.Planet_Position) ([]float64, error) {
	if source == nil {
		return nil, errors.New("Planet position is nil")
	}
	parsed, err := source.AsPosition()
	if err != nil {
		return nil, err
	}
	if parsed.X == nil || parsed.Y == nil {
		return nil, errors.New("Planet position coordinates are nil")
	}
	return []float64{*parsed.X, *parsed.Y}, nil
}

// MustPlanetBiome implements a converter for a planet biome.
func MustPlanetBiome(c Converter, in *api.Planet_Biome) (gen.Biome, error) {
	if in == nil {
		return gen.Biome{}, errors.New("Planet biome is nil")
	}
	biome, err := in.AsBiome()
	if err != nil {
		return gen.Biome{}, fmt.Errorf("parse planet biome: %w", err)
	}
	return c.ConvertPlanetBiome(biome)
}

// MustPlanetBiomeName implements a converter for a planet biome name.
func MustPlanetBiomeName(c Converter, in *api.Planet_Biome) (string, error) {
	biome, err := MustPlanetBiome(c, in)
	if err != nil {
		return "", err
	}
	return biome.Name, nil
}
