// Package config handles configuration via environment variables
package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "HELL"

// Config contains configuration values
type Config struct {
	MongoURI   string `envconfig:"MONGODB_URI" required:"true"`
	APIRootURL string `envconfig:"API_URL" required:"true"`
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
