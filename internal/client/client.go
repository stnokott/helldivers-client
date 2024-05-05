// Package client wraps the API specs into a client
package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/config"
)

// Client wraps the generated OpenAPI client
type Client struct {
	api *api.ClientWithResponses
	log *log.Logger
}

const _maxHTTPRetries = 3

// New creates a new client instance
func New(cfg *config.Config, logger *log.Logger) (*Client, error) {
	options := api.WithHTTPClient(
		newRateLimitHTTPClient(_maxHTTPRetries, logger),
	)
	c, err := api.NewClientWithResponses(cfg.APIRootURL, options)
	if err != nil {
		return nil, fmt.Errorf("client initialization: %w", err)
	}

	return &Client{
		api: c,
		log: logger,
	}, nil
}

// Connect implements main.ConnectWaiter.
func (c *Client) Connect(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if _, err := c.WarID(ctx); err == nil {
				return nil
			} else {
				c.log.Printf("API query: %v", err)
			}
		}
	}
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
