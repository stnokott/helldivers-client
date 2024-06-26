//go:build !goverter

package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db"
)

func mustDispatchMessage(from api.DispatchMessage0) *api.Dispatch_Message {
	dispatchMessage := new(api.Dispatch_Message)
	if err := dispatchMessage.FromDispatchMessage0(from); err != nil {
		panic(err)
	}
	return dispatchMessage
}

var validDispatch = api.Dispatch{
	Id:        ptr(int32(678)),
	Message:   mustDispatchMessage("A dispatch message"),
	Published: ptr(time.Date(2025, 1, 2, 3, 4, 5, 6, time.UTC)),
	Type:      ptr(int32(111)),
}

func TestDispatch(t *testing.T) {
	type modifier func(*api.Dispatch)
	tests := []struct {
		name     string
		modifier modifier
		want     []db.EntityMerger
		wantErr  bool
	}{
		{
			name: "valid",
			modifier: func(d *api.Dispatch) {
				// keep valid
			},
			want: []db.EntityMerger{
				&db.Dispatch{
					ID:         678,
					Message:    "A dispatch message",
					CreateTime: db.PGTimestamp(time.Date(2025, 1, 2, 3, 4, 5, 6, time.UTC)),
					Type:       111,
				},
			},
			wantErr: false,
		},
		{
			name: "empty required ID",
			modifier: func(d *api.Dispatch) {
				d.Id = nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty required message",
			modifier: func(d *api.Dispatch) {
				d.Message = nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty create time",
			modifier: func(d *api.Dispatch) {
				d.Published = nil
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dispatch api.Dispatch
			if err := copytest.DeepCopy(&dispatch, &validDispatch); err != nil {
				t.Errorf("failed to create dispatch struct copy: %v", err)
				return
			}
			// call modifiers on valid copies
			tt.modifier(&dispatch)
			data := APIData{
				Dispatches: &[]api.Dispatch{dispatch},
			}
			converter := &ConverterImpl{}
			got, err := Dispatches(converter, data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dispatches() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dispatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
