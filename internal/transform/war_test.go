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

func int32Ptr(x int32) *int32 {
	return &x
}

func float64Ptr(x float64) *float64 {
	return &x
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestWarTransform(t *testing.T) {
	type args struct {
		data warRequestData
	}
	tests := []struct {
		name    string
		w       War
		args    args
		want    *db.DocsProvider
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: warRequestData{
					WarID: &api.WarId{Id: int32Ptr(2)},
					War: &api.War{
						Started:          timePtr(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local)),
						Ended:            timePtr(time.Date(2030, 01, 01, 23, 59, 0, 0, time.Local)),
						ImpactMultiplier: float64Ptr(2.5),
						Factions:         &[]string{"Humans", "Automatons"},
					},
				},
			},
			want: &db.DocsProvider{
				CollectionName: "wars",
				Docs: []db.DocWrapper{
					{
						DocID: int32(2),
						Document: structs.War{
							ID:               2,
							StartTime:        primitive.Timestamp{T: uint32(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local).Unix())},
							EndTime:          primitive.Timestamp{T: uint32(time.Date(2030, 01, 01, 23, 59, 0, 0, time.Local).Unix())},
							ImpactMultiplier: 2.5,
							Factions:         []string{"Humans", "Automatons"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil war ID",
			args: args{
				data: warRequestData{
					WarID: &api.WarId{Id: nil},
					War:   &api.War{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil end time",
			args: args{
				data: warRequestData{
					WarID: &api.WarId{Id: int32Ptr(2)},
					War: &api.War{
						Started:          timePtr(time.Date(2024, 01, 01, 23, 59, 0, 0, time.Local)),
						Ended:            nil,
						ImpactMultiplier: float64Ptr(2.5),
						Factions:         &[]string{"Humans", "Automatons"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Transform(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("War.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("War.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
