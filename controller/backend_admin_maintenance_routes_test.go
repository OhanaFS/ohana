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
	"github.com/OhanaFS/stitch"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetAllJobs(t *testing.T) {

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
		Port:       "5555",
		CaCert:     certPaths.CaCertPath,
		PublicCert: certPaths.PublicCertPath,
		PrivateKey: certPaths.PrivateKeyPath,
	}

	configFile := &config.Config{Stitch: stitchConfig, Inc: incConfig}
	logger := config.NewLogger(configFile)
	db := testutil.NewMockDB(t)

	session := testutil.NewMockSession(t)
	sessionId, err := session.Create(nil, "superuser", time.Hour)
	Inc := inc.NewInc(configFile, db)
	inc.RegisterIncServices(Inc)

	// Wait for inc to start
	time.Sleep(time.Second * 3)

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

	// Create jobs

	job := dbfs.Job{
		StartTime:                    time.Now(),
		EndTime:                      time.Now(),
		TotalTimeTaken:               0,
		TotalShardsScanned:           0,
		TotalFilesScanned:            0,
		MissingShardsCheck:           true,
		MissingShardsProgress:        nil,
		OrphanedShardsCheck:          false,
		OrphanedShardsProgress:       nil,
		QuickShardsHealthCheck:       false,
		QuickShardsHealthProgress:    nil,
		AllFilesShardsHealthCheck:    false,
		AllFilesShardsHealthProgress: nil,
		PermissionCheck:              false,
		PermissionResults:            nil,
		DeleteFragments:              false,
		DeleteFragmentsResults:       nil,
		Progress:                     100,
		StatusMsg:                    "done pog",
		Status:                       3,
	}

	jpms := dbfs.JobProgressMissingShard{
		JobId:      1,
		StartTime:  job.StartTime,
		EndTime:    time.Time{},
		ServerId:   "",
		InProgress: false,
		Msg:        "",
	}

	// save the entries

	assert.NoError(t, db.Save(&job).Error) // This save is f-ed up. Saves as 1,
	assert.NoError(t, db.Save(&jpms).Error)

	var newJob dbfs.Job

	data1 := "Hello this is the first version of the file. Hopefully it won't get corrupted poggies"
	data2 := "Hello this is the second version of the file. " +
		"It won't get corrupted like the other one because this file isn't cringe."

	t.Run("GetAllJobs", func(t *testing.T) {

		Assert := assert.New(t)

		// Get the jobs

		req := httptest.NewRequest("GET", "/api/v1/cluster/maintenance/all", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("start_num", "0")
		w := httptest.NewRecorder()
		bc.GetAllJobs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		var jobs []dbfs.Job
		Assert.NoError(json.Unmarshal([]byte(body), &jobs))
		Assert.Equal(1, len(jobs))
		Assert.Equal(job.JobId, jobs[0].JobId)
		Assert.Equal(len(jobs[0].MissingShardsProgress), 1)
		Assert.Equal(len(jobs[0].OrphanedShardsProgress), 0)

	})

	t.Run("GetJob", func(t *testing.T) {

		Assert := assert.New(t)

		// Get the jobs

		req := httptest.NewRequest("GET", "/api/v1/cluster/maintenance/job/1", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(job.JobId)),
		})
		w := httptest.NewRecorder()
		bc.GetJob(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		var job dbfs.Job
		Assert.NoError(json.Unmarshal([]byte(body), &job))
		Assert.Equal(job.JobId, uint(1))
		Assert.Equal(len(job.MissingShardsProgress), 1)
		Assert.Equal(len(job.OrphanedShardsProgress), 0)

	})

	t.Run("DeleteJob", func(t *testing.T) {

		Assert := assert.New(t)

		// Delete the job
		req := httptest.NewRequest("DELETE", "/api/v1/cluster/maintenance/job/1", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(job.JobId)),
		})
		w := httptest.NewRecorder()
		bc.DeleteJob(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

		// Checking the db itself to ensure it got deleted

		req = httptest.NewRequest("GET", "/api/v1/cluster/maintenance/job/1", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(job.JobId)),
		})
		w = httptest.NewRecorder()
		bc.GetJob(w, req)
		Assert.Equal(http.StatusNotFound, w.Code)

	})

	var testFileID string

	t.Run("StartJob", func(t *testing.T) {

		// Before starting our full shards check, we want to create a few files so that it
		// can be tested

		// Get root folder

		Assert := assert.New(t)

		rootFolder, err := dbfs.GetRootFolder(db)
		Assert.NoError(err)

		testFile, err := dbfstestutils.EXAMPLECreateFile(db, user, dbfstestutils.ExampleFile{
			FileName:       "Test123",
			ParentFolderId: rootFolder.FileId,
			Server:         incConfig.ServerName,
			FragmentPath:   stitchConfig.ShardsLocation,
			FileData:       data1,
			Size:           50,
			ActualSize:     50,
		})
		Assert.NoError(err)

		testFileID = testFile.FileId

		// Turn on versioning
		Assert.NoError(testFile.UpdateMetaData(db, dbfs.FileMetadataModification{
			VersioningMode: dbfs.VersioningOnVersions,
		}, user))

		// Get the shards for the file to corrupt it later pepelaugh
		shards, err := testFile.GetFileFragments(db, user)
		Assert.NoError(err)
		Assert.True(len(shards) > 2, shards) // This is a sanity check to make sure we have at least 3 shards

		// Update the file

		Assert.NoError(dbfstestutils.EXAMPLEUpdateFile(db, testFile, dbfstestutils.ExampleUpdate{
			NewSize:       70,
			NewActualSize: 70,
			FragmentPath:  stitchConfig.ShardsLocation,
			FileData:      data2,
			Server:        incConfig.ServerName,
		}, user))

		shards2, err := testFile.GetFileFragments(db, user)
		Assert.NoError(err)
		Assert.True(len(shards) > 2, shards) // This is a sanity check to make sure we have at least 3 shards

		// Corrupt one of the shards from both versions
		corruptPath := path.Join(stitchConfig.ShardsLocation, shards[0].FileFragmentPath)
		Assert.NoError(dbfstestutils.EXAMPLECorruptFragments(corruptPath))

		corruptPath = path.Join(stitchConfig.ShardsLocation, shards2[0].FileFragmentPath)
		Assert.NoError(dbfstestutils.EXAMPLECorruptFragments(corruptPath))

		// Starting a job with a full shards check
		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("full_shards_check", "true")
		w := httptest.NewRecorder()
		bc.StartJob(w, req)

		// Starting a job with a full shards check
		req = httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("full_shards_check", "true")
		w = httptest.NewRecorder()
		bc.StartJob(w, req)

		var jobs []dbfs.Job
		Assert.NoError(db.Find(&jobs).Error)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(2))
		Assert.Equal(newJob.AllFilesShardsHealthCheck, true)

	})

	var resultsAFSHC []dbfs.ResultsAFSHC

	t.Run("Get Full Shards Result", func(t *testing.T) {

		time.Sleep(1 * time.Second)
		Assert := assert.New(t)

		// Get the resultsAFSHC
		req := httptest.NewRequest("GET",
			"/api/v1/cluster/maintenance/job/2/full_shards",
			nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		// mux
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})
		w := httptest.NewRecorder()
		bc.GetFullShardsResult(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &resultsAFSHC))
		Assert.Equal(2, len(resultsAFSHC))
	})

	t.Run("Repair said shards", func(t *testing.T) {

		// This test is reliant on the previous test running successfully
		// As we are using the results from resultsAFSHC which is populated from the previous test

		Assert := assert.New(t)

		// Repair said shards

		// Creating the request to fix all shards
		fixAFSHC := make([]dbfs.FixAFSHC, len(resultsAFSHC))
		for i, result := range resultsAFSHC {
			fixAFSHC[i] = dbfs.FixAFSHC{
				DataId: result.DataId,
				Fix:    true,
			}
		}

		fixAFSHCBytes, err := json.Marshal(fixAFSHC)
		Assert.NoError(err)

		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/2/full_shards",
			bytes.NewReader(fixAFSHCBytes)).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("Content-Type", "application/json") // technically not required but good practice
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})

		w := httptest.NewRecorder()
		bc.FixFullShardsResult(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

		// In theory all the shards should be fixed. Let's check that

		req = httptest.NewRequest("GET",
			strings.Replace(inc.FragmentHealthCheckPath,
				inc.PathReplaceShardString, resultsAFSHC[0].FragPath, -1),
			nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"shardPath": resultsAFSHC[0].FragPath,
		})
		w = httptest.NewRecorder()
		Inc.ShardHealthCheckRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		var verificationResult stitch.ShardVerificationResult
		Assert.NoError(json.Unmarshal([]byte(body), &verificationResult))
		Assert.Equal(true, verificationResult.IsAvailable)
		Assert.Equal(0, len(verificationResult.BrokenBlocks))

		fmt.Println(body)

		// Check that fragments are updated properly

		var fragments []dbfs.Fragment
		Assert.Nil(db.Where("file_version_data_id = ?", resultsAFSHC[0].DataId).Find(&fragments).Error)
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

		// Try to download the older version of the file

		req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID":    testFileID,
			"versionID": "1",
		})
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		Assert.Equal(http.StatusOK, w.Code, w.Body.String())
		Assert.Contains(w.Body.String(), data1)

		// Checking the cron job to see it being marked as done.

		// Get the resultsAFSHC
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
		bc.GetFullShardsResult(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &resultsAFSHC))
		Assert.Equal(2, len(resultsAFSHC))
		Assert.Equal(dbfs.CronErrorTypeSolved, resultsAFSHC[0].ErrorType)
		Assert.Equal(dbfs.CronErrorTypeSolved, resultsAFSHC[1].ErrorType)

	})

	Inc.HttpServer.Close()

}
