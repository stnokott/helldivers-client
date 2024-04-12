package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

func mustPlanetEvent(from api.Event) *api.Planet_Event {
	planetEvent := new(api.Planet_Event)
	if err := planetEvent.FromEvent(from); err != nil {
		panic(err)
	}
	return planetEvent
}

func TestEventsTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		e       Events
		args    args
		want    *db.DocsProvider[structs.Event]
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					Planets: &[]api.Planet{
						{
							Event: mustPlanetEvent(
								api.Event{
									Id:        ptr(int32(5)),
									EventType: ptr(int32(6)),
									Faction:   ptr("Foo"),
									MaxHealth: ptr(int64(10000)),
									StartTime: ptr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local)),
									EndTime:   ptr(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)),
								},
							),
						},
						{
							Event: mustPlanetEvent(
								api.Event{
									Id:        ptr(int32(6)),
									EventType: ptr(int32(7)),
									Faction:   ptr("Bar"),
									MaxHealth: ptr(int64(10000)),
									StartTime: ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)),
									EndTime:   ptr(time.Date(2027, 1, 1, 0, 0, 0, 0, time.Local)),
								},
							),
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Event]{
				CollectionName: "events",
				Docs: []db.DocWrapper[structs.Event]{
					{
						DocID: int32(5),
						Document: structs.Event{
							ID:        5,
							Type:      6,
							Faction:   "Foo",
							MaxHealth: 10000,
							StartTime: db.PrimitiveTime(time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local)),
							EndTime:   db.PrimitiveTime(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)),
						},
					},
					{
						DocID: int32(6),
						Document: structs.Event{
							ID:        6,
							Type:      7,
							Faction:   "Bar",
							MaxHealth: 10000,
							StartTime: db.PrimitiveTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)),
							EndTime:   db.PrimitiveTime(time.Date(2027, 1, 1, 0, 0, 0, 0, time.Local)),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil faction",
			args: args{
				data: APIData{
					Planets: &[]api.Planet{
						{
							Event: mustPlanetEvent(
								api.Event{
									Id:        ptr(int32(5)),
									EventType: ptr(int32(6)),
									Faction:   nil,
									MaxHealth: ptr(int64(10000)),
									StartTime: ptr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local)),
									EndTime:   ptr(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)),
								},
							),
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
			got, err := tt.e.Transform(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Events.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Events.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
