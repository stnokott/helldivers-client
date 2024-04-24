-- name: GetCampaign :one
SELECT id FROM campaigns
WHERE id = $1;

-- name: CampaignExists :one
SELECT EXISTS(SELECT * FROM campaigns WHERE id = $1);

-- name: MergeCampaign :execrows
INSERT INTO campaigns (
    id, planet_id, type, count
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
    SET planet_id=$2, type=$3, count=$4
WHERE FALSE IN (
    EXCLUDED.planet_id=$2, EXCLUDED.type=$3, EXCLUDED.count=$4
);
