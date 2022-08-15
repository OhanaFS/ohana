package dbfs_test

import (
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"sync"
	"testing"
	"time"
)

func TestDeletion(t *testing.T) {

	// init db

	db := testutil.NewMockDB(t)

	debugPrint := true

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	// Making folders
	newFolder, err := dbfs.CreateFolderByPath(db, "/TestFakeFiles", &superUser, "testserver", false)
	assert.NoError(t, err)

	// Creating fake files
	file1, err := EXAMPLECreateFile(db, &superUser, "somefile1.txt", newFolder.FileId)
	assert.NoError(t, err)

	file2, err := EXAMPLECreateFile(db, &superUser, "somefile2.txt", newFolder.FileId)
	assert.NoError(t, err)

	file3, err := EXAMPLECreateFile(db, &superUser, "somefile3.txt", newFolder.FileId)
	assert.NoError(t, err)

	// Create copy of file2
	// gonna put it in root because I laze.

	rootFolder, err := dbfs.GetRootFolder(db)
	assert.NoError(t, err)

	err = file2.Copy(db, rootFolder, &superUser, "testserver")
	assert.NoError(t, err)

	// get file2 copy

	file2Copy, err := dbfs.GetFileByPath(db, "somefile2.txt", &superUser, false)
	assert.NoError(t, err)

	// Attempting to run cron job without any deletion function
	t.Run("NoDeletion", func(t *testing.T) {

		Assert := assert.New(t)

		fragments, err := dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)
		Assert.Equal(0, len(fragments))

	})

	// Deleting files

	t.Run("Delete Files", func(t *testing.T) {

		Assert := assert.New(t)

		// Ensuring it works with multiple versions
		err = EXAMPLEUpdateFile(db, file1, "", &superUser)
		Assert.NoError(err)

		// Deletion on file2 copy
		err := file1.Delete(db, &superUser, "deleteServer")
		Assert.NoError(err)
		err = file2.Delete(db, &superUser, "deleteServer")
		Assert.NoError(err)
		err = file3.Delete(db, &superUser, "deleteServer")
		Assert.NoError(err)

		if debugPrint {
			fmt.Println(file2Copy.FileName)
		}

	})

	// Getting fragments that should be deleted

	t.Run("Get To Be Deleted Fragments", func(t *testing.T) {

		Assert := assert.New(t)

		// We should expect 3 file's fragments to be collected.

		fragments, err := dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)
		Assert.Equal(ExampleTotalShards*3, len(fragments))

		// Let's try deleting file2 copy

		err = file2Copy.Delete(db, &superUser, "deleteServer")
		Assert.NoError(err)

		// We should expect totalShards to be 3 x default
		fragments, err = dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)
		Assert.Equal(ExampleTotalShards*4, len(fragments))

	})

	t.Run("Delete Fragments", func(t *testing.T) {

		Assert := assert.New(t)

		fragments, err := dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)

		dataIdFragmentMap := make(map[string][]dbfs.Fragment)

		for _, fragment := range fragments {
			dataIdFragmentMap[fragment.FileVersionDataId] = append(dataIdFragmentMap[fragment.FileVersionDataId], fragment)
		}

		// "Deletion Code"

		const maxGoroutines = 10
		input := make(chan string, len(dataIdFragmentMap))
		output := make(chan string, len(dataIdFragmentMap))

		// Worker function

		for i := 0; i < maxGoroutines; i++ {
			go deleteWorker(dataIdFragmentMap, input, output)
		}

		for dataId, _ := range dataIdFragmentMap {
			input <- dataId
		}
		close(input)

		err = db.Transaction(func(tx *gorm.DB) error {

			for i := 0; i < len(dataIdFragmentMap); i++ {
				dataIdProcessed := <-output

				// Create transaction
				err2 := dbfs.FinishDeleteDataId(tx, dataIdProcessed)
				Assert.NoError(err2)
			}
			return nil
		})
		Assert.NoError(err)

		// Checking to ensure that no other fragments are left

		fragments, err = dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)
		Assert.Equal(0, len(fragments))

	})

	t.Run("Test MarkOldFileVersions", func(t *testing.T) {

		Assert := assert.New(t)

		// Checking that attribute is not marked
		rows, err := dbfs.MarkOldFileVersions(db)
		Assert.NoError(err)
		Assert.Equal(int64(0), rows)

		// Setting bad attribute
		err = dbfs.SetHowLongToKeepFileVersions(db, -1)
		Assert.ErrorIs(err, dbfs.ErrInvalidCronJobProperty)

		// Setting good attribute
		err = dbfs.SetHowLongToKeepFileVersions(db, 1)
		Assert.NoError(err)

		// Checking that attribute is marked
		rows, err = dbfs.MarkOldFileVersions(db)
		Assert.NoError(err)
		Assert.Equal(int64(0), rows)

		// Creating 3 files.

		// Get root folder
		rootFolder, err := dbfs.GetRootFolder(db)

		file1, err := EXAMPLECreateFile(db, &superUser, "deletefile1.txt", rootFolder.FileId)
		Assert.NoError(err)
		file2, err := EXAMPLECreateFile(db, &superUser, "deletefile2.txt", rootFolder.FileId)
		Assert.NoError(err)
		_, err = EXAMPLECreateFile(db, &superUser, "dontDeleteFile3.txt", rootFolder.FileId)
		Assert.NoError(err)

		// Turning on versioning for file 1 and file 2

		err = file1.UpdateMetaData(db, dbfs.FileMetadataModification{VersioningMode: dbfs.VersioningOnVersions}, &superUser)
		Assert.NoError(err)
		err = file2.UpdateMetaData(db, dbfs.FileMetadataModification{VersioningMode: dbfs.VersioningOnVersions}, &superUser)
		Assert.NoError(err)

		// Updating file 1 and file 2

		err = EXAMPLEUpdateFile(db, file1, "", &superUser)
		Assert.NoError(err)
		err = EXAMPLEUpdateFile(db, file2, "", &superUser)
		Assert.NoError(err)

		// Checking that file 1 and file 2 are not marked for deletion
		rows, err = dbfs.MarkOldFileVersions(db)
		Assert.NoError(err)
		Assert.Equal(int64(0), rows)

		// Updating file 1's fileversion to be old

		// Getting the latest file version
		file1, err = dbfs.GetFileById(db, file1.FileId, &superUser)

		// Manually updating the file version to be old

		result := db.Model(&dbfs.FileVersion{}).Where("file_id = ? AND version_no <> ?",
			file1.FileId, file1.VersionNo).Update("modified_time", time.Now().Add(time.Hour*24*-2))
		Assert.NoError(result.Error)

		// Checking that file 1 is marked for deletion
		rows, err = dbfs.MarkOldFileVersions(db)
		Assert.NoError(err)
		Assert.Equal(int64(2), rows) // should be 2 because the change from no versioning to versioning also counts

		rows, err = dbfs.MarkOldFileVersions(db)
		Assert.NoError(err)
		Assert.Equal(int64(0), rows) // should be 2 because the change from no versioning to versioning also counts

	})

}

