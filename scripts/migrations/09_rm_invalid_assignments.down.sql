-- undo references
UPDATE assignments
SET task_ids = array_replace(task_ids, -1, 1)
;

-- delete dummy task
DELETE FROM assignment_tasks
WHERE id = -1
;

-- reinsert tasks from backup table
INSERT INTO assignment_tasks
OVERRIDING SYSTEM VALUE
SELECT * FROM __invalid_assignment_tasks;

-- drop backup table
DROP TABLE __invalid_assignment_tasks;
