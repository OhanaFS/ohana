package dbfs_test

import (
	"context"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"sync"
	"testing"
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
		err := file1.Delete(db, &superUser)
		Assert.NoError(err)
		err = file2.Delete(db, &superUser)
		Assert.NoError(err)
		err = file3.Delete(db, &superUser)
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

		err = file2Copy.Delete(db, &superUser)
		Assert.NoError(err)

		// We should expect totalShards to be 3 x default
		fragments, err = dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)
		Assert.Equal(ExampleTotalShards*4, len(fragments))

	})

	t.Run("Delete Fragments", func(t *testing.T) {

		Assert := assert.New(t)

		// This doesn't work with multiple goroutines still :/
		db := ctxutil.GetTransaction(context.Background(), db)

		fragments, err := dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)

		dataIdFragmentMap := make(map[string][]dbfs.Fragment)

		for _, fragment := range fragments {
			dataIdFragmentMap[fragment.FileVersionDataId] = append(dataIdFragmentMap[fragment.FileVersionDataId], fragment)
		}

		// "Deletion Code"

		// dataWg is the wait group for all the goroutines / threads
		var dataWg sync.WaitGroup

		// each dataId has its own goroutine
		for dataId, dataIdFragments := range dataIdFragmentMap {
			dataWg.Add(1)

			go func(dataId2 string, dataIdFragments2 []dbfs.Fragment, db2 *gorm.DB) {

				defer dataWg.Done()

				// fragWg is the waitgroup for the dataId
				var fragWg sync.WaitGroup

				// each fragment has its own goroutine
				for _, fragment := range dataIdFragments2 {
					fragWg.Add(1)
					go func(path, server, dataId string) {
						// Assume this is a delete fragment call.
						fmt.Println("Deleting fragment:", path, server, dataId)
						defer fragWg.Done()
					}(fragment.FileFragmentPath, fragment.ServerName, fragment.FileVersionDataId)
				}
				fragWg.Wait()

				fmt.Println("Deleted fragments for dataId:", dataId2)

				// This doesn't work. Table gets locked still and returns no such table.
				/*
					db3 := ctxutil.GetTransaction(context.Background(), db2)
					err = dbfs.FinishDeleteDataId(db3, dataId2)
					Assert.NoError(err)

				*/

			}(dataId, dataIdFragments, db)
		}

		dataWg.Wait()

		// THIS IS HERE BECAUSE IT GETS LOCKED IN THE GOROUTINES AND IDK WHY.
		for dataId, _ := range dataIdFragmentMap {
			err = dbfs.FinishDeleteDataId(db, dataId)
		}

		// Checking to ensure that no other fragments are left

		fragments, err = dbfs.GetToBeDeletedFragments(db)
		Assert.NoError(err)
		Assert.Equal(0, len(fragments))

	})

}
