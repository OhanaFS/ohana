package controller_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

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
)

var certsConfigured = false
var tempDirConfigured = false
var tempDir string

func TestAdminClusterHistoricalRoutes(t *testing.T) {

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

	// Load the db with HistoricalStats that can be used for testing

	fakeStartLoadDataDate := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	fakeEndLoadDataDate := time.Now()

	// days between fakeStartLoadDataDate and fakeEndLoadDataDate
	daysBetween := int(fakeEndLoadDataDate.Sub(fakeStartLoadDataDate).Hours() / 24)
	tempLoadData := make([]dbfs.HistoricalStats, daysBetween)
	num := 0

	for (fakeStartLoadDataDate.Unix() - fakeEndLoadDataDate.Unix()) < 0 {
		// load fake data in the db

		tempD := fakeStartLoadDataDate.Day()
		tempM := int(fakeStartLoadDataDate.Month())
		tempY := fakeStartLoadDataDate.Year()

		if num >= len(tempLoadData) {
			tempLoadData = append(tempLoadData, dbfs.HistoricalStats{})
		}

		tempLoadData[num] = dbfs.HistoricalStats{
			Day:            tempD,
			Month:          tempM,
			Year:           tempY,
			NonReplicaUsed: int64(tempD + tempM + tempY),
			ReplicaUsed:    int64(tempD+tempM+tempY) * 2,
			NumOfFiles:     int64(tempD+tempM+tempY) * 3,
		}
		num++

		fakeStartLoadDataDate = fakeStartLoadDataDate.Add(time.Hour * 24)
	}

	// dump data in
	assert.NoError(t, db.Save(&tempLoadData).Error)

	t.Run("GetHistoricalFiles", func(t *testing.T) {

		type testRoute func(http.ResponseWriter, *http.Request)

		testHistorical := func(t *testing.T, functionToTest testRoute, testType int) {

			urlString := ""

			switch testType {
			case dbfs.HistoricalNonReplicaUsed:
				{
					urlString = "/api/v1/cluster/stats/historical/non_replica_used_historical"
				}
			case dbfs.HistoricalReplicaUsed:
				{
					urlString = "/api/v1/cluster/stats/historical/replica_used_historical"
				}
			case dbfs.HistoricalNumOfFiles:
				{
					urlString = "/api/v1/cluster/stats/historical/num_of_files_historical"
				}
			}

			Assert := assert.New(t)

			// Get number of files via route (day)

			req := httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "1")
			w := httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body := w.Body.String()

			var numOfFilesHistorical []dbfs.DateInt64Value
			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(10, len(numOfFilesHistorical))
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Today's date vs 10 days ago
			todayDate := time.Now()
			tenDaysAgo := todayDate.AddDate(0, 0, -9)

			Assert.Equal(todayDate.Year(), numOfFilesHistorical[len(numOfFilesHistorical)-1].Year())
			Assert.Equal(int(todayDate.Month()), numOfFilesHistorical[len(numOfFilesHistorical)-1].Month())
			Assert.Equal(todayDate.Day(), numOfFilesHistorical[len(numOfFilesHistorical)-1].Day())

			Assert.Equal(tenDaysAgo.Year(), numOfFilesHistorical[0].Year())
			Assert.Equal(int(tenDaysAgo.Month()), numOfFilesHistorical[0].Month())
			Assert.Equal(tenDaysAgo.Day(), numOfFilesHistorical[0].Day())
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get number of files via route (week)

			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "2")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(10, len(numOfFilesHistorical))
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get number of files via route (month)

			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "3")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(10, len(numOfFilesHistorical))
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get two days of data
			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "1")
			req.Header.Add("start_date", "2020-01-01")
			req.Header.Add("end_date", "2020-01-02")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(2, len(numOfFilesHistorical))
			Assert.Equal(1, numOfFilesHistorical[0].Day())
			Assert.Equal(2, numOfFilesHistorical[1].Day())

			// Get 1 days of data
			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "1")
			req.Header.Add("start_date", "2020-01-01")
			req.Header.Add("end_date", "2020-01-01")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(1, len(numOfFilesHistorical))
			Assert.Equal(1, numOfFilesHistorical[0].Day())
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get two weeks of data
			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "2")
			req.Header.Add("start_date", "2022-07-31")
			req.Header.Add("end_date", "2022-08-12")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(2, len(numOfFilesHistorical))
			Assert.Equal(7, numOfFilesHistorical[0].Month())
			Assert.Equal(31, numOfFilesHistorical[0].Day())
			Assert.Equal(8, numOfFilesHistorical[1].Month())
			Assert.Equal(7, numOfFilesHistorical[1].Day())
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get one week of data
			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "2")
			req.Header.Add("start_date", "2022-07-31")
			req.Header.Add("end_date", "2022-07-31")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(1, len(numOfFilesHistorical))
			Assert.Equal(7, numOfFilesHistorical[0].Month())
			Assert.Equal(31, numOfFilesHistorical[0].Day())
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get two months of data
			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "3")
			req.Header.Add("start_date", "2020-01-01")
			req.Header.Add("end_date", "2020-02-01")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(2, len(numOfFilesHistorical))
			Assert.Equal(1, numOfFilesHistorical[0].Month())
			Assert.Equal(2, numOfFilesHistorical[1].Month())
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

			// Get one month of data
			req = httptest.NewRequest("GET", urlString,
				nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req.Header.Add("range_type", "3")
			req.Header.Add("start_date", "2020-01-01")
			req.Header.Add("end_date", "2020-01-01")
			w = httptest.NewRecorder()

			functionToTest(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			Assert.NoError(json.Unmarshal([]byte(body), &numOfFilesHistorical))
			Assert.Equal(1, len(numOfFilesHistorical))
			Assert.Equal(1, numOfFilesHistorical[0].Month())
			Assert.True(AssertCorrectEntry(numOfFilesHistorical, testType))

		}

		testHistorical(t, bc.GetNumOfFilesHistorical, dbfs.HistoricalNumOfFiles)

		testHistorical(t, bc.GetStorageUsedHistorical, dbfs.HistoricalNonReplicaUsed)

		testHistorical(t, bc.GetStorageUsedReplicaHistorical, dbfs.HistoricalReplicaUsed)

	})

	// Close everything
	Inc.HttpServer.Close()

}

