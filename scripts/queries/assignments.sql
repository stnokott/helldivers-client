-- name: GetAssignment :one
SELECT id FROM assignments
WHERE id = $1;

-- name: AssignmentExists :one
SELECT EXISTS(SELECT * FROM assignments WHERE id = $1);

-- name: MergeAssignment :execrows
INSERT INTO assignments (
    id, title, briefing, description, expiration, task_ids, reward_type, reward_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (id) DO UPDATE
    SET title=$2, briefing=$3, description=$4, expiration=$5, task_ids=$6, reward_type=$7, reward_amount=$8
;

-- name: InsertAssignmentTask :one
INSERT INTO assignment_tasks (
    task_type, values, value_types
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: DeleteAssignmentTasks :exec
DELETE FROM assignment_tasks
WHERE id IN (
    SELECT assignment_tasks.id
    FROM assignments
    JOIN assignment_tasks
        ON assignment_tasks.id = ANY(task_ids)
    WHERE assignments.id = sqlc.arg(assignment_id)
);
