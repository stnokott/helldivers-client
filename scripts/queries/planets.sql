-- TODO: check if INSERT ... ON CONFLICT DO UPDATE is possible

-- name: GetPlanet :one
SELECT id FROM planets
WHERE id = $1;

-- name: InsertPlanet :one
INSERT INTO planets (
    id, name, sector, position, waypoint_ids, disabled, biome_name, hazard_names, max_health, initial_owner
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id;

-- name: UpdatePlanet :one
UPDATE planets
    SET name=$2, sector=$3, position=$4, waypoint_ids=$5, disabled=$6, biome_name=$7, hazard_names=$8, max_health=$9, initial_owner=$10
WHERE id = $1
RETURNING id;

-- name: GetBiome :one
SELECT name FROM biomes
WHERE name = $1;

-- name: InsertBiome :one
INSERT INTO biomes (
    name, description
) VALUES (
    $1, $2
)
RETURNING name;

-- name: UpdateBiome :one
UPDATE biomes
    SET description=$2
WHERE name=$1
RETURNING name;

-- name: GetHazard :one
SELECT name FROM hazards
WHERE name = $1;

-- name: InsertHazard :one
INSERT INTO hazards (
    name, description
) VALUES (
    $1, $2
)
RETURNING name;

-- name: UpdateHazard :one
UPDATE hazards
    SET description=$2
WHERE name=$1
RETURNING name;
