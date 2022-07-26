package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
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
	Redis          RedisConfig    `yaml:"redis"`
	Stitch         StitchConfig   `yaml:"stitch"`
	SPA            SPAConfig      `yaml:"-"`
}

type HttpConfig struct {
	// Bind is the Address:port to bind the HTTP server to.
	Bind string `yaml:"bind"`
	// BaseURL is the publicly accessible URL for the HTTP server.
	BaseURL string `yaml:"base_url"`
}

type DatabaseConfig struct {
	// ConnectionString is the connection string for the database.
	ConnectionString string `yaml:"connection_string"`
}

type AuthConfig struct {
	// URL for the SSO authenticating server
	ConfigURL string `yaml:"config_url"`
	// Client ID required for the SSO authenticating server
	ClientID string `yaml:"client_id"`
	// Client secret required for the SSO authenticating server
	ClientSecret string `yaml:"client_secret"`
	// URL for the callback after authentication
	RedirectURL string `yaml:"redirect_url"`
}

type RedisConfig struct {
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Db       int    `yaml:"db"`
}

// SPAConfig is the configuration for the SPA router. It is not exposed to the
// configuration file.
type SPAConfig struct {
	StaticPath           string
	IndexPath            string
	UseDevelopmentServer bool
	DevelopmentServerURL string
}

// StitchConfig is the configuration for the Stitch router.
// It contains the path of where the shards are located on the server
type StitchConfig struct {
	ShardsLocation string `yaml:"shards_location"`
}

// CertsConfig is the configuration for the certs, stored via flags
type FlagsConfig struct {
	GenCA        *bool
	GenCerts     *bool
	GenCAPath    *string
	GenCertsPath *string
	CsrPath      *string
	CertPath     *string
	PkPath       *string
	AllHosts     *string
	NumOfCerts   *int
}

// LoadFlagsConfig tries to load the configuration from flags
func LoadFlagsConfig() *FlagsConfig {
	var config FlagsConfig
	config.GenCA = flag.Bool("gen-ca", false, "Generate CA certs")
	config.GenCerts = flag.Bool("gen-certs", false, "Generate certs")
	config.GenCAPath = flag.String("gen-ca-path", "certificates/main", "Filename of output CA certs")
	config.GenCertsPath = flag.String("gen-certs-path", "certificates/output", "Filename of output certs")
	config.CsrPath = flag.String("csr-path", "certificates/main_csr.json", "Filename of input CSR")
	config.CertPath = flag.String("cert-path", "certificates/main_GLOBAL_CERTIFICATE.pem", "Filename of CA cert")
	config.PkPath = flag.String("pk-path", "certificates/main_PRIVATE_KEY.pem", "Filename of CA private key")
	config.AllHosts = flag.String("hosts", "certhosts.yaml", "yaml file of hosts")
	config.NumOfCerts = flag.Int("num-of-certs", 1, "Number of certs to generate")

	flag.Parse()
	return &config
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

	config.SPA = SPAConfig{
		StaticPath:           "web/dist",
		IndexPath:            "index.html",
		UseDevelopmentServer: config.Environment == EnvironmentDevelopment,
		DevelopmentServerURL: "http://localhost:3000",
	}

	// If on development mode, allow override of the HTTP bind using the PORT
	// environment variable.
	if config.Environment == EnvironmentDevelopment {
		portOverride := os.Getenv("PORT")
		if portOverride != "" {
			config.HTTP.Bind = "127.0.0.1:" + portOverride
			log.Println("Overriding HTTP bind using the PORT environment variable during development mode.")
		}
	}

	return &config, nil
}
