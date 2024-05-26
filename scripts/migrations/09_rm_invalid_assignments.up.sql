/*
This migration is required due to a bug (fixed in this commit) which didn't set task IDs correctly
when merging assignments (all task IDs were set to 1 instead of the correct ones).

Unfortunately, these invalid relations can't be fixed now since it's impossible to
deduce them.

This migration thus replaces all existing task ID FKs with references to a dummy task to
indicate that these task references are invalid.
*/

-- create backup of existing tasks
CREATE TABLE __invalid_assignment_tasks AS (SELECT * FROM assignment_tasks);

-- truncate table since it contains invalid tasks
TRUNCATE TABLE assignment_tasks;

-- insert dummy task
INSERT INTO assignment_tasks (
    id, task_type, values, value_types
) OVERRIDING SYSTEM VALUE
VALUES (
    -1, 0, '{}', '{}'
);

-- update references
UPDATE assignments
SET task_ids = array_replace(task_ids, 1, -1)
;
