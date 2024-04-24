package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db"
)

func mustCampaignPlanet(from api.Planet) *api.Campaign2_Planet {
	campaignPlanet := new(api.Campaign2_Planet)
	if err := campaignPlanet.FromPlanet(from); err != nil {
		panic(err)
	}
	return campaignPlanet
}

func mustPlanetEvent(from api.Event) *api.Planet_Event {
	planetEvent := new(api.Planet_Event)
	if err := planetEvent.FromEvent(from); err != nil {
		panic(err)
	}
	return planetEvent
}

var validEventPlanet = api.Planet{
	Index: ptr(int32(888)),
	Name:  ptr("A planet"),
	Event: mustPlanetEvent(
		api.Event{
			Id:                ptr(int32(997)),
			StartTime:         ptr(time.Date(2024, 4, 5, 6, 7, 8, 9, time.UTC)),
			EndTime:           ptr(time.Date(2025, 4, 5, 6, 7, 8, 9, time.UTC)),
			EventType:         ptr(int32(667)),
			Faction:           ptr("Terminids"),
			MaxHealth:         ptr(int64(4455667788)),
			CampaignId:        ptr(int64(123)),
			Health:            nil, // not required, persisted in dynamic snapshots
			JointOperationIds: nil, // not mapped currently
		},
	),
}

var validEventCampaign = api.Campaign2{
	Id: ptr(int32(123)),
	Planet: mustCampaignPlanet(api.Planet{
		Index: ptr(int32(888)),
		Name:  ptr("Foo"),
	}),
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
					ID:         997,
					StartTime:  db.PGTimestamp(time.Date(2024, 4, 5, 6, 7, 8, 9, time.UTC)),
					EndTime:    db.PGTimestamp(time.Date(2025, 4, 5, 6, 7, 8, 9, time.UTC)),
					Type:       667,
					Faction:    "Terminids",
					MaxHealth:  4455667788,
					CampaignID: 123,
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
				event, _ := p.Event.AsEvent()
				event.Id = nil
				p.Event = mustPlanetEvent(event)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var planet api.Planet
			if err := copytest.DeepCopy(&planet, &validEventPlanet); err != nil {
				t.Errorf("failed to create event struct copy: %v", err)
				return
			}

			planets := []api.Planet{validEventPlanet}

			tt.modifier(&planets[0])

			data := APIData{
				Planets: &planets,
				Campaigns: &[]api.Campaign2{
					validEventCampaign, // campaign remains static
				},
			}
			got, err := Events(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Events() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Events() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
