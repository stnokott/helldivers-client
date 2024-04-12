package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Dispatches implements worker.docTransformer
type Dispatches struct{}

func (_ Dispatches) Transform(data APIData) (*db.DocsProvider, error) {
	if data.Dispatches == nil {
		return nil, errors.New("got nil dispatches slice")
	}

	dispatches := *data.Dispatches
	dispatchDocs := make([]db.DocWrapper, len(dispatches))

	for i, dispatch := range dispatches {
		if dispatch.Id == nil ||
			dispatch.Published == nil ||
			dispatch.Type == nil ||
			dispatch.Message == nil {
			return nil, errFromNils(&dispatch)
		}

		dispatchDocs[i] = db.DocWrapper{
			DocID: *dispatch.Id,
			Document: structs.Dispatch{
				ID:         *dispatch.Id,
				CreateTime: db.PrimitiveTime(*dispatch.Published),
				Type:       *dispatch.Type,
				Message:    *dispatch.Message,
			},
		}
	}
	return &db.DocsProvider{
		CollectionName: db.CollDispatches,
		Docs:           dispatchDocs,
	}, nil
}
