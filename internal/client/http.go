package client

import (
	"errors"
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

func newRateLimitHTTPClient(rl *rate.Limiter, logger *log.Logger) *rateLimitHTTPClient {
	return &rateLimitHTTPClient{
		client: http.DefaultClient,
		rl:     rl,
		log:    logger,
	}
}

func (c *rateLimitHTTPClient) Do(req *http.Request) (*http.Response, error) {
	reserved := c.rl.Reserve()
	if !reserved.OK() {
		return nil, errors.New("rate limiter configured incorrectly")
	}
	delay := reserved.Delay()
	if delay > 0 {
		c.log.Print("WARN: rate limit exceeded, waiting")
		time.Sleep(delay)
	}
	c.log.Printf("requesting %s", req.URL.Path)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
