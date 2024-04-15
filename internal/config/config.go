// Package config handles configuration via environment variables
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "HELL"

// Config contains configuration values
type Config struct {
	MongoURI             string   `envconfig:"MONGODB_URI" required:"true" desc:"URI to MongoDB host. Example: mongodb://user:pass@localhost:27017"`
	APIRootURL           string   `envconfig:"API_URL" required:"true" desc:"Root URL of Helldivers 2 API. Example: http://localhost:4000"`
	APIRateLimitInterval Interval `envconfig:"API_RATE_LIMIT_INTERVAL" default:"10s" desc:"Interval of API rate limit. Requests will wait if internal rate limit is exceeded."`
	APIRateLimitCount    int      `envconfig:"API_RATE_LIMIT_COUNT" default:"5" desc:"Allowed requests per interval of API rate limit. Requests will wait if internal rate limit is exceeded."`
	WorkerInterval       Interval `envconfig:"WORKER_INTERVAL" default:"5m" desc:"Interval at which data will be queried from the API and written to the database."`
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() *Config {
	c := new(Config)
	if err := envconfig.Process(envPrefix, c); err != nil {
		fmt.Printf("failed to read config from ENV: %v\n\n", err)
		envconfig.Usage(envPrefix, c) // nolint:errcheck
		os.Exit(1)
	}
	return c
}

// Interval is a duration implementing the envconfig.Decoder interface
type Interval time.Duration

// Decode implements the envconfig.Decoder interface
func (id *Interval) Decode(value string) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*id = Interval(d)
	return nil
}

// String implement the Stringer interface
func (id *Interval) String() string {
	return time.Duration(*id).String()
}
