-- name: GetCampaign :one
SELECT id FROM campaigns
WHERE id = $1;

-- name: MergeCampaign :one
INSERT INTO campaigns (
    id, planet_id, type, count
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
    SET planet_id=$2, type=$3, count=$4
RETURNING id;
