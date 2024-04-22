package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func Snapshot(data APIData) (mergers []db.EntityMerger, err error) {
	s := &db.Snapshot{
		Snapshot: gen.Snapshot{
			WarSnapshotID: -1, // will be filled during insert
			StatisticsID:  -1, // will be filled during insert
		},
		WarSnapshot:     gen.WarSnapshot{},
		PlanetSnapshots: []db.PlanetSnapshot{},
		Statistics:      gen.SnapshotStatistic{},
	}

	if err = snapshotSetWarID(s, data.WarID); err != nil {
		return
	}
	if err = snapshotSetWar(s, data.War); err != nil {
		return
	}
	if err = snapshotSetAssignments(s, data.Assignments); err != nil {
		return
	}
	if err = snapshotSetCampaigns(s, data.Campaigns); err != nil {
		return
	}
	if err = snapshotSetDispatches(s, data.Dispatches); err != nil {
		return
	}
	if err = snapshotSetPlanets(s, data.Planets); err != nil {
		return
	}

	mergers = []db.EntityMerger{s}
	return
}

func snapshotSetWarID(snap *db.Snapshot, warID *api.WarId) error {
	if warID == nil || warID.Id == nil {
		return errors.New("got nil War ID, will be omitted")
	}
	snap.WarSnapshot.WarID = *warID.Id
	return nil
}

func snapshotSetWar(snap *db.Snapshot, warPtr *api.War) error {
	if warPtr == nil {
		return errors.New("got nil War, snapshot timestamp will be omitted")
	}
	war := *warPtr
	if war.ImpactMultiplier == nil {
		return errFromNils(warPtr)
	}
	snap.WarSnapshot.ImpactMultiplier = *war.ImpactMultiplier

	stats, err := makeWarStatistics(war.Statistics)
	if err != nil {
		return err
	}
	snap.Statistics = *stats
	return nil
}

func snapshotSetAssignments(snap *db.Snapshot, assignmentsPtr *[]api.Assignment2) error {
	if assignmentsPtr == nil {
		return errors.New("got nil Assignments slice, will be omitted")
	}
	assignments := *assignmentsPtr
	assignmentIDs := make([]int64, len(assignments))
	for i, assignment := range assignments {
		if assignment.Id == nil {
			return errors.New("got nil Assignment ID, will be omitted")
		}
		assignmentIDs[i] = *assignment.Id
	}
	snap.Snapshot.AssignmentIds = assignmentIDs
	return nil
}

func snapshotSetCampaigns(snap *db.Snapshot, campaignsPtr *[]api.Campaign2) error {
	if campaignsPtr == nil {
		return errors.New("got nil Campaigns slice, will be omitted")
	}
	campaigns := *campaignsPtr
	campaignIDs := make([]int32, len(campaigns))
	for i, campaign := range campaigns {
		if campaign.Id == nil {
			return errors.New("got nil Campaign ID, will be omitted")
		}
		campaignIDs[i] = *campaign.Id
	}
	snap.Snapshot.CampaignIds = campaignIDs
	return nil
}

func snapshotSetDispatches(snap *db.Snapshot, dispatchesPtr *[]api.Dispatch) error {
	if dispatchesPtr == nil {
		return errors.New("got nil Dispatches slice, will be omitted")
	}
	dispatches := *dispatchesPtr
	dispatchIDs := make([]int32, len(dispatches))
	for i, dispatch := range dispatches {
		if dispatch.Id == nil {
			return errors.New("got nil Dispatch ID, will be omitted")
		}
		dispatchIDs[i] = *dispatch.Id
	}
	snap.Snapshot.DispatchIds = dispatchIDs
	return nil
}

func snapshotSetPlanets(snap *db.Snapshot, planetsPtr *[]api.Planet) error {
	if planetsPtr == nil {
		return errors.New("got nil Planets slice, will be omitted")
	}
	planets := *planetsPtr
	planetSnapshots := make([]db.PlanetSnapshot, len(planets))
	for i, planet := range planets {
		if planet.Index == nil ||
			planet.Health == nil ||
			planet.CurrentOwner == nil ||
			planet.RegenPerSecond == nil {
			return errFromNils(&planet)
		}
		eventSnapshot, err := makeEventSnapshot(planet.Event)
		if err != nil {
			return err
		}

		planetStatistics, err := makePlanetStatistics(planet.Statistics)
		if err != nil {
			return err
		}
		var attacking []int32
		if planet.Attacking != nil {
			attacking = *planet.Attacking
		}
		planetSnap := db.PlanetSnapshot{
			PlanetSnapshot: gen.PlanetSnapshot{
				PlanetID:           *planet.Index,
				Health:             *planet.Health,
				CurrentOwner:       *planet.CurrentOwner,
				EventSnapshotID:    nil, // will be filled later from DB
				AttackingPlanetIds: attacking,
				StatisticsID:       -1, // will be filled later from DB
				RegenPerSecond:     *planet.RegenPerSecond,
			},
			Event:      eventSnapshot,
			Statistics: *planetStatistics,
		}
		planetSnapshots[i] = planetSnap
	}
	snap.PlanetSnapshots = planetSnapshots
	return nil
}

func makeEventSnapshot(eventPtr *api.Planet_Event) (*gen.EventSnapshot, error) {
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
	return &gen.EventSnapshot{
		EventID: *planetEvent.Id,
		Health:  *planetEvent.Health,
	}, nil
}

func makeStatistics(stats api.Statistics) (*gen.SnapshotStatistic, error) {
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

	return &gen.SnapshotStatistic{
		MissionsWon:     db.PGUint64(*stats.MissionsWon),
		MissionsLost:    db.PGUint64(*stats.MissionsLost),
		MissionTime:     db.PGUint64(*stats.MissionTime),
		TerminidKills:   db.PGUint64(*stats.TerminidKills),
		AutomatonKills:  db.PGUint64(*stats.AutomatonKills),
		IlluminateKills: db.PGUint64(*stats.IlluminateKills),
		BulletsFired:    db.PGUint64(*stats.BulletsFired),
		BulletsHit:      db.PGUint64(*stats.BulletsHit),
		TimePlayed:      db.PGUint64(*stats.TimePlayed),
		Deaths:          db.PGUint64(*stats.Deaths),
		Revives:         db.PGUint64(*stats.Revives),
		Friendlies:      db.PGUint64(*stats.Friendlies),
		PlayerCount:     db.PGUint64(*stats.PlayerCount),
	}, nil
}

func makePlanetStatistics(statsPtr *api.Planet_Statistics) (*gen.SnapshotStatistic, error) {
	if statsPtr == nil {
		return nil, errors.New("got nil Planet Statistics")
	}
	stats, err := statsPtr.AsStatistics()
	if err != nil {
		return nil, fmt.Errorf("cannot parse Planet Statistics: %w", err)
	}
	return makeStatistics(stats)
}

func makeWarStatistics(statsPtr *api.War_Statistics) (*gen.SnapshotStatistic, error) {
	if statsPtr == nil {
		return nil, errors.New("got nil War Statistics")
	}
	stats, err := statsPtr.AsStatistics()
	if err != nil {
		return nil, fmt.Errorf("cannot parse War Statistics: %w", err)
	}
	return makeStatistics(stats)
}
