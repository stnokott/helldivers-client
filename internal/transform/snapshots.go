package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Snapshots implements worker.DocTransformer
type Snapshots struct{}

// Transform implements the worker.DocTransformer interface
func (Snapshots) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.Snapshot] {
	provider := &db.DocsProvider[structs.Snapshot]{
		CollectionName: db.CollSnapshots,
		Docs:           []db.DocWrapper[structs.Snapshot]{},
	}

	doc := db.DocWrapper[structs.Snapshot]{
		DocID:    nil,
		Document: structs.Snapshot{},
	}
	snapshotSetWarID(&doc.Document, data.WarID, errFunc)
	snapshotSetWar(&doc, data.War, errFunc)
	snapshotSetAssignments(&doc.Document, data.Assignments, errFunc)
	snapshotSetCampaigns(&doc.Document, data.Campaigns, errFunc)
	snapshotSetDispatches(&doc.Document, data.Dispatches, errFunc)
	snapshotSetPlanets(&doc.Document, data.Planets, errFunc)

	provider.Docs = append(provider.Docs, doc)
	return provider
}

func snapshotSetWarID(snap *structs.Snapshot, warID *api.WarId, errFunc func(error)) {
	if warID == nil || warID.Id == nil {
		errFunc(errors.New("got nil War ID, will be omitted"))
		return
	}
	snap.WarSnapshot.WarID = *warID.Id
}

func snapshotSetWar(doc *db.DocWrapper[structs.Snapshot], warPtr *api.War, errFunc func(error)) {
	if warPtr == nil {
		errFunc(errors.New("got nil War, snapshot timestamp will be omitted"))
		return
	}
	war := *warPtr
	if war.Now == nil || war.ImpactMultiplier == nil {
		errFunc(errFromNils(warPtr))
		return
	}
	doc.DocID = primitive.NewDateTimeFromTime(*war.Now)
	doc.Document.Timestamp = primitive.NewDateTimeFromTime(*war.Now)
	doc.Document.WarSnapshot.ImpactMultiplier = *war.ImpactMultiplier

	stats, err := makeWarStatistics(war.Statistics)
	if err != nil {
		errFunc(err)
		return
	}
	doc.Document.Statistics = *stats
}

func snapshotSetAssignments(snap *structs.Snapshot, assignmentsPtr *[]api.Assignment2, errFunc func(error)) {
	if assignmentsPtr == nil {
		errFunc(errors.New("got nil Assignments slice, will be omitted"))
		return
	}
	assignments := *assignmentsPtr
	assignmentIds := []int64{}
	for _, assignment := range assignments {
		if assignment.Id == nil {
			errFunc(errors.New("got nil Assignment ID, will be omitted"))
			continue
		}
		assignmentIds = append(assignmentIds, *assignment.Id)
	}
	snap.AssignmentIDs = assignmentIds
}

func snapshotSetCampaigns(snap *structs.Snapshot, campaignsPtr *[]api.Campaign2, errFunc func(error)) {
	if campaignsPtr == nil {
		errFunc(errors.New("got nil Campaigns slice, will be omitted"))
		return
	}
	campaigns := *campaignsPtr
	campaignIDs := []int32{}
	for _, campaign := range campaigns {
		if campaign.Id == nil {
			errFunc(errors.New("got nil Campaign ID, will be omitted"))
			continue
		}
		campaignIDs = append(campaignIDs, *campaign.Id)
	}
	snap.CampaignIDs = campaignIDs
}

func snapshotSetDispatches(snap *structs.Snapshot, dispatchesPtr *[]api.Dispatch, errFunc func(error)) {
	if dispatchesPtr == nil {
		errFunc(errors.New("got nil Dispatches slice, will be omitted"))
		return
	}
	dispatches := *dispatchesPtr
	dispatchIDs := []int32{}
	for _, dispatch := range dispatches {
		if dispatch.Id == nil {
			errFunc(errors.New("got nil Dispatch ID, will be omitted"))
			continue
		}
		dispatchIDs = append(dispatchIDs, *dispatch.Id)
	}
	snap.DispatchIDs = dispatchIDs
}

