DROP TRIGGER validate_planet_snapshot_refs ON planet_snapshots;
DROP FUNCTION validate_planet_snapshot_refs;

DROP TRIGGER validate_snapshot_refs ON snapshots;
DROP FUNCTION validate_snapshot_refs;

DROP TABLE IF EXISTS snapshots;

DROP TABLE IF EXISTS planet_snapshots;

DROP TABLE IF EXISTS event_snapshots;

DROP TABLE IF EXISTS assignment_snapshots;

DROP TABLE IF EXISTS war_snapshots;

DROP TABLE IF EXISTS snapshot_statistics;
