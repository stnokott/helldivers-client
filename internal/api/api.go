// Package api contains interfaces for the Helldivers 2 API, generated from OpenAPI spec
package api

import (
	"fmt"
	"io"
	"net/http"
)

//go:generate go run github.com/deepmap/oapi-codegen@v2.1.0 --config=oapi-codegen.cfg.yaml http://localhost:4000/Helldivers-2-API.json

// Data implements the interface required for client processing
func (resp *GetRawApiWarSeasonCurrentWarIDResponse) Data() (*WarId, error) {
	if resp.StatusCode() == 200 {
		return resp.JSON200, nil
	}
	return nil, respErr(resp.HTTPResponse)
}

// Data implements the interface required for client processing
func (resp *GetApiV1WarResponse) Data() (*War, error) {
	if resp.StatusCode() == 200 {
		return resp.JSON200, nil
	}
	return nil, respErr(resp.HTTPResponse)
}

// Data implements the interface required for client processing
func (resp *GetApiV1AssignmentsAllResponse) Data() (*[]Assignment2, error) {
	if resp.StatusCode() == 200 {
		return resp.JSON200, nil
	}
	return nil, respErr(resp.HTTPResponse)
}

// Data implements the interface required for client processing
func (resp *GetApiV1CampaignsAllResponse) Data() (*[]Campaign2, error) {
	if resp.StatusCode() == 200 {
		return resp.JSON200, nil
	}
	return nil, respErr(resp.HTTPResponse)
}

// Data implements the interface required for client processing
func (resp *GetApiV1DispatchesAllResponse) Data() (*[]Dispatch, error) {
	if resp.StatusCode() == 200 {
		return resp.JSON200, nil
	}
	return nil, respErr(resp.HTTPResponse)
}

// Data implements the interface required for client processing
func (resp *GetApiV1PlanetsAllResponse) Data() (*[]Planet, error) {
	if resp.StatusCode() == 200 {
		return resp.JSON200, nil
	}
	return nil, respErr(resp.HTTPResponse)
}

func respErr(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		body = []byte("n/a")
	}
	_ = resp.Body.Close()
	return fmt.Errorf("HTTP status %s: '%s'", resp.Status, string(body))
}
