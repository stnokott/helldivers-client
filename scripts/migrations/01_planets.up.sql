CREATE TABLE IF NOT EXISTS biomes
(
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    PRIMARY KEY (name)
);

COMMENT ON TABLE biomes
    IS 'Represents information about a biomes of a planet.';



CREATE TABLE IF NOT EXISTS hazards
(
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    PRIMARY KEY (name)
);

COMMENT ON TABLE hazards
    IS 'Describes an environmental hazards that can be present on a planet';



CREATE TABLE IF NOT EXISTS planets
(
    id integer NOT NULL UNIQUE,
    name text NOT NULL UNIQUE CONSTRAINT name_not_empty CHECK (name <> ''),
    sector text NOT NULL CONSTRAINT sector_not_empty CHECK (sector <> ''),
    position double precision[2] NOT NULL CONSTRAINT position_exactly_two_values CHECK (array_length(position, 1) = 2),
    waypoint_ids integer[] NOT NULL,
    disabled boolean NOT NULL,
    biome_name text NOT NULL REFERENCES biomes,
    hazard_names text[] NOT NULL, -- reference check is performed in trigger function
    max_health bigint NOT NULL CONSTRAINT max_health_not_negative CHECK (max_health > 0),
    initial_owner text NOT NULL CONSTRAINT initial_owner_not_empty CHECK (initial_owner <> ''),
    PRIMARY KEY (id)
);

CREATE OR REPLACE FUNCTION validate_planet_hazard_refs() RETURNS TRIGGER AS $validate_planet_hazard_refs$
	DECLARE
		new_hazard_name text;
    BEGIN
		-- check hazard refs
		FOREACH new_hazard_name IN ARRAY NEW.hazard_names LOOP
			IF NOT EXISTS (SELECT 1 FROM hazards WHERE name = new_hazard_name) THEN
				RAISE EXCEPTION 'planet "%" has non-existent hazard "%"', NEW.name, new_hazard_name;
			END IF;
		END LOOP;

        RETURN NEW;
    END;
$validate_planet_hazard_refs$ LANGUAGE plpgsql;

CREATE TRIGGER validate_planet_hazard_refs BEFORE INSERT OR UPDATE ON planets
    FOR EACH ROW EXECUTE FUNCTION validate_planet_hazard_refs();

COMMENT ON TABLE planets
    IS 'Represents information of a planet from the "WarInfo" endpoint returned by ArrowHead''s API';

COMMENT ON COLUMN planets.id
    IS 'The unique identifier ArrowHead assigned to this planet';

COMMENT ON COLUMN planets.name
    IS 'The name of the planet, as shown in game';

COMMENT ON COLUMN planets.sector
    IS 'The name of the sector the planet is in, as shown in game';

COMMENT ON COLUMN planets.position
    IS 'The coordinates of this planet on the galactic war map in format [X, Y]';

COMMENT ON COLUMN planets.waypoint_ids
    IS 'List of indexes of all the planets to which this planet is connected';

COMMENT ON COLUMN planets.disabled
    IS 'Whether or not this planet is disabled, as assigned by ArrowHead';

COMMENT ON COLUMN planets.biome_name
    IS 'The biomes this planet has.';

COMMENT ON COLUMN planets.hazard_names
    IS 'All hazardss that are applicable to this planet.';

COMMENT ON COLUMN planets.max_health
    IS 'The maximum health pool of this planet';

COMMENT ON COLUMN planets.initial_owner
    IS 'The faction that originally owned the plane';
