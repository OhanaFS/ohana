package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type Config struct {
	Environment string         `yaml:"environment"`
	HTTP        HttpConfig     `yaml:"http"`
	Database    DatabaseConfig `yaml:"database"`
	Redis       RedisConfig    `yaml:"redis"`
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

type RedisConfig struct {
	password string `yaml:"password"`
	address  string `yaml:"address"`
	db       int    `yaml:"db"`
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

	return &config, nil
}
