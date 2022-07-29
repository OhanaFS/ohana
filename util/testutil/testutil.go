package testutil

import (
	"testing"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewMockDB(t *testing.T) *gorm.DB {
	assert := assert.New(t)
	appDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

func NewMockSession(t *testing.T) service.Session {
	assert := assert.New(t)
	session, err := service.NewSession(&config.Config{}, zap.NewNop())
	assert.Nil(err)

	return session
}
