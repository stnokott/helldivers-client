// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package gen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Represents an assignment given by Super Earth to the community. This is also known as "Major Order"s in the game
type Assignment struct {
	ID int64
	// The title of the assignment
	Title string
	// A long form description of the assignment, usually contains context
	Briefing string
	// A very short summary of the description
	Description string
	// The date when the assignment will expire.
	Expiration pgtype.Timestamp
	// A list of tasks that need to be completed for this assignment
	TaskIds []int64
	// The type of reward (medals, super credits, ...)
	RewardType int32
	// The amount of Type that will be awarded
	RewardAmount pgtype.Numeric
}

type AssignmentSnapshot struct {
	ID           int64
	AssignmentID int64
	// A list of numbers, how they represent progress is unknown.
	Progress []pgtype.Numeric
}

// Represents a task in an Assignment that needs to be completed to finish the assignment
type AssignmentTask struct {
	// Auto-generated by sequence
	ID int64
	// The type of task this represents
	TaskType int32
	// A list of numbers, purpose unknown
	Values []pgtype.Numeric
	// A list of numbers, purpose unknown
	ValueTypes []pgtype.Numeric
}

// Represents information about a biomes of a planet.
type Biome struct {
	Name        string
	Description string
}

// Represents an ongoing campaign on a planet
type Campaign struct {
	// The unique identifier of this campaign
	ID int32
	// The type of campaign, this should be mapped onto an enum
	Type int32
	// Indicates how many campaigns have already been fought on this Planet
	Count pgtype.Numeric
}

// Represents a message from high command to the players, usually updates on the status of the war effort.
type Dispatch struct {
	// The unique identifier of this dispatch
	ID int32
	// When the dispatch was published
	CreateTime pgtype.Timestamp
	// The type of dispatch, purpose unknown
	Type int32
	// The message this dispatch represents
	Message string
}

// Represents an ongoing event on a Planet.
type Event struct {
	ID         int32
	CampaignID int32
	// The type of event
	Type int32
	// The faction that initiated the event
	Faction string
	// The maximum health of the Event at the time of snapshot
	MaxHealth int64
	// When the event started
	StartTime pgtype.Timestamp
	// When the event will end (or has ended).
	EndTime pgtype.Timestamp
}

// Contains dynamic data about a currently-ongoing event
type EventSnapshot struct {
	// Auto-generated by sequence
	ID      int64
	EventID int32
	Health  int64
}

// Describes an environmental hazards that can be present on a planet
type Hazard struct {
	Name        string
	Description string
}

// Represents information of a planet from the "WarInfo" endpoint returned by ArrowHead's API
type Planet struct {
	// The unique identifier ArrowHead assigned to this planet
	ID int32
	// The name of the planet, as shown in game
	Name string
	// The name of the sector the planet is in, as shown in game
	Sector string
	// The coordinates of this planet on the galactic war map in format [X, Y]
	Position []float64
	// List of indexes of all the planets to which this planet is connected
	WaypointIds []int32
	// Whether or not this planet is disabled, as assigned by ArrowHead
	Disabled bool
	// The biomes this planet has.
	BiomeName string
	// All hazardss that are applicable to this planet.
	HazardNames []string
	// The maximum health pool of this planet
	MaxHealth int64
	// The faction that originally owned the plane
	InitialOwner string
}

// Contains dynamic data about a planet currently part of this war
type PlanetSnapshot struct {
	// Auto-generated by sequence
	ID int64
	// ID of the planet this snapshot captures.
	PlanetID int32
	// The current health this planet has
	Health int64
	// The faction that currently controls the planet
	CurrentOwner string
	// Information on the active event ongoing on this planet, if one is active
	EventSnapshotID *int64
	// A list of Index integers that this planet is currently attacking.
	AttackingPlanetIds []int32
	// How much the planet regenerates per second if left alone
	RegenPerSecond float64
	// A set of statistics scoped to this planet.
	StatisticsID int64
}

// Contains the dynamic data of any metrics changing over time.
type Snapshot struct {
	// The time the snapshot of the war was taken, auto-generated as current timestamp
	CreateTime pgtype.Timestamp
	// Dynamic data about current war
	WarSnapshotID int64
	// Snapshots for currently active assignments
	AssignmentSnapshotIds []int64
	// Currently active campaigns
	CampaignIds []int32
	// Currently active dispatches
	DispatchIds []int32
	// Dynamic data about planets at point of snapshot
	PlanetSnapshotIds []int64
	// Global statistics for the current war
	StatisticsID int64
}

// Contains statistics of missions, kills, success rate etc
type SnapshotStatistic struct {
	// Auto-generated by sequence
	ID           int64
	MissionsWon  pgtype.Numeric
	MissionsLost pgtype.Numeric
	// The total amount of time spent planetside (in seconds)
	MissionTime     pgtype.Numeric
	TerminidKills   pgtype.Numeric
	AutomatonKills  pgtype.Numeric
	IlluminateKills pgtype.Numeric
	BulletsFired    pgtype.Numeric
	BulletsHit      pgtype.Numeric
	// The total amount of time played (including off-planet) in seconds
	TimePlayed pgtype.Numeric
	// The amount of casualties on the side of humanity
	Deaths pgtype.Numeric
	// The amount of revives(?)
	Revives pgtype.Numeric
	// The amount of friendly fire casualties
	Friendlies pgtype.Numeric
	// The total amount of players present (at the time of the snapshot)
	PlayerCount pgtype.Numeric
}

// Represents the global information of the ongoing war
type War struct {
	ID int32
	// When this war was started
	StartTime pgtype.Timestamp
	// When this war will end (or has ended)
	EndTime pgtype.Timestamp
	// A list of factions currently involved in the war
	Factions []string
}

// Contains the dynamic data about a war.
type WarSnapshot struct {
	// Auto-generated by sequence
	ID    int64
	WarID int32
	// A fraction used to calculate the impact of a mission on the war effort
	ImpactMultiplier float64
}
