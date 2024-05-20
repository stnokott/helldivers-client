ALTER TABLE assignment_snapshots
ALTER COLUMN progress TYPE numeric[];

ALTER TABLE campaigns
ALTER COLUMN count TYPE numeric;

ALTER TABLE assignments
ALTER COLUMN reward_amount TYPE numeric;

ALTER TABLE assignment_tasks
ALTER COLUMN values TYPE numeric[],
ALTER COLUMN value_types TYPE numeric[];
