package inc_test

import (
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/selfsign"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRegisterServer(t *testing.T) {

	// Initialize the config

	Assert := assert.New(t)

	db := testutil.NewMockDB(t)

	stitchConfig := config.StitchConfig{
		ShardsLocation: "shards/",
	}

	// Setting up certs for configs

	ogc := config.LoadFlagsConfig()
	trueBool := true
	ogc.GenCA = &trueBool
	ogc.GenCerts = &trueBool

	fakeHosts := selfsign.Hosts{Hosts: []string{"localhost", "localhost2"}}

	hostFile, err := os.Create("certhosts.yaml")
	Assert.Nil(err)
	defer hostFile.Close()

	encoder := yaml.NewEncoder(hostFile)
	Assert.Nil(encoder.Encode(fakeHosts))

	//err = selfsign.ProcessFlags(ogc)
	Assert.Nil(err)

	configFile := &config.Config{Stitch: stitchConfig,
		Database: config.DatabaseConfig{
			ConnectionString: "ItShouldn'tUseThisBecauseSQLiteInTesting",
			ServerName:       "testServer",
			HostName:         "localhost",
			Port:             "5555",
			Ca_Cert:          "certificates/main_GLOBAL_CERTIFICATE.pem",
			Public_Cert:      "certificates/output_cert.pem",
			Private_Key:      "certificates/output_key.pem",
		},
	}

	t.Run("Running a Server for ping test", func(t *testing.T) {

		mux := http.NewServeMux()
		mux.HandleFunc("/inc/ping", inc.Pong)

		server := &http.Server{
			Addr:    ":" + configFile.Database.Port,
			Handler: mux,
		}

		go server.ListenAndServe()

	})

	t.Run("Test Pong", func(t *testing.T) {

		time.Sleep(1 * time.Second)

		//req := httptest.NewRequest("GET", "/inc/ping", nil)
		//w := httptest.NewRecorder()
		//inc.Pong(w, req)
		//Assert.Equal(http.StatusOK, w.Code)

		Assert := assert.New(t)

		Assert.Equal(inc.Ping(configFile.Database.HostName, configFile.Database.Port), true)

	})

	t.Run("Register Server", func(t *testing.T) {

		Assert := assert.New(t)

		err := inc.RegisterServer(db, configFile, true)
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
		err := inc.RegisterServer(db, configFile, true)
		Assert.NoError(err)

	})

	t.Run("Mark a server as offline", func(t *testing.T) {

		Assert := assert.New(t)

		err := inc.MarkServerOffline(db, configFile, "server3")
		Assert.NoError(err)

		// Check if the server is in the database
		_, err = dbfs.GetServerAddress(db, "server3")
		Assert.Error(dbfs.ErrServerOffline, err)

	})

	t.Run("Cert Cleanup", func(t *testing.T) {
		err = os.RemoveAll("certificates")
		Assert.NoError(err)
		err = os.RemoveAll("certhosts.yaml")
		Assert.NoError(err)

	})

}
