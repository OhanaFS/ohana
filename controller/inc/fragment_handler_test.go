package inc_test

import (
	"context"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

		_, err := Inc.LocalOrphanedShardsCheck()
		Assert.NoError(err)

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
		_, err = Inc.LocalOrphanedShardsCheck()
		Assert.ErrorIs(err, inc.ErrOrphanedShardsFound)

		// Delete that weird file
		Assert.NoError(os.Remove(shardPath))

	})

	t.Run("Testing LocalMissingShards", func(t *testing.T) {

		Assert := assert.New(t)

		_, err := Inc.LocalMissingShardsCheck()
		Assert.NoError(err)

		dir, err := os.ReadDir(Inc.ShardsLocation)
		Assert.NoError(err)
		Assert.NoError(os.Remove(path.Join(Inc.ShardsLocation, dir[0].Name())))

		_, err = Inc.LocalMissingShardsCheck()
		Assert.ErrorIs(err, inc.ErrMissingShardsFound)

	})

	t.Run("Cleanup", func(t *testing.T) {
		err = os.RemoveAll(Inc.ShardsLocation)
		Assert.NoError(err)
	})

	defer Inc.HttpServer.Shutdown(context.Background())

}

type testFileFrag struct {
	File      dbfs.File
	Fragments []dbfs.Fragment
}
