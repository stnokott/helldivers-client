// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: snapshots.sql

package gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getLatestSnapshot = `-- name: GetLatestSnapshot :one
SELECT create_time, war_snapshot_id, assignment_snapshot_ids, campaign_ids, dispatch_ids, planet_snapshot_ids, statistics_id FROM snapshots
ORDER BY create_time desc
LIMIT 1
`

func (q *Queries) GetLatestSnapshot(ctx context.Context) (Snapshot, error) {
	row := q.db.QueryRow(ctx, getLatestSnapshot)
	var i Snapshot
	err := row.Scan(
		&i.CreateTime,
		&i.WarSnapshotID,
		&i.AssignmentSnapshotIds,
		&i.CampaignIds,
		&i.DispatchIds,
		&i.PlanetSnapshotIds,
		&i.StatisticsID,
	)
	return i, err
}

const insertAssignmentSnapshot = `-- name: InsertAssignmentSnapshot :one
INSERT INTO assignment_snapshots (
    assignment_id, progress
) VALUES (
    $1, $2
)
RETURNING id
`

type InsertAssignmentSnapshotParams struct {
	AssignmentID int64
	Progress     []int32
}

func (q *Queries) InsertAssignmentSnapshot(ctx context.Context, arg InsertAssignmentSnapshotParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertAssignmentSnapshot, arg.AssignmentID, arg.Progress)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertEventSnapshot = `-- name: InsertEventSnapshot :one
INSERT INTO event_snapshots (
    event_id, health
) VALUES (
    $1, $2
)
RETURNING id
`

type InsertEventSnapshotParams struct {
	EventID int32
	Health  int64
}

func (q *Queries) InsertEventSnapshot(ctx context.Context, arg InsertEventSnapshotParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertEventSnapshot, arg.EventID, arg.Health)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertPlanetSnapshot = `-- name: InsertPlanetSnapshot :one
INSERT INTO planet_snapshots (
    planet_id, health, current_owner, event_snapshot_id, attacking_planet_ids, regen_per_second, statistics_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id
`

type InsertPlanetSnapshotParams struct {
	PlanetID           int32
	Health             int64
	CurrentOwner       string
	EventSnapshotID    *int64
	AttackingPlanetIds []int32
	RegenPerSecond     float64
	StatisticsID       int64
}

func (q *Queries) InsertPlanetSnapshot(ctx context.Context, arg InsertPlanetSnapshotParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertPlanetSnapshot,
		arg.PlanetID,
		arg.Health,
		arg.CurrentOwner,
		arg.EventSnapshotID,
		arg.AttackingPlanetIds,
		arg.RegenPerSecond,
		arg.StatisticsID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertSnapshot = `-- name: InsertSnapshot :one
INSERT INTO snapshots (
    war_snapshot_id, assignment_snapshot_ids, campaign_ids, dispatch_ids, planet_snapshot_ids, statistics_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING create_time
`

type InsertSnapshotParams struct {
	WarSnapshotID         int64
	AssignmentSnapshotIds []int64
	CampaignIds           []int32
	DispatchIds           []int32
	PlanetSnapshotIds     []int64
	StatisticsID          int64
}

func (q *Queries) InsertSnapshot(ctx context.Context, arg InsertSnapshotParams) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, insertSnapshot,
		arg.WarSnapshotID,
		arg.AssignmentSnapshotIds,
		arg.CampaignIds,
		arg.DispatchIds,
		arg.PlanetSnapshotIds,
		arg.StatisticsID,
	)
	var create_time pgtype.Timestamp
	err := row.Scan(&create_time)
	return create_time, err
}

const insertSnapshotStatistics = `-- name: InsertSnapshotStatistics :one
INSERT INTO snapshot_statistics (
    missions_won, missions_lost, mission_time, terminid_kills, automaton_kills, illuminate_kills, bullets_fired, bullets_hit, time_played, deaths, revives, friendlies, player_count
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING id
`

type InsertSnapshotStatisticsParams struct {
	MissionsWon     pgtype.Numeric
	MissionsLost    pgtype.Numeric
	MissionTime     pgtype.Numeric
	TerminidKills   pgtype.Numeric
	AutomatonKills  pgtype.Numeric
	IlluminateKills pgtype.Numeric
	BulletsFired    pgtype.Numeric
	BulletsHit      pgtype.Numeric
	TimePlayed      pgtype.Numeric
	Deaths          pgtype.Numeric
	Revives         pgtype.Numeric
	Friendlies      pgtype.Numeric
	PlayerCount     pgtype.Numeric
}

func (q *Queries) InsertSnapshotStatistics(ctx context.Context, arg InsertSnapshotStatisticsParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertSnapshotStatistics,
		arg.MissionsWon,
		arg.MissionsLost,
		arg.MissionTime,
		arg.TerminidKills,
		arg.AutomatonKills,
		arg.IlluminateKills,
		arg.BulletsFired,
		arg.BulletsHit,
		arg.TimePlayed,
		arg.Deaths,
		arg.Revives,
		arg.Friendlies,
		arg.PlayerCount,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertWarSnapshot = `-- name: InsertWarSnapshot :one
INSERT INTO war_snapshots (
    war_id, impact_multiplier
) VALUES (
    $1, $2
)
RETURNING id
`

type InsertWarSnapshotParams struct {
	WarID            int32
	ImpactMultiplier float64
}

func (q *Queries) InsertWarSnapshot(ctx context.Context, arg InsertWarSnapshotParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertWarSnapshot, arg.WarID, arg.ImpactMultiplier)
	var id int64
	err := row.Scan(&id)
	return id, err
}