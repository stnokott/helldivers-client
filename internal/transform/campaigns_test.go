package transform

import (
	"reflect"
	"testing"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

func mustCampaign2Planet(from api.Planet) *api.Campaign2_Planet {
	planet := new(api.Campaign2_Planet)
	if err := planet.FromPlanet(from); err != nil {
		panic(err)
	}
	return planet
}

func TestCampaignsTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		c       Campaigns
		args    args
		want    *db.DocsProvider[structs.Campaign]
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					Campaigns: &[]api.Campaign2{
						{
							Id: ptr(int32(5)),
							Planet: mustCampaign2Planet(api.Planet{
								Index: ptr(int32(99)),
							}),
							Type:  ptr(int32(7)),
							Count: ptr(int32(10)),
						},
						{
							Id: ptr(int32(6)),
							Planet: mustCampaign2Planet(api.Planet{
								Index: ptr(int32(42)),
							}),
							Type:  ptr(int32(8)),
							Count: ptr(int32(2)),
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Campaign]{
				CollectionName: "campaigns",
				Docs: []db.DocWrapper[structs.Campaign]{
					{
						DocID: int32(5),
						Document: structs.Campaign{
							ID:       5,
							PlanetID: 99,
							Type:     7,
							Count:    10,
						},
					},
					{
						DocID: int32(6),
						Document: structs.Campaign{
							ID:       6,
							PlanetID: 42,
							Type:     8,
							Count:    2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil planet",
			args: args{
				data: APIData{
					Campaigns: &[]api.Campaign2{
						{
							Id:     ptr(int32(5)),
							Planet: nil,
							Type:   ptr(int32(7)),
							Count:  ptr(int32(10)),
						},
					},
				},
			},
			want:    &db.DocsProvider[structs.Campaign]{
				CollectionName: db.CollCampaigns,
				Docs: []db.DocWrapper[structs.Campaign]{},
			},
			wantErr: true,
		},
		{
			name: "nil embedded",
			args: args{
				data: APIData{
					Campaigns: &[]api.Campaign2{
						{
							Id: ptr(int32(5)),
							Planet: mustCampaign2Planet(api.Planet{
								Index: nil,
							}),
							Type:  ptr(int32(7)),
							Count: ptr(int32(10)),
						},
					},
				},
			},
			want:    &db.DocsProvider[structs.Campaign]{
				CollectionName: db.CollCampaigns,
				Docs: []db.DocWrapper[structs.Campaign]{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := false
			errFunc := func(err error) {
				if !tt.wantErr {
					t.Logf("Campaigns.Transform() error: %v", err)
				}
				gotErr = true
			}
			got := tt.c.Transform(tt.args.data, errFunc)
			if gotErr != tt.wantErr {
				t.Errorf("Campaigns.Transform() returned error, wantErr %v", tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Campaigns.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
