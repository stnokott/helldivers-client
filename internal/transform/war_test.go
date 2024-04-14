// Package transform converts API structs to DB structs
package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestWarTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		w       War
		args    args
		want    *db.DocsProvider[structs.War]
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					WarID: &api.WarId{Id: ptr(int32(2))},
					War: &api.War{
						Started:  ptr(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local)),
						Ended:    ptr(time.Date(2030, 01, 01, 23, 59, 0, 0, time.Local)),
						Factions: &[]string{"Humans", "Automatons"},
					},
				},
			},
			want: &db.DocsProvider[structs.War]{
				CollectionName: "wars",
				Docs: []db.DocWrapper[structs.War]{
					{
						DocID: int32(2),
						Document: structs.War{
							ID:        2,
							StartTime: primitive.NewDateTimeFromTime(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local)),
							EndTime:   primitive.NewDateTimeFromTime(time.Date(2030, 01, 01, 23, 59, 0, 0, time.Local)),
							Factions:  []string{"Humans", "Automatons"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil war ID",
			args: args{
				data: APIData{
					WarID: &api.WarId{Id: nil},
					War: &api.War{
						Started:  ptr(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local)),
						Ended:    ptr(time.Date(2030, 01, 01, 23, 59, 0, 0, time.Local)),
						Factions: &[]string{"Humans", "Automatons"},
					},
				},
			},
			want: &db.DocsProvider[structs.War]{
				CollectionName: db.CollWars,
				Docs:           []db.DocWrapper[structs.War]{},
			},
			wantErr: true,
		},
		{
			name: "nil end time",
			args: args{
				data: APIData{
					WarID: &api.WarId{Id: ptr(int32(2))},
					War: &api.War{
						Started:  ptr(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local)),
						Ended:    nil,
						Factions: &[]string{"Humans", "Automatons"},
					},
				},
			},
			want: &db.DocsProvider[structs.War]{
				CollectionName: db.CollWars,
				Docs:           []db.DocWrapper[structs.War]{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := false
			errFunc := func(err error) {
				if !tt.wantErr {
					t.Logf("War.Transform() error: %v", err)
				}
				gotErr = true
			}
			got := tt.w.Transform(tt.args.data, errFunc)
			if gotErr != tt.wantErr {
				t.Errorf("War.Transform() returned error, wantErr %v", tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("War.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
