package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

func TestDispatchesTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		d       Dispatches
		args    args
		want    *db.DocsProvider
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					Dispatches: &[]api.Dispatch{
						{
							Id:        ptr(int32(5)),
							Published: ptr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local)),
							Type:      ptr(int32(7)),
							Message:   ptr("Foo"),
						},
						{
							Id:        ptr(int32(6)),
							Published: ptr(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)),
							Type:      ptr(int32(8)),
							Message:   ptr("Bar"),
						},
					},
				},
			},
			want: &db.DocsProvider{
				CollectionName: "dispatches",
				Docs: []db.DocWrapper{
					{
						DocID: int32(5),
						Document: structs.Dispatch{
							ID:         5,
							CreateTime: db.PrimitiveTime(time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local)),
							Type:       7,
							Message:    "Foo",
						},
					},
					{
						DocID: int32(6),
						Document: structs.Dispatch{
							ID:         6,
							CreateTime: db.PrimitiveTime(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)),
							Type:       8,
							Message:    "Bar",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil type",
			args: args{
				data: APIData{
					Dispatches: &[]api.Dispatch{
						{
							Id:        ptr(int32(5)),
							Published: ptr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local)),
							Type:      nil,
							Message:   ptr("Foo"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Transform(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dispatches.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dispatches.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
