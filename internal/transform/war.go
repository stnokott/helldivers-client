package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// War implements worker.docTransformer
type War struct{}

func (_ War) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.War] {
	provider := &db.DocsProvider[structs.War]{
		CollectionName: db.CollWars,
		Docs:           []db.DocWrapper[structs.War]{},
	}

	warID := data.WarID
	if warID.Id == nil {
		errFunc(errors.New("got nil war ID"))
		return provider
	}

	war := data.War
	if war.Started == nil ||
		war.Ended == nil ||
		war.Factions == nil {
		errFunc(errFromNils(war))
	} else {
		provider.Docs = append(provider.Docs, db.DocWrapper[structs.War]{
			DocID: *warID.Id,
			Document: structs.War{
				ID:        *warID.Id,
				StartTime: primitive.NewDateTimeFromTime(*war.Started),
				EndTime:   primitive.NewDateTimeFromTime(*war.Ended),
				Factions:  *war.Factions,
			},
		})
	}
	return provider
}
