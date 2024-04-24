// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: campaigns.sql

package gen

import (
	"context"
)

const campaignExists = `-- name: CampaignExists :one
SELECT EXISTS(SELECT id, type, count FROM campaigns WHERE id = $1)
`

func (q *Queries) CampaignExists(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRow(ctx, campaignExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getCampaign = `-- name: GetCampaign :one
SELECT id FROM campaigns
WHERE id = $1
`

func (q *Queries) GetCampaign(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, getCampaign, id)
	err := row.Scan(&id)
	return id, err
}

const mergeCampaign = `-- name: MergeCampaign :execrows
INSERT INTO campaigns (
    id, type, count
) VALUES (
    $1, $2, $3
)
ON CONFLICT (id) DO UPDATE
    SET type=$2, count=$3
WHERE FALSE IN (
    EXCLUDED.type=$2, EXCLUDED.count=$3
)
`

type MergeCampaignParams struct {
	ID    int32
	Type  int32
	Count int32
}

func (q *Queries) MergeCampaign(ctx context.Context, arg MergeCampaignParams) (int64, error) {
	result, err := q.db.Exec(ctx, mergeCampaign, arg.ID, arg.Type, arg.Count)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}