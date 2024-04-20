package db

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

var validPlanet = Planet{
	Planet: gen.Planet{
		ID:           1,
		Name:         "Foo",
		Sector:       "Bar",
		Position:     []float64{1, 2},
		WaypointIds:  []int32{1, 2, 3},
		Disabled:     false,
		BiomeName:    "FooBiome",
		HazardNames:  []string{"BarHazard"},
		MaxHealth:    1000,
		InitialOwner: "Super Humans",
	},
	Biome: Biome{
		Name:        "FooBiome",
		Description: "This biome contains a lot of spaghetti",
	},
	Hazards: []Hazard{
		{
			Name:        "BarHazard",
			Description: "This hazard contains a lot of bugs",
		},
	},
}

func TestPlanetsSchema(t *testing.T) {
	// modifier applies a change to the valid struct, based on the test
	type modifier func(*Planet)
	tests := []struct {
		name          string
		modifier      modifier
		wantBiomeErr  bool
		wantHazardErr bool
		wantPlanetErr bool
	}{
		{
			name:          "valid",
			modifier:      func(p *Planet) {},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: false,
		},
		{
			name: "empty required sector",
			modifier: func(p *Planet) {
				p.Sector = ""
			},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: true,
		},
		{
			name: "empty foreign key biome name",
			modifier: func(p *Planet) {
				p.BiomeName = ""
			},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: true,
		},
		{
			name: "position Y coordinate missing",
			modifier: func(p *Planet) {
				p.Position = []float64{5}
			},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: true,
		},
		{
			name: "position has too many coordinates",
			modifier: func(p *Planet) {
				p.Position = []float64{3, 4, 5}
			},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: true,
		},
		{
			name: "biome foreign key not existing",
			modifier: func(p *Planet) {
				p.Biome.Name = "a different biome"
			},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: true,
		},
		{
			name: "hazard foreign key not existing",
			modifier: func(p *Planet) {
				p.Hazards = []Hazard{}
			},
			wantBiomeErr:  false,
			wantHazardErr: false,
			wantPlanetErr: true,
		},
		{
			name: "biome name empty",
			modifier: func(p *Planet) {
				p.Biome.Name = ""
			},
			wantBiomeErr: true,
		},
		{
			name: "biome description empty",
			modifier: func(p *Planet) {
				p.Biome.Description = ""
			},
			wantBiomeErr: true,
		},
		{
			name: "hazard name empty",
			modifier: func(p *Planet) {
				p.Hazards[0].Name = ""
			},
			wantBiomeErr:  false,
			wantHazardErr: true,
		},
		{
			name: "hazard description empty",
			modifier: func(p *Planet) {
				p.Hazards[0].Description = ""
			},
			wantBiomeErr:  false,
			wantHazardErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withClient(t, func(client *Client, migration *migrate.Migrate) {
				if err := migration.Up(); err != nil {
					t.Errorf("failed to migrate up: %v", err)
					return
				}
				planet := validPlanet
				tt.modifier(&planet)

				_, err := client.queries.InsertBiome(context.Background(), gen.InsertBiomeParams(planet.Biome))
				if (err != nil) != tt.wantBiomeErr {
					t.Logf("InsertBiome() error = %v, wantBiomeErr = %v", err, tt.wantBiomeErr)
					return
				}

				for _, hazard := range planet.Hazards {
					_, err = client.queries.InsertHazard(context.Background(), gen.InsertHazardParams(hazard))
					if (err != nil) != tt.wantHazardErr {
						t.Logf("InsertHazard() error = %v, wantHazardErr = %v", err, tt.wantHazardErr)
						return
					}
				}

				_, err = client.queries.InsertPlanet(context.Background(), gen.InsertPlanetParams(planet.Planet))
				if (err != nil) != tt.wantPlanetErr {
					t.Errorf("InsertPlanet() error = %v, wantErr = %v", err, tt.wantPlanetErr)
					return
				}
				if err != nil {
					return
				}
				fetchedResult, err := client.queries.GetPlanet(context.Background(), planet.ID)
				if err != nil {
					t.Errorf("failed to fetch inserted planet: %v", err)
					return
				}
				if fetchedResult != planet.ID {
					t.Errorf("failed to validate INSERT: inserted data has ID %d, DB returned %d", planet.ID, fetchedResult)
				}
			})
		})
	}
}
