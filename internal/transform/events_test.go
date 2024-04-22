package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

func mustPlanetEvent(from api.Event) *api.Planet_Event {
	planetEvent := new(api.Planet_Event)
	if err := planetEvent.FromEvent(from); err != nil {
		panic(err)
	}
	return planetEvent
}

var validEvent = api.Event{
	Id:                ptr(int32(997)),
	StartTime:         ptr(time.Date(2024, 4, 5, 6, 7, 8, 9, time.Local)),
	EndTime:           ptr(time.Date(2025, 4, 5, 6, 7, 8, 9, time.Local)),
	EventType:         ptr(int32(667)),
	Faction:           ptr("Terminids"),
	MaxHealth:         ptr(int64(4455667788)),
	CampaignId:        nil, // not required, linked through Planet
	Health:            nil, // not required, persisted in dynamic snapshots
	JointOperationIds: nil, // not required
}

func TestEvent(t *testing.T) {
	// modifier for planet to allow nulling event
	type modifier func(*api.Planet)
	tests := []struct {
		name     string
		modifier modifier
		want     []db.EntityMerger
		wantErr  bool
	}{
		{
			name: "valid",
			modifier: func(p *api.Planet) {
				// keep valid
			},
			want: []db.EntityMerger{
				&db.Event{
					ID:        997,
					StartTime: db.PGTimestamp(time.Date(2024, 4, 5, 6, 7, 8, 9, time.Local)),
					EndTime:   db.PGTimestamp(time.Date(2025, 4, 5, 6, 7, 8, 9, time.Local)),
					Type:      667,
					Faction:   "Terminids",
					MaxHealth: 4455667788,
				},
			},
			wantErr: false,
		},
		{
			name: "nil event",
			modifier: func(p *api.Planet) {
				p.Event = nil
			},
			want:    []db.EntityMerger{},
			wantErr: false,
		},
		{
			name: "empty ID",
			modifier: func(p *api.Planet) {
				event := validEvent
				event.Id = nil
				p.Event = mustPlanetEvent(event)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planets := []api.Planet{
				{
					Event: mustPlanetEvent(validEvent),
				},
			}
			// call modifiers on valid copies
			tt.modifier(&planets[0])
			data := APIData{
				Planets: &planets,
			}
			got, err := Events(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Events() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Events() = %v, want %v", got, tt.want)
				t.Log(got[0])
				t.Log(tt.want[0])
				return
			}
		})
	}
}
