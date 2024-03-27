// Package client wraps the API specs into a client
package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
)

// Client wraps the generated OpenAPI client
type Client struct {
	c *api.ClientWithResponses
}

// New creates a new client instance
func New(host string) (*Client, error) {
	c, err := api.NewClientWithResponses(host)
	if err != nil {
		return nil, fmt.Errorf("client initialization failed: %w", err)
	}

	return &Client{
		c: c,
	}, nil
}

func respBodyErr(body []byte) error {
	return errors.New("unknown error occured: " + string(body))
}

// Seasons returns the overview of current and past seasons
func (c *Client) Seasons(ctx context.Context) (*api.WarSeasonOverview, error) {
	resp, err := c.c.Helldivers2WebApiWarSeasonControllerIndexWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("request initialization failed: %w", err)
	}
	if resp.JSON200 != nil {
		return resp.JSON200, nil
	}
	if resp.JSON429 != nil && resp.JSON429.Error != nil {
		return nil, fmt.Errorf("rate limit exceeded: %s", *resp.JSON429.Error)
	}
	return nil, respBodyErr(resp.Body)
}
