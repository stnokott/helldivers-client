CREATE TABLE IF NOT EXISTS dispatches
(
    id integer NOT NULL UNIQUE,
    create_time timestamp without time zone NOT NULL,
    type integer NOT NULL,
    message text NOT NULL CHECK (message <> ''),
    PRIMARY KEY (id)
);

COMMENT ON TABLE dispatches
    IS 'Represents a message from high command to the players, usually updates on the status of the war effort.';

COMMENT ON COLUMN dispatches.id
    IS 'The unique identifier of this dispatch';

COMMENT ON COLUMN dispatches.create_time
    IS 'When the dispatch was published';

COMMENT ON COLUMN dispatches.type
    IS 'The type of dispatch, purpose unknown';

COMMENT ON COLUMN dispatches.message
    IS 'The message this dispatch represents';

