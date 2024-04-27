-- name: GetAssignment :one
SELECT id FROM data.assignments
WHERE id = $1;

-- name: AssignmentExists :one
SELECT EXISTS(SELECT * FROM data.assignments WHERE id = $1);

-- name: MergeAssignment :execrows
INSERT INTO data.assignments (
    id, title, briefing, description, expiration, task_ids, reward_type, reward_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (id) DO UPDATE
    SET title=$2, briefing=$3, description=$4, expiration=$5, task_ids=$6, reward_type=$7, reward_amount=$8
WHERE FALSE IN (
    EXCLUDED.title=$2, EXCLUDED.briefing=$3, EXCLUDED.description=$4, EXCLUDED.expiration=$5, EXCLUDED.reward_type=$7, EXCLUDED.reward_amount=$8
);

-- name: InsertAssignmentTask :execrows
INSERT INTO data.assignment_tasks (
    task_type, values, value_types
) VALUES (
    $1, $2, $3
);

-- name: DeleteAssignmentTasks :exec
DELETE FROM data.assignment_tasks
WHERE id = ANY(sqlc.arg(ids)::bigint[]);
