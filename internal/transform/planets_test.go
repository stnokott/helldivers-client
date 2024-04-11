package transform

import (
	"reflect"
	"testing"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

func mustPlanetPosition(rawJSON []byte) *api.Planet_Position {
	planetPosition := new(api.Planet_Position)
	if err := planetPosition.UnmarshalJSON(rawJSON); err != nil {
		panic(err)
	}
	return planetPosition
}

func TestPlanetsTransform(t *testing.T) {
	type args struct {
		data planetsRequestData
	}
	tests := []struct {
		name    string
		p       Planets
		args    args
		want    *db.DocsProvider
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: &[]api.Planet{
					{
						Index:          ptr(int32(5)),
						Name:           ptr("Foo"),
						Sector:         ptr("Bar"),
						Position:       mustPlanetPosition([]byte(`{"X": 3, "Y": 5}`)),
						Waypoints:      &[]int32{4, 5, 6},
						Disabled:       ptr(false),
						MaxHealth:      ptr(int64(1000)),
						InitialOwner:   ptr("Automatons"),
						RegenPerSecond: ptr(float64(0.6)),
					},
				},
			},
			want: &db.DocsProvider{
				CollectionName: "planets",
				Docs: []db.DocWrapper{
					{
						DocID: int32(5),
						Document: structs.Planet{
							ID:             5,
							Name:           "Foo",
							Sector:         "Bar",
							Position:       structs.PlanetPosition{X: 3, Y: 5},
							Waypoints:      []int32{4, 5, 6},
							Disabled:       false,
							MaxHealth:      1000,
							InitialOwner:   "Automatons",
							RegenPerSecond: 0.6,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil disabled",
			args: args{
				data: &[]api.Planet{
					{
						Index:          ptr(int32(5)),
						Name:           ptr("Foo"),
						Sector:         ptr("Bar"),
						Position:       mustPlanetPosition([]byte(`{"X": 3, "Y": 5}`)),
						Waypoints:      &[]int32{4, 5, 6},
						Disabled:       nil,
						MaxHealth:      ptr(int64(1000)),
						InitialOwner:   ptr("Automatons"),
						RegenPerSecond: ptr(float64(0.6)),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid position",
			args: args{
				data: &[]api.Planet{
					{
						Index:          ptr(int32(5)),
						Name:           ptr("Foo"),
						Sector:         ptr("Bar"),
						Position:       mustPlanetPosition([]byte(`{"X": "3.5", "Y": 5}`)),
						Waypoints:      &[]int32{4, 5, 6},
						Disabled:       ptr(false),
						MaxHealth:      ptr(int64(1000)),
						InitialOwner:   ptr("Automatons"),
						RegenPerSecond: ptr(float64(0.6)),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Transform(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Planets.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planets.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
