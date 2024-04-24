CREATE TABLE IF NOT EXISTS events
(
    id integer NOT NULL UNIQUE,
    type integer NOT NULL,
    faction text NOT NULL CHECK (faction <> ''),
    max_health bigint NOT NULL CHECK (max_health >= 0),
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL CHECK (end_time > start_time),
    PRIMARY KEY (id)
);

COMMENT ON TABLE events
    IS 'Represents an ongoing event on a Planet.';

COMMENT ON COLUMN events.type
    IS 'The type of event';

COMMENT ON COLUMN events.faction
    IS 'The faction that initiated the event';

COMMENT ON COLUMN events.max_health
    IS 'The maximum health of the Event at the time of snapshot';

COMMENT ON COLUMN events.start_time
    IS 'When the event started';

COMMENT ON COLUMN events.end_time
    IS 'When the event will end (or has ended).';
