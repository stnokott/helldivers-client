package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Planet)(nil)

// Planet implements EntityMerger
type Planet struct {
	gen.Planet
	Biome   gen.Biome
	Hazards []gen.Hazard
}

func (p *Planet) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	biomeName, err := mergeBiome(ctx, tx, p.Biome, stats)
	if err != nil {
		return err
	}
	p.BiomeName = biomeName

	hazardNames, err := mergeHazards(ctx, tx, p.Hazards, stats)
	if err != nil {
		return err
	}
	p.HazardNames = hazardNames

	exists, err := tx.PlanetExists(ctx, p.ID)
	if err != nil {
		return fmt.Errorf("failed to check if planet '%s' exists: %v", p.Name, err)
	}
	rows, err := tx.MergePlanet(ctx, gen.MergePlanetParams(p.Planet))
	if err != nil {
		return fmt.Errorf("failed to merge planet '%s': %v", p.Name, err)
	}
	stats.Incr("Planets", exists, rows)
	return nil
}

func mergeBiome(ctx context.Context, tx *gen.Queries, biome gen.Biome, stats tableMergeStats) (string, error) {
	exists, err := tx.BiomeExists(ctx, biome.Name)
	if err != nil {
		return "", fmt.Errorf("failed to check if biome '%s' exists: %v", biome.Name, err)
	}

	rows, err := tx.MergeBiome(ctx, gen.MergeBiomeParams(biome))
	if err != nil {
		return "", fmt.Errorf("failed to merge biome '%s': %v", biome.Name, err)
	}
	stats.Incr("Biomes", exists, rows)
	return biome.Name, nil
}

func mergeHazards(ctx context.Context, tx *gen.Queries, hazards []gen.Hazard, stats tableMergeStats) ([]string, error) {
	hazardNames := make([]string, len(hazards))
	for i, hazard := range hazards {
		exists, err := tx.HazardExists(ctx, hazard.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check if hazard '%s' exists: %v", hazard.Name, err)
		}

		rows, err := tx.MergeHazard(ctx, gen.MergeHazardParams(hazard))
		if err != nil {
			return nil, fmt.Errorf("failed to merge hazard '%s': %v", hazard.Name, err)
		}
		hazardNames[i] = hazard.Name
		stats.Incr("Hazards", exists, rows)
	}
	return hazardNames, nil
}
