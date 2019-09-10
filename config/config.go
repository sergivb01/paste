package config

import (
	"os"

	"github.com/joho/godotenv"
)

// config defines the configuration for the application
type Config struct {
	Listen      string
	PostgresURI string
}

// loadConfig loads the configuration from provied yml file path
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Listen:      getEnv("LISTEN", ":8080"),
		PostgresURI: getEnv("POSTGRES_URI", "none"),
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}
