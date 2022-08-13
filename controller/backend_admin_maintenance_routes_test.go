package controller_test

import (
	"context"
	"encoding/json"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
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

	t.Run("StartJob", func(t *testing.T) {

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
		assert.NoError(t, db.Find(&jobs).Error)

		Assert := assert.New(t)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &newJob))
		Assert.Equal(newJob.JobId, uint(2))
		Assert.Equal(newJob.AllFilesShardsHealthCheck, true)

	})

	t.Run("Get Full Shards Result", func(t *testing.T) {

		time.Sleep(1 * time.Second)
		Assert := assert.New(t)

		// Get the result
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
		var result []dbfs.ResultsAFSHC
		Assert.NoError(json.Unmarshal([]byte(body), &result))

		Assert.Equal(len(result), 0)
	})

}
