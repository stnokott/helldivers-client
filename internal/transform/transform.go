// Package transform converts API structs to DB structs
package transform

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

type War struct{}

type warRequestData struct {
	WarID *api.WarId
	War   *api.War
}

func (_ War) Request(api *client.Client, ctx context.Context) (data warRequestData, err error) {
	data.WarID, err = apiWithTimeout(api.WarID, 1*time.Second)
	if err != nil {
		return
	}
	data.War, err = apiWithTimeout(api.War, 5*time.Second)
	if err != nil {
		return
	}
	return
}

func (_ War) Transform(data warRequestData) (db.DocProvider, error) {
	warID := data.WarID
	if warID.Id == nil {
		return nil, errors.New("got nil war ID")
	}

	war := data.War
	if war.Started == nil ||
		war.Ended == nil ||
		war.ImpactMultiplier == nil ||
		war.Factions == nil {
		return nil, errFromNils(war)
	}
	return &warDocProvider{
		ID:               *warID.Id,
		StartTime:        db.PrimitiveTime(*war.Started),
		EndTime:          db.PrimitiveTime(*war.Ended),
		ImpactMultiplier: *war.ImpactMultiplier,
		Factions:         *war.Factions,
	}, nil
}

type warDocProvider structs.War

func (p *warDocProvider) DocID() any {
	return p.ID
}

func (p *warDocProvider) Document() any {
	return p
}

func (p *warDocProvider) CollectionName() db.CollectionName {
	return db.CollWars
}

func apiWithTimeout[T any](apiFunc func(context.Context) (T, error), timeout time.Duration) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return apiFunc(ctx)
}

// errFromNils returns an error containing the list of nil fields in v.
func errFromNils(v any) error {
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
