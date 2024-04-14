package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Dispatches implements worker.docTransformer
type Dispatches struct{}

func (_ Dispatches) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.Dispatch] {
	provider := &db.DocsProvider[structs.Dispatch]{
		CollectionName: db.CollDispatches,
		Docs:           []db.DocWrapper[structs.Dispatch]{},
	}

	if data.Dispatches == nil {
		errFunc(errors.New("got nil dispatches slice"))
		return provider
	}

	dispatches := *data.Dispatches

	for _, dispatch := range dispatches {
		if dispatch.Id == nil ||
			dispatch.Published == nil ||
			dispatch.Type == nil ||
			dispatch.Message == nil {
			errFunc(errFromNils(&dispatch))
			continue
		}

		provider.Docs = append(provider.Docs, db.DocWrapper[structs.Dispatch]{
			DocID: *dispatch.Id,
			Document: structs.Dispatch{
				ID:         *dispatch.Id,
				CreateTime: primitive.NewDateTimeFromTime(*dispatch.Published),
				Type:       *dispatch.Type,
				Message:    *dispatch.Message,
			},
		})
	}
	return provider
}
