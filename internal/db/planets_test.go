package db

import (
	"context"
	"math"
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
		name     string
		modifier modifier
		wantErr  bool
	}{
		{
			name:     "valid",
			modifier: func(p *Planet) {},
			wantErr:  false,
		},
		{
			name: "empty required sector",
			modifier: func(p *Planet) {
				p.Sector = ""
			},
			wantErr: true,
		},
		{
			name: "negative max health",
			modifier: func(p *Planet) {
				p.MaxHealth = -1
			},
			wantErr: true,
		},
		{
			name: "high max health",
			modifier: func(p *Planet) {
				p.MaxHealth = math.MaxInt64
			},
			wantErr: false,
		},
		{
			name: "position Y coordinate missing",
			modifier: func(p *Planet) {
				p.Position = []float64{5}
			},
			wantErr: true,
		},
		{
			name: "position has too many coordinates",
			modifier: func(p *Planet) {
				p.Position = []float64{3, 4, 5}
			},
			wantErr: true,
		},
		{
			name: "empty foreign key biome name",
			modifier: func(p *Planet) {
				p.BiomeName = ""
			},
			// biome name in planet will be changed to inserted planet, so should not produce error
			wantErr: false,
		},
		{
			name: "biome foreign key different",
			modifier: func(p *Planet) {
				p.Biome.Name = "a different biome"
			},
			// same as above
			wantErr: false,
		},
		{
			name: "hazard foreign key not existing",
			modifier: func(p *Planet) {
				p.Hazards = []Hazard{}
			},
			// same as for biomes
			wantErr: false,
		},
		{
			name: "biome name empty",
			modifier: func(p *Planet) {
				p.Biome.Name = ""
			},
			wantErr: true,
		},
		{
			name: "biome description empty",
			modifier: func(p *Planet) {
				p.Biome.Description = ""
			},
			wantErr: true,
		},
		{
			name: "hazard name empty",
			modifier: func(p *Planet) {
				p.Hazards[0].Name = ""
			},
			wantErr: true,
		},
		{
			name: "hazard description empty",
			modifier: func(p *Planet) {
				p.Hazards[0].Description = ""
			},
			wantErr: true,
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

				err := planet.Merge(context.Background(), client.queries, tableMergeStats{})
				if (err != nil) != tt.wantErr {
					t.Errorf("Planet.Merge() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if err != nil {
					// any subsequent tests don't make sense if error encountered
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
