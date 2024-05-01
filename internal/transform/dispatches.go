package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
)

func Dispatches(c Converter, data APIData) ([]db.EntityMerger, error) {
	if data.Dispatches == nil {
		return nil, errors.New(("got nil dispatches slice"))
	}

	src := *data.Dispatches
	mergers := make([]db.EntityMerger, len(src))
	for i, dispatch := range src {
		merger, err := c.ConvertDispatch(dispatch)
		if err != nil {
			return nil, err
		}
		mergers[i] = merger
	}
	return mergers, nil
}
