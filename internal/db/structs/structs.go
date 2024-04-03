// Package structs contains the types required for MongoDB mapping
package structs

import "time"

// WarSeason represents the global overview of the war, it's planets, capitals etc
type WarSeason struct {
	// The identifier for this war, taken directly from the API
	ID int `bson:"_id"`
	// Empty, not mapped yet
	Capitals []any
	// Empty, not mapped yet
	PlanetPermanentEffects []any `bson:"planet_permanent_effects"`
	// When this war season was started
	StartDate time.Time `bson:"start_date"`
	// When this war season is scheduled to end
	EndDate time.Time `bson:"end_date"`
	// All planets present in this war season
	Planets []Planet
}

// Planet is a planet on the galaxy map
type Planet struct {
	// The index of this planet, for convenience kept the same as in the official API
	ID int `bson:"_id"`
	// The human readable name of the planet, or unknown if it's not a known name
	Name string
	// Whether or not this planet is currently playable (enabled)
	Disabled bool
	// Which faction originally claimed this planet
	InitialOwner string `bson:"initial_owner"`
	// Maximum health of this planet, used in conflict states
	MaxHealth float64 `bson:"max_health"`
	// The coordinates in the galaxy where this planet is located
	Position Position
	// The name of the sector this planet resides in (or the identifier as a string if it's not a known sector)
	Sector string
	// Waypoints, seems to link planets together but purpose unclear
	Waypoints []int
	// Planet status change over time during a war
	History []PlanetHistory
}

// Position is a 2D-coordinate on the galaxy map
type Position struct {
	X int
	Y int
}

// PlanetHistory captures the current war status of a planet at a point in time.
type PlanetHistory struct {
	// Time at which this status was retrieved from the API
	Timestamp time.Time `bson:"_id"`
	// The current 'health' of this planet
	Health float64
	// The progression of liberation on this planet, presented as a %
	Liberation float64
	// The faction that owns the planet at this moment
	Owner string
	// The amount of helldivers currently on this planet
	PlayerCount int `bson:"players"`
	// At which rate this planet will regenerate it's health
	RegenPerSecond float64 `bson:"regen_per_second"`
}

// WarNews represents a message in the Helldivers 2 newsfeed
type WarNews struct {
	// The identifier of this campaign
	ID int `bson:"_id"`
	// Localized versions of a newsfeed message
	Message WarNewsMessage
	// When this message was published
	Published time.Time
	// A type identifier, haven't figured out what they mean (seems to be 0 mostly)
	Type int
}

// WarNewsMessage contains localized versions of a newsfeed message
type WarNewsMessage struct {
	// The message from Super Earth about the news in de-DE
	DE string
	// The message from Super Earth about the news in en-US
	EN string
	// The message from Super Earth about the news in es-ES
	ES string
	// The message from Super Earth about the news in fr-FR
	FR string
	// The message from Super Earth about the news in it-IT
	IT string
	// The message from Super Earth about the news in pl-PL
	PL string
	// The message from Super Earth about the news in ru-RU
	RU string
	// The message from Super Earth about the news in zh-Hans
	ZH string
}
