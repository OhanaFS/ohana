package dbfs_test

import (
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T) {
	db := testutil.NewMockDB(t)

	t.Run("Getting root folder", func(t *testing.T) {

		Assert := assert.New(t)

		file, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		Assert.Equal("00000000-0000-0000-0000-000000000000", file.FileID)
		Assert.Equal("root", file.FileName)

		json, err := json.Marshal(file)

		Assert.Nil(err)

		fmt.Println(string(json))

	})

}
