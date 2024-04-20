package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
)

func Wars(data APIData) ([]db.EntityMerger, error) {
	warID := data.WarID
	if warID.Id == nil {
		return nil, errors.New("got nil war ID")
	}

	war := data.War
	if war.Started == nil ||
		war.Ended == nil ||
		war.Factions == nil {
		return nil, errFromNils(war)
	}
	return []db.EntityMerger{
		&db.War{
			ID:        *warID.Id,
			StartTime: db.PGTimestamp(*war.Started),
			EndTime:   db.PGTimestamp(*war.Ended),
			Factions:  *war.Factions,
		},
	}, nil
}
