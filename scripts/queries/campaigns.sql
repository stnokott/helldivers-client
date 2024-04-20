-- name: GetCampaign :one
SELECT id FROM campaigns
WHERE id = $1;

-- name: InsertCampaign :one
INSERT INTO campaigns (
    id, planet_id, type, count
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: UpdateCampaign :one
UPDATE campaigns
    SET planet_id=$2, type=$3, count=$4
WHERE id = $1
RETURNING id;
