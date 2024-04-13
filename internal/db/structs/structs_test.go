// Package structs contains the types required for MongoDB mapping

package structs

import (
	"math"
	"testing"
)

func TestBSONLong(t *testing.T) {
	tests := []struct {
		name string
		long BSONLong
	}{
		{"regular int", BSONLong(999)},
		{"max int64", BSONLong(math.MaxInt64)},
		{"max uint32", BSONLong(math.MaxUint32)},
		{"max uint64", BSONLong(math.MaxUint64)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ, marshalled, err := tt.long.MarshalBSONValue()
			if err != nil {
				t.Errorf("BSONLong.MarshalBSONValue() = error %v, want nil", err)
				return
			}

			var got BSONLong
			if err = got.UnmarshalBSONValue(typ, marshalled); err != nil {
				t.Errorf("BSONLong.UnmarshalBSONValue() = error %v, want nil", err)
			}

			if got != tt.long {
				t.Errorf("marshal<->unmarshal got %d, want %d", got, tt.long)
			}
		})
	}
}
