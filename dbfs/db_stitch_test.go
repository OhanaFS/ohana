package dbfs_test

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestSettingStitchParams(t *testing.T) {

	db := testutil.NewMockDB(t)

	fakeLogger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	t.Run("Testing default", func(t *testing.T) {
		Assert := assert.New(t)

		params, err := dbfs.GetStitchParams(db, fakeLogger)
		Assert.NoError(err)

		Assert.Equal(params.DataShards, 2)
		Assert.Equal(params.ParityShards, 1)
		Assert.Equal(params.KeyThreshold, 2)

	})

	t.Run("Testing setting params", func(t *testing.T) {

		Assert := assert.New(t)

		err := dbfs.SetStitchParams(db, 3, 4, 5)
		Assert.NoError(err)

		params, err := dbfs.GetStitchParams(db, fakeLogger)
		Assert.NoError(err)

		Assert.Equal(params.DataShards, 3)
		Assert.Equal(params.ParityShards, 4)
		Assert.Equal(params.KeyThreshold, 5)

		// Trying invalid values

		err = dbfs.SetStitchParams(db, 0, 0, 0)
		Assert.Error(err)
		Assert.Contains(err.Error(), "dataShards")
		Assert.Contains(err.Error(), "parityShards")
		Assert.Contains(err.Error(), "keyThreshold")

		// Trying invalid values

		err = dbfs.SetStitchParams(db, 11, 11, 11)
		Assert.Error(err)
		Assert.Contains(err.Error(), "dataShards")
		Assert.Contains(err.Error(), "parityShards")
		Assert.Contains(err.Error(), "keyThreshold")

		// Test that it goes back to defaults

		params, err = dbfs.GetStitchParams(db, fakeLogger)
		Assert.NoError(err)

		Assert.Equal(params.DataShards, 2)
		Assert.Equal(params.ParityShards, 1)
		Assert.Equal(params.KeyThreshold, 2)

	})

}
