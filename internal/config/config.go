// Package config handles configuration via environment variables
package config

import (
	"log"
	"os"
)

// Config contains configuration values
type Config struct {
	MongoURI   string
	APIRootURL string
}

// Get reads environment variables and parses them into a Config struct.
//
// Any required environment variables which are not provided will cause the application to exit.
func Get() Config {
	mongo := mustGetEnv("MONGODB_URI")
	api := mustGetEnv("API_URL")
	return Config{
		MongoURI:   mongo,
		APIRootURL: api,
	}
}

func mustGetEnv(env string) string {
	v := os.Getenv(env)
	if v == "" {
		log.Fatalf("environment variable %s is required, but not available", env)
	}
	return v
}
