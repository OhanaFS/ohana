package inc_test

import (
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/selfsign"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var certsConfigured = false
var tempDirConfigured = false
var tempDir string

func TestFragmentHandler(t *testing.T) {

	// Set params here
	filesFrags := make([]testFileFrag, 10)
	filesToBeVersionedAndDeleted := []int{0, 1, 2}
	filesToBeDeleted := []int{3, 4, 5}
	filesToBeVersioned := []int{6, 7, 8}
	//filesLeftAlone := []int{9}

	cleanupBefore := true
	cleanupOn := true

	db := testutil.NewMockDB(t)

	debugPrint := true
	if debugPrint {
		fmt.Println("Debug print on")
	}

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	// making Inc

	Assert := assert.New(t)
	tempdir, err := getTempDir()
	Assert.Nil(err)

	err = genCerts(tempdir)
	Assert.Nil(err)

	sharddir := filepath.Join(tempdir, "shards")

	stitchConfig := config.StitchConfig{
		ShardsLocation: sharddir,
	}

	if cleanupBefore {
		err := os.RemoveAll(tempdir)
		Assert.NoError(err)
		err = os.RemoveAll(stitchConfig.ShardsLocation)
		Assert.NoError(err)
	}

	if w, err := os.Stat(stitchConfig.ShardsLocation); os.IsNotExist(err) {
		err := os.MkdirAll(stitchConfig.ShardsLocation, 0755)
		if err != nil {
			panic("ERROR. CANNOT CREATE SHARDS FOLDER.")
		}
	} else if !w.IsDir() {
		panic("ERROR. SHARDS FOLDER IS NOT A DIRECTORY.")
	}

	configFile := &config.Config{Stitch: stitchConfig,
		Inc: config.IncConfig{
			ServerName: "testServer",
			HostName:   "localhost",
			Port:       "5555",
			CaCert:     tempdir + "/certificates/main_GLOBAL_CERTIFICATE.pem",
			PublicCert: tempdir + "/certificates/output_cert.pem",
			PrivateKey: tempdir + "/certificates/output_key.pem",
		},
	}

	// new zap logger
	logger, _ := zap.NewDevelopment()

	stitchParams, err := dbfs.GetStitchParams(db, logger)
	Assert.NoError(err)
	dataShards, parityShards, keyThreshold := stitchParams.DataShards, stitchParams.ParityShards, stitchParams.KeyThreshold
	totalShards := dataShards + parityShards

	Inc := inc.NewInc(configFile, db)

	t.Run("Running a Server for ping test", func(t *testing.T) {

		mux := http.NewServeMux()
		mux.HandleFunc("/inc/ping", inc.Pong)

		server := &http.Server{
			Addr:    ":" + configFile.Inc.Port,
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

		Assert.Equal(inc.Ping(configFile.Inc.HostName, configFile.Inc.Port), true)

	})

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
		err, s := Inc.CronJobDeleteFragments(true)
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

		err, s := Inc.CronJobDeleteFragments(true)
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

		err, s := Inc.CronJobDeleteFragments(true)
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

	if cleanupOn {
		t.Run("Cleanup", func(t *testing.T) {
			err := os.RemoveAll(tempdir)
			Assert.NoError(err)
			err = os.RemoveAll(Inc.ShardsLocation)
			Assert.NoError(err)
		})

	}

}

func genCerts(tempdir string) error {
	// Setting up certs for configs

	if certsConfigured {
		return nil
	} else {
		certsConfigured = true
	}

	ogc := config.LoadFlagsConfig()
	trueBool := true
	ogc.GenCA = &trueBool
	ogc.GenCerts = &trueBool
	tempDirCA := filepath.Join(tempdir, "certificates/main")
	ogc.GenCAPath = &tempDirCA
	tempDirCerts := filepath.Join(tempdir, "certificates/output")
	ogc.GenCertsPath = &tempDirCerts
	tempCertPath := filepath.Join(tempdir, "certificates/main_GLOBAL_CERTIFICATE.pem")
	ogc.CertPath = &tempCertPath
	tempPkPath := filepath.Join(tempdir, "certificates/main_PRIVATE_KEY.pem")
	ogc.PkPath = &tempPkPath
	tempCsrPath := filepath.Join(tempdir, "certificates/main_csr.json")
	ogc.CsrPath = &tempCsrPath
	tempHostsFile := filepath.Join(tempdir, "certhosts.yaml")
	ogc.AllHosts = &tempHostsFile

	fakeHosts := selfsign.Hosts{Hosts: []string{"localhost", "localhost2"}}

	hostFile, err := os.Create(filepath.Join(tempdir, "certhosts.yaml"))
	if err != nil {
		return err
	}
	defer hostFile.Close()

	encoder := yaml.NewEncoder(hostFile)
	if err := encoder.Encode(fakeHosts); err != nil {
		return err
	}

	return selfsign.ProcessFlags(ogc)
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

type testFileFrag struct {
	File      dbfs.File
	Fragments []dbfs.Fragment
}