func snapshotSetPlanets(snap *structs.Snapshot, planetsPtr *[]api.Planet, errFunc func(error)) {
	if planetsPtr == nil {
		errFunc(errors.New("got nil Planets slice, will be omitted"))
		return
	}
	planets := *planetsPtr
	planetSnapshots := []structs.PlanetSnapshot{}
	for _, planet := range planets {
		if planet.Index == nil ||
			planet.Health == nil ||
			planet.CurrentOwner == nil {
			errFunc(errFromNils(&planet))
			continue
		}
		eventSnapshot, err := makeEventSnapshot(planet.Event)
		if err != nil {
			errFunc(err)
			continue
		}
		planetStatistics, err := makePlanetStatistics(planet.Statistics)
		if err != nil {
			errFunc(err)
			continue
		}
		var attacking []int32
		if planet.Attacking != nil {
			attacking = *planet.Attacking
		}
		planetSnap := structs.PlanetSnapshot{
			ID:           *planet.Index,
			Health:       *planet.Health,
			CurrentOwner: *planet.CurrentOwner,
			Event:        eventSnapshot,
			Statistics:   *planetStatistics,
			Attacking:    attacking,
		}
		planetSnapshots = append(planetSnapshots, planetSnap)
	}
	snap.Planets = planetSnapshots
}

func makeEventSnapshot(eventPtr *api.Planet_Event) (*structs.EventSnapshot, error) {
	if eventPtr == nil {
		// events are optional, so nil is ok
		return nil, nil
	}
	planetEvent, err := eventPtr.AsEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Planet Event: %w", err)
	}
	if planetEvent.Id == nil || planetEvent.Health == nil {
		return nil, errFromNils(&planetEvent)
	}
	return &structs.EventSnapshot{
		EventID: *planetEvent.Id,
		Health:  *planetEvent.Health,
	}, nil
}

func makeStatistics(stats api.Statistics) (*structs.Statistics, error) {
	if stats.MissionsWon == nil ||
		stats.MissionsLost == nil ||
		stats.MissionTime == nil ||
		stats.TerminidKills == nil ||
		stats.AutomatonKills == nil ||
		stats.IlluminateKills == nil ||
		stats.BulletsFired == nil ||
		stats.BulletsHit == nil ||
		stats.TimePlayed == nil ||
		stats.Deaths == nil ||
		stats.Revives == nil ||
		stats.Friendlies == nil ||
		stats.PlayerCount == nil {
		return nil, errFromNils(&stats)
	}

	return &structs.Statistics{
		MissionsWon:  structs.BSONLong(*stats.MissionsWon),
		MissionsLost: structs.BSONLong(*stats.MissionsLost),
		MissionTime:  structs.BSONLong(*stats.MissionTime),
		Kills: structs.StatisticsKills{
			Terminid:   structs.BSONLong(*stats.TerminidKills),
			Automaton:  structs.BSONLong(*stats.AutomatonKills),
			Illuminate: structs.BSONLong(*stats.IlluminateKills),
		},
		BulletsFired: structs.BSONLong(*stats.BulletsFired),
		BulletsHit:   structs.BSONLong(*stats.BulletsHit),
		TimePlayed:   structs.BSONLong(*stats.TimePlayed),
		Deaths:       structs.BSONLong(*stats.Deaths),
		Revives:      structs.BSONLong(*stats.Revives),
		Friendlies:   structs.BSONLong(*stats.Friendlies),
		PlayerCount:  structs.BSONLong(*stats.PlayerCount),
	}, nil
}

func makePlanetStatistics(statsPtr *api.Planet_Statistics) (*structs.Statistics, error) {
	if statsPtr == nil {
		return nil, errors.New("got nil Planet Statistics")
	}
	stats, err := statsPtr.AsStatistics()
	if err != nil {
		return nil, fmt.Errorf("cannot parse Planet Statistics: %w", err)
	}
	return makeStatistics(stats)
}

func makeWarStatistics(statsPtr *api.War_Statistics) (*structs.Statistics, error) {
	if statsPtr == nil {
		return nil, errors.New("got nil War Statistics")
	}
	stats, err := statsPtr.AsStatistics()
	if err != nil {
		return nil, fmt.Errorf("cannot parse War Statistics: %w", err)
	}
	return makeStatistics(stats)
}
