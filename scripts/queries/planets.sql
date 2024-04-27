-- name: GetPlanet :one
SELECT id FROM planets
WHERE id = $1;

-- name: PlanetExists :one
SELECT EXISTS(SELECT * FROM planets WHERE id = $1);

-- name: MergePlanet :execrows
INSERT INTO planets (
    id, name, sector, position, waypoint_ids, disabled, biome_name, hazard_names, max_health, initial_owner
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (id) DO UPDATE
    SET name=$2, sector=$3, position=$4, waypoint_ids=$5, disabled=$6, biome_name=$7, hazard_names=$8, max_health=$9, initial_owner=$10
WHERE FALSE IN (
    EXCLUDED.name=$2, EXCLUDED.sector=$3, EXCLUDED.position=$4,EXCLUDED. waypoint_ids=$5, EXCLUDED.disabled=$6, EXCLUDED.biome_name=$7, EXCLUDED.hazard_names=$8, EXCLUDED.max_health=$9, EXCLUDED.initial_owner=$10
);

-- name: GetBiome :one
SELECT name FROM biomes
WHERE name = $1;

-- name: BiomeExists :one
SELECT EXISTS(SELECT * FROM biomes WHERE name = $1);

-- name: MergeBiome :execrows
INSERT INTO biomes (
    name, description
) VALUES (
    $1, $2
)
ON CONFLICT (name) DO UPDATE
    SET description=$2
WHERE FALSE IN (
    EXCLUDED.description=$2
);

-- name: GetHazard :one
SELECT name FROM hazards
WHERE name = $1;

-- name: HazardExists :one
SELECT EXISTS(SELECT * FROM hazards WHERE name = $1);

-- name: MergeHazard :execrows
INSERT INTO hazards (
    name, description
) VALUES (
    $1, $2
)
ON CONFLICT (name) DO UPDATE
    SET description=$2
WHERE FALSE IN (
    EXCLUDED.description=$2
);

