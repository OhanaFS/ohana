package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type Config struct {
	Environment    string         `yaml:"environment"`
	HTTP           HttpConfig     `yaml:"http"`
	Database       DatabaseConfig `yaml:"database"`
	Authentication AuthConfig     `yaml:"authentication"`
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

type AuthConfig struct {
	//URL for the SSO authenticating server
	ConfigURL string `yaml:"config_url"`
	//Client ID required for the SSO authenticating server
	ClientID string `yaml:"client_id"`
	//Client secret required for the SSO authenticating server
	ClientSecret string `yaml:"client_secret"`
	//URL for the callback after authentication
	RedirectURL string `yaml:"redirect_url"`
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
