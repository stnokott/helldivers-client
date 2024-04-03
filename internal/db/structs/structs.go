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
	// War season change over time
	History []WarSeasonHistory
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
	// Planets being currently attacked from this planet
	AttackTargets []int `bson:"attack_targets"`
	// active campaign on this planet (optional)
	Campaign *PlanetCampaign
}

// WarSeasonHistory captures the historic status of the Helldivers offensive in the galactic war
type WarSeasonHistory struct {
	// Time at which this status was retrieved from the API
	Timestamp time.Time `bson:"_id"`
	// Always empty AFAIK, haven't figured this out
	ActiveElectionPolicyEffects []int `bson:"active_election_policy_effects"`
	// Always empty AFAIK, haven't figured this out
	CommunityTargets []int `bson:"community_targets"`
	// I don't fully understand what this does, feel free to ping me if you know
	ImpactMultiplier float64 `bson:"impact_multiplier"`
	// Currently active global event, past and present
	GlobalEvents []WarSeasonHistoryGlobalEvent `bson:"global_events"`
}

// PlanetCampaign contains information about a currently active campaign
type PlanetCampaign struct {
	// not sure what this counts, it's generally a low number
	Count int
	// The type of this campaign, haven't found out what they mean yet
	Type int
}

// WarSeasonHistoryGlobalEvent contains information about a global event, past and present
type WarSeasonHistoryGlobalEvent struct {
	// The title of the global event, appears to be more a status than an actual title
	Title string
	// A list of effects, usually strategems or bonuses
	Effects []string
	// Planets affected by this event
	PlanetIDs []int `bson:"planet_ids"`
	// The race involved in this campaign (so far seems to always be 'Human')
	Race string
	// The localized message from Super Earth about the global event
	Message WarNewsMessage
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
