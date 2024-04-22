package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

var dispatchValid = api.Dispatch{
	Id:        ptr(int32(678)),
	Message:   ptr("A dispatch message"),
	Published: ptr(time.Date(2025, 1, 2, 3, 4, 5, 6, time.Local)),
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
					CreateTime: db.PGTimestamp(time.Date(2025, 1, 2, 3, 4, 5, 6, time.Local)),
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
			dispatch := dispatchValid
			// call modifiers on valid copies
			tt.modifier(&dispatch)
			data := APIData{
				Dispatches: &[]api.Dispatch{dispatch},
			}
			got, err := Dispatches(data)
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
