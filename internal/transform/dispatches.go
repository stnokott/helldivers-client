package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
)

func Dispatches(data APIData) ([]db.EntityMerger, error) {
	if data.Dispatches == nil {
		return nil, errors.New(("got nil dispatches slice"))
	}

	src := *data.Dispatches
	dispatches := make([]db.EntityMerger, len(src))
	for i, dispatch := range src {
		if dispatch.Id == nil ||
			dispatch.Published == nil ||
			dispatch.Type == nil ||
			dispatch.Message == nil {
			return nil, errFromNils(&dispatch)
		}

		dispatches[i] = &db.Dispatch{
			ID:         *dispatch.Id,
			CreateTime: db.PGTimestamp(*dispatch.Published),
			Type:       *dispatch.Type,
			Message:    *dispatch.Message,
		}
	}
	return dispatches, nil
}
