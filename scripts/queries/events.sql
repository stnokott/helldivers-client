-- name: GetEvent :one
SELECT id FROM data.events
WHERE id = $1;

-- name: EventExists :one
SELECT EXISTS(SELECT * FROM data.events WHERE id = $1);

-- name: MergeEvent :execrows
INSERT INTO data.events (
    id, campaign_id, type, faction, max_health, start_time, end_time
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (id) DO UPDATE
    SET campaign_id=$2, type=$3, faction=$4, max_health=$5, start_time=$6, end_time=$7
WHERE FALSE IN (
    EXCLUDED.campaign_id=$2, EXCLUDED.type=$3, EXCLUDED.faction=$4, EXCLUDED.max_health=$5, EXCLUDED.start_time=$6, EXCLUDED.end_time=$7
);
