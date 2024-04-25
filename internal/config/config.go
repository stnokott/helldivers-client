// Package config handles configuration via environment variables
package config

import (
	"fmt"
	"os"
	"time"

	env "go-simpler.org/env"
)

const envPrefix = "HELL"

// Config contains configuration values
type Config struct {
	PostgresURI    string        `env:"POSTGRES_URI,required" usage:"URI to MongoDB host. Example: postgresql://user:pass@localhost:5432/database"`
	APIRootURL     string        `env:"API_URL,required" usage:"Root URL of Helldivers 2 API. Example: http://localhost:4000"`
	WorkerInterval time.Duration `env:"WORKER_INTERVAL" default:"5m" usage:"Interval at which data will be queried from the API and written to the database."`
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() *Config {

	c := new(Config)
	if err := env.Load(c, nil); err != nil {
		fmt.Printf("failed to read config from ENV: %v\n", err)
		fmt.Println("Usage:")
		env.Usage(c, os.Stdout)
		os.Exit(1)
	}
	return c
}
