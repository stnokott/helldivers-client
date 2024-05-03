package transform

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

//go:generate go run github.com/jmattheis/goverter/cmd/goverter@v1.4.0 gen github.com/stnokott/helldivers-client/internal/transform

// Converter converts API structs into mergable DB structs.
//
// goverter:converter
// goverter:wrapErrors
//
// goverter:extend Must.*
//
// goverter:output:package github.com/stnokott/helldivers-client/internal/transform
// goverter:output:file ./generated.go
type Converter interface {
	ConvertAssignment(source api.Assignment2) (*db.Assignment, error)
	// goverter:map Id ID
	// goverter:ignore TaskIds
	// goverter:map Reward RewardType | parseAssignmentRewardType
	// goverter:map Reward RewardAmount | parseAssignmentRewardAmount
	ConvertSingleAssignment(source api.Assignment2) (*gen.Assignment, error)
	// goverter:ignore ID
	// goverter:map Type TaskType
	ConvertAssignmentTask(source api.Task2) (gen.AssignmentTask, error)
	ConvertAssignmentTasks(source []api.Task2) ([]gen.AssignmentTask, error)

	// goverter:map Id ID
	ConvertCampaign(source api.Campaign2) (*db.Campaign, error)

	// goverter:map Id ID
	// goverter:map Published CreateTime
	ConvertDispatch(source api.Dispatch) (*db.Dispatch, error)

	// goverter:map Id ID
	// goverter:map CampaignId CampaignID
	// goverter:map EventType Type
	ConvertEvent(source api.Event) (*db.Event, error)

	// goverter:map . Planet
	ConvertPlanet(source api.Planet) (*db.Planet, error)
	ConvertPlanetBiome(source api.Biome) (gen.Biome, error)
	ConvertPlanetHazard(source api.Hazard) (gen.Hazard, error)
	// goverter:map Index ID
	// goverter:map Waypoints WaypointIds
	// goverter:map Biome BiomeName
	// goverter:map Hazards HazardNames
	ConvertSinglePlanet(source api.Planet) (gen.Planet, error)

	// goverter:map WarID ID
	// goverter:autoMap War
	// goverter:map War.Started StartTime
	// goverter:map War.Ended EndTime
	ConvertWar(source APIData) (*db.War, error)

	// goverter:default DefaultSnapshot
	// goverter:ignore CreateTime
	// goverter:ignore WarSnapshotID
	// goverter:ignore AssignmentSnapshotIds
	// goverter:ignore PlanetSnapshotIds
	// goverter:ignore StatisticsID
	// goverter:map Campaigns CampaignIds
	// goverter:map Dispatches DispatchIds
	ConvertSnapshot(source APIData) (gen.Snapshot, error)
	// goverter:ignore ID
	// goverter:autoMap War
	ConvertWarSnapshot(source APIData) (*gen.WarSnapshot, error)
	// goverter:default DefaultAssignmentSnapshot
	// goverter:ignore ID
	// goverter:map Id AssignmentID
	ConvertAssignmentSnapshot(source api.Assignment2) (gen.AssignmentSnapshot, error)
	ConvertAssignmentSnapshots(source []api.Assignment2) ([]gen.AssignmentSnapshot, error)
	// goverter:default DefaultPlanetSnapshot
	// goverter:ignore ID
	// goverter:map Index PlanetID
	// goverter:ignore EventSnapshotID
	// goverter:ignore StatisticsID
	// goverter:map Attacking AttackingPlanetIds
	ConvertPlanetSnapshotOnly(source api.Planet) (gen.PlanetSnapshot, error)
	ConvertPlanetSnapshotsOnly(source []api.Planet) ([]gen.PlanetSnapshot, error)
	// goverter:map . PlanetSnapshot
	ConvertPlanetSnapshot(source api.Planet) (db.PlanetSnapshot, error)
	ConvertPlanetSnapshots(source []api.Planet) ([]db.PlanetSnapshot, error)
	// goverter:ignore ID
	// goverter:map Id EventID
	ConvertEventSnapshot(source api.Event) (*gen.EventSnapshot, error)
	// goverter:ignore ID
	ConvertStatistics(source api.Statistics) (*gen.SnapshotStatistic, error)
	ConvertWarStatistics(source *api.War_Statistics) (*gen.SnapshotStatistic, error)
}

// MustBool dereferences a boolean or returns an error if nil.
func MustBool(ptr *bool) (bool, error) {
	return mustPtr(ptr)
}

// MustInt32Ptr dereferences an int32 or returns an error if nil.
func MustInt32Ptr(ptr *int32) (int32, error) {
	return mustPtr(ptr)
}

// MustInt64Ptr dereferences an int64 or returns an error if nil.
func MustInt64Ptr(ptr *int64) (int64, error) {
	return mustPtr(ptr)
}

// MustFloat64Ptr dereferences a float64 or returns an error if nil.
func MustFloat64Ptr(ptr *float64) (float64, error) {
	return mustPtr(ptr)
}

// MustString dereferences a string or returns an error if nil.
func MustString(ptr *string) (string, error) {
	return mustPtr(ptr)
}

// MustInt32Slice dereferences an int32 or returns an error if nil.
func MustInt32Slice(ptr *[]int32) ([]int32, error) {
	return mustPtr(ptr)
}

// MustStringSlice dereferences a string slice or returns an error if nil.
func MustStringSlice(ptr *[]string) ([]string, error) {
	return mustPtr(ptr)
}

// MustNumeric converts a uint64 into a pgx-compatible type or an error if nil.
func MustNumeric(ptr *uint64) (pgtype.Numeric, error) {
	x, err := mustPtr(ptr)
	if err != nil {
		return pgtype.Numeric{}, err
	}
	return db.PGUint64(x), nil
}

// MustTimestamp converts a timestamp into a pgx-compatible type or an error if nil.
func MustTimestamp(ptr *time.Time) (pgtype.Timestamp, error) {
	t, err := mustPtr(ptr)
	if err != nil {
		return pgtype.Timestamp{}, err
	}
	return db.PGTimestamp(t), nil
}

// nolint: ireturn
func mustPtr[T any](ptr *T) (T, error) {
	if ptr == nil {
		var zero T
		return zero, fmt.Errorf("%T is nil", ptr)
	}
	return *ptr, nil
}
