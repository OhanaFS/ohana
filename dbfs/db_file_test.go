package dbfs_test

import (
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T) {
	db := testutil.NewMockDB(t)

	debugPrint := true

	FAKEUSER := dbfs.User{}

	t.Run("Getting root folder", func(t *testing.T) {

		Assert := assert.New(t)

		file, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		Assert.Equal("00000000-0000-0000-0000-000000000000", file.FileID)
		Assert.Equal("root", file.FileName)

		json, err := json.Marshal(file)

		Assert.Nil(err)

		if debugPrint {
			fmt.Println(string(json))
		}

	})

	t.Run("Creating a few folders at the root folder", func(t *testing.T) {
		Assert := assert.New(t)

		rootFolder, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		// Folder 1

		newFolder1, err := dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test1", FAKEUSER)

		Assert.Nil(err)

		Assert.Equal("Test1", newFolder1.FileName)

		// Folder 2

		newFolder2, err := dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test2", FAKEUSER)

		Assert.Nil(err)

		Assert.Equal("Test2", newFolder2.FileName)

		// Folder 3

		newFolder3, err := dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test3", FAKEUSER)

		Assert.Nil(err)

		Assert.Equal("Test3", newFolder3.FileName)

		// Folder 4 (nested inside Folder 3)

		newFolder4, err := dbfs.CreateFolderByParentID(db, newFolder3.FileID, "Test4", FAKEUSER)

		Assert.Nil(err)

		Assert.Equal("Test4", newFolder4.FileName)

		// Folder 5 (same name as Folder 2, should fail due to existing folder)

		_, err = dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test3", FAKEUSER)

		Assert.Error(dbfs.ErrFolderExists, err)

		// Listing files

		ls, err := dbfs.ListFilesByFolderID(db, rootFolder.FileID, FAKEUSER)

		Assert.Nil(err)

		Assert.Equal(3, len(ls))

		ls, err = dbfs.ListFilesByFolderID(db, newFolder3.FileID, FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))
		Assert.Equal("Test4", ls[0].FileName)

	})

	t.Run("GetFileByPath", func(t *testing.T) {
		// Pre-requisite: Creating a few folders at the root folder already ran

		// Checking for right file

		Assert := assert.New(t)

		file, err := dbfs.GetFileByPath(db, "/Test3/Test4", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "/Test3/Test4/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/Test4", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/Test4/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal("Test3", file.FileName)

		// Checking for non-existent file

		file, err = dbfs.GetFileByPath(db, "/Test3/Test5", FAKEUSER)
		Assert.Error(dbfs.ErrFileNotFound)

	})

	t.Run("ListFilesByPath", func(t *testing.T) {
		// Pre-requisite: Creating a few folders at the root folder already ran

		// Checking for right file

		Assert := assert.New(t)

		files, err := dbfs.ListFilesByPath(db, "/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(3, len(files))

		files, err = dbfs.ListFilesByPath(db, "Test3/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(1, len(files))

		files, err = dbfs.ListFilesByPath(db, "Test3/Test4", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(0, len(files))

		// Checking for non-existent file

		files, err = dbfs.ListFilesByPath(db, "/Test3/Test2", FAKEUSER)
		Assert.Error(dbfs.ErrFileNotFound)

	})

	t.Run("CreateFolderByPath", func(t *testing.T) {
		// Creating a folder at /Test2/Test_2

		Assert := assert.New(t)

		ls, err := dbfs.ListFilesByPath(db, "Test2/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

		innerFolder, err := dbfs.CreateFolderByPath(db, "/Test2/Test_2", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal("Test_2", innerFolder.FileName)

		ls, err = dbfs.ListFilesByPath(db, "Test2/", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))

	})

	t.Run("Delete Folder By ID", func(t *testing.T) {

		db := testutil.NewMockDB(t)

		// DeleteFolderByID

		Assert := assert.New(t)

		// Creating /Test2/Test2_2

		folder, err := dbfs.CreateFolderByPath(db, "/Test2", FAKEUSER)
		Assert.Nil(err)
		folder, err = dbfs.CreateFolderByPath(db, "/Test2/Test_2", FAKEUSER)
		Assert.Nil(err)

		// Delete single folder.

		ls, err := dbfs.ListFilesByPath(db, "/Test2", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))

		folder, err = dbfs.GetFileByPath(db, "/Test2/Test_2", FAKEUSER)
		Assert.Nil(err)

		err = dbfs.DeleteFolderByID(db, folder.FileID, FAKEUSER)
		Assert.Nil(err)

		ls, err = dbfs.ListFilesByPath(db, "/Test2", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

		// Delete non-existent folder

		err = dbfs.DeleteFolderByID(db, "adjslkfjakldfjkl", FAKEUSER)
		Assert.Error(dbfs.ErrFileNotFound, err)

		// Creating folders for next test

		_, err = dbfs.CreateFolderByPath(db, "/Test3", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_1", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_2", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_1/POG", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB/TestB_1", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB/TestB_2", FAKEUSER)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestC", FAKEUSER)
		Assert.Nil(err)

		// Delete Folder that has contents inside "/Test3"

		ls, err = dbfs.ListFilesByPath(db, "/Test3", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(3, len(ls))

		parentFolder, err := dbfs.GetFileByPath(db, "/Test3", FAKEUSER)
		Assert.Nil(err)

		err = dbfs.DeleteFolderByID(db, parentFolder.FileID, FAKEUSER)

		Assert.Error(dbfs.ErrFolderNotEmpty, err)

		// TODO: Delete folder which one has no permissions to

		// DeleteFolderByIDCascade

		err = dbfs.DeleteFolderByIDCascade(db, parentFolder.FileID, FAKEUSER)
		Assert.Nil(err)

		ls, err = dbfs.ListFilesByPath(db, "/Test3", FAKEUSER)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

	})
}
