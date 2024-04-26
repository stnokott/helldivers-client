package client

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var defaultBackoff = 5 * time.Second

// rateLimitHTTPClient defines a non-thread safe HTTP client which can automatically deal with rate limit headers.
type rateLimitHTTPClient struct {
	client   *http.Client
	maxRetry int
	log      *log.Logger
}

func newRateLimitHTTPClient(maxRetry int, logger *log.Logger) *rateLimitHTTPClient {
	return &rateLimitHTTPClient{
		client:   http.DefaultClient,
		maxRetry: maxRetry,
		log:      logger,
	}
}

func (c *rateLimitHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.log.Printf("requesting %s", req.URL.Path)

	// start retry loop
	for n := c.maxRetry; n >= 0; n-- {
		resp, err := c.doIter(req)
		if resp != nil {
			return resp, nil
		}
		if err != nil {
			return nil, err
		}
		if n > 0 {
			c.log.Printf("will retry %d more times", n)
		}
		// no response, but also no err, so we go to the next iteration
	}
	return nil, fmt.Errorf("no valid response after %d retries", c.maxRetry)
}

// doIter performs a retry iteration.
// It can have three different outcomes:
//  1. successful request -> returns response, nil err
//  2. unsuccessful request -> waits if 429, then returns nil response, nil err
//  3. context cancelled while waiting -> returns nil response, err
func (c *rateLimitHTTPClient) doIter(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		return resp, nil
	}
	if err != nil {
		c.log.Printf("** HTTP error: %v", err)
		return nil, err
	}
	c.log.Printf("** HTTP error: %s", resp.Status)
	if resp.StatusCode != http.StatusTooManyRequests {
		// no wait needed, immediately return
		return nil, nil
	}

	retryAfter := c.retryAfter(resp.Header)
	c.log.Printf("** retrying after %s", retryAfter.String())
	backoffTimer := time.NewTimer(retryAfter)
	ctx := req.Context()
	select {
	case <-backoffTimer.C:
		return nil, nil
	case <-ctx.Done():
		if !backoffTimer.Stop() {
			// drain channel
			<-backoffTimer.C
		}
		return nil, ctx.Err()
	}
}

func (c *rateLimitHTTPClient) retryAfter(header http.Header) time.Duration {
	var backoff time.Duration
	retryAfter := header.Get("Retry-After")
	if retryAfter == "" {
		c.log.Println("got no 'Retry-After' response header, using default backoff")
		backoff = defaultBackoff
	} else {
		parsed, err := strconv.Atoi(retryAfter)
		if err != nil {
			c.log.Printf("got invalid 'Retry-After' response header (%v), using default backoff", err)
			backoff = defaultBackoff
		} else {
			backoff = time.Duration(parsed) * time.Second
		}
	}
	return backoff
}
