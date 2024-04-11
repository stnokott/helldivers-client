// Package structs contains the types required for MongoDB mapping
package structs

import (
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
	Waypoints []int32 `bson:"waypoints,omitempty"`
	// Whether or not this planet is disabled, as assigned by ArrowHead
	Disabled bool `bson:"disabled"`
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

type Campaign struct {
	// The unique identifier of this Campaign
	ID int `bson:"_id"`
	// The planet on which this campaign is being fought
	PlanetID int `bson:"planet_id"`
	// The type of campaign, this should be mapped onto an enum
	Type int
	// Indicates how many campaigns have already been fought on this Planet
	Count int
}

type Dispatch struct {
	// The unique identifier of this dispatch
	ID int `bson:"_id"`
	// When the dispatch was published
	CreateTime primitive.Timestamp `bson:"create_time,omitempty"`
	// The type of dispatch, purpose unknown
	Type int
	// The message this dispatch represents
	Message string `bson:"message,omitempty"`
}

type Event struct {
	// The unique identifier of this event
	ID int `bson:"_id"`
	// The type of event
	Type int
	// The faction that initiated the event
	Faction string `bson:"faction,omitempty"`
	// The maximum health of the Event at the time of snapshot
	MaxHealth int `bson:"max_health"`
	// When the event started
	StartTime primitive.Timestamp `bson:"start_time,omitempty"`
	// When the event will end
	EndTime primitive.Timestamp `bson:"end_time,omitempty"`
}

type Assignment struct {
	// The unique identifier of this assignment
	ID int `bson:"_id"`
	// The title of the assignment
	Title string `bson:"title,omitempty"`
	// A long form description of the assignment, usually contains context
	Briefing string `bson:"briefing,omitempty"`
	// A very short summary of the description
	Description string `bson:"description,omitempty"`
	// A list of tasks that need to be completed for this assignment
	Tasks []AssignmentTask `bson:"tasks,omitempty"`
	// The reward for completing the assignment
	Reward AssignmentReward `bson:"reward,omitempty"`
}

// AssignmentTask represents a task in an Assignment that needs to be completed to finish the assignment
type AssignmentTask struct {
	// The type of task this represents
	Type int
	// A list of numbers, purpose unknown
	Values []int `bson:"values,omitempty"`
	// A list of numbers, purpose unknown
	ValueTypes []int `bson:"value_types,omitempty"`
}

type AssignmentReward struct {
	// The type of reward (medals, super credits, ...)
	Type int
	// The amount of Type that will be awarded
	Amount int
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

type Snapshot struct {
	// The time the snapshot of the war was taken
	Timestamp primitive.Timestamp `bson:"_id"`
	// FK ID of war
	WarID int `bson:"war_id"`
	// Currently active assignments
	AssignmentIDs []int `bson:"assignment_ids,omitempty"`
	// Currently active campaigns
	CampaignIDs []int `bson:"campaign_ids,omitempty"`
	// Currently active dispatches
	DispatchIDs []int `bson:"dispatch_ids,omitempty"`
	// Dynamic data about planets
	Planets []PlanetSnapshot `bson:"planets,omitempty"`
}

// PlanetSnapshot contains information about planets currently part of this war
type PlanetSnapshot struct {
	ID int `bson:"planet_id"`
	// The current health this planet has
	Health int
	// The faction that currently controls the planet
	CurrentOwner string `bson:"current_owner,omitempty"`
	// Information on the active event ongoing on this planet, if one is active
	Event      *EventSnapshot    `bson:"event,omitempty"`
	Statistics *PlanetStatistics `bson:"statistics,omitempty"`
}

type EventSnapshot struct {
	// FK ID of event
	EventID int `bson:"event_id"`
	// The health of the Event at the time of snapshot
	Health int
}

type PlanetStatistics struct {
	// The amount of missions won
	MissionsWon int64 `bson:"missions_won"`
	// The amount of missions lost
	MissionsLost int64 `bson:"missions_lost"`
	// The total amount of time spent planetside (in seconds)
	MissionTime int64 `bson:"mission_time"`
	Kills       StatisticsKills
	// The total amount of bullets fired
	BulletsFired int64 `bson:"bullets_fired"`
	// The total amount of bullets hit
	BulletsHit int64 `bson:"bullets_hit"`
	// The total amount of time played (including off-planet) in seconds
	TimePlayed int64 `bson:"time_played"`
	// The amount of casualties on the side of humanity
	Deaths int64
	// The amount of revives(?)
	Revives int64
	// The amount of friendly fire casualties
	Friendlies int64
	// The total amount of players present (at the time of the snapshot)
	PlayerCount int64 `bson:"player_count"`
}

type StatisticsKills struct {
	// The total amount of bugs killed since start of the season
	Terminid int64
	// The total amount of automatons killed since start of the season
	Automaton int64
	// The total amount of Illuminate killed since start of the season
	Illuminate int64
}
