-- name: GetWar :one
SELECT id FROM wars
WHERE id = $1;

-- name: WarExists :one
SELECT EXISTS(SELECT * FROM wars WHERE id = $1);

-- name: MergeWar :execrows
INSERT INTO wars (
    id, start_time, end_time, factions
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
    SET start_time=$2, end_time=$3, factions=$4
WHERE FALSE IN (
    EXCLUDED.start_time=$2, EXCLUDED.end_time=$3, EXCLUDED.factions=$4
);
