ALTER TABLE data.war_snapshots SET SCHEMA public;
ALTER TABLE data.event_snapshots SET SCHEMA public;
ALTER TABLE data.assignment_snapshots SET SCHEMA public;
ALTER TABLE data.snapshot_statistics SET SCHEMA public;
ALTER TABLE data.planet_snapshots SET SCHEMA public;
ALTER TABLE data.snapshots SET SCHEMA public;

ALTER TABLE data.dispatches SET SCHEMA public;

ALTER TABLE data.assignments SET SCHEMA public;
ALTER TABLE data.assignment_tasks SET SCHEMA public;

ALTER TABLE data.hazards SET SCHEMA public;
ALTER TABLE data.biomes SET SCHEMA public;
ALTER TABLE data.planets SET SCHEMA public;

ALTER TABLE data.events SET SCHEMA public;

ALTER TABLE data.campaigns SET SCHEMA public;

ALTER TABLE data.wars SET SCHEMA public;



DROP SCHEMA IF EXISTS data;
