// Package structs contains the types required for MongoDB mapping
package structs

import (
	"encoding/binary"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	// The unique identifier ArrowHead assigned to this planet
	ID int32 `bson:"_id"`
	// The name of the planet, as shown in game
	Name string `bson:"name,omitempty"`
	// The name of the sector the planet is in, as shown in game
	Sector string `bson:"sector,omitempty"`
	// The coordinates of this planet on the galactic war map
	Position PlanetPosition `bson:"position,omitempty"`
	// A list of Index of all the planets to which this planet is connected
	Waypoints []int32 `bson:"waypoints"`
	// Whether or not this planet is disabled, as assigned by ArrowHead
	Disabled bool `bson:"disabled"`
	// The biome this planet has.
	Biome Biome `bson:"biome,omitempty"`
	// All Hazards that are applicable to this planet.
	Hazards []Hazard `bson:"hazards"`
	// The maximum health pool of this planet
	MaxHealth int64 `bson:"max_health"`
	// The faction that originally owned the planet
	InitialOwner string `bson:"initial_owner,omitempty"`
	// How much the planet regenerates per second if left alone
	RegenPerSecond float64 `bson:"regen_per_second"`
}

type PlanetPosition struct {
	X float64
	Y float64
}

type Biome struct {
	Description string `bson:"description,omitempty"`
	Name        string `bson:"name,omitempty"`
}

type Hazard struct {
	Description string `bson:"description,omitempty"`
	Name        string `bson:"name,omitempty"`
}

type Campaign struct {
	// The unique identifier of this Campaign
	ID int32 `bson:"_id"`
	// The planet on which this campaign is being fought
	PlanetID int32 `bson:"planet_id"`
	// The type of campaign, this should be mapped onto an enum
	Type int32
	// Indicates how many campaigns have already been fought on this Planet
	Count int32
}

type Dispatch struct {
	// The unique identifier of this dispatch
	ID int32 `bson:"_id"`
	// When the dispatch was published
	CreateTime primitive.Timestamp `bson:"create_time,omitempty"`
	// The type of dispatch, purpose unknown
	Type int32
	// The message this dispatch represents
	Message string `bson:"message,omitempty"`
}

type Event struct {
	// The unique identifier of this event
	ID int32 `bson:"_id"`
	// The type of event
	Type int32
	// The faction that initiated the event
	Faction string `bson:"faction,omitempty"`
	// The maximum health of the Event at the time of snapshot
	MaxHealth int64 `bson:"max_health"`
	// When the event started
	StartTime primitive.Timestamp `bson:"start_time,omitempty"`
	// When the event will end
	EndTime primitive.Timestamp `bson:"end_time,omitempty"`
}

type Assignment struct {
	// The unique identifier of this assignment
	ID int64 `bson:"_id"`
	// The title of the assignment
	Title string `bson:"title,omitempty"`
	// A long form description of the assignment, usually contains context
	Briefing string `bson:"briefing,omitempty"`
	// A very short summary of the description
	Description string `bson:"description,omitempty"`
	// The date when the assignment will expire.
	Expiration primitive.Timestamp `bson:"expiration,omitempty"`
	// A list of numbers, how they represent progress is unknown.
	Progress []int32 `bson:"progress"`
	// A list of tasks that need to be completed for this assignment
	Tasks []AssignmentTask `bson:"tasks,omitempty"`
	// The reward for completing the assignment
	Reward AssignmentReward `bson:"reward,omitempty"`
}

// AssignmentTask represents a task in an Assignment that needs to be completed to finish the assignment
type AssignmentTask struct {
	// The type of task this represents
	Type int32
	// A list of numbers, purpose unknown
	Values []int32 `bson:"values,omitempty"`
	// A list of numbers, purpose unknown
	ValueTypes []int32 `bson:"value_types,omitempty"`
}

type AssignmentReward struct {
	// The type of reward (medals, super credits, ...)
	Type int32
	// The amount of Type that will be awarded
	Amount int32
}

