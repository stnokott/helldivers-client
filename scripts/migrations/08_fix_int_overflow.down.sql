ALTER TABLE assignment_tasks
ALTER COLUMN values TYPE integer[],
ALTER COLUMN value_types TYPE integer[];

ALTER TABLE assignments
ALTER COLUMN reward_amount TYPE integer;

ALTER TABLE campaigns
ALTER COLUMN count TYPE integer;

ALTER TABLE assignment_snapshots
ALTER COLUMN progress TYPE integer[];
