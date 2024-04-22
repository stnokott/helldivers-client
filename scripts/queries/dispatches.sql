-- name: GetDispatch :one
SELECT id FROM dispatches
WHERE id = $1;

-- name: InsertDispatch :one
INSERT INTO dispatches (
    id, create_time, type, message
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: UpdateDispatch :one
UPDATE dispatches
    SET create_time=$2, type=$3, message=$4
WHERE id = $1
RETURNING id;
