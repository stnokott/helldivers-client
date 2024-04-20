-- name: GetAssignment :one
SELECT id FROM assignments
WHERE id = $1;

-- name: InsertAssignment :one
INSERT INTO assignments (
    id, title, briefing, description, expiration, progress, task_ids, reward_type, reward_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id;

-- name: UpdateAssignment :one
UPDATE assignments
    SET title=$2, briefing=$3, description=$4, expiration=$5, progress=$6, task_ids=$7, reward_type=$8, reward_amount=$9
WHERE id = $1
RETURNING id;

-- name: GetAssignmentTask :one
SELECT id FROM assignment_tasks
WHERE id = $1;

-- name: InsertAssignmentTask :one
INSERT INTO assignment_tasks (
    type, values, value_types
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: UpdateAssignmentTask :one
UPDATE assignment_tasks
    SET type=$2, values=$3, value_types=$4
WHERE id = $1
RETURNING id;
