// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: wars.sql

package gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getWar = `-- name: GetWar :one
SELECT id FROM wars
WHERE id = $1
`

func (q *Queries) GetWar(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, getWar, id)
	err := row.Scan(&id)
	return id, err
}

const mergeWar = `-- name: MergeWar :one
INSERT INTO wars (
    id, start_time, end_time, factions
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
    SET start_time=$2, end_time=$3, factions=$4
RETURNING id
`

type MergeWarParams struct {
	ID        int32
	StartTime pgtype.Timestamp
	EndTime   pgtype.Timestamp
	Factions  []string
}

func (q *Queries) MergeWar(ctx context.Context, arg MergeWarParams) (int32, error) {
	row := q.db.QueryRow(ctx, mergeWar,
		arg.ID,
		arg.StartTime,
		arg.EndTime,
		arg.Factions,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}
