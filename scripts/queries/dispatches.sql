-- name: GetDispatch :one
SELECT id FROM data.dispatches
WHERE id = $1;

-- name: DispatchExists :one
SELECT EXISTS(SELECT * FROM data.dispatches WHERE id = $1);

-- name: MergeDispatch :execrows
INSERT INTO data.dispatches (
    id, create_time, type, message
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
    SET create_time=$2, type=$3, message=$4
WHERE FALSE IN (
    EXCLUDED.create_time=$2, EXCLUDED.type=$3, EXCLUDED.message=$4
);
