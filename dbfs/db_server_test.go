package dbfs_test

import (
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetServers(t *testing.T) {
	db := testutil.NewMockDB(t)

	debugPrint := true

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	t.Run("Create Servers", func(t *testing.T) {

		Assert := assert.New(t)

		server := dbfs.Server{
			Name:      "server1",
			HostName:  "localhost",
			Port:      "5555",
			Status:    dbfs.ServerStarting,
			FreeSpace: 23,
		}
		err = db.Save(&server).Error
		Assert.NoError(err)

		server = dbfs.Server{
			Name:      "server2",
			HostName:  "localhost.lo",
			Port:      "5555",
			Status:    dbfs.ServerStarting,
			FreeSpace: 25,
		}

		err = db.Save(&server).Error
		Assert.NoError(err)

		server = dbfs.Server{
			Name:      "server3",
			HostName:  "localhost.lol",
			Port:      "5555",
			Status:    dbfs.ServerStarting,
			FreeSpace: 24,
		}

		err = db.Save(&server).Error
		Assert.NoError(err)

		// duplicate server, updating params. Checking to see that it works properly.

		server = dbfs.Server{
			Name:      "server3",
			HostName:  "localhost.change",
			Port:      "5555",
			Status:    dbfs.ServerStarting,
			FreeSpace: 24,
		}

		err = db.Save(&server).Error
		Assert.NoError(err)

	})

	t.Run("Get Servers", func(t *testing.T) {

		Assert := assert.New(t)

		servers, err := dbfs.GetServers(db)
		Assert.NoError(err)

		if debugPrint {
			for _, server := range servers {
				fmt.Println(server)
			}
		}

		Assert.Equal(3, len(servers))

	})

	t.Run("Get Servers Address", func(t *testing.T) {

		Assert := assert.New(t)

		address, err := dbfs.GetServerAddress(db, "server1")
		Assert.NoError(err)
		Assert.Equal("localhost:5555", address)

		// checking that duplicates are fine

		address, err = dbfs.GetServerAddress(db, "server3")
		Assert.NoError(err)
		Assert.Equal("localhost.change:5555", address)

	})

}
