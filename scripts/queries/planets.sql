-- name: GetPlanet :one
SELECT id FROM planets
WHERE id = $1;

-- name: MergePlanet :one
INSERT INTO planets (
    id, name, sector, position, waypoint_ids, disabled, biome_name, hazard_names, max_health, initial_owner
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (id) DO UPDATE
    SET name=$2, sector=$3, position=$4, waypoint_ids=$5, disabled=$6, biome_name=$7, hazard_names=$8, max_health=$9, initial_owner=$10
RETURNING id;

-- name: GetBiome :one
SELECT name FROM biomes
WHERE name = $1;

-- name: MergeBiome :one
INSERT INTO biomes (
    name, description
) VALUES (
    $1, $2
)
ON CONFLICT (name) DO UPDATE
    SET description=$2
RETURNING name;

-- name: GetHazard :one
SELECT name FROM hazards
WHERE name = $1;

-- name: MergeHazard :one
INSERT INTO hazards (
    name, description
) VALUES (
    $1, $2
)
ON CONFLICT (name) DO UPDATE
    SET description=$2
RETURNING name;