func TestCheckOrphanedFiles(t *testing.T) {
	db := testutil.NewMockDB(t)

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	// Setting up the environment

	// Get root folder
	rootFolder, err := dbfs.GetRootFolder(db)
	assert.Nil(t, err)

	// Making a few folders

	folders := make([]*dbfs.File, 5)

	for i := 0; i < 5; i++ {
		folders[i], err = rootFolder.CreateSubFolder(db, fmt.Sprintf("folder%d", i),
			&superUser, "testServer")
		assert.Nil(t, err)
	}

	// Making a few files
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			tempFolder, err := folders[i].CreateSubFolder(db, fmt.Sprintf("folder_%d", i),
				&superUser, "testServer")
			assert.Nil(t, err)
			_, err = EXAMPLECreateFile(db, &superUser, "file_lvl3"+fmt.Sprintf("%d", i), tempFolder.FileId)
			assert.Nil(t, err)
		}
		_, err = EXAMPLECreateFile(db, &superUser, "file_lvl2"+fmt.Sprintf("%d", i), folders[i].FileId)
		assert.Nil(t, err)
	}

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

	var orphanedFilesResult []dbfs.ResultsOrphanedFile

	t.Run("CheckOrphanedFiles", func(t *testing.T) {

		Assert := assert.New(t)
		orphanedFilesResult, err = dbfs.CheckOrphanedFiles(db, 2, true)
		Assert.NoError(err)
		Assert.Equal(2, len(orphanedFilesResult), orphanedFilesResult)

	})

	t.Run("FixOrphanedFiles", func(t *testing.T) {

		Assert := assert.New(t)

		// Make the fixes to be passed to FixOrphanedFiles
		// we are going to do a different fix for each file

		// Make a list of fixes

		fixes := make([]dbfs.OrphanedFilesActions, len(orphanedFilesResult))
		fixes[0] = dbfs.OrphanedFilesActions{
			ParentFolderId: orphanedFilesResult[0].ParentFolderId,
			Delete:         true,
		}
		fixes[1] = dbfs.OrphanedFilesActions{
			ParentFolderId: orphanedFilesResult[1].ParentFolderId,
			Move:           true,
		}

		// Send the results
		err := dbfs.FixOrphanedFiles(db, 2, fixes)
		Assert.NoError(err)

		// There should be a new folder in the database

		rootContents, err := rootFolder.ListContents(db, &superUser)
		Assert.NoError(err)
		fmt.Println(rootContents)

		checkFolder, err := dbfs.GetFileByPath(db, "/Orphaned Files", &superUser, false)
		Assert.NoError(err)
		Assert.NotNil(checkFolder)

		// There should be 1 folder inside the folder
		contents, err := checkFolder.ListContents(db, &superUser)
		Assert.NoError(err)
		Assert.Equal(1, len(contents))
		contentsFile, err := contents[0].ListContents(db, &superUser)

		Assert.Equal(file2.FileId, contentsFile[0].FileId)
		Assert.Contains(orphanedFilesResult[1].Contents, file2.FileName)

	})

}

func deleteWorker(dataIdFragmentMap map[string][]dbfs.Fragment, input <-chan string, output chan<- string) {

	for j := range input {
		var fragWg sync.WaitGroup

		for _, fragment := range dataIdFragmentMap[j] {
			fragWg.Add(1)
			go func(path, server, dataId string) {
				// Assume this is a delete fragment call.
				fmt.Println("Deleting fragment:", path, server, dataId)
				defer fragWg.Done()
			}(fragment.FileFragmentPath, fragment.ServerName, fragment.FileVersionDataId)
		}

		fragWg.Wait()
		output <- j
	}

}
