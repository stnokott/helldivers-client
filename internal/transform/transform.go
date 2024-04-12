// Package transform converts API structs to DB structs
package transform

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stnokott/helldivers-client/internal/api"
)

type APIData struct {
	Planets     *[]api.Planet
	WarID       *api.WarId
	War         *api.War
	Campaigns   *[]api.Campaign2
	Dispatches  *[]api.Dispatch
	Assignments *[]api.Assignment2
}

// errFromNils returns an error containing the list of nil fields in v.
func errFromNils[T any](v *T) error {
	names := []string{}
	value := reflect.ValueOf(v).Elem()
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.IsNil() {
			names = append(names, valueType.Field(i).Name)
		}
	}

	return fmt.Errorf("nil fields in %s struct: %s", valueType.Name(), strings.Join(names, ", "))
}
