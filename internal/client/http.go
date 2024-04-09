package client

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type rateLimitHTTPClient struct {
	client *http.Client
	rl     *rate.Limiter
	log    *log.Logger
}

func newRateLimitHTTPClient(d time.Duration, n int, logger *log.Logger) *rateLimitHTTPClient {
	return &rateLimitHTTPClient{
		client: http.DefaultClient,
		rl:     rate.NewLimiter(rate.Every(d), n),
		log:    logger,
	}
}

func (c *rateLimitHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	if !c.rl.Allow() {
		c.log.Print("WARN: rate limit exceeded, waiting")
	}
	err := c.rl.Wait(ctx)
	if err != nil {
		return nil, err
	}
	c.log.Println(req.URL.Path)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
