package config

import "go.uber.org/zap"

func NewLogger(c *Config) *zap.Logger {
	var logger *zap.Logger

	if c.Environment == EnvironmentDevelopment {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	return logger
}
