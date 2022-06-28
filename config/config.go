package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type Config struct {
	Environment string         `yaml:"environment"`
	HTTP        HttpConfig     `yaml:"http"`
	Database    DatabaseConfig `yaml:"database"`
}

type HttpConfig struct {
	// Bind is the address:port to bind the HTTP server to.
	Bind string `yaml:"bind"`
	// BaseURL is the publicly accessible URL for the HTTP server.
	BaseURL string `yaml:"base_url"`
}

type DatabaseConfig struct {
	// ConnectionString is the connection string for the database.
	ConnectionString string `yaml:"connection_string"`
}

// LoadConfig tries to load the configuration from the file specified in the
// CONFIG_FILE environment variable. If the variable is not set, it defaults to
// "config.yaml".
func LoadConfig() (*Config, error) {
	configFileName := os.Getenv("CONFIG_FILE")
	if configFileName == "" {
		configFileName = "config.yaml"
	}

	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	// If on development mode, allow override of the HTTP bind using the PORT
	// environment variable.
	if config.Environment == EnvironmentDevelopment {
		portOverride := os.Getenv("PORT")
		if portOverride != "" {
			config.HTTP.Bind = ":" + portOverride
			log.Println("Overriding HTTP bind using the PORT environment variable during development mode.")
		}
	}

	return &config, nil
}
