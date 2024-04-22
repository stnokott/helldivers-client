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

-- name: InsertAssignmentTask :one
INSERT INTO assignment_tasks (
    type, values, value_types
) VALUES (
    $1, $2, $3
)
RETURNING id;
