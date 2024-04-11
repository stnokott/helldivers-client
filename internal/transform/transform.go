// Package transform converts API structs to DB structs
package transform

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

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
