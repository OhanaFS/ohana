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

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("username = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	t.Run("Getting root folder", func(t *testing.T) {

		Assert := assert.New(t)

		file, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		Assert.Equal("00000000-0000-0000-0000-000000000000", file.FileID)
		Assert.Equal("root", file.FileName)

		jsonOutput, err := json.Marshal(file)

		Assert.Nil(err)

		if debugPrint {
			fmt.Println(string(jsonOutput))
		}

	})

	t.Run("Creating a few folders at the root folder", func(t *testing.T) {
		Assert := assert.New(t)

		rootFolder, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		// Folder 1

		newFolder1, err := dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test1", &superUser)

		Assert.Nil(err)

		Assert.Equal("Test1", newFolder1.FileName)

		// Folder 2

		newFolder2, err := dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test2", &superUser)

		Assert.Nil(err)

		Assert.Equal("Test2", newFolder2.FileName)

		// Folder 3

		newFolder3, err := dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test3", &superUser)

		Assert.Nil(err)

		Assert.Equal("Test3", newFolder3.FileName)

		// Folder 4 (nested inside Folder 3)

		newFolder4, err := dbfs.CreateFolderByParentID(db, newFolder3.FileID, "Test4", &superUser)

		Assert.Nil(err)

		Assert.Equal("Test4", newFolder4.FileName)

		// Folder 5 (same name as Folder 2, should fail due to existing folder)

		_, err = dbfs.CreateFolderByParentID(db, rootFolder.FileID, "Test3", &superUser)

		Assert.Error(dbfs.ErrFileFolderExists, err)

		// Listing files

		ls, err := dbfs.ListFilesByFolderID(db, rootFolder.FileID, &superUser)

		Assert.Nil(err)

		Assert.Equal(3, len(ls))

		ls, err = dbfs.ListFilesByFolderID(db, newFolder3.FileID, &superUser)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))
		Assert.Equal("Test4", ls[0].FileName)

	})

	t.Run("GetFileByPath", func(t *testing.T) {
		// Pre-requisite: Creating a few folders at the root folder already ran

		// Checking for right file

		Assert := assert.New(t)

		file, err := dbfs.GetFileByPath(db, "/Test3/Test4", &superUser)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "/Test3/Test4/", &superUser)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/Test4", &superUser)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/Test4/", &superUser)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/", &superUser)
		Assert.Nil(err)
		Assert.Equal("Test3", file.FileName)

		// Checking for non-existent file

		file, err = dbfs.GetFileByPath(db, "/Test3/Test5", &superUser)
		Assert.Error(dbfs.ErrFileNotFound)

	})

	t.Run("ListFilesByPath", func(t *testing.T) {
		// Pre-requisite: Creating a few folders at the root folder already ran

		// Checking for right file

		Assert := assert.New(t)

		files, err := dbfs.ListFilesByPath(db, "/", &superUser)
		Assert.Nil(err)
		Assert.Equal(3, len(files))

		files, err = dbfs.ListFilesByPath(db, "Test3/", &superUser)
		Assert.Nil(err)
		Assert.Equal(1, len(files))

		files, err = dbfs.ListFilesByPath(db, "Test3/Test4", &superUser)
		Assert.Nil(err)
		Assert.Equal(0, len(files))

		// Checking for non-existent file

		files, err = dbfs.ListFilesByPath(db, "/Test3/Test2", &superUser)
		Assert.Error(dbfs.ErrFileNotFound)

	})

	t.Run("CreateFolderByPath", func(t *testing.T) {
		// Creating a folder at /Test2/Test_2

		Assert := assert.New(t)

		ls, err := dbfs.ListFilesByPath(db, "Test2/", &superUser)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

		innerFolder, err := dbfs.CreateFolderByPath(db, "/Test2/Test_2", &superUser)
		Assert.Nil(err)
		Assert.Equal("Test_2", innerFolder.FileName)

		ls, err = dbfs.ListFilesByPath(db, "Test2/", &superUser)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))

	})

	t.Run("Delete Folder By ID", func(t *testing.T) {

		db := testutil.NewMockDB(t)
		superUser := dbfs.User{}

		// Getting superuser account
		err := db.Where("username = ?", "superuser").First(&superUser).Error
		assert.Nil(t, err)

		// DeleteFolderByID

		Assert := assert.New(t)

		// Creating /Test2/Test2_2

		folder, err := dbfs.CreateFolderByPath(db, "/Test2", &superUser)
		Assert.Nil(err)
		folder, err = dbfs.CreateFolderByPath(db, "/Test2/Test_2", &superUser)
		Assert.Nil(err)

		// Delete single folder.

		ls, err := dbfs.ListFilesByPath(db, "/Test2", &superUser)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))

		folder, err = dbfs.GetFileByPath(db, "/Test2/Test_2", &superUser)
		Assert.Nil(err)

		err = dbfs.DeleteFolderByID(db, folder.FileID, &superUser)
		Assert.Nil(err)

		ls, err = dbfs.ListFilesByPath(db, "/Test2", &superUser)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

		// Delete non-existent folder

		err = dbfs.DeleteFolderByID(db, "100c09f5-e163-468b-ad34-80944bbf8dfa", &superUser)
		Assert.Error(dbfs.ErrFileNotFound, err)

		// Creating folders for next test

		_, err = dbfs.CreateFolderByPath(db, "/Test3", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_1", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_2", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_1/POG", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB/TestB_1", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB/TestB_2", &superUser)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestC", &superUser)
		Assert.Nil(err)

		// Delete Folder that has contents inside "/Test3"

		ls, err = dbfs.ListFilesByPath(db, "/Test3", &superUser)
		Assert.Nil(err)
		Assert.Equal(3, len(ls))

		parentFolder, err := dbfs.GetFileByPath(db, "/Test3", &superUser)
		Assert.Nil(err)

		err = dbfs.DeleteFolderByID(db, parentFolder.FileID, &superUser)

		Assert.Error(dbfs.ErrFolderNotEmpty, err)

		// TODO: Delete folder which one has no permissions to

		// DeleteFolderByIDCascade

		err = dbfs.DeleteFolderByIDCascade(db, parentFolder.FileID, &superUser)
		Assert.Nil(err)

		ls, err = dbfs.ListFilesByPath(db, "/Test3", &superUser)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

	})

	t.Run("Checking permissions for users", func(t *testing.T) {

		Assert := assert.New(t)

		// Making folders

		parentFolder, err := dbfs.CreateFolderByPath(db, "/TestPerms", &superUser)
		Assert.Nil(err)

		subFolder1, err := dbfs.CreateFolderByPath(db, "/TestPerms/perm1", &superUser)
		Assert.Nil(err)

		subFolder2, err := dbfs.CreateFolderByPath(db, "/TestPerms/perm2", &superUser)
		Assert.Nil(err)

		// Creating a user
		user1, err := dbfs.CreateNewUser(db, "permissionCheckUser", "user1Name", 1, "permissionCheckUser")
		Assert.Nil(err)

		// giving said user permissions to /TestPerms/perm1

		err = subFolder1.AddPermissionUsers(db, &dbfs.PermissionNeeded{
			Read: true,
		}, &superUser, *user1)
		Assert.Nil(err)

		// Check that user has permissions to /TestPerms/perm1 but not /TestPerms/perm2 and /TestPerms
		hasPermission, err := user1.HasPermission(db, subFolder1, &dbfs.PermissionNeeded{Read: true})
		Assert.Equal(true, hasPermission)
		Assert.Nil(err)
		hasPermission, err = user1.HasPermission(db, parentFolder, &dbfs.PermissionNeeded{Read: true})
		Assert.Error(dbfs.ErrFileNotFound, err)
		Assert.Equal(false, hasPermission)
		hasPermission, err = user1.HasPermission(db, subFolder2, &dbfs.PermissionNeeded{Read: true})
		Assert.Error(dbfs.ErrFileNotFound, err)
		Assert.Equal(false, hasPermission)

		// Creating a folder under /TestPerms/perm1

		subFolder3, err := dbfs.CreateFolderByPath(db, "/TestPerms/perm1/perm3", user1)
		Assert.Error(err, dbfs.ErrNoPermission)
		subFolder3, err = dbfs.CreateFolderByPath(db, "/TestPerms/perm1/perm3", &superUser)
		Assert.Nil(err)

		// Check that user has permissions to /TestPerms/perm1/perm3
		hasPermission, err = user1.HasPermission(db, subFolder3, &dbfs.PermissionNeeded{Read: true})
		Assert.Equal(true, hasPermission)
		Assert.Nil(err)
		hasPermission, err = superUser.HasPermission(db, subFolder3, &dbfs.PermissionNeeded{Read: true})
		Assert.Equal(true, hasPermission)
		Assert.Nil(err)

	})

	// Creating fake files
	t.Run("Creating fake files", func(t *testing.T) {

		Assert := assert.New(t)

		// Making folders
		newFolder, err := dbfs.CreateFolderByPath(db, "/TestFakeFiles", &superUser)
		Assert.Nil(err)

		// Creating fake files
		file, err := dbfs.EXAMPLECreateFile(db, &superUser, "somefile.txt", newFolder.FileID)

		Assert.Nil(err)
		// Check if you can get da fragments
		fragments, err := file.GetFileFragments(db, &superUser)
		Assert.Nil(err)
		Assert.Equal(5, len(fragments))

		if debugPrint {
			fmt.Println(fragments)
		}

	})

}
