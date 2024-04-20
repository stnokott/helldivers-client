CREATE TABLE IF NOT EXISTS campaigns
(
    id integer NOT NULL UNIQUE,
    planet_id integer NOT NULL,
    type integer NOT NULL,
    count integer NOT NULL CHECK (count >= 0),
    PRIMARY KEY (id)
);

COMMENT ON TABLE campaigns
    IS 'Represents an ongoing campaign on a planet';

COMMENT ON COLUMN campaigns.id
    IS 'The unique identifier of this campaign';

COMMENT ON COLUMN campaigns.planet_id
    IS 'The planet on which this campaign is being fought';

COMMENT ON COLUMN campaigns.type
    IS 'The type of campaign, this should be mapped onto an enum';

COMMENT ON COLUMN campaigns.count
    IS 'Indicates how many campaigns have already been fought on this Planet';