type War struct {
	ID int32 `bson:"_id"`
	// When this war was started
	StartTime primitive.Timestamp `bson:"start_time,omitempty"`
	// When this war will end (or has ended)
	EndTime primitive.Timestamp `bson:"end_time,omitempty"`
	// A fraction used to calculate the impact of a mission on the war effort
	ImpactMultiplier float64 `bson:"impact_multiplier,omitempty"`
	// A list of factions currently involved in the war
	Factions []string `bson:"factions,omitempty"`
}

// TODO: add global statistics
type Snapshot struct {
	// The time the snapshot of the war was taken
	Timestamp primitive.Timestamp `bson:"_id"`
	// FK ID of war
	WarID int32 `bson:"war_id"`
	// Currently active assignments
	AssignmentIDs []int64 `bson:"assignment_ids"`
	// Currently active campaigns
	CampaignIDs []int32 `bson:"campaign_ids"`
	// Currently active dispatches
	DispatchIDs []int32 `bson:"dispatch_ids"`
	// Dynamic data about planets
	Planets []PlanetSnapshot `bson:"planets"`
}

// PlanetSnapshot contains information about planets currently part of this war
type PlanetSnapshot struct {
	// ID of the planet this snapshot captures.
	ID int32 `bson:"planet_id"`
	// The current health this planet has
	Health int64
	// The faction that currently controls the planet
	CurrentOwner string `bson:"current_owner,omitempty"`
	// Information on the active event ongoing on this planet, if one is active
	Event *EventSnapshot `bson:"event,omitempty"`
	// A set of statistics scoped to this planet.
	Statistics PlanetStatistics `bson:"statistics,omitempty"`
	// A list of Index integers that this planet is currently attacking.
	Attacking []int32 `bson:"attacking"`
}

type EventSnapshot struct {
	// FK ID of event
	EventID int32 `bson:"event_id"`
	// The health of the Event at the time of snapshot
	Health int64
}

// BSONLong implements custom BSON marshallong for uint64.
//
// It is required since MongoDB natively only supports signed 64-bit values (long).
type BSONLong uint64

// MarshalBSONValue implements bson.ValueMarshaler by converting to int64 which is natively supported by MongoDB.
func (long *BSONLong) MarshalBSONValue() (bsontype.Type, []byte, error) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(*long))
	return bson.TypeInt64, bytes, nil
}

// UnmarshalBSONValue implements bson.ValueUnmarshaler by converting from int64 which is natively supported by MongoDB.
func (long *BSONLong) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	if t != bson.TypeInt64 {
		return fmt.Errorf("BSONLong was encoded as %s, needs %s", t.String(), bson.TypeInt64.String())
	}
	*long = BSONLong(binary.LittleEndian.Uint64(b))
	return nil
}

type PlanetStatistics struct {
	// The amount of missions won
	MissionsWon BSONLong `bson:"missions_won"`
	// The amount of missions lost
	MissionsLost BSONLong `bson:"missions_lost"`
	// The total amount of time spent planetside (in seconds)
	MissionTime BSONLong `bson:"mission_time"`
	Kills       StatisticsKills
	// The total amount of bullets fired
	BulletsFired BSONLong `bson:"bullets_fired"`
	// The total amount of bullets hit
	BulletsHit BSONLong `bson:"bullets_hit"`
	// The total amount of time played (including off-planet) in seconds
	TimePlayed BSONLong `bson:"time_played"`
	// The amount of casualties on the side of humanity
	Deaths BSONLong
	// The amount of revives(?)
	Revives BSONLong
	// The amount of friendly fire casualties
	Friendlies BSONLong
	// The total amount of players present (at the time of the snapshot)
	PlayerCount BSONLong `bson:"player_count"`
}

type StatisticsKills struct {
	// The total amount of bugs killed since start of the season
	Terminid BSONLong
	// The total amount of automatons killed since start of the season
	Automaton BSONLong
	// The total amount of Illuminate killed since start of the season
	Illuminate BSONLong
}
