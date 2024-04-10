// Package client wraps the API specs into a client
package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
)

// Client wraps the generated OpenAPI client
type Client struct {
	api *api.ClientWithResponses
}

const (
	rateLimitDuration = 10 * time.Second
	rateLimitCount    = 5
)

// New creates a new client instance
func New(rootURI string, logger *log.Logger) (*Client, error) {
	options := api.WithHTTPClient(
		newRateLimitHTTPClient(rateLimitDuration, rateLimitCount, logger),
	)
	c, err := api.NewClientWithResponses(rootURI, options)
	if err != nil {
		return nil, fmt.Errorf("client initialization failed: %w", err)
	}

	return &Client{
		api: c,
	}, nil
}

func processResp[
	T any,
	PT interface{ Data() (*T, error) },
](
	ctx context.Context,
	requestFunc func(context.Context, ...api.RequestEditorFn) (PT, error),
) (*T, error) {
	resp, err := requestFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	data, err := resp.Data()
	if err != nil {
		return nil, fmt.Errorf("request response unavailable: %w", err)
	}

	return data, nil
}

// WarID returns the ID of the current war
func (c *Client) WarID(ctx context.Context) (*api.WarId, error) {
	return processResp(ctx, c.api.GetRawApiWarSeasonCurrentWarIDWithResponse)
}

// War returns the current war
func (c *Client) War(ctx context.Context) (*api.War, error) {
	return processResp(ctx, c.api.GetApiV1WarWithResponse)
}

// Assignments returns all currently active assignments
func (c *Client) Assignments(ctx context.Context) (*[]api.Assignment2, error) {
	return processResp(ctx, c.api.GetApiV1AssignmentsAllWithResponse)
}

// Campaigns returns all currently active campaigns
func (c *Client) Campaigns(ctx context.Context) (*[]api.Campaign2, error) {
	return processResp(ctx, c.api.GetApiV1CampaignsAllWithResponse)
}

// Dispatches returns all currently active dispatches
func (c *Client) Dispatches(ctx context.Context) (*[]api.Dispatch, error) {
	return processResp(ctx, c.api.GetApiV1DispatchesAllWithResponse)
}

// Planets returns all planets in the current war
func (c *Client) Planets(ctx context.Context) (*[]api.Planet, error) {
	return processResp(ctx, c.api.GetApiV1PlanetsAllWithResponse)
}
