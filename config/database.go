package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(config *Config) (*gorm.DB, error) {
	dbString := config.Database.ConnectionString
	var db *gorm.DB
	var err error

	loggerConfig := logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Silent,
	}

	if config.Environment == EnvironmentDevelopment {
		loggerConfig =
			logger.Config{
				LogLevel: logger.Info,
				Colorful: true,
			}
	}

	lg := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		loggerConfig,
	)

	const maxTries = 10
	for try := 0; try < maxTries; try++ {
		db, err = gorm.Open(postgres.Open(dbString), &gorm.Config{
			Logger:                                   lg,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			break
		}
		if try < maxTries-1 {
			time.Sleep(time.Second)
		}
	}

	if err != nil {
		return nil, err
	}

	// Test connection
	if testDbPing, err := db.DB(); err != nil {
		return nil, err
	} else {
		if err := testDbPing.Ping(); err != nil {
			return nil, err
		}
	}

	return db, nil

}
