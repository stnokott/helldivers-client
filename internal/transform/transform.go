// Package transform converts API structs to DB structs
package transform

import (
	"github.com/stnokott/helldivers-client/internal/api"
)

// APIData contains all relevant (unprocessed) data from the API, used for further processing.
type APIData struct {
	WarID       *api.WarId
	War         *api.War
	Planets     *[]api.Planet
	Campaigns   *[]api.Campaign2
	Dispatches  *[]api.Dispatch
	Assignments *[]api.Assignment2
}
