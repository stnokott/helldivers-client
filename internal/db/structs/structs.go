// Package structs contains the types required for MongoDB mapping
package structs

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
