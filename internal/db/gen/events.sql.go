// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: events.sql

package gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const eventExists = `-- name: EventExists :one
SELECT EXISTS(SELECT id, campaign_id, type, faction, max_health, start_time, end_time FROM data.events WHERE id = $1)
`

func (q *Queries) EventExists(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRow(ctx, eventExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getEvent = `-- name: GetEvent :one
SELECT id FROM data.events
WHERE id = $1
`

func (q *Queries) GetEvent(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, getEvent, id)
	err := row.Scan(&id)
	return id, err
}

const mergeEvent = `-- name: MergeEvent :execrows
INSERT INTO data.events (
    id, campaign_id, type, faction, max_health, start_time, end_time
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (id) DO UPDATE
    SET campaign_id=$2, type=$3, faction=$4, max_health=$5, start_time=$6, end_time=$7
WHERE FALSE IN (
    EXCLUDED.campaign_id=$2, EXCLUDED.type=$3, EXCLUDED.faction=$4, EXCLUDED.max_health=$5, EXCLUDED.start_time=$6, EXCLUDED.end_time=$7
)
`

type MergeEventParams struct {
	ID         int32
	CampaignID int32
	Type       int32
	Faction    string
	MaxHealth  int64
	StartTime  pgtype.Timestamp
	EndTime    pgtype.Timestamp
}

func (q *Queries) MergeEvent(ctx context.Context, arg MergeEventParams) (int64, error) {
	result, err := q.db.Exec(ctx, mergeEvent,
		arg.ID,
		arg.CampaignID,
		arg.Type,
		arg.Faction,
		arg.MaxHealth,
		arg.StartTime,
		arg.EndTime,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
