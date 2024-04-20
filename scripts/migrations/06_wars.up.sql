CREATE TABLE IF NOT EXISTS wars
(
    id integer NOT NULL UNIQUE,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL CHECK (end_time > start_time),
    factions text[] NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE wars
    IS 'Represents the global information of the ongoing war';

COMMENT ON COLUMN wars.start_time
    IS 'When this war was started';

COMMENT ON COLUMN wars.end_time
    IS 'When this war will end (or has ended)';

COMMENT ON COLUMN wars.factions
    IS 'A list of factions currently involved in the war';

