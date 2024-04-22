-- name: GetEvent :one
SELECT id FROM events
WHERE id = $1;

-- name: MergeEvent :one
INSERT INTO events (
    id, type, faction, max_health, start_time, end_time
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (id) DO UPDATE
    SET type=$2, faction=$3, max_health=$4, start_time=$5, end_time=$6
RETURNING id;
