// Package config handles configuration via environment variables
package config

import (
	"log"
	"os"
)

// Config contains configuration values
type Config struct {
	MongoURI string
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() Config {
	uri := mustGetEnv("MONGODB_URI")
	return Config{
		MongoURI: uri,
	}
}

func mustGetEnv(env string) string {
	v := os.Getenv(env)
	if v == "" {
		log.Fatalf("environment variable %s is required, but not available", env)
	}
	return v
}
