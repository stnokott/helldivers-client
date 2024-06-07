// Package config handles configuration via environment variables
package config

import (
	"fmt"
	"os"

	"go-simpler.org/env"
)

// Config contains configuration values
type Config struct {
	PostgresURI     string `env:"POSTGRES_URI,required" usage:"URI to MongoDB host. Example: postgresql://user:pass@localhost:5432/database"`
	APIRootURL      string `env:"API_URL,required" usage:"Root URL of Helldivers 2 API. Example: http://localhost:4000"`
	WorkerCron      string `env:"WORKER_CRON" default:"*/5 * * * *" usage:"Cron expression defining the interval at which data will be queried from the API and written to the database."`
	HealthchecksURL string `env:"HEALTHCHECKS_URL" default:"" usage:"Root URL of healthchecks.io endpoint."`
}

// MustGet reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the
// application to print usage and exit.
func MustGet() *Config {
	c, err := Get()
	if err != nil {
		fmt.Printf("failed to read config from ENV: %v\n", err)
		fmt.Println("Usage:")
		env.Usage(&Config{}, os.Stdout, nil)
		os.Exit(1)
	}
	return c
}

// Get reads environment variables and parses them into a Config struct.
func Get() (*Config, error) {
	c := new(Config)
	if err := env.Load(c, nil); err != nil {
		return nil, err
	}
	return c, nil
}
