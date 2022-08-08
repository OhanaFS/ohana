package inc_test

import (
	"context"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRegisterServer(t *testing.T) {

	// Initialize the config

	Assert := assert.New(t)
	tempdir := t.TempDir()

	db := testutil.NewMockDB(t)

	stitchConfig := config.StitchConfig{
		ShardsLocation: "shards/",
	}

	// Setting up certs for configs

	certsPath, err := selfsigntestutils.GenCertsTest(tempdir)

	Assert.Nil(err)

	configFile := &config.Config{Stitch: stitchConfig,
		Inc: config.IncConfig{
			ServerName: "testServer",
			HostName:   "localhost",
			Port:       "5555",
			CaCert:     certsPath.CaCertPath,
			PublicCert: certsPath.PublicCertPath,
			PrivateKey: certsPath.PrivateKeyPath,
		},
	}

	var incServer *inc.Inc

	// Create inc server

	t.Run("Register Server", func(t *testing.T) {

		Assert := assert.New(t)

		incServer = inc.NewInc(configFile, db)

		err := incServer.RegisterServer(true)
		Assert.NoError(err)

		// Check if the server is in the database
		server, err := dbfs.GetServerAddress(db, "testServer")
		Assert.NoError(err)
		Assert.Equal("localhost:5555", server)

		// Get All Servers
		servers, err := dbfs.GetServers(db)
		Assert.NoError(err)
		for _, server := range servers {
			fmt.Println(server)
		}

	})

	t.Run("Register server while another server is already registered", func(t *testing.T) {

		Assert := assert.New(t)

		//manually registering a server. setting it as "in process" for 10 sec
		go func() {

			db := ctxutil.GetTransaction(context.Background(), db)

			testServer := dbfs.Server{
				Name:      "server3",
				HostName:  "127.0.0.1",
				Port:      "5555",
				Status:    dbfs.ServerStarting,
				FreeSpace: uint64(24),
			}
			err := db.Save(&testServer).Error
			assert.NoError(t, err)
			fmt.Println("Registered server3 as Pending")
			time.Sleep(time.Second * 2)
			testServer.Status = dbfs.ServerOnline
			err = db.Save(&testServer).Error
			fmt.Println("Registered server3 as Online")
			Assert.NoError(err)

		}()

		/*
			Note: Gorm's sqlite3 driver doesn't like concurrent threads accessing the same database. (fair)
			but unlike the generic sqlite3 driver, we can't set a max concurrent to 1 to force it to retry.
			I'm just using sleep to force it to happen at a different timing.

			Will not happen in production according to google. Just a gorm sqlite3 thing.

			I've been debugging this for like a hr and a half as if you don't mark it as shared anyway it just gives
			you the error "Table not found" which drove me nuts. (googling fixed it ofc but I was stubborn)
		*/
		time.Sleep(time.Second * 1)
		err := incServer.RegisterServer(true)
		Assert.NoError(err)

	})

	t.Run("Mark a server as offline", func(t *testing.T) {

		Assert := assert.New(t)

		err := inc.MarkServerOffline(db, incServer.ServerName, "server3")
		Assert.NoError(err)

		// Check if the server is in the database
		_, err = dbfs.GetServerAddress(db, "server3")
		Assert.Error(dbfs.ErrServerOffline, err)

	})

	defer incServer.HttpServer.Shutdown(context.Background())

}
