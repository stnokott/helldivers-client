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
	MongoURI       string         `envconfig:"MONGODB_URI" required:"true"`
	APIRootURL     string         `envconfig:"API_URL" required:"true"`
	WorkerInterval WorkerInterval `envconfig:"WORKER_INTERVAL" default:"5m"`
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() Config {
	var c Config
	if err := envconfig.Process(envPrefix, &c); err != nil {
		log.Fatalf("failed to read config from ENV: %v", err)
	}
	return c
}

type WorkerInterval time.Duration

func (id *WorkerInterval) Decode(value string) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*id = WorkerInterval(d)
	return nil
}

func (id *WorkerInterval) String() string {
	return time.Duration(*id).String()
}
