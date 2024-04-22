package transform

import (
	"reflect"
	"testing"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

func mustCampaignPlanet(from api.Planet) *api.Campaign2_Planet {
	campaignPlanet := new(api.Campaign2_Planet)
	if err := campaignPlanet.FromPlanet(from); err != nil {
		panic(err)
	}
	return campaignPlanet
}

var validCampaign = api.Campaign2{
	Id:    ptr(int32(987)),
	Count: ptr(int32(123)),
	Type:  ptr(int32(7)),
	Planet: mustCampaignPlanet(api.Planet{
		Index: ptr(int32(3)),
	}),
}

func TestCampaigns(t *testing.T) {
	// modifier changes the valid assignment to one that is suited for the test
	type modifier func(*api.Campaign2)
	tests := []struct {
		name     string
		modifier modifier
		want     []db.EntityMerger
		wantErr  bool
	}{
		{
			name: "valid",
			modifier: func(a *api.Campaign2) {
				// keep valid
			},
			want: []db.EntityMerger{
				&db.Campaign{
					ID:       987,
					PlanetID: 3,
					Type:     7,
					Count:    123,
				},
			},
			wantErr: false,
		},
		{
			name: "missing required planet ID",
			modifier: func(c *api.Campaign2) {
				c.Planet = mustCampaignPlanet(api.Planet{
					Index: nil,
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing required type",
			modifier: func(c *api.Campaign2) {
				c.Type = nil
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := validCampaign
			// call modifier on valid assignment copy
			tt.modifier(&campaign)
			data := APIData{
				Campaigns: &[]api.Campaign2{
					campaign,
				},
			}
			got, err := Campaigns(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Campaigns() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Campaigns() = %v, want %v", got, tt.want)
			}
		})
	}
}
