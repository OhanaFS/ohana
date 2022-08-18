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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFixFullShardResult(t *testing.T) {

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
		fixAFSHC := make([]dbfs.ShardActions, len(resultsAFSHC))
		for i, result := range resultsAFSHC {
			fixAFSHC[i] = dbfs.ShardActions{
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

func TestFixQuickShardResult(t *testing.T) {

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

		// Get the shards for the file to corrupt it later
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

		// Starting a job with a quick shards check
		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("quick_shards_check", "true")
		w := httptest.NewRecorder()
		bc.StartJob(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(1))
		Assert.Equal(newJob.QuickShardsHealthCheck, true)

	})

	var resultsCFSHC []dbfs.ResultsCFSHC

	t.Run("Get Full Shards Result", func(t *testing.T) {

		time.Sleep(1 * time.Second)
		Assert := assert.New(t)

		// Get the resultsCFSHC
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
		bc.GetQuickShardsResult(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &resultsCFSHC))
		Assert.Equal(1, len(resultsCFSHC)) // should not be 2 as this only checks for current versions
	})

	t.Run("Repair said shards", func(t *testing.T) {

		// This test is reliant on the previous test running successfully
		// As we are using the results from resultsCFSHC which is populated from the previous test

		Assert := assert.New(t)

		// Repair said shards

		// Creating the request to fix all shards
		shardActions := make([]dbfs.ShardActions, len(resultsCFSHC))
		for i, result := range resultsCFSHC {
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
		bc.FixQuickShardsResult(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

		// In theory all the shards should be fixed. Let's check that

		req = httptest.NewRequest("GET",
			strings.Replace(inc.FragmentHealthCheckPath,
				inc.PathReplaceShardString, resultsCFSHC[0].FragPath, -1),
			nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"shardPath": resultsCFSHC[0].FragPath,
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
		Assert.Nil(db.Where("file_version_data_id = ?", resultsCFSHC[0].DataId).Find(&fragments).Error)
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

		// Get the resultsCFSHC
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
		bc.GetQuickShardsResult(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &resultsCFSHC))
		Assert.Equal(1, len(resultsCFSHC))
		Assert.Equal(dbfs.CronErrorTypeSolved, resultsCFSHC[0].ErrorType)

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

func TestFixOrphanedFilesResult(t *testing.T) {

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
		Port:       "5557",
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

	// Creating files to test with

	// Get root folder
	//rootFolder, err := dbfs.GetRootFolder(db)
	//assert.Nil(t, err)

	// Making two file that is like not to be in any folder at all
	parent := "See! I'm empty"
	file := dbfs.File{
		FileId:             "randomfileidlmao",
		FileName:           "I have no parents",
		MIMEType:           "please adopt me",
		EntryType:          dbfs.IsFile,
		ParentFolderFileId: &parent,
		VersionNo:          0,
		DataId:             "I don't have any data as well",
		DataIdVersion:      0,
		Size:               0,
		ActualSize:         0,
		CreatedTime:        time.Time{},
		ModifiedUser:       nil,
		ModifiedUserUserId: nil,
		ModifiedTime:       time.Time{},
		VersioningMode:     0,
		Checksum:           "",
		TotalShards:        0,
		DataShards:         0,
		ParityShards:       0,
		KeyThreshold:       0,
		EncryptionKey:      "",
		EncryptionIv:       "",
		PasswordProtected:  false,
		LinkFile:           nil,
		LinkFileFileId:     nil,
		LastChecked:        time.Time{},
		Status:             0,
		HandledServer:      "",
	}

	err = db.Create(&file).Error
	assert.Nil(t, err)

	parent2 := "lonely inside rip"
	file2 := dbfs.File{
		FileId:             "whocares",
		FileName:           "sadge",
		MIMEType:           "I do tho",
		EntryType:          dbfs.IsFile,
		ParentFolderFileId: &parent2,
		VersionNo:          0,
		DataId:             "you are still loved",
		DataIdVersion:      0,
		Size:               0,
		ActualSize:         0,
		CreatedTime:        time.Time{},
		ModifiedUser:       nil,
		ModifiedUserUserId: nil,
		ModifiedTime:       time.Time{},
		VersioningMode:     0,
		Checksum:           "",
		TotalShards:        0,
		DataShards:         0,
		ParityShards:       0,
		KeyThreshold:       0,
		EncryptionKey:      "",
		EncryptionIv:       "",
		PasswordProtected:  false,
		LinkFile:           nil,
		LinkFileFileId:     nil,
		LastChecked:        time.Time{},
		Status:             0,
		HandledServer:      "",
	}
	err = db.Create(&file2).Error
	assert.Nil(t, err)

	var newJob dbfs.Job
	var orphanedFilesResults []dbfs.ResultsOrphanedFile

	t.Run("Check if the files are orphaned", func(t *testing.T) {

		Assert := assert.New(t)

		// Starting a job
		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("orphaned_files_check", "true")
		w := httptest.NewRecorder()
		bc.StartJob(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(1))
		Assert.Equal(newJob.OrphanedFilesCheck, true)

		time.Sleep(time.Second * 3)

		// Checking if the files are orphaned

		var test []dbfs.JobProgressOrphanedShard

		err = bc.Db.Find(&test).Error
		Assert.Nil(err)

		req = httptest.NewRequest("GET",
			"/api/v1/cluster/maintenance/job/{id}/orphaned_files",
			nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(newJob.JobId))})
		w = httptest.NewRecorder()
		bc.GetOrphanedFilesResult(w, req)

		body = w.Body.String()
		Assert.Equal(http.StatusOK, w.Code, body)

		Assert.NoError(json.Unmarshal([]byte(body), &orphanedFilesResults))
		Assert.Equal(len(orphanedFilesResults), 2)

	})

	t.Run("Fixing Orphaned Files", func(t *testing.T) {

		Assert := assert.New(t)

		fixes := make([]dbfs.OrphanedFilesActions, len(orphanedFilesResults))
		fixes[0] = dbfs.OrphanedFilesActions{
			ParentFolderId: orphanedFilesResults[0].ParentFolderId,
			Delete:         true,
		}
		fixes[1] = dbfs.OrphanedFilesActions{
			ParentFolderId: orphanedFilesResults[1].ParentFolderId,
			Move:           true,
		}

		b, err := json.Marshal(fixes)
		Assert.NoError(err)

		req := httptest.NewRequest("POST",
			"/api/v1/cluster/maintenance/job/{id}/orphaned_files",
			bytes.NewReader(b)).WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(newJob.JobId))})
		w := httptest.NewRecorder()
		bc.FixOrphanedFilesResult(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

	})

	t.Run("Ensuring that the job is marked as completed", func(t *testing.T) {
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

func TestOrphanedShardsResult(t *testing.T) {

	// Set up env

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

	// Create a few files that are normal

	data1 := "Hello this is the first version of the file."
	data2 := "Hello this is the second version of the file."

	// Getting Superuser to use with testing
	user, err := dbfs.GetUser(db, "superuser")
	assert.NoError(t, err)

	t.Run("Create normal files", func(t *testing.T) {

		time.Sleep(time.Second * 10)

		Assert := assert.New(t)

		testFile, err := dbfstestutils.EXAMPLECreateFile(db, user, dbfstestutils.ExampleFile{
			FileName:       "Test123",
			ParentFolderId: "00000000-0000-0000-0000-000000000000",
			Server:         incConfig.ServerName,
			FragmentPath:   stitchConfig.ShardsLocation,
			FileData:       data1,
			Size:           50,
			ActualSize:     50,
		})
		Assert.NoError(err)

		// Turn on versioning
		Assert.NoError(testFile.UpdateMetaData(db, dbfs.FileMetadataModification{
			VersioningMode: dbfs.VersioningOnVersions,
		}, user))

		// Update the file

		Assert.NoError(dbfstestutils.EXAMPLEUpdateFile(db, testFile, dbfstestutils.ExampleUpdate{
			NewSize:       70,
			NewActualSize: 70,
			FragmentPath:  stitchConfig.ShardsLocation,
			FileData:      data2,
			Server:        incConfig.ServerName,
		}, user))
	})

	var newJob dbfs.Job

	t.Run("Test orphaned Shards", func(t *testing.T) {

		// Run inital tests, should be all good
		Assert := assert.New(t)

		// Starting a job with a quick shards check
		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("orphaned_shards_check", "true")
		w := httptest.NewRecorder()
		bc.StartJob(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(1))
		Assert.Equal(newJob.OrphanedShardsCheck, true)

		// Get the job
		time.Sleep(time.Second * 3)

		// Get the resultsMissingShards
		req = httptest.NewRequest("GET",
			"/api/v1/cluster/maintenance/job/2/orphaned_shards",
			nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		// mux
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})
		w = httptest.NewRecorder()
		bc.GetOrphanedShardsResult(w, req)
		body = w.Body.String()
		Assert.Equal(http.StatusOK, w.Code, body)

		var results []dbfs.ResultsOrphanedShard

		Assert.NoError(json.Unmarshal([]byte(body), &results))
		Assert.Equal(0, len(results))

	})

	t.Run("Add a random fragment to the shards", func(t *testing.T) {

		Assert := assert.New(t)

		// Create a random file in wd
		randomFile, err := os.Create(filepath.Join(bc.Inc.ShardsLocation, "random.txt"))
		Assert.NoError(err)

		// Write some random data to the file
		_, err = randomFile.Write([]byte("Hello this is a random file"))
		Assert.NoError(err)

		// Close the file
		Assert.NoError(randomFile.Close())

	})

	var results []dbfs.ResultsOrphanedShard

	t.Run("Test orphaned Shards", func(t *testing.T) {

		Assert := assert.New(t)

		// Starting a job with a quick shards check
		req := httptest.NewRequest("POST", "/api/v1/cluster/maintenance/job/start", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("orphaned_shards_check", "true")
		w := httptest.NewRecorder()
		bc.StartJob(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(2))
		Assert.Equal(newJob.OrphanedShardsCheck, true)

		// Get the job
		time.Sleep(time.Second * 3)

		// Get the resultsMissingShards
		req = httptest.NewRequest("GET",
			"/api/v1/cluster/maintenance/job/2/orphaned_shards",
			nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		// mux
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})
		w = httptest.NewRecorder()
		bc.GetOrphanedShardsResult(w, req)
		body = w.Body.String()
		Assert.Equal(http.StatusOK, w.Code, body)

		Assert.NoError(json.Unmarshal([]byte(body), &results))
		Assert.Equal(1, len(results))

	})

	t.Run("Delete Orphaned Shards", func(t *testing.T) {

		Assert := assert.New(t)

		orphanedShardActions := make([]dbfs.OrphanedShardActions, len(results))

		for i, result := range results {
			orphanedShardActions[i] = dbfs.OrphanedShardActions{
				ServerId: result.ServerId,
				FileName: result.FileName, // TODO: Make it consistent
				Delete:   true,
			}
		}

		actionsBody, err := json.Marshal(orphanedShardActions)
		Assert.NoError(err)

		// Send the request to delete the orphaned shards
		req := httptest.NewRequest("POST",
			"/api/v1/cluster/maintenance/job/2/orphaned_shards",
			bytes.NewReader(actionsBody)).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(newJob.JobId)),
		})

		w := httptest.NewRecorder()
		bc.FixOrphanedShardsResult(w, req)

		body := w.Body.String()
		Assert.Equal(http.StatusOK, w.Code, body)
		Assert.Contains(body, "true")

		// Test to check that the file is deleted :)

		_, err = os.Stat(filepath.Join(bc.Inc.ShardsLocation, "random.txt"))
		Assert.Error(err)
		Assert.True(os.IsNotExist(err))

	})

	Inc.HttpServer.Close()

}

