// Package config handles configuration via environment variables
package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "HELL"

// Config contains configuration values
type Config struct {
	MongoURI             string   `envconfig:"MONGODB_URI" required:"true"`
	APIRootURL           string   `envconfig:"API_URL" required:"true"`
	APIRateLimitInterval Interval `envconfig:"API_RATE_LIMIT_INTERVAL" default:"10s"`
	APIRateLimitCount    int      `envconfig:"API_RATE_LIMIT_COUNT" default:"5"`
	WorkerInterval       Interval `envconfig:"WORKER_INTERVAL" default:"5m"`
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() *Config {
	c := new(Config)
	if err := envconfig.Process(envPrefix, c); err != nil {
		log.Fatalf("failed to read config from ENV: %v", err)
	}
	return c
}

type Interval time.Duration

func (id *Interval) Decode(value string) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*id = Interval(d)
	return nil
}

func (id *Interval) String() string {
	return time.Duration(*id).String()
}
