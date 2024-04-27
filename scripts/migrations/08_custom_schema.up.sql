CREATE SCHEMA IF NOT EXISTS data;



ALTER TABLE wars SET SCHEMA data;

ALTER TABLE campaigns SET SCHEMA data;

ALTER TABLE events SET SCHEMA data;

ALTER TABLE planets SET SCHEMA data;
ALTER TABLE biomes SET SCHEMA data;
ALTER TABLE hazards SET SCHEMA data;

ALTER TABLE assignment_tasks SET SCHEMA data;
ALTER TABLE assignments SET SCHEMA data;

ALTER TABLE dispatches SET SCHEMA data;

ALTER TABLE snapshots SET SCHEMA data;
ALTER TABLE planet_snapshots SET SCHEMA data;
ALTER TABLE snapshot_statistics SET SCHEMA data;
ALTER TABLE assignment_snapshots SET SCHEMA data;
ALTER TABLE event_snapshots SET SCHEMA data;
ALTER TABLE war_snapshots SET SCHEMA data;
