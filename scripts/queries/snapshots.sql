-- name: GetLatestSnapshot :one
SELECT * FROM data.snapshots
ORDER BY create_time desc
LIMIT 1;

-- name: InsertSnapshot :one
INSERT INTO data.snapshots (
    war_snapshot_id, assignment_snapshot_ids, campaign_ids, dispatch_ids, planet_snapshot_ids, statistics_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING create_time;

-- name: InsertWarSnapshot :one
INSERT INTO data.war_snapshots (
    war_id, impact_multiplier
) VALUES (
    $1, $2
)
RETURNING id;

-- name: InsertAssignmentSnapshot :one
INSERT INTO data.assignment_snapshots (
    assignment_id, progress
) VALUES (
    $1, $2
)
RETURNING id;

-- name: InsertPlanetSnapshot :one
INSERT INTO data.planet_snapshots (
    planet_id, health, current_owner, event_snapshot_id, attacking_planet_ids, regen_per_second, statistics_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: InsertEventSnapshot :one
INSERT INTO data.event_snapshots (
    event_id, health
) VALUES (
    $1, $2
)
RETURNING id;

-- name: InsertSnapshotStatistics :one
INSERT INTO data.snapshot_statistics (
    missions_won, missions_lost, mission_time, terminid_kills, automaton_kills, illuminate_kills, bullets_fired, bullets_hit, time_played, deaths, revives, friendlies, player_count
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING id;
