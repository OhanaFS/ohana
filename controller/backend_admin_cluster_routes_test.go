package controller_test

import (
	"context"
	"encoding/json"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	dbfstestutils "github.com/OhanaFS/ohana/dbfs/testutils"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var certsConfigured = false
var tempDirConfigured = false
var tempDir string

func TestAdminClusterRoutes(t *testing.T) {

	// Setting up env

	tempDir, err := getTempDir()
	assert.NoError(t, err)
	shardDir := filepath.Join(tempDir, "shards")

	// TODO: Clean this up.

	//Set up mock Db and session store
	stitchConfig := config.StitchConfig{
		ShardsLocation: shardDir,
	}
	incConfig := config.IncConfig{
		ServerName: "localServer",
		HostName:   "localhost",
		BindIp:     "127.0.0.1",
		Port:       "5555",
		CaCert:     "",
		PublicCert: "",
		PrivateKey: "",
	}

	configFile := &config.Config{Stitch: stitchConfig, Inc: incConfig}
	logger := config.NewLogger(configFile)
	db := testutil.NewMockDB(t)

	session := testutil.NewMockSession(t)
	sessionId, err := session.Create(nil, "superuser", time.Hour)
	Inc := inc.NewInc(configFile, db)

	// Setting up controller
	bc := &controller.BackendController{
		Db:         db,
		Logger:     logger,
		Path:       configFile.Stitch.ShardsLocation,
		ServerName: configFile.Inc.ServerName,
	}

	bc.InitialiseShardsFolder()

	// Getting Superuser to use with testing
	user, err := dbfs.GetUser(db, "superuser")
	assert.NoError(t, err)

	// Setting up files and logs

	// get root folder
	rootFolder, err := dbfs.GetRootFolder(db)

	// some consts
	const (
		oldSize       = 500
		oldActualSize = 800
		newSize       = 600
		newActualSize = 900
		filesToCreate = 20
		fatalCount    = 5
		errorCount    = 10
		warningCount  = 15
		infoCount     = 20
		debugCount    = 25
		traceCount    = 30
	)

	// create some files

	for i := 0; i < filesToCreate; i++ {
		f, err := dbfstestutils.EXAMPLECreateFile(db, user, dbfstestutils.ExampleFile{
			FileName:       "blank" + strconv.Itoa(i) + ".txt",
			ParentFolderId: rootFolder.FileId,
			Server:         Inc.ServerName,
			FragmentPath:   Inc.ShardsLocation,
			FileData:       "blank 123 " + strconv.Itoa(i),
			Size:           oldSize,
			ActualSize:     oldActualSize,
		})
		assert.NoError(t, err)
		err = dbfstestutils.EXAMPLEUpdateFile(db, f, dbfstestutils.ExampleUpdate{
			NewSize:       newSize,
			NewActualSize: newActualSize,
			Server:        Inc.ServerName,
			Password:      "",
		}, user)
	}

	// Creating some logs

	logsToCreate := map[int8]int{
		dbfs.LogServerFatal:   fatalCount,
		dbfs.LogServerError:   errorCount,
		dbfs.LogServerWarning: warningCount,
		dbfs.LogServerInfo:    infoCount,
		dbfs.LogServerDebug:   debugCount,
		dbfs.LogServerTrace:   traceCount,
	}

	alertCount := fatalCount + errorCount + warningCount

	dbfsLogger := dbfs.NewLogger(db, Inc.ServerName)

	for logType, count := range logsToCreate {
		for i := 0; i < count; i++ {
			switch logType {
			case dbfs.LogServerFatal:
				dbfsLogger.LogFatal("fatal log " + strconv.Itoa(i))
			case dbfs.LogServerError:
				dbfsLogger.LogError("error log " + strconv.Itoa(i))
			case dbfs.LogServerWarning:
				dbfsLogger.LogWarning("warning log " + strconv.Itoa(i))
			case dbfs.LogServerInfo:
				dbfsLogger.LogInfo("info log " + strconv.Itoa(i))
			case dbfs.LogServerDebug:
				dbfsLogger.LogDebug("debug log " + strconv.Itoa(i))
			case dbfs.LogServerTrace:
				dbfsLogger.LogTrace("trace log	" + strconv.Itoa(i))
			}
		}
	}

	t.Run("GetNumOfFiles", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/numOfFiles", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetNumOfFiles(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, strconv.Itoa(filesToCreate))

	})

	t.Run("GetStorageUsed", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/nonReplicaUsed", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetStorageUsed(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, strconv.Itoa(newSize*filesToCreate))

	})

	t.Run("GetStorageUsedReplica", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/replica", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetStorageUsedReplica(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, strconv.Itoa(newActualSize*filesToCreate+oldActualSize*filesToCreate))

	})

	// Need to be able to get alert or warning per server... maybe patch it into the get of /servers/servername/alerts?

	var alertID int

	t.Run("GetAllAlerts", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/alerts", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetAllAlerts(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		var alerts []dbfs.Alert
		Assert.NoError(json.Unmarshal([]byte(body), &alerts))
		Assert.Equal(alertCount, len(alerts))

		alertID = alerts[0].LogId

	})

	t.Run("GetAlert", func(t *testing.T) {

		Assert := assert.New(t)

		getNum := strconv.Itoa(alertID)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/alerts/"+getNum, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": getNum,
		})
		w := httptest.NewRecorder()
		bc.GetAlert(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		var alert dbfs.Alert
		Assert.NoError(json.Unmarshal([]byte(body), &alert))
		Assert.Equal(alertID, alert.LogId)

	})

	t.Run("ClearAlert", func(t *testing.T) {

		Assert := assert.New(t)

		getNum := strconv.Itoa(alertID)

		req := httptest.NewRequest("DELETE", "/api/v1/cluster/stats/alerts/"+getNum, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": getNum,
		})
		w := httptest.NewRecorder()
		bc.ClearAlert(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

		// Checking that it is gone (should be false)

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/alerts/"+getNum, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": getNum,
		})
		w = httptest.NewRecorder()
		bc.GetAlert(w, req)
		Assert.Equal(http.StatusNotFound, w.Code)

	})

	t.Run("ClearAllAlerts", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("DELETE", "/api/v1/cluster/stats/alerts", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.ClearAllAlerts(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, "true")

		// Checking that it got deleted

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/alerts", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetAllAlerts(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		var alerts []dbfs.Alert
		Assert.NoError(json.Unmarshal([]byte(body), &alerts))
		Assert.Equal(len(alerts), 0)
	})

}

func getTempDir() (string, error) {

	if tempDirConfigured {
		return tempDir, nil
	} else {
		tempDirConfigured = true
	}

	tempDir, err := ioutil.TempDir("", "ohana-test")
	if err != nil {
		return "", err
	}
	return tempDir, nil
}
