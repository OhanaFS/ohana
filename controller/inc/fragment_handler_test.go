package inc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	dbfstestutils "github.com/OhanaFS/ohana/dbfs/test_utils"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func TestFragmentHandler(t *testing.T) {

	// Set params here
	filesFrags := make([]testFileFrag, 10)
	filesToBeVersionedAndDeleted := []int{0, 1, 2}
	filesToBeDeleted := []int{3, 4, 5}
	filesToBeVersioned := []int{6, 7, 8}
	//filesLeftAlone := []int{9}

	db := testutil.NewMockDB(t)

	debugPrint := true
	if debugPrint {
		fmt.Println("Debug print on")
	}

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	Assert := assert.New(t)
	tempdir := t.TempDir()
	Assert.Nil(err)

	sharddir := filepath.Join(tempdir, "shards")

	stitchConfig := config.StitchConfig{
		ShardsLocation: sharddir,
	}

	if w, err := os.Stat(stitchConfig.ShardsLocation); os.IsNotExist(err) {
		err := os.MkdirAll(stitchConfig.ShardsLocation, 0755)
		if err != nil {
			panic("ERROR. CANNOT CREATE SHARDS FOLDER.")
		}
	} else if !w.IsDir() {
		panic("ERROR. SHARDS FOLDER IS NOT A DIRECTORY.")
	}

	certsPaths, err := selfsigntestutils.GenCertsTest(tempdir)
	Assert.Nil(err)

	configFile := &config.Config{Stitch: stitchConfig,
		Inc: config.IncConfig{
			ServerName: "testServer",
			HostName:   "localhost",
			Port:       "5555",
			CaCert:     certsPaths.CaCertPath,
			PublicCert: certsPaths.PublicCertPath,
			PrivateKey: certsPaths.PrivateKeyPath,
		},
	}

	// new zap logger
	logger, _ := zap.NewDevelopment()

	stitchParams, err := dbfs.GetStitchParams(db, logger)
	Assert.NoError(err)
	dataShards, parityShards, keyThreshold := stitchParams.DataShards, stitchParams.ParityShards, stitchParams.KeyThreshold
	totalShards := dataShards + parityShards

	Inc := inc.NewInc(configFile, db)
	inc.RegisterIncServices(Inc)

	t.Run("Creating random files", func(t *testing.T) {

		Assert := assert.New(t)

		// Creating files

		newFileFunc := func(filename string) (dbfs.File, []dbfs.Fragment) {

			// root folder
			rootFolder, err := dbfs.GetRootFolder(db)
			Assert.NoError(err)

			// This is the fileKey and fileIV for the passwordProtect
			fileKey, fileIv, err := dbfs.GenerateKeyIV()
			Assert.NoError(err)

			// This is the key and IV for the pipeline
			dataKey, dataIv, err := dbfs.GenerateKeyIV()
			Assert.NoError(err)

			// File, PasswordProtect entries for dbfs.
			dbfsFile := dbfs.File{
				FileId:             uuid.New().String(),
				FileName:           filename,
				MIMEType:           "",
				ParentFolderFileId: &rootFolder.FileId, // root folder for now
				Size:               int(5),
				VersioningMode:     dbfs.VersioningOff,
				TotalShards:        totalShards,
				DataShards:         dataShards,
				ParityShards:       parityShards,
				KeyThreshold:       keyThreshold,
				PasswordProtected:  false,
				HandledServer:      Inc.ServerName,
			}

			passwordProtect := dbfs.PasswordProtect{
				FileId:         dbfsFile.FileId,
				FileKey:        fileKey,
				FileIv:         fileIv,
				PasswordActive: false,
			}

			err = db.Transaction(func(tx *gorm.DB) error {

				err := dbfs.CreateInitialFile(tx, &dbfsFile, fileKey, fileIv, dataKey, dataIv, &superUser)
				if err != nil {
					return err
				}

				err = tx.Create(&passwordProtect).Error
				if err != nil {
					return err
				}

				return dbfs.CreatePermissions(tx, &dbfsFile)

			})
			Assert.NoError(err)

			err = dbfs.FinishFile(db, &dbfsFile, &superUser, 40, "")
			Assert.NoError(err)

			// Fragment creation

			for i := 1; i <= int(dbfsFile.TotalShards); i++ {

				shardName := dbfsFile.DataId + ".shard" + strconv.Itoa(i)
				shardPath := path.Join(Inc.ShardsLocation, shardName)
				shardFile, err := os.Create(shardPath)
				Assert.NoError(err)
				// write some crap
				_, err = shardFile.Write([]byte(fmt.Sprintf("Hello World %d" + strconv.Itoa(i))))
				Assert.NoError(err)
				err = shardFile.Close()
				Assert.NoError(err)

				err = dbfs.CreateFragment(db, dbfsFile.FileId, dbfsFile.DataId, dbfsFile.VersionNo, i, Inc.ServerName, shardName)
				Assert.NoError(err)
			}

			// Get Fragments
			fragments, err := dbfsFile.GetFileFragments(db, &superUser)

			return dbfsFile, fragments

		}

		for i := 0; i < 10; i++ {
			testFile, frag := newFileFunc("testFile" + strconv.Itoa(i))
			filesFrags[i] = testFileFrag{testFile, frag}
		}

		dir, err := os.ReadDir(Inc.ShardsLocation)
		Assert.NoError(err)
		Assert.Equal(len(filesFrags)*totalShards, len(dir))

		err = dbfs.SetHowLongToKeepFileVersions(db, 1)
		Assert.NoError(err)
		s, err := Inc.CronJobDeleteFragments(true)
		Assert.NoError(err)
		fmt.Println(s)

	})

	t.Run("Marking a few files to be versioned", func(t *testing.T) {

		Assert := assert.New(t)

		for _, i := range append(filesToBeVersionedAndDeleted, filesToBeVersioned...) {
			err := filesFrags[i].File.UpdateMetaData(db,
				dbfs.FileMetadataModification{VersioningMode: dbfs.VersioningOnVersions}, &superUser)
			Assert.NoError(err)
		}

	})

	t.Run("Updating all files sans those marked to be left alone", func(t *testing.T) {

		Assert := assert.New(t)

		updateFileTestFunc := func(file *dbfs.File) {

			err := file.UpdateFile(db, file.Size, file.Size, "ignore checksum",
				Inc.ServerName, "", "", "", &superUser)
			Assert.NoError(err)

			for i := 1; i <= int(file.TotalShards); i++ {
				shardName := file.DataId + ".shard" + strconv.Itoa(i)
				shardPath := path.Join(Inc.ShardsLocation, shardName)
				shardFile, err := os.Create(shardPath)
				Assert.NoError(err)

				// write some crap
				_, err = shardFile.Write([]byte(fmt.Sprintf("Hello World Updated %d" + strconv.Itoa(i))))
				Assert.NoError(err)
				err = shardFile.Close()
				Assert.NoError(err)

				err = file.UpdateFragment(db, i, shardName, "ignore checksum", Inc.ServerName)
				Assert.NoError(err)

			}

			err = file.FinishUpdateFile(db, "blah")
			Assert.NoError(err)

		}

		for _, i := range append(append(filesToBeVersionedAndDeleted, filesToBeVersioned...), filesToBeDeleted...) {
			updateFileTestFunc(&filesFrags[i].File)
		}

		countOfUpdatedFiles := len(filesToBeVersioned) + len(filesToBeVersionedAndDeleted) + len(filesToBeDeleted)

		dir, err := os.ReadDir(Inc.ShardsLocation)
		Assert.NoError(err)
		Assert.Equal(len(filesFrags)*totalShards+countOfUpdatedFiles*totalShards, len(dir))

	})

	t.Run("Deleting all files sans those marked to be left alone", func(t *testing.T) {

		Assert := assert.New(t)

		for _, i := range append(filesToBeVersionedAndDeleted, filesToBeDeleted...) {
			err := filesFrags[i].File.Delete(db, &superUser, Inc.ServerName)
			Assert.NoError(err)
		}

		// At this stage, if we run the cron job, it should delete filesToBeDeleted only which

		Inc.Db = db.Begin()

		s, err := Inc.CronJobDeleteFragments(true)
		Assert.NoError(err)
		fmt.Println("*******")
		fmt.Println(s)

		Inc.Db.Commit()
		Inc.Db = db // reset the db connection

		dir, err := os.ReadDir(Inc.ShardsLocation)
		Assert.NoError(err)

		initialShards := len(filesFrags) * totalShards
		countOfUpdatedFiles := len(filesToBeVersioned) + len(filesToBeVersionedAndDeleted) + len(filesToBeDeleted)
		countOfUpdatedShards := countOfUpdatedFiles * totalShards
		previousShardCount := initialShards + countOfUpdatedShards

		expectedRecords := previousShardCount - len(filesToBeVersionedAndDeleted)*totalShards*2 -
			len(filesToBeDeleted)*totalShards*2

		Assert.Equal(expectedRecords, len(dir))

	})

	t.Run("Changing old file version to be old enough to be cleaned up", func(t *testing.T) {

		Assert := assert.New(t)

		// Going to mark it indivisually, just easier

		for _, i := range filesToBeVersioned {

			result := db.Model(&dbfs.FileVersion{}).Where("file_id = ? AND version_no <> ?",
				filesFrags[i].File.FileId, filesFrags[i].File.VersionNo).Update("modified_time", time.Now().Add(time.Hour*24*-2))

			Assert.Equal(int64(2), result.RowsAffected)
			Assert.NoError(result.Error)
		}

		rows, err := dbfs.MarkOldFileVersions(db)
		Assert.Equal(int64(2*len(filesToBeVersioned)), rows)

		Inc.Db = db.Begin()

		s, err := Inc.CronJobDeleteFragments(true)
		Assert.NoError(err)
		fmt.Println("*******")
		fmt.Println(s)

		Inc.Db.Commit()
		Inc.Db = db // reset the db connection

		dir, err := os.ReadDir(Inc.ShardsLocation)
		Assert.NoError(err)

		initialShards := len(filesFrags) * totalShards
		countOfUpdatedFiles := len(filesToBeVersioned) + len(filesToBeVersionedAndDeleted) + len(filesToBeDeleted)
		countOfUpdatedShards := countOfUpdatedFiles * totalShards
		countOfDeletedShards := len(filesToBeVersionedAndDeleted)*totalShards*2 +
			len(filesToBeDeleted)*totalShards*2
		previousShardCount := initialShards + countOfUpdatedShards - countOfDeletedShards

		expectedRecords := previousShardCount - len(filesToBeVersioned)*totalShards

		Assert.Equal(expectedRecords, len(dir))

	})

	t.Run("Testing LocalOrphanedShards", func(t *testing.T) {

		Assert := assert.New(t)

		results, err := Inc.LocalOrphanedShardsCheck(-1, false)
		Assert.NoError(err)
		Assert.Len(results, 0)

		// Test via routes

		req := httptest.NewRequest(http.MethodGet,
			inc.FragmentOrphanedPath, nil)
		w := httptest.NewRecorder()
		req.Header.Add("job_id", "0") // TODO: Create a value incrementer
		Inc.OrphanedShardsRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)
		results2, err := dbfs.GetResultsOrphanedShard(db, 0)
		Assert.Nil(err)
		Assert.Equal(0, len(results2), results2)

		// add a weird file into the shards folder

		shardName := "weird.shard"
		shardPath := path.Join(Inc.ShardsLocation, shardName)
		shardFile, err := os.Create(shardPath)
		Assert.NoError(err)

		// write some crap
		_, err = shardFile.Write([]byte("Hello World"))
		Assert.NoError(err)
		Assert.NoError(shardFile.Close())

		// check to see if it is in the list of orphaned shards
		results, err = Inc.LocalOrphanedShardsCheck(-1, false)
		Assert.NoError(err)
		Assert.Len(results, 1)

		// Test via routes

		req = httptest.NewRequest(http.MethodGet,
			inc.FragmentOrphanedPath, nil)
		w = httptest.NewRecorder()
		req.Header.Add("job_id", "1") // TODO: Create a value incrementer
		Inc.OrphanedShardsRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)
		results2, err = dbfs.GetResultsOrphanedShard(db, 1)
		Assert.Nil(err)
		Assert.Equal(1, len(results2), results2)

		// Delete that weird file
		Assert.NoError(os.Remove(shardPath))

	})

	t.Run("Testing LocalMissingShards", func(t *testing.T) {

		Assert := assert.New(t)

		results, err := Inc.LocalMissingShardsCheck(-2, false)
		Assert.NoError(err)
		Assert.Len(results, 0)

		// Testing via route
		req := httptest.NewRequest(http.MethodGet,
			inc.FragmentMissingPath, nil)
		w := httptest.NewRecorder()
		req.Header.Add("job_id", "-1")
		Inc.MissingShardsRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)
		results, err = dbfs.GetResultsMissingShard(db, -1)
		Assert.Nil(err)
		Assert.Equal(0, len(results), results)

		// Delete
		dir, err := os.ReadDir(Inc.ShardsLocation)
		Assert.NoError(err)
		Assert.NoError(os.Remove(path.Join(Inc.ShardsLocation, dir[0].Name())))

		results, err = Inc.LocalMissingShardsCheck(-1, false)
		Assert.NoError(err)
		Assert.Len(results, 1)

		// Testing via route
		req = httptest.NewRequest(http.MethodGet,
			inc.FragmentMissingPath, nil)
		w = httptest.NewRecorder()
		req.Header.Add("job_id", "0")
		Inc.MissingShardsRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)
		results, err = dbfs.GetResultsMissingShard(db, 0)
		Assert.Nil(err)
		Assert.Equal(1, len(results), results)

	})

	t.Run("Cleanup", func(t *testing.T) {
		err = os.RemoveAll(Inc.ShardsLocation)
		Assert.NoError(err)
	})

	defer Inc.HttpServer.Shutdown(context.Background())

}

