package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// config defines the configuration for the application
type Config struct {
	Listen string `yaml:"listen"`
	
	PostgresURI string `yaml:"postgresURI"`
}

// loadConfig loads the configuration from provied yml file path
func LoadFromFile(path string) (Config, error) {
	var c Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return c, err
	}

	return c, nil
}
