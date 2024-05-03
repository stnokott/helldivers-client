package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

// Wars converts API data into mergable DB entities.
func Wars(c Converter, data APIData) ([]db.EntityMerger, error) {
	war, err := c.ConvertWar(data)
	if err != nil {
		return nil, err
	}
	return []db.EntityMerger{war}, nil
}

// MustWarID implements a converter for a war ID.
func MustWarID(source *api.WarId) (int32, error) {
	if source == nil || source.Id == nil {
		return -1, errors.New("WarID is nil")
	}
	return *source.Id, nil
}