func TestBackendController_SetStitchParameters(t *testing.T) {

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

	//fake zapper

	t.Run("Checking inital params", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/maintenance/stitch", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetStitchParameters(w, req)

		var stitchParams dbfs.StitchParams
		body := w.Body.String()
		Assert.NoError(json.Unmarshal([]byte(body), &stitchParams))
		Assert.Equal(http.StatusOK, w.Result().StatusCode)

		Assert.Equal(stitchParams.DataShards, 2)
		Assert.Equal(stitchParams.ParityShards, 1)
		Assert.Equal(stitchParams.KeyThreshold, 2)

	})

	t.Run("Setting params", func(t *testing.T) {

		// Run inital tests, should be all good
		Assert := assert.New(t)

		// Starting a job with a quick shards check
		req := httptest.NewRequest("POST", "/api/v1/maintenance/stitch", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("data_shards", "4")
		req.Header.Add("key_threshold", "5")
		req.Header.Add("parity_shards", "6")
		w := httptest.NewRecorder()
		bc.SetStitchParameters(w, req)

		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second)

		req = httptest.NewRequest("GET", "/api/v1/maintenance/stitch", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetStitchParameters(w, req)

		var stitchParams dbfs.StitchParams
		body := w.Body.String()
		Assert.NoError(json.Unmarshal([]byte(body), &stitchParams))
		Assert.Equal(http.StatusOK, w.Result().StatusCode)

		Assert.Equal(stitchParams.DataShards, 4)
		Assert.Equal(stitchParams.ParityShards, 5)
		Assert.Equal(stitchParams.KeyThreshold, 6)

	})

	t.Run("Invalid params", func(t *testing.T) {

		// Run inital tests, should be all good
		Assert := assert.New(t)

		// Starting a job with a quick shards check
		req := httptest.NewRequest("POST", "/api/v1/maintenance/stitch", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("data_shards", "11")
		req.Header.Add("key_threshold", "11")
		req.Header.Add("parity_shards", "11")
		w := httptest.NewRecorder()
		bc.SetStitchParameters(w, req)

		Assert.Equal(http.StatusBadRequest, w.Code)

		time.Sleep(time.Second)

		req = httptest.NewRequest("GET", "/api/v1/maintenance/stitch", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetStitchParameters(w, req)

		var stitchParams dbfs.StitchParams
		body := w.Body.String()
		Assert.NoError(json.Unmarshal([]byte(body), &stitchParams))
		Assert.Equal(http.StatusOK, w.Result().StatusCode)

		Assert.Equal(stitchParams.DataShards, 2)
		Assert.Equal(stitchParams.ParityShards, 1)
		Assert.Equal(stitchParams.KeyThreshold, 2)
	})

}

func TestBackendController_RotateKey(t *testing.T) {

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
		Port:       "5559",
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

	data1 := "Hello this is the first version of the file. Hopefully it won't get corrupted poggies"
	data2 := "Hello this is the second version of the file. " +
		"It won't get corrupted like the other one because this file isn't cringe."

	var testFileID string

	t.Run("Create Files", func(t *testing.T) {

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

		// Update the file

		Assert.NoError(dbfstestutils.EXAMPLEUpdateFile(db, testFile, dbfstestutils.ExampleUpdate{
			NewSize:       70,
			NewActualSize: 70,
			FragmentPath:  stitchConfig.ShardsLocation,
			FileData:      data2,
			Server:        incConfig.ServerName,
		}, user))

	})

	t.Run("test that you can download the file", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/files/"+testFileID, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID":    testFileID,
			"versionID": "0",
		})
		w := httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		Assert.Equal(data1, w.Body.String())

	})

	t.Run("RotateKey", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/maintenance/stitch", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Body = ioutil.NopCloser(strings.NewReader(`{"file_id": "` + testFileID + `"}`))
		w := httptest.NewRecorder()
		bc.RotateKey(w, req)

		Assert.Equal(http.StatusOK, w.Code)

		// Check that the file is still downloadable

		req = httptest.NewRequest("GET", "/api/v1/files/"+testFileID, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": testFileID,
		})
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		Assert.Equal(data2, w.Body.String())

		// Check that the file is still downloadable (original version)

		req = httptest.NewRequest("GET", "/api/v1/files/"+testFileID, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID":    testFileID,
			"versionID": "0",
		})
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		Assert.Equal(data1, w.Body.String())

	})

	Inc.HttpServer.Close()

}
