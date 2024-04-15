// Package config handles configuration via environment variables
package config

import (
	"log"
	"os"
	"time"

	env "go-simpler.org/env"
)

const envPrefix = "HELL"

// Config contains configuration values
type Config struct {
	MongoURI             string        `env:"MONGO_URI,required" usage:"URI to MongoDB host. Example: mongodb://user:pass@localhost:27017"`
	APIRootURL           string        `env:"API_URL,required" usage:"Root URL of Helldivers 2 API. Example: http://localhost:4000"`
	APIRateLimitInterval time.Duration `env:"API_RATE_LIMIT_INTERVAL" default:"10s" usage:"Interval of API rate limit. Requests will wait if internal rate limit is exceeded."`
	APIRateLimitCount    int           `env:"API_RATE_LIMIT_COUNT" default:"5" usage:"Allowed requests per interval of API rate limit. Requests will wait if internal rate limit is exceeded."`
	WorkerInterval       time.Duration `env:"WORKER_INTERVAL" default:"5m" usage:"Interval at which data will be queried from the API and written to the database."`
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() *Config {

	c := new(Config)
	if err := env.Load(c, nil); err != nil {
		log.Printf("failed to read config from ENV: %v", err)
		log.Println("Usage:")
		env.Usage(c, log.Default().Writer())
		os.Exit(1)
	}
	return c
}
