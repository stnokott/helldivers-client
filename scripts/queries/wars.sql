-- name: GetWar :one
SELECT id FROM wars
WHERE id = $1;

-- name: MergeWar :one
INSERT INTO wars (
    id, start_time, end_time, factions
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
    SET start_time=$2, end_time=$3, factions=$4
RETURNING id;
