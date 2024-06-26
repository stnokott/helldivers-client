CREATE TABLE IF NOT EXISTS assignment_tasks
(
    id bigint NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
    task_type integer NOT NULL,
    values integer[] NOT NULL,
    value_types integer[] NOT NULL,
    CONSTRAINT equal_value_lengths CHECK (array_length(values, 1) = array_length(value_types, 1)),
    PRIMARY KEY (id)
);

COMMENT ON TABLE assignment_tasks
    IS 'Represents a task in an Assignment that needs to be completed to finish the assignment';

COMMENT ON COLUMN assignment_tasks.id
    IS 'Auto-generated by sequence';

COMMENT ON COLUMN assignment_tasks.task_type
    IS 'The type of task this represents';

COMMENT ON COLUMN assignment_tasks.values
    IS 'A list of numbers, purpose unknown';

COMMENT ON COLUMN assignment_tasks.value_types
    IS 'A list of numbers, purpose unknown';



CREATE TABLE IF NOT EXISTS assignments
(
    id bigint NOT NULL UNIQUE,
    title text NOT NULL CONSTRAINT title_not_empty CHECK (title <> ''),
    briefing text NOT NULL CONSTRAINT briefing_not_empty CHECK (briefing <> ''),
    description text NOT NULL CONSTRAINT description_not_empty CHECK (description <> ''),
    expiration timestamp without time zone NOT NULL CONSTRAINT expiration_not_default CHECK (expiration > '1900-01-01'),
    task_ids bigint[] NOT NULL,
    reward_type integer NOT NULL,
    reward_amount integer NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE assignments
    IS 'Represents an assignment given by Super Earth to the community. This is also known as "Major Order"s in the game';

COMMENT ON COLUMN assignments.title
    IS 'The title of the assignment';

COMMENT ON COLUMN assignments.briefing
    IS 'A long form description of the assignment, usually contains context';

COMMENT ON COLUMN assignments.description
    IS 'A very short summary of the description';

COMMENT ON COLUMN assignments.expiration
    IS 'The date when the assignment will expire.';

COMMENT ON COLUMN assignments.task_ids
    IS 'A list of tasks that need to be completed for this assignment';

COMMENT ON COLUMN assignments.reward_type
    IS 'The type of reward (medals, super credits, ...)';

COMMENT ON COLUMN assignments.reward_amount
    IS 'The amount of Type that will be awarded';
