package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	dbfstestutils "github.com/OhanaFS/ohana/dbfs/test_utils"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestFixMissingShardResult(t *testing.T) {

	// Setting up env

	tempDir = t.TempDir()
	shardDir := filepath.Join(tempDir, "shards")
	certPaths, err := selfsigntestutils.GenCertsTest(tempDir)
	assert.NoError(t, err)

	//Set up mock Db and session store
	stitchConfig := config.StitchConfig{
		ShardsLocation: shardDir,
	}
	incConfig := config.IncConfig{
		ServerName: "localServer",
		HostName:   "localhost",
		BindIp:     "127.0.0.1",
		Port:       "5556",
		CaCert:     certPaths.CaCertPath,
		PublicCert: certPaths.PublicCertPath,
		PrivateKey: certPaths.PrivateKeyPath,
	}

	configFile := &config.Config{Stitch: stitchConfig, Inc: incConfig}
	logger := config.NewLogger(configFile)
	db := testutil.NewMockDB(t)

	session := testutil.NewMockSession(t)
	sessionId, err := session.Create(nil, "superuser", time.Hour)
	Inc := inc.NewInc(configFile, db, logger)
	inc.RegisterIncServices(Inc)

	// Wait for inc to start
	time.Sleep(time.Second * 10)

	// Setting up controller
	bc := &controller.BackendController{
		Db:         db,
		Logger:     logger,
		Path:       configFile.Stitch.ShardsLocation,
		ServerName: configFile.Inc.ServerName,
		Inc:        Inc,
	}

	bc.InitialiseShardsFolder()

	// Getting Superuser to use with testing
	user, err := dbfs.GetUser(db, "superuser")
	assert.NoError(t, err)

	// Testing

	data1 := "Hello this is the first version of the file. Hopefully it won't get corrupted poggies"
	data2 := "Hello this is the second version of the file. " +
		"It won't get corrupted like the other one because this file isn't cringe."

	var newJob dbfs.Job
	var testFileID string

	t.Run("StartJob", func(t *testing.T) {

		// Before starting our full shards check, we want to create a few files so that it
		// can be tested

		// Get root folder

		Assert := assert.New(t)

		time.Sleep(time.Second * 2)

		testFile, err := dbfstestutils.EXAMPLECreateFile(bc.Db, user, dbfstestutils.ExampleFile{
			FileName:       "Test123",
			ParentFolderId: "00000000-0000-0000-0000-000000000000",
			Server:         incConfig.ServerName,
			FragmentPath:   stitchConfig.ShardsLocation,
			FileData:       data1,
			Size:           50,
			ActualSize:     50,
		})
		Assert.NoError(err)

		testFileID = testFile.FileId

		// Turn on versioning
		Assert.NoError(testFile.UpdateMetaData(bc.Db, dbfs.FileMetadataModification{
			VersioningMode: dbfs.VersioningOnVersions,
		}, user))

		// Get the shards for the file to corrupt it later pepelaugh
		shards, err := testFile.GetFileFragments(bc.Db, user)
		Assert.NoError(err)
		Assert.True(len(shards) > 2, shards) // This is a sanity check to make sure we have at least 3 shards

		// Update the file

		Assert.NoError(dbfstestutils.EXAMPLEUpdateFile(bc.Db, testFile, dbfstestutils.ExampleUpdate{
			NewSize:       70,
			NewActualSize: 70,
			FragmentPath:  stitchConfig.ShardsLocation,
			FileData:      data2,
			Server:        incConfig.ServerName,
		}, user))

		shards2, err := testFile.GetFileFragments(bc.Db, user)
		Assert.NoError(err)
		Assert.True(len(shards) > 2, shards) // This is a sanity check to make sure we have at least 3 shards

		// Remove one of the shards from both versions
		Assert.NoError(os.Remove(path.Join(stitchConfig.ShardsLocation, shards[0].FileFragmentPath)))
		Assert.NoError(os.Remove(path.Join(stitchConfig.ShardsLocation, shards2[0].FileFragmentPath)))

		// Starting a job with a quick shards check
		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("missing_shards_check", "true")
		w := httptest.NewRecorder()
		bc.StartJob(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(1))
		Assert.Equal(newJob.MissingShardsCheck, true)

	})

	var resultsMissingShards []dbfs.ResultsMissingShard

	t.Run("Get Missing Shards Result", func(t *testing.T) {

		time.Sleep(1 * time.Second)
		Assert := assert.New(t)

		// Get the resultsMissingShards
		req := httptest.NewRequest("GET",
			"/api/v1/cluster/maintenance/job/2/missing_shards",
			nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		// mux
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})
		w := httptest.NewRecorder()
		bc.GetMissingShardsResult(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &resultsMissingShards))
		Assert.Equal(2, len(resultsMissingShards))
	})

	t.Run("Repair said shards", func(t *testing.T) {

		// This test is reliant on the previous test running successfully
		// As we are using the results from resultsMissingShards which is populated from the previous test

		Assert := assert.New(t)

		// Repair said shards

		// Creating the request to fix all shards
		shardActions := make([]dbfs.ShardActions, len(resultsMissingShards))
		for i, result := range resultsMissingShards {
			shardActions[i] = dbfs.ShardActions{
				DataId: result.DataId,
				Fix:    true,
			}
		}

		fixAFSHCBytes, err := json.Marshal(shardActions)
		Assert.NoError(err)

		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/1/quick_shards",
			bytes.NewReader(fixAFSHCBytes)).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("Content-Type", "application/json") // technically not required but good practice
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})

		w := httptest.NewRecorder()
		bc.FixMissingShardsResult(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

		// In theory all the shards should be fixed. Let's check that TODO

		fmt.Println(body)

		// Check that fragments are updated properly

		var fragments []dbfs.Fragment
		Assert.Nil(bc.Db.Where("file_version_data_id = ?", resultsMissingShards[0].DataId).Find(&fragments).Error)
		Assert.Equal(3, len(fragments))

		// Try to download the file

		req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": testFileID,
		})
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		Assert.Equal(http.StatusOK, w.Code, w.Body.String())
		Assert.Contains(w.Body.String(), data2)

		// Checking the cron job to see it being marked as done.

		// Get the resultsMissingShards
		req = httptest.NewRequest("GET",
			"/api/v1/cluster/maintenance/job/2/full_shards",
			nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		// mux
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})
		w = httptest.NewRecorder()
		bc.GetMissingShardsResult(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &resultsMissingShards))
		Assert.Equal(1, len(resultsMissingShards))
		Assert.Equal(dbfs.CronErrorTypeSolved, resultsMissingShards[0].ErrorType)

	})

	t.Run("Get Job again to see if it is marked as done", func(t *testing.T) {

		Assert := assert.New(t)

		// Get the Job
		req := httptest.NewRequest("GET", "/api/v1/cluster/maintenance/job/2", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})
		w := httptest.NewRecorder()
		bc.GetJob(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		var job dbfs.Job
		Assert.NoError(json.Unmarshal([]byte(body), &job))

		Assert.Equal(100, job.Progress)

	})

	Inc.HttpServer.Close()

}
