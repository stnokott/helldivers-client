-- name: GetCampaign :one
SELECT id FROM campaigns
WHERE id = $1;

-- name: CampaignExists :one
SELECT EXISTS(SELECT * FROM campaigns WHERE id = $1);

-- name: MergeCampaign :execrows
INSERT INTO campaigns (
    id, type, count
) VALUES (
    $1, $2, $3
)
ON CONFLICT (id) DO UPDATE
    SET type=$2, count=$3
WHERE FALSE IN (
    EXCLUDED.type=$2, EXCLUDED.count=$3
);
