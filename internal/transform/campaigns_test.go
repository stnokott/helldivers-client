//go:build !goverter

package transform

import (
	"reflect"
	"testing"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/copytest"
	"github.com/stnokott/helldivers-client/internal/db"
)

var validCampaign = api.Campaign2{
	Id:    ptr(int32(987)),
	Count: ptr(uint64(123)),
	Type:  ptr(int32(7)),
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
					ID:    987,
					Type:  7,
					Count: db.PGUint64(123),
				},
			},
			wantErr: false,
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
			var campaign api.Campaign2
			if err := copytest.DeepCopy(&campaign, &validCampaign); err != nil {
				t.Errorf("failed to create campaign struct copy: %v", err)
				return
			}
			// call modifier on valid assignment copy
			tt.modifier(&campaign)
			data := APIData{
				Campaigns: &[]api.Campaign2{
					campaign,
				},
			}
			converter := &ConverterImpl{}
			got, err := Campaigns(converter, data)
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
