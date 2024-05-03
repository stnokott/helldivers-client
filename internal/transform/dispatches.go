package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

// Dispatches converts API data into mergable DB entities.
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

// MustDispatchMessage returns the default locale representation of a localized dispatch message.
func MustDispatchMessage(source *api.Dispatch_Message) (string, error) {
	if source == nil {
		return "", errors.New("Dispatch message is nil")
	}
	return source.AsDispatchMessage0()
}
