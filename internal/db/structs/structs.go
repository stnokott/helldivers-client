// Package structs contains the types required for MongoDB mapping
package structs

import "time"

// Planet is a planet on the galaxy map
type Planet struct {
	ID           int `bson:"_id"`
	Name         string
	Disabled     bool
	InitialOwner string  `bson:"initial_owner"`
	MaxHealth    float64 `bson:"max_health"`
	Position     Position
	Sector       string
	Waypoints    []int
}

// Position is a 2D-coordinate on the galaxy map
type Position struct {
	X int
	Y int
}

// PlanetStatus captures the current war status of a planet at a point in time.
type PlanetStatus struct {
	Timestamp      time.Time `bson:"_id"`
	PlanetID       int       `bson:"planet_id"`
	WarID          int       `bson:"war_id"`
	Health         float64
	Liberation     float64
	Owner          string
	PlayerCount    int     `bson:"players"`
	RegenPerSecond float64 `bson:"regen_per_second"`
}
