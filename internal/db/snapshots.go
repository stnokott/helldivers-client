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

type PlanetSnapshot struct {
	gen.PlanetSnapshot
	Event      *gen.EventSnapshot
	Statistics gen.SnapshotStatistic
}

func (s *Snapshot) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	warSnapID, err := insertWarSnapshot(ctx, tx, s.WarSnapshot, stats)
	if err != nil {
		return err
	}

	assignmentSnapIDs, err := insertAssignmentSnapshots(ctx, tx, s.AssignmentSnapshots, stats)
	if err != nil {
		return err
	}

	planetSnapIDs, err := insertPlanetSnapshots(ctx, tx, s.PlanetSnapshots, stats)
	if err != nil {
		return err
	}

	statsID, err := insertSnapshotStatistics(ctx, tx, s.Statistics, stats)
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
	stats.Incr("Snapshots", false, 1)
	return nil
}

func insertWarSnapshot(ctx context.Context, tx *gen.Queries, warSnap gen.WarSnapshot, stats tableMergeStats) (int64, error) {
	id, err := tx.InsertWarSnapshot(ctx, gen.InsertWarSnapshotParams{
		WarID:            warSnap.WarID,
		ImpactMultiplier: warSnap.ImpactMultiplier,
	})
	if err != nil {
		return -1, fmt.Errorf("failed to insert war snapshot: %v", err)
	}
	stats.Incr("War Snapshots", false, 1)
	return id, nil
}

func insertAssignmentSnapshots(ctx context.Context, tx *gen.Queries, assignmentSnaps []gen.AssignmentSnapshot, stats tableMergeStats) ([]int64, error) {
	ids := make([]int64, len(assignmentSnaps))
	for i, snap := range assignmentSnaps {
		id, err := tx.InsertAssignmentSnapshot(ctx, gen.InsertAssignmentSnapshotParams{
			AssignmentID: snap.AssignmentID,
			Progress:     snap.Progress,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to insert assignment snapshot: %v", err)
		}
		stats.Incr("Assignment Snapshots", false, 1)
		ids[i] = id
	}
	return ids, nil
}

func insertPlanetSnapshots(ctx context.Context, tx *gen.Queries, planetSnaps []PlanetSnapshot, stats tableMergeStats) ([]int64, error) {
	ids := make([]int64, len(planetSnaps))
	for i, snap := range planetSnaps {
		eventSnapID, err := insertEventSnapshot(ctx, tx, snap.Event, stats)
		if err != nil {
			return nil, err
		}
		statsID, err := insertSnapshotStatistics(ctx, tx, snap.Statistics, stats)
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
			return nil, fmt.Errorf("failed to insert planet snapshot: %w", err)
		}
		ids[i] = id
		stats.Incr("Planet Snapshots", false, 1)
	}
	return ids, nil
}

func insertEventSnapshot(ctx context.Context, tx *gen.Queries, eventSnap *gen.EventSnapshot, stats tableMergeStats) (*int64, error) {
	if eventSnap == nil {
		// event is optional
		return nil, nil
	}
	id, err := tx.InsertEventSnapshot(ctx, gen.InsertEventSnapshotParams{
		EventID: eventSnap.EventID,
		Health:  eventSnap.Health,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert event snapshot: %w", err)
	}
	stats.Incr("Event Snapshots", false, 1)
	return &id, nil
}

func insertSnapshotStatistics(ctx context.Context, tx *gen.Queries, snapshotStats gen.SnapshotStatistic, stats tableMergeStats) (int64, error) {
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
		return -1, fmt.Errorf("failed to insert snapshot statistics: %w", err)
	}
	stats.Incr("Snapshots Statistics", false, 1)
	return id, nil
}