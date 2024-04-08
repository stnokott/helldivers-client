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

// CurrentWar returns the current war season
func (c *Client) CurrentWar(ctx context.Context) (*api.War, error) {
	resp, err := c.c.GetApiV1WarWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("request initialization failed: %w", err)
	}
	if resp.JSON200 != nil {
		return resp.JSON200, nil
	}
	return nil, respBodyErr(resp.Body)
}
