package transform

import (
	"reflect"
	"testing"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

func mustPlanetPosition(from api.Position) *api.Planet_Position {
	planetPosition := new(api.Planet_Position)
	if err := planetPosition.FromPosition(from); err != nil {
		panic(err)
	}
	return planetPosition
}

// TODO: generic
func mustPlanetBiome(from api.Biome) *api.Planet_Biome {
	planetBiome := new(api.Planet_Biome)
	if err := planetBiome.FromBiome(from); err != nil {
		panic(err)
	}
	return planetBiome
}

func TestPlanetsTransform(t *testing.T) {
	type args struct {
		data APIData
	}
	tests := []struct {
		name    string
		p       Planets
		args    args
		want    *db.DocsProvider[structs.Planet]
		wantErr bool
	}{
		{
			name: "complete",
			args: args{
				data: APIData{
					Planets: &[]api.Planet{
						{
							Index:  ptr(int32(5)),
							Name:   ptr("Foo"),
							Sector: ptr("Bar"),
							Position: mustPlanetPosition(api.Position{
								X: ptr(float64(3)),
								Y: ptr(float64(5)),
							}),
							Waypoints: &[]int32{4, 5, 6},
							Disabled:  ptr(false),
							Biome: mustPlanetBiome(api.Biome{
								Name:        ptr("Foobiome"),
								Description: ptr("Bardescription"),
							}),
							Hazards: &[]api.Hazard{
								{
									Name:        ptr("Foohazard"),
									Description: ptr("Barhazard"),
								},
							},
							MaxHealth:      ptr(int64(1000)),
							InitialOwner:   ptr("Automatons"),
							RegenPerSecond: ptr(float64(0.6)),
						},
					},
				},
			},
			want: &db.DocsProvider[structs.Planet]{
				CollectionName: "planets",
				Docs: []db.DocWrapper[structs.Planet]{
					{
						DocID: int32(5),
						Document: structs.Planet{
							ID:        5,
							Name:      "Foo",
							Sector:    "Bar",
							Position:  structs.PlanetPosition{X: 3, Y: 5},
							Waypoints: []int32{4, 5, 6},
							Disabled:  false,
							Biome: structs.Biome{
								Name:        "Foobiome",
								Description: "Bardescription",
							},
							Hazards: []structs.Hazard{
								{
									Name:        "Foohazard",
									Description: "Barhazard",
								},
							},
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
				data: APIData{
					Planets: &[]api.Planet{
						{
							Index:  ptr(int32(5)),
							Name:   ptr("Foo"),
							Sector: ptr("Bar"),
							Position: mustPlanetPosition(api.Position{
								X: ptr(float64(3)),
								Y: ptr(float64(5)),
							}),
							Biome: mustPlanetBiome(api.Biome{
								Name:        ptr("Foobiome"),
								Description: ptr("Bardescription"),
							}),
							Hazards: &[]api.Hazard{
								{
									Name:        ptr("Foohazard"),
									Description: ptr("Barhazard"),
								},
							},
							Waypoints:      &[]int32{4, 5, 6},
							Disabled:       nil,
							MaxHealth:      ptr(int64(1000)),
							InitialOwner:   ptr("Automatons"),
							RegenPerSecond: ptr(float64(0.6)),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil embedded",
			args: args{
				data: APIData{
					Planets: &[]api.Planet{
						{
							Index:  ptr(int32(5)),
							Name:   ptr("Foo"),
							Sector: ptr("Bar"),
							Position: mustPlanetPosition(api.Position{
								X: nil,
								Y: ptr(float64(5)),
							}),
							Biome: mustPlanetBiome(api.Biome{
								Name:        ptr("Foobiome"),
								Description: ptr("Bardescription"),
							}),
							Hazards: &[]api.Hazard{
								{
									Name:        ptr("Foohazard"),
									Description: ptr("Barhazard"),
								},
							},
							Waypoints:      &[]int32{4, 5, 6},
							Disabled:       ptr(false),
							MaxHealth:      ptr(int64(1000)),
							InitialOwner:   ptr("Automatons"),
							RegenPerSecond: ptr(float64(0.6)),
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
