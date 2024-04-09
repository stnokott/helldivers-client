// Package structs contains the types required for MongoDB mapping
package structs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID             int            `bson:"_id"`
	Name           string         `bson:"name,omitempty"`
	Sector         string         `bson:"sector,omitempty"`
	Position       PlanetPosition `bson:"position,omitempty"`
	Waypoints      []int          `bson:"waypoints,omitempty"`
	Disabled       bool           `bson:"disabled"`
	MaxHealth      int            `bson:"max_health"`
	InitialOwner   string         `bson:"initial_owner,omitempty"`
	RegenPerSecond float64        `bson:"regen_per_second"`
}

type PlanetPosition struct {
	X int
	Y int
}

type Campaign struct {
	ID       int `bson:"_id"`
	PlanetID int `bson:"planet_id"`
	Type     int
	Count    int
}

type Dispatch struct {
	ID         int                 `bson:"_id"`
	CreateTime primitive.Timestamp `bson:"create_time,omitempty"`
	Type       int
	Message    string `bson:"message,omitempty"`
}

type Event struct {
	ID        int `bson:"_id"`
	Type      int
	Faction   string              `bson:"faction,omitempty"`
	MaxHealth int                 `bson:"max_health"`
	StartTime primitive.Timestamp `bson:"start_time,omitempty"`
	EndTime   primitive.Timestamp `bson:"end_time,omitempty"`
}

type Assignment struct {
	ID          int              `bson:"_id"`
	Title       string           `bson:"title,omitempty"`
	Briefing    string           `bson:"briefing,omitempty"`
	Description string           `bson:"description,omitempty"`
	Tasks       []AssignmentTask `bson:"tasks,omitempty"`
	Reward      AssignmentReward `bson:"reward,omitempty"`
}

type AssignmentTask struct {
	Type       int
	Values     []int `bson:"values,omitempty"`
	ValueTypes []int `bson:"value_types,omitempty"`
}

type AssignmentReward struct {
	Type   int
	Amount int
}

type War struct {
	ID               int                 `bson:"_id"`
	StartTime        primitive.Timestamp `bson:"start_time,omitempty"`
	EndTime          primitive.Timestamp `bson:"end_time,omitempty"`
	ImpactMultiplier float64             `bson:"impact_multiplier,omitempty"`
	Factions         []string            `bson:"factions,omitempty"`
}

type Snapshot struct {
	ID            primitive.Timestamp `bson:"_id"`
	WarID         int                 `bson:"war_id"`
	AssignmentIDs []int               `bson:"assignment_ids,omitempty"`
	CampaignIDs   []int               `bson:"campaign_ids,omitempty"`
	DispatchIDs   []int               `bson:"dispatch_ids,omitempty"`
	Planets       []PlanetSnapshot    `bson:"planets,omitempty"`
}

type PlanetSnapshot struct {
	ID           int `bson:"planet_id"`
	Health       int
	CurrentOwner string            `bson:"current_owner,omitempty"`
	Event        *EventSnapshot    `bson:"event,omitempty"`
	Statistics   *PlanetStatistics `bson:"statistics,omitempty"`
}

type EventSnapshot struct {
	EventID int `bson:"event_id"`
	Health  int
}

type PlanetStatistics struct {
	MissionsWon  int64 `bson:"missions_won"`
	MissionsLost int64 `bson:"missions_lost"`
	MissionTime  int64 `bson:"mission_time"`
	Kills        StatisticsKills
	BulletsFired int64 `bson:"bullets_fired"`
	BulletsHit   int64 `bson:"bullets_hit"`
	TimePlayed   int64 `bson:"time_played"`
	Deaths       int64
	Revives      int64
	Friendlies   int64
	PlayerCount  int64 `bson:"player_count"`
}

type StatisticsKills struct {
	Terminid   int64
	Automaton  int64
	Illuminate int64
}
