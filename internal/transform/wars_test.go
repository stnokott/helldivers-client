package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db"
)

var validWarID = api.WarId{
	Id: ptr(int32(999)),
}

var validWar = api.War{
	Started: ptr(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
	Ended:   ptr(time.Date(2025, 1, 1, 1, 1, 1, 1, time.UTC)),
	Factions: &[]string{
		"Humans",
		"Automatons",
	},
}

func TestWar(t *testing.T) {
	// warIDModifier changes the valid war to one that is suited for the test
	type warIDModifier func(*api.WarId)
	type warModifier func(*api.War)
	tests := []struct {
		name          string
		warIDModifier warIDModifier
		warModifier   warModifier
		want          []db.EntityMerger
		wantErr       bool
	}{
		{
			name: "valid",
			warIDModifier: func(wi *api.WarId) {
				// keep valid
			},
			warModifier: func(a *api.War) {
				// keep valid
			},
			want: []db.EntityMerger{
				&db.War{
					ID:        999,
					StartTime: db.PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
					EndTime:   db.PGTimestamp(time.Date(2025, 1, 1, 1, 1, 1, 1, time.UTC)),
					Factions:  []string{"Humans", "Automatons"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty factions",
			warIDModifier: func(wi *api.WarId) {
				// keep valid
			},
			warModifier: func(w *api.War) {
				w.Factions = &[]string{}
			},
			want: []db.EntityMerger{
				&db.War{
					ID:        999,
					StartTime: db.PGTimestamp(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
					EndTime:   db.PGTimestamp(time.Date(2025, 1, 1, 1, 1, 1, 1, time.UTC)),
					Factions:  []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "empty required war ID",
			warIDModifier: func(wi *api.WarId) {
				wi.Id = nil
			},
			warModifier: func(w *api.War) {
				// keep valid
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty required start time",
			warIDModifier: func(wi *api.WarId) {
				// keep valid
			},
			warModifier: func(w *api.War) {
				w.Started = nil
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				warID api.WarId
				war   api.War
			)
			if err := copytest.DeepCopy(
				&warID, &validWarID,
				&war, &validWar,
			); err != nil {
				t.Errorf("failed to create struct copies: %v", err)
				return
			}
			// call modifiers on valid copies
			tt.warIDModifier(&warID)
			tt.warModifier(&war)
			data := APIData{
				WarID: &warID,
				War:   &war,
			}
			got, err := Wars(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("War() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("War() = %v, want %v", got, tt.want)
			}
		})
	}
}