func TestAdminClusterRoutes(t *testing.T) {

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
			FragmentPath:  Inc.ShardsLocation,
			FileData:      "new 123 " + strconv.Itoa(i),
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
	logCount := fatalCount + errorCount + warningCount + infoCount + debugCount + traceCount

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

		// Get number of files via route

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/numOfFiles", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetNumOfFiles(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()
		Assert.Contains(body, strconv.Itoa(filesToCreate))

		// Get number of files via DumpDailyStats
		Assert.NoError(dbfs.DumpDailyStats(db))
		stats, err := dbfs.GetTodayStat(db)
		Assert.NoError(err)
		Assert.Equal(int64(filesToCreate), stats.NumOfFiles)

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

		// Get via DumpDailyStats
		Assert.NoError(dbfs.DumpDailyStats(db))
		stats, err := dbfs.GetTodayStat(db)
		Assert.NoError(err)
		Assert.Equal(int64(newSize*filesToCreate), stats.NonReplicaUsed)

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

		// Get via DumpDailyStats
		Assert.NoError(dbfs.DumpDailyStats(db))
		stats, err := dbfs.GetTodayStat(db)
		Assert.NoError(err)
		Assert.Equal(int64(newActualSize*filesToCreate+oldActualSize*filesToCreate), stats.ReplicaUsed)

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

	t.Run("GetAllLogs", func(t *testing.T) {

		Assert := assert.New(t)

		// Next page

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		req.Header.Add("start_num", "0")
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		var logs []dbfs.Log
		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(50, len(logs))

		// Page 2

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		req.Header.Add("start_num", "1")
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(50, len(logs))

		// Page 3

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		req.Header.Add("start_num", "2")
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(logCount-100, len(logs))

		// setting a start date

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		req.Header.Add("start_num", "0")
		req.Header.Add("start_date", "2002-10-02T15:00:00Z")
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(50, len(logs))

		// Setting an end date (should return none)

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		req.Header.Add("start_num", "2")
		req.Header.Add("start_num", "0")
		req.Header.Add("end_date", "2002-10-02T15:00:00Z")
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(0, len(logs))

		// server_filter

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		req.Header.Add("start_num", "0")
		req.Header.Add("server_filter", "localServer")
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(50, len(logs))

		// server_filter and type_filter

		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		req.Header.Add("start_num", "0")
		req.Header.Add("server_filter", "localServer")
		req.Header.Add("type_filter", strconv.Itoa(int(dbfs.LogServerFatal)))
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(fatalCount, len(logs))

	})

	t.Run("GetLog", func(t *testing.T) {
		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/logs/1", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})
		w := httptest.NewRecorder()
		bc.GetLog(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		var log dbfs.Log
		Assert.NoError(json.Unmarshal([]byte(body), &log))
		Assert.Equal(1, log.LogId)
	})

	t.Run("ClearLog", func(t *testing.T) {
		Assert := assert.New(t)

		req := httptest.NewRequest("DELETE", "/api/v1/cluster/stats/logs/1", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})
		w := httptest.NewRecorder()
		bc.ClearLog(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.Contains(body, "true")
	})

	t.Run("ClearAllLogs", func(t *testing.T) {
		Assert := assert.New(t)

		req := httptest.NewRequest("DELETE", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.ClearAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.Contains(body, "true")

		// Check that the logs are actually cleared
		req = httptest.NewRequest("GET", "/api/v1/cluster/stats/logs", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetAllLogs(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		var logs []dbfs.Log

		Assert.NoError(json.Unmarshal([]byte(body), &logs))
		Assert.Equal(0, len(logs))

	})

	t.Run("GetServerStatuses", func(t *testing.T) {
		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/server_statuses", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetServerStatuses(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		var serverStatuses []dbfs.Server

		Assert.NoError(json.Unmarshal([]byte(body), &serverStatuses))
		Assert.Equal(1, len(serverStatuses)) //
		fmt.Println(serverStatuses)
		Assert.Equal("localServer", serverStatuses[0].Name)
		Assert.Equal(Inc.HostName, serverStatuses[0].HostName)
		Assert.Equal(Inc.Port, serverStatuses[0].Port)
		Assert.Equal(dbfs.ServerOnline, serverStatuses[0].Status)

	})

	t.Run("GetSpecificServerStatus", func(t *testing.T) {
		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/server_statuses/localServer", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"serverName": "localServer",
		})
		w := httptest.NewRecorder()
		bc.GetSpecificServerStatus(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		var serverStatus dbfs.Server

		fmt.Println(body)

		Assert.NoError(json.Unmarshal([]byte(body), &serverStatus))
		Assert.Equal("localServer", serverStatus.Name)
		Assert.Equal(Inc.HostName, serverStatus.HostName)
		Assert.Equal(Inc.Port, serverStatus.Port)
		Assert.Equal(dbfs.ServerOnline, serverStatus.Status)

	})

	t.Run("CronDeleteFragments", func(t *testing.T) {

		Assert := assert.New(t)

		req := httptest.NewRequest("GET", "/api/v1/cluster/stats/cron/delete_fragments", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()

		// before this can start we need to start a transaction because stupid gorm and sqlite

		bc.Inc.Db = db.Begin()

		bc.CronDeleteFragments(w, req)

		bc.Inc.Db.Commit()
		bc.Inc.Db = db

		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		Assert.Contains(body, "20")

	})

	/*
		t.Run("DeleteServer", func(t *testing.T) {
			Assert := assert.New(t)

			req := httptest.NewRequest("DELETE", "/api/v1/cluster/stats/server_statuses/localServer", nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			req = mux.SetURLVars(req, map[string]string{
				"serverName": "localServer",
			})
			w := httptest.NewRecorder()
			bc.DeleteServer(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body := w.Body.String()

			Assert.Contains(body, "true")

			// Check that the server is actually deleted
			req = httptest.NewRequest("GET", "/api/v1/cluster/stats/server_statuses", nil).WithContext(
				ctxutil.WithUser(context.Background(), user))
			req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
			w = httptest.NewRecorder()
			bc.GetServerStatuses(w, req)
			Assert.Equal(http.StatusOK, w.Code)
			body = w.Body.String()

			var serverStatuses []dbfs.Server

			Assert.NoError(json.Unmarshal([]byte(body), &serverStatuses))
			Assert.Equal("localServer", serverStatuses[0].Name)
			Assert.Equal(Inc.HostName, serverStatuses[0].HostName)
			Assert.Equal(Inc.Port, serverStatuses[0].Port)
			Assert.Equal(dbfs.ServerOffline, serverStatuses[0].Status)
		})
	*/
	os.RemoveAll(tempDir)

	Inc.HttpServer.Close()

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

func AssertCorrectEntry(values []dbfs.DateInt64Value, typeOfRequest int) bool {

	/*
		NonReplicaUsed: int64(tempD + tempM + tempY),
		ReplicaUsed:    int64(tempD+tempM+tempY) * 2,
		NumOfFiles:     int64(tempD+tempM+tempY) * 3,
	*/

	if len(values) == 0 {
		return true
	}

	for _, v := range values {

		switch typeOfRequest {
		case dbfs.HistoricalNonReplicaUsed:
			{
				if v.Value != int64(v.Day()+v.Month()+v.Year()) {
					return false
				}
			}
		case dbfs.HistoricalReplicaUsed:
			{
				if v.Value != int64(v.Day()+v.Month()+v.Year())*2 {
					return false
				}
			}
		case dbfs.HistoricalNumOfFiles:
			{
				if v.Value != int64(v.Day()+v.Month()+v.Year())*3 {
					return false
				}
			}

		}
	}

	return true

}
