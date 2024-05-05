package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Snapshot)(nil)

// Snapshot implements EntityMerger
type Snapshot struct {
	gen.Snapshot
	WarSnapshot         gen.WarSnapshot
	AssignmentSnapshots []gen.AssignmentSnapshot
	PlanetSnapshots     []PlanetSnapshot
	Statistics          gen.SnapshotStatistic
}

// PlanetSnapshot wraps all snapshots relevant for a planet.
type PlanetSnapshot struct {
	gen.PlanetSnapshot
	Event      *gen.EventSnapshot
	Statistics gen.SnapshotStatistic
}

// Merge implements EntityMerger.
func (s *Snapshot) Merge(ctx context.Context, tx *gen.Queries, onMerge onMergeFunc) error {
	warSnapID, err := insertWarSnapshot(ctx, tx, s.WarSnapshot, onMerge)
	if err != nil {
		return err
	}

	assignmentSnapIDs, err := insertAssignmentSnapshots(ctx, tx, s.AssignmentSnapshots, onMerge)
	if err != nil {
		return err
	}

	planetSnapIDs, err := insertPlanetSnapshots(ctx, tx, s.PlanetSnapshots, onMerge)
	if err != nil {
		return err
	}

	statsID, err := insertSnapshotStatistics(ctx, tx, s.Statistics, onMerge)
	if err != nil {
		return err
	}

	// perform INSERT
	if _, err = tx.InsertSnapshot(ctx, gen.InsertSnapshotParams{
		WarSnapshotID:         warSnapID,
		AssignmentSnapshotIds: assignmentSnapIDs,
		CampaignIds:           s.CampaignIds,
		DispatchIds:           s.DispatchIds,
		PlanetSnapshotIds:     planetSnapIDs,
		StatisticsID:          statsID,
	}); err != nil {
		return fmt.Errorf("failed to insert snapshot: %v", err)
	}
	onMerge(gen.TableSnapshots, false, 1)
	return nil
}

func insertWarSnapshot(ctx context.Context, tx *gen.Queries, warSnap gen.WarSnapshot, onMerge onMergeFunc) (int64, error) {
	id, err := tx.InsertWarSnapshot(ctx, gen.InsertWarSnapshotParams{
		WarID:            warSnap.WarID,
		ImpactMultiplier: warSnap.ImpactMultiplier,
	})
	if err != nil {
		return -1, fmt.Errorf("failed to insert war snapshot: %v", err)
	}
	onMerge(gen.TableWarSnapshots, false, 1)
	return id, nil
}

func insertAssignmentSnapshots(ctx context.Context, tx *gen.Queries, assignmentSnaps []gen.AssignmentSnapshot, onMerge onMergeFunc) ([]int64, error) {
	ids := make([]int64, len(assignmentSnaps))
	for i, snap := range assignmentSnaps {
		id, err := tx.InsertAssignmentSnapshot(ctx, gen.InsertAssignmentSnapshotParams{
			AssignmentID: snap.AssignmentID,
			Progress:     snap.Progress,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to insert assignment snapshot: %v", err)
		}
		onMerge(gen.TableAssignmentSnapshots, false, 1)
		ids[i] = id
	}
	return ids, nil
}

func insertPlanetSnapshots(ctx context.Context, tx *gen.Queries, planetSnaps []PlanetSnapshot, onMerge onMergeFunc) ([]int64, error) {
	ids := make([]int64, len(planetSnaps))
	for i, snap := range planetSnaps {
		eventSnapID, err := insertEventSnapshot(ctx, tx, snap.Event, onMerge)
		if err != nil {
			return nil, err
		}
		statsID, err := insertSnapshotStatistics(ctx, tx, snap.Statistics, onMerge)
		if err != nil {
			return nil, err
		}

		id, err := tx.InsertPlanetSnapshot(ctx, gen.InsertPlanetSnapshotParams{
			PlanetID:           snap.PlanetID,
			Health:             snap.Health,
			CurrentOwner:       snap.CurrentOwner,
			EventSnapshotID:    eventSnapID,
			AttackingPlanetIds: snap.AttackingPlanetIds,
			RegenPerSecond:     snap.RegenPerSecond,
			StatisticsID:       statsID,
		})
		if err != nil {
			return nil, fmt.Errorf("insert planet snapshot: %w", err)
		}
		ids[i] = id
		onMerge(gen.TablePlanetSnapshots, false, 1)
	}
	return ids, nil
}

func insertEventSnapshot(ctx context.Context, tx *gen.Queries, eventSnap *gen.EventSnapshot, onMerge onMergeFunc) (*int64, error) {
	if eventSnap == nil {
		// event is optional
		return nil, nil
	}
	id, err := tx.InsertEventSnapshot(ctx, gen.InsertEventSnapshotParams{
		EventID: eventSnap.EventID,
		Health:  eventSnap.Health,
	})
	if err != nil {
		return nil, fmt.Errorf("insert event snapshot: %w", err)
	}
	onMerge(gen.TableEventSnapshots, false, 1)
	return &id, nil
}

func insertSnapshotStatistics(ctx context.Context, tx *gen.Queries, snapshotStats gen.SnapshotStatistic, onMerge onMergeFunc) (int64, error) {
	id, err := tx.InsertSnapshotStatistics(ctx, gen.InsertSnapshotStatisticsParams{
		MissionsWon:     snapshotStats.MissionsWon,
		MissionsLost:    snapshotStats.MissionsLost,
		MissionTime:     snapshotStats.MissionTime,
		TerminidKills:   snapshotStats.TerminidKills,
		AutomatonKills:  snapshotStats.AutomatonKills,
		IlluminateKills: snapshotStats.IlluminateKills,
		BulletsFired:    snapshotStats.BulletsFired,
		BulletsHit:      snapshotStats.BulletsHit,
		TimePlayed:      snapshotStats.TimePlayed,
		Deaths:          snapshotStats.Deaths,
		Revives:         snapshotStats.Revives,
		Friendlies:      snapshotStats.Friendlies,
		PlayerCount:     snapshotStats.PlayerCount,
	})
	if err != nil {
		return -1, fmt.Errorf("insert snapshot statistics: %w", err)
	}
	onMerge(gen.TableSnapshotStatistics, false, 1)
	return id, nil
}
