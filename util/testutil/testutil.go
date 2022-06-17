package testutil

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
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
