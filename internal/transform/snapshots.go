package transform

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// Snapshot converts API data into a mergable DB entity.
func Snapshot(c Converter, data APIData) (mergers []db.EntityMerger, err error) {
	snapshot, err := c.ConvertSnapshot(data)
	if err != nil {
		return nil, err
	}
	warSnap, err := c.ConvertWarSnapshot(data)
	if err != nil {
		return nil, err
	}
	if data.Assignments == nil {
		return nil, errors.New("Assignments is nil")
	}
	assignmentSnapshots, err := c.ConvertAssignmentSnapshots(*data.Assignments)
	if err != nil {
		return nil, err
	}
	if data.Planets == nil {
		return nil, errors.New("Planets is nil")
	}
	planetSnapshots, err := c.ConvertPlanetSnapshots(*data.Planets)
	if err != nil {
		return nil, err
	}
	warStats, err := c.ConvertWarStatistics(data.War.Statistics)
	if err != nil {
		return nil, err
	}

	s := &db.Snapshot{
		Snapshot:            snapshot,
		WarSnapshot:         *warSnap,
		AssignmentSnapshots: assignmentSnapshots,
		PlanetSnapshots:     planetSnapshots,
		Statistics:          *warStats,
	}

	mergers = []db.EntityMerger{s}
	return
}

// DefaultSnapshot generates a snapshot with default values for FK IDs and identity columns.
func DefaultSnapshot() gen.Snapshot {
	return gen.Snapshot{
		CreateTime:            pgtype.Timestamp{Valid: false}, // identity column
		WarSnapshotID:         -1,                             // will be filled during insert
		AssignmentSnapshotIds: nil,                            // see above
		PlanetSnapshotIds:     nil,                            // see above
		StatisticsID:          -1,                             // see above
	}
}

// DefaultPlanetSnapshot generates a planet snapshot with default values for FK IDs and identity columns.
func DefaultPlanetSnapshot() gen.PlanetSnapshot {
	return gen.PlanetSnapshot{
		ID:              -1,  // identity column
		EventSnapshotID: nil, // will be filled later from DB
		StatisticsID:    -1,  // will be filled later from DB
	}
}

// DefaultAssignmentSnapshot returns a assignment snapshot with a default identity column.
func DefaultAssignmentSnapshot() gen.AssignmentSnapshot {
	return gen.AssignmentSnapshot{
		ID: -1, // identity column
	}
}

// MustSnapshotCampaignIDs implements a converter for campaign IDs.
func MustSnapshotCampaignIDs(source *[]api.Campaign2) ([]int32, error) {
	if source == nil {
		return nil, errors.New("Campaigns is nil")
	}
	campaigns := *source
	campaignIDs := make([]int32, len(campaigns))
	for i, campaign := range campaigns {
		if campaign.Id == nil {
			return nil, errors.New("Campaign ID is nil")
		}
		campaignIDs[i] = *campaign.Id
	}
	return campaignIDs, nil
}

// MustSnapshotDispatchIDs implements a converter for dispatch IDs.
func MustSnapshotDispatchIDs(source *[]api.Dispatch) ([]int32, error) {
	if source == nil {
		return nil, errors.New("Dispatches is nil")
	}
	dispatches := *source
	dispatchIDs := make([]int32, len(dispatches))
	for i, dispatch := range dispatches {
		if dispatch.Id == nil {
			return nil, errors.New("Dispatch ID is nil")
		}
		dispatchIDs[i] = *dispatch.Id
	}
	return dispatchIDs, nil
}

// MustEventSnapshot implements a converter for event snapshots.
func MustEventSnapshot(c Converter, source *api.Planet_Event) (*gen.EventSnapshot, error) {
	if source == nil {
		// events are optional, so nil is ok
		return nil, nil
	}
	planetEvent, err := source.AsEvent()
	if err != nil {
		return nil, fmt.Errorf("parse Planet Event: %w", err)
	}
	return c.ConvertEventSnapshot(planetEvent)
}

// MustPlanetStatistics implements a converter for planet statistics.
func MustPlanetStatistics(c Converter, source *api.Planet_Statistics) (gen.SnapshotStatistic, error) {
	if source == nil {
		return gen.SnapshotStatistic{}, errors.New("got nil Planet Statistics")
	}
	parsed, err := source.AsStatistics()
	if err != nil {
		return gen.SnapshotStatistic{}, fmt.Errorf("parse Planet Statistics: %w", err)
	}
	stats, err := c.ConvertStatistics(parsed)
	if err != nil {
		return gen.SnapshotStatistic{}, err
	}
	return *stats, nil
}

// MustWarStatistics implements a converter for war statistics.
func MustWarStatistics(c Converter, source *api.War_Statistics) (*gen.SnapshotStatistic, error) {
	if source == nil {
		return nil, errors.New("got nil War Statistics")
	}
	parsed, err := source.AsStatistics()
	if err != nil {
		return nil, fmt.Errorf("parse War Statistics: %w", err)
	}
	return c.ConvertStatistics(parsed)
}
