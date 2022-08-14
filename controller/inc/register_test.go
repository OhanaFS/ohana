package inc_test

import (
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
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

		fakeLogger, _ := zap.NewDevelopment()
		incServer = inc.NewInc(configFile, db, fakeLogger)
		inc.RegisterIncServices(incServer)

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

	t.Run("Mark a server as offline", func(t *testing.T) {

		Assert := assert.New(t)

		err := inc.MarkServerOffline(db, incServer.ServerName, "server3")
		Assert.NoError(err)

		// Check if the server is in the database
		_, err = dbfs.GetServerAddress(db, "server3")
		Assert.Error(dbfs.ErrServerOffline, err)

	})

	incServer.HttpServer.Close()

}
