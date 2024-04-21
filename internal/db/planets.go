package db

import (
	"context"
	"fmt"
	"log"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Planet)(nil)

// Planet implements EntityMerger
type Planet struct {
	gen.Planet
	Biome   Biome
	Hazards []Hazard
}

type Biome gen.Biome

type Hazard gen.Hazard

func (p *Planet) Merge(ctx context.Context, tx *gen.Queries, stats *MergeStats, logger *log.Logger) error {
	logger.Printf("** merging planet '%s'", p.Name)
	biomeName, err := mergeBiome(ctx, tx, p.Biome, logger)
	if err != nil {
		return err
	}
	p.BiomeName = biomeName

	hazardNames, err := mergeHazards(ctx, tx, p.Hazards, logger)
	if err != nil {
		return err
	}
	p.HazardNames = hazardNames

	id, err := tx.GetPlanet(ctx, p.ID)
	exists, err := entityExistsByPK(id, err, p.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing planet: %v", err)
	}
	if exists {
		// perform UPDATE
		if _, err = tx.UpdatePlanet(ctx, gen.UpdatePlanetParams(p.Planet)); err != nil {
			return fmt.Errorf("failed to update planet ('%s'): %v", p.Name, err)
		}
		stats.Updates++
	} else {
		// perform INSERT
		if _, err = tx.InsertPlanet(ctx, gen.InsertPlanetParams(p.Planet)); err != nil {
			return fmt.Errorf("failed to insert planet ('%s'): %v", p.Name, err)
		}
		stats.Inserts++
	}
	return nil
}

func mergeBiome(ctx context.Context, tx *gen.Queries, biome Biome, logger *log.Logger) (string, error) {
	id, err := tx.GetBiome(ctx, biome.Name)
	exists, err := entityExistsByPK(id, err, biome.Name)
	var biomeName string
	if exists {
		// perform UPDATE
		biomeName, err = tx.UpdateBiome(ctx, gen.UpdateBiomeParams(biome))
		if err != nil {
			return "", fmt.Errorf("failed to update biome ('%s'): %v", biome.Name, err)
		}
	} else {
		// perform INSERT
		biomeName, err = tx.InsertBiome(ctx, gen.InsertBiomeParams(biome))
		if err != nil {
			return "", fmt.Errorf("failed to insert biome ('%s'): %v", biome.Name, err)
		}
	}
	return biomeName, nil
}

func mergeHazards(ctx context.Context, tx *gen.Queries, hazards []Hazard, logger *log.Logger) ([]string, error) {
	hazardNames := make([]string, len(hazards))
	for i, hazard := range hazards {
		id, err := tx.GetHazard(ctx, hazard.Name)
		exists, err := entityExistsByPK(id, err, hazard.Name)
		var hazardName string
		if exists {
			// perform UPDATE
			hazardName, err = tx.UpdateHazard(ctx, gen.UpdateHazardParams(hazard))
			if err != nil {
				return nil, fmt.Errorf("failed to update hazard ('%s'): %v", hazard.Name, err)
			}
		} else {
			// perform INSERT
			hazardName, err = tx.InsertHazard(ctx, gen.InsertHazardParams(hazard))
			if err != nil {
				return nil, fmt.Errorf("failed to insert hazard ('%s'): %v", hazard.Name, err)
			}
		}
		hazardNames[i] = hazardName
	}
	return hazardNames, nil
}
