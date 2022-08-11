package testutil

import (
	"gorm.io/driver/postgres"
	"log"
	"os"
	"testing"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMockDB(t *testing.T) *gorm.DB {
	assert := assert.New(t)
	appDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: true,
			}),
	})
	assert.NoError(err)
	assert.NoError(dbfs.InitDB(appDB))
	t.Log("Migrated in-memory database")

	var tables []string
	appDB.
		Raw("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;").
		Scan(&tables)
	t.Logf("Tables: %v", tables)

	return appDB
}

// NewPostgresDB returns a new database connection based on the dbString
// e.g. use it like
// db := testutil.NewPostgresDB(t, "postgres://ohanaAdmin:ohanaMeansFamily@127.0.0.1:5432/ohana")
func NewPostgresDB(t *testing.T, dbString string) *gorm.DB {
	assert := assert.New(t)

	var db *gorm.DB
	var err error

	const maxTries = 10
	for try := 0; try < maxTries; try++ {
		db, err = gorm.Open(postgres.Open(dbString), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\n", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				}),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			break
		}
		if try < maxTries-1 {
			time.Sleep(time.Second)
		}
	}

	assert.NoError(err)

	// Test connection
	testDbPing, err := db.DB()
	assert.NoError(err)
	assert.NoError(testDbPing.Ping())

	// Init
	assert.NoError(dbfs.InitDB(db))
	t.Log("Migrated database")

	// TODO Clear the database

	var tables []string
	db.
		Raw("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;").
		Scan(&tables)
	t.Logf("Tables: %v", tables)

	return db
}

func NewMockSession(t *testing.T) service.Session {
	assert := assert.New(t)
	session, err := service.NewSession(&config.Config{}, zap.NewNop())
	assert.Nil(err)

	return session
}
