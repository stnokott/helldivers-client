package transform

import (
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func mustPlanetPosition(from api.Position) *api.Planet_Position {
	planetPosition := new(api.Planet_Position)
	if err := planetPosition.FromPosition(from); err != nil {
		panic(err)
	}
	return planetPosition
}

func mustPlanetBiome(from api.Biome) *api.Planet_Biome {
	planetBiome := new(api.Planet_Biome)
	if err := planetBiome.FromBiome(from); err != nil {
		panic(err)
	}
	return planetBiome
}

func mustPlanetStatistics(from api.Statistics) *api.Planet_Statistics {
	planetStats := new(api.Planet_Statistics)
	if err := planetStats.FromStatistics(from); err != nil {
		panic(err)
	}
	return planetStats
}

var validPlanet = api.Planet{
	Index: ptr(int32(3)),
	Name:  ptr("A planet"),
	Biome: mustPlanetBiome(api.Biome{
		Name:        ptr("Foobiome"),
		Description: ptr("Foodescription"),
	}),
	Hazards: &[]api.Hazard{
		{
			Name:        ptr("Barhazard"),
			Description: ptr("Bardescription"),
		},
	},
	Disabled:     ptr(false),
	InitialOwner: ptr("Humans"),
	MaxHealth:    ptr(int64(112233445566)),
	Position: mustPlanetPosition(api.Position{
		X: ptr(float64(38)), Y: ptr(float64(6)),
	}),
	Sector:    ptr("A sector"),
	Waypoints: &[]int32{3, 4, 5},
}

func TestPlanets(t *testing.T) {
	// modifier changes the valid assignment to one that is suited for the test
	type modifier func(*api.Planet)
	tests := []struct {
		name     string
		modifier modifier
		want     []db.EntityMerger
		wantErr  bool
	}{
		{
			name: "valid",
			modifier: func(a *api.Planet) {
				// keep valid
			},
			want: []db.EntityMerger{
				&db.Planet{
					Planet: gen.Planet{
						ID:           3,
						Name:         "A planet",
						BiomeName:    "Foobiome",
						HazardNames:  []string{"Barhazard"},
						Disabled:     false,
						InitialOwner: "Humans",
						MaxHealth:    112233445566,
						Position:     []float64{38, 6},
						Sector:       "A sector",
						WaypointIds:  []int32{3, 4, 5},
					},
					Biome: gen.Biome{
						Name:        "Foobiome",
						Description: "Foodescription",
					},
					Hazards: []gen.Hazard{
						{
							Name:        "Barhazard",
							Description: "Bardescription",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty waypoints",
			modifier: func(p *api.Planet) {
				p.Waypoints = &[]int32{}
			},
			want: []db.EntityMerger{
				&db.Planet{
					Planet: gen.Planet{
						ID:           3,
						Name:         "A planet",
						BiomeName:    "Foobiome",
						HazardNames:  []string{"Barhazard"},
						Disabled:     false,
						InitialOwner: "Humans",
						MaxHealth:    112233445566,
						Position:     []float64{38, 6},
						Sector:       "A sector",
						WaypointIds:  []int32{},
					},
					Biome: gen.Biome{
						Name:        "Foobiome",
						Description: "Foodescription",
					},
					Hazards: []gen.Hazard{
						{
							Name:        "Barhazard",
							Description: "Bardescription",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty required name",
			modifier: func(p *api.Planet) {
				p.Name = nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty required biome",
			modifier: func(p *api.Planet) {
				p.Biome = nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty required hazards",
			modifier: func(p *api.Planet) {
				p.Hazards = nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty embedded position coordinate",
			modifier: func(p *api.Planet) {
				p.Position = mustPlanetPosition(api.Position{
					X: ptr(float64(7)), Y: nil,
				})
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var planet api.Planet
			// deep copy will copy values behind pointers instead of the pointers themselves
			copyOption := copier.Option{DeepCopy: true}
			if err := copier.CopyWithOption(&planet, &validPlanet, copyOption); err != nil {
				t.Errorf("failed to create planet struct copy: %v", err)
				return
			}
			// call modifier on valid planet copy
			tt.modifier(&planet)
			data := APIData{
				Planets: &[]api.Planet{
					planet,
				},
			}
			got, err := Planets(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Planets() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planets() = %v, want %v", got, tt.want)
			}
		})
	}
}