func TestStitchFragment(t *testing.T) {

	db := testutil.NewMockDB(t)

	debugPrint := true
	if debugPrint {
		fmt.Println("Debug print on")
	}

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	tempdir := t.TempDir()
	assert.Nil(t, err)

	sharddir := filepath.Join(tempdir, "shards")

	stitchConfig := config.StitchConfig{
		ShardsLocation: sharddir,
	}

	if w, err := os.Stat(stitchConfig.ShardsLocation); os.IsNotExist(err) {
		err := os.MkdirAll(stitchConfig.ShardsLocation, 0755)
		if err != nil {
			panic("ERROR. CANNOT CREATE SHARDS FOLDER.")
		}
	} else if !w.IsDir() {
		panic("ERROR. SHARDS FOLDER IS NOT A DIRECTORY.")
	}

	certsPaths, err := selfsigntestutils.GenCertsTest(tempdir)
	assert.Nil(t, err)

	configFile := &config.Config{Stitch: stitchConfig,
		Inc: config.IncConfig{
			ServerName: "testServer",
			HostName:   "localhost",
			Port:       "5555",
			CaCert:     certsPaths.CaCertPath,
			PublicCert: certsPaths.PublicCertPath,
			PrivateKey: certsPaths.PrivateKeyPath,
		},
	}

	Inc := inc.NewInc(configFile, db)
	inc.RegisterIncServices(Inc)

	logger := config.NewLogger(configFile)
	bc := &controller.BackendController{
		Db:         db,
		Logger:     logger,
		Path:       configFile.Stitch.ShardsLocation,
		ServerName: configFile.Inc.ServerName,
		Inc:        Inc,
	}
	bc.InitialiseShardsFolder()

	// create a new file

	rootFolder, err := dbfs.GetRootFolder(db)

	file, err := dbfstestutils.EXAMPLECreateFile(db, &superUser, dbfstestutils.ExampleFile{
		FileName:       "test123",
		ParentFolderId: rootFolder.FileId,
		Server:         configFile.Inc.ServerName,
		FragmentPath:   configFile.Stitch.ShardsLocation,
		FileData:       "Blah123",
		Size:           50,
		ActualSize:     50,
	})

	jobIdNo := 0

	t.Run("Checking health of fragments", func(t *testing.T) {
		Assert := assert.New(t)

		// Test that all fragments are healthy

		jobIdNo = jobIdNo + 1

		Assert.Nil(err)
		Assert.NotNil(file)
		err := Inc.LocalCurrentFilesFragmentsHealthCheck(jobIdNo)

		results, err := dbfs.GetResultsCffhc(db, jobIdNo)
		Assert.NoError(err)
		Assert.NotNil(results)
		Assert.Equal(0, len(results))

		// Testing via route

		jobIdNo = jobIdNo + 1

		req := httptest.NewRequest(http.MethodGet,
			inc.CurrentFilesHealthPath, nil)
		w := httptest.NewRecorder()
		req.Header.Add("job_id", strconv.Itoa(jobIdNo))
		Inc.CurrentFilesFragmentsHealthCheckRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)

		results, err = dbfs.GetResultsCffhc(db, jobIdNo)
		Assert.NoError(err)
		Assert.NotNil(results)
		Assert.Equal(0, len(results))

		// damange the file
		fragments, err := file.GetFileFragments(db, &superUser)
		Assert.NoError(err)
		Assert.NotNil(fragments)

		// Check again.

		jobIdNo = jobIdNo + 1

		err = dbfstestutils.EXAMPLECorruptFragments(
			path.Join(configFile.Stitch.ShardsLocation, fragments[0].FileFragmentPath))
		Assert.Nil(err)

		err = Inc.LocalCurrentFilesFragmentsHealthCheck(jobIdNo)
		Assert.NoError(err)

		results, err = dbfs.GetResultsCffhc(db, jobIdNo)
		Assert.Nil(err)
		Assert.NotNil(results)
		Assert.Equal(1, len(results))

		// Checking via route

		jobIdNo = jobIdNo + 1

		req = httptest.NewRequest(http.MethodGet,
			inc.CurrentFilesHealthPath, nil)
		w = httptest.NewRecorder()
		req.Header.Add("job_id", strconv.Itoa(jobIdNo))
		Inc.CurrentFilesFragmentsHealthCheckRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)

		results, err = dbfs.GetResultsCffhc(db, jobIdNo)
		Assert.NoError(err)
		Assert.NotNil(results)
		Assert.Equal(1, len(results))

		// Check using InvidiaulFragmentHealthCheck
		result, err := Inc.IndividualFragHealthCheck(fragments[0])
		Assert.Nil(err)
		Assert.NotNil(result)
		Assert.True(len(result.BrokenBlocks) > 0)

		result, err = Inc.IndividualFragHealthCheck(fragments[1])
		Assert.Nil(err)
		Assert.NotNil(result)
		Assert.True(len(result.BrokenBlocks) == 0)

		// Check bad fragment via route

		req = httptest.NewRequest("GET",
			strings.Replace(inc.FragmentHealthCheckPath,
				"{fragmentPath}", fragments[0].FileFragmentPath, -1), nil)
		w = httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{
			"fragmentPath": fragments[0].FileFragmentPath,
		})
		fmt.Println(fragments[0].FileFragmentPath)
		Inc.FragmentHealthCheckRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body := w.Body.String()

		err = json.Unmarshal([]byte(body), &result)
		Assert.Nil(err)

		Assert.True(len(result.BrokenBlocks) > 0)

		// Check good fragment via route

		req = httptest.NewRequest("GET",
			strings.Replace(inc.FragmentHealthCheckPath,
				"{fragmentPath}", fragments[1].FileFragmentPath, -1), nil)
		w = httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{
			"fragmentPath": fragments[1].FileFragmentPath,
		})
		Inc.FragmentHealthCheckRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		err = json.Unmarshal([]byte(body), &result)
		Assert.Nil(err)

		Assert.True(len(result.BrokenBlocks) == 0)

		err = db.Find(&fragments).Error
		fmt.Println(err)
	})

	t.Run("Checking health of all fragments", func(t *testing.T) {

		Assert := assert.New(t)
		strings := make([]string, 0)

		// We are going to update the file we created and corrupted earlier.
		// This should cause the health check to fail for LocalAllFilesFragmentsHealthCheck
		// but not for LocalCurrentFilesFragmentsHealthCheck

		// We are also going to create the edge case, where a file that is corrupted is coppied
		// (thus the new file is linked to the same data ID and fragments)
		// and then update the new file

		// Create new folder for copied file
		newFolder, err := rootFolder.CreateSubFolder(db, "blah", &superUser, "omgServer")
		Assert.Nil(err)

		// Copy file
		err = file.Copy(db, newFolder, &superUser, "omgServer")
		Assert.Nil(err)

		// Get the copied file
		_, err = dbfs.GetFileByPath(db, "/blah/test123", &superUser, false)

		jobIdNo = jobIdNo + 1

		Assert.NoError(Inc.LocalCurrentFilesFragmentsHealthCheck(jobIdNo))

		results, err := dbfs.GetResultsCffhc(db, jobIdNo)
		Assert.NoError(err)
		Assert.Equal(1, len(results), results)
		Assert.NoError(json.Unmarshal([]byte(results[0].FileId), &strings))
		Assert.Equal(2, len(strings))
		Assert.True(strings[0] != strings[1])

		Assert.NoError(dbfstestutils.EXAMPLEUpdateFile(db, file, dbfstestutils.ExampleUpdate{
			NewSize:       50,
			NewActualSize: 50,
			FragmentPath:  configFile.Stitch.ShardsLocation,
			FileData:      "New file data pog",
			Server:        configFile.Inc.ServerName,
			Password:      "",
		}, &superUser))

		jobIdNo = jobIdNo + 1

		Assert.Nil(Inc.LocalCurrentFilesFragmentsHealthCheck(jobIdNo))

		results, err = dbfs.GetResultsCffhc(db, jobIdNo)
		Assert.NoError(err)
		Assert.Equal(1, len(results), results)
		Assert.NoError(json.Unmarshal([]byte(results[0].FileId), &strings))
		Assert.Equal(1, len(strings))

		jobIdNo = jobIdNo + 1

		// Checking All File Fragments

		Assert.NoError(Inc.LocalAllFilesFragmentsHealthCheck(jobIdNo))

		results2, err := dbfs.GetResultsAffhc(db, jobIdNo)
		Assert.Nil(err)
		Assert.NotNil(results2)
		Assert.Equal(1, len(results2))
		Assert.NoError(json.Unmarshal([]byte(results[0].FileId), &strings))
		Assert.Equal(1, len(strings))

		// Testing via route

		jobIdNo = jobIdNo + 1

		req := httptest.NewRequest(http.MethodGet,
			inc.AllFilesHealthPath, nil)
		w := httptest.NewRecorder()
		req.Header.Add("job_id", strconv.Itoa(jobIdNo))
		Inc.AllFilesFragmentsHealthCheckRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)

		results2, err = dbfs.GetResultsAffhc(db, jobIdNo)
		Assert.Nil(err)
		Assert.NotNil(results2)
		Assert.Equal(1, len(results2))
		Assert.NoError(json.Unmarshal([]byte(results[0].FileId), &strings))
		Assert.Equal(1, len(strings))

	})

	t.Run("Deleteing a fragment with the Route", func(t *testing.T) {

		Assert := assert.New(t)

		// Get the number of files in the directory
		files, err := ioutil.ReadDir(configFile.Stitch.ShardsLocation)
		Assert.NoError(err)
		Assert.True(len(files) > 0, files)

		// original count of files
		originalCount := len(files)
		firstFileName := files[0].Name()

		req := httptest.NewRequest("DELETE",
			strings.Replace(inc.FragmentPath,
				"{fragmentPath}", firstFileName, -1), nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{
			"fragmentPath": firstFileName,
		})
		Inc.DeleteFragmentRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		Assert.True(strings.Contains(w.Body.String(), "true"), w.Body.String())

		// Make sure it's actually deleted

		files, err = ioutil.ReadDir(configFile.Stitch.ShardsLocation)
		Assert.NoError(err)
		Assert.True(len(files) == originalCount-1, files)
		Assert.NotEqualf(firstFileName, files[0].Name(), "File not deleted")

	})

	t.Run("Missing shards check", func(t *testing.T) {

		// In the previous test, we should have deleted a fragment, thus it should be missing
		// now we are going to check if the missing fragment check works
		Assert := assert.New(t)

		jobIdNo = jobIdNo + 1

		req := httptest.NewRequest(http.MethodGet,
			inc.FragmentMissingPath, nil)
		w := httptest.NewRecorder()
		req.Header.Add("job_id", strconv.Itoa(jobIdNo))
		Inc.MissingShardsRoute(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		time.Sleep(time.Second / 2)
		results, err := dbfs.GetResultsMissingShard(db, jobIdNo)
		Assert.Nil(err)
		Assert.Equal(1, len(results), results)

	})

}

type testFileFrag struct {
	File      dbfs.File
	Fragments []dbfs.Fragment
}
