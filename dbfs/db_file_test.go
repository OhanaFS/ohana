package dbfs_test

import (
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

const (
	ExampleTotalShards  = 5
	ExampleDataShards   = 3
	ExampleParityShards = 2
	ExampleKeyThreshold = 2
)

func TestFile(t *testing.T) {
	db := testutil.NewMockDB(t)

	debugPrint := true

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.Nil(t, err)

	t.Run("Getting root folder", func(t *testing.T) {

		Assert := assert.New(t)

		file, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		Assert.Equal("00000000-0000-0000-0000-000000000000", file.FileId)
		Assert.Equal("root", file.FileName)

		//jsonOutput, err := json.Marshal(file)
		//
		//Assert.Nil(err)
		//
		//if debugPrint {
		//	fmt.Println(string(jsonOutput))
		//}

	})

	t.Run("Creating a few folders at the root folder", func(t *testing.T) {
		Assert := assert.New(t)

		rootFolder, err := dbfs.GetRootFolder(db)

		Assert.Nil(err)

		// Folder 1

		newFolder1, err := dbfs.CreateFolderByParentId(db, rootFolder.FileId, "Test1", &superUser, "testserver")

		Assert.Nil(err)

		Assert.Equal("Test1", newFolder1.FileName)

		// Folder 2

		newFolder2, err := dbfs.CreateFolderByParentId(db, rootFolder.FileId, "Test2", &superUser, "testserver")

		Assert.Nil(err)

		Assert.Equal("Test2", newFolder2.FileName)

		// Folder 3

		newFolder3, err := dbfs.CreateFolderByParentId(db, rootFolder.FileId, "Test3", &superUser, "testserver")

		Assert.Nil(err)

		Assert.Equal("Test3", newFolder3.FileName)

		// Folder 4 (nested inside Folder 3)

		newFolder4, err := dbfs.CreateFolderByParentId(db, newFolder3.FileId, "Test4", &superUser, "testserver")

		Assert.Nil(err)

		Assert.Equal("Test4", newFolder4.FileName)

		// Folder 5 (same name as Folder 2, should fail due to existing folder)

		_, err = dbfs.CreateFolderByParentId(db, rootFolder.FileId, "Test3", &superUser, "testserver")

		Assert.Error(dbfs.ErrFileFolderExists, err)

		// Listing files

		ls, err := dbfs.ListFilesByFolderId(db, rootFolder.FileId, &superUser)

		Assert.Nil(err)

		Assert.Equal(5, len(ls))

		ls, err = dbfs.ListFilesByFolderId(db, newFolder3.FileId, &superUser)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))
		Assert.Equal("Test4", ls[0].FileName)

	})

	t.Run("GetFileByPath", func(t *testing.T) {
		// Pre-requisite: Creating a few folders at the root folder already ran

		// Checking for right file

		Assert := assert.New(t)

		file, err := dbfs.GetFileByPath(db, "/Test3/Test4", &superUser, false)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "/Test3/Test4/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/Test4", &superUser, false)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/Test4/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal("Test4", file.FileName)

		file, err = dbfs.GetFileByPath(db, "Test3/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal("Test3", file.FileName)

		// Checking for non-existent file

		file, err = dbfs.GetFileByPath(db, "/Test3/Test5", &superUser, false)
		Assert.Error(dbfs.ErrFileNotFound)

	})

	t.Run("ListFilesByPath", func(t *testing.T) {
		// Pre-requisite: Creating a few folders at the root folder already ran

		// Checking for right file

		Assert := assert.New(t)

		files, err := dbfs.ListFilesByPath(db, "/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(5, len(files))

		files, err = dbfs.ListFilesByPath(db, "Test3/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(1, len(files))

		files, err = dbfs.ListFilesByPath(db, "Test3/Test4", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(0, len(files))

		// Checking for non-existent file

		files, err = dbfs.ListFilesByPath(db, "/Test3/Test2", &superUser, false)
		Assert.Error(dbfs.ErrFileNotFound)

	})

	t.Run("CreateFolderByPath", func(t *testing.T) {
		// Creating a folder at /Test2/Test_2

		Assert := assert.New(t)

		ls, err := dbfs.ListFilesByPath(db, "Test2/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

		innerFolder, err := dbfs.CreateFolderByPath(db, "/Test2/Test_2", &superUser, "testserver", false)
		Assert.Nil(err)
		Assert.Equal("Test_2", innerFolder.FileName)

		ls, err = dbfs.ListFilesByPath(db, "Test2/", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))

		// json marshall test
		jsonOutput, err := json.Marshal(innerFolder)
		fmt.Println(string(jsonOutput))

	})

	t.Run("Delete Folder By Id", func(t *testing.T) {

		db := testutil.NewMockDB(t)
		superUser := dbfs.User{}

		// Getting superuser account
		err := db.Where("email = ?", "superuser").First(&superUser).Error
		assert.Nil(t, err)

		// DeleteFolderById

		Assert := assert.New(t)

		// Creating /Test2/Test2_2

		folder, err := dbfs.CreateFolderByPath(db, "/Test2", &superUser, "testserver", false)
		Assert.Nil(err)
		folder, err = dbfs.CreateFolderByPath(db, "/Test2/Test_2", &superUser, "testserver", false)

		Assert.Nil(err)

		// Delete single folder.

		ls, err := dbfs.ListFilesByPath(db, "/Test2", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(1, len(ls))

		folder, err = dbfs.GetFileByPath(db, "/Test2/Test_2", &superUser, false)
		Assert.Nil(err)

		err = dbfs.DeleteFolderById(db, folder.FileId, &superUser)
		Assert.Nil(err)

		ls, err = dbfs.ListFilesByPath(db, "/Test2", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

		// Delete non-existent folder

		err = dbfs.DeleteFolderById(db, "100c09f5-e163-468b-ad34-80944bbf8dfa", &superUser)
		Assert.Error(dbfs.ErrFileNotFound, err)

		// Creating folders for next test

		_, err = dbfs.CreateFolderByPath(db, "/Test3", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_1", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_2", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestA/TestA_1/POG", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB/TestB_1", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestB/TestB_2", &superUser, "testserver", false)
		Assert.Nil(err)
		_, err = dbfs.CreateFolderByPath(db, "/Test3/TestC", &superUser, "testserver", false)

		Assert.Nil(err)

		// Delete Folder that has contents inside "/Test3"

		ls, err = dbfs.ListFilesByPath(db, "/Test3", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(3, len(ls))

		parentFolder, err := dbfs.GetFileByPath(db, "/Test3", &superUser, false)
		Assert.Nil(err)

		err = dbfs.DeleteFolderById(db, parentFolder.FileId, &superUser)

		Assert.Error(dbfs.ErrFolderNotEmpty, err)

		// DeleteFolderByIdCascade

		err = dbfs.DeleteFolderByIdCascade(db, parentFolder.FileId, &superUser, "deleteServer")
		Assert.Nil(err)

		ls, err = dbfs.ListFilesByPath(db, "/Test3", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(0, len(ls))

	})

	t.Run("Checking permissions for users", func(t *testing.T) {

		Assert := assert.New(t)

		// Making folders

		parentFolder, err := dbfs.CreateFolderByPath(db, "/TestPerms", &superUser, "testserver", false)
		Assert.Nil(err)

		subFolder1, err := dbfs.CreateFolderByPath(db, "/TestPerms/perm1", &superUser, "testserver", false)
		Assert.Nil(err)

		subFolder2, err := dbfs.CreateFolderByPath(db, "/TestPerms/perm2", &superUser, "testserver", false)

		Assert.Nil(err)

		// Creating a user
		user1, err := dbfs.CreateNewUser(db, "permissionCheckUser", "user1Name", 1,
			"permissionCheckUser", "refreshToken", "accessToken", "idToken", "testServer")
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

		subFolder3, err := dbfs.CreateFolderByPath(db, "/TestPerms/perm1/perm3", user1, "testserver", false)
		Assert.Error(err, dbfs.ErrNoPermission)
		subFolder3, err = dbfs.CreateFolderByPath(db, "/TestPerms/perm1/perm3", &superUser, "testserver", false)
		Assert.Nil(err)

		// Check that user has permissions to /TestPerms/perm1/perm3
		hasPermission, err = user1.HasPermission(db, subFolder3, &dbfs.PermissionNeeded{Read: true})
		Assert.Equal(true, hasPermission)
		Assert.Nil(err)
		hasPermission, err = superUser.HasPermission(db, subFolder3, &dbfs.PermissionNeeded{Read: true})
		Assert.Equal(true, hasPermission)
		Assert.Nil(err)

		// See if the user has his own folder created.
		userFolder, err := dbfs.GetHomeFolder(db, user1)
		Assert.Nil(err)
		Assert.Equal(userFolder.FileId, user1.HomeFolderId)

	})

	// Creating fake files
	t.Run("Creating fake files", func(t *testing.T) {

		Assert := assert.New(t)

		// Making folders
		newFolder, err := dbfs.CreateFolderByPath(db, "/TestFakeFiles", &superUser, "testserver", false)
		Assert.Nil(err)

		// Creating fake files
		file, err := EXAMPLECreateFile(db, &superUser, "somefile.txt", newFolder.FileId)

		Assert.Nil(err)
		// Check if you can get da fragments
		fragments, err := file.GetFileFragments(db, &superUser)
		Assert.Nil(err)
		Assert.Equal(5, len(fragments))

		if debugPrint {
			fmt.Println(fragments)
		}

		file2, err := EXAMPLECreateFile(db, &superUser, "somefile2.txt", newFolder.FileId)
		Assert.Nil(err)
		// Check if you can get da fragments
		fragments, err = file2.GetFileFragments(db, &superUser)
		Assert.Nil(err)
		Assert.Equal(5, len(fragments))

		if debugPrint {
			fmt.Println(fragments)
		}

		file, err = EXAMPLECreateFile(db, &superUser, "somefile2.txt", newFolder.FileId)
		Assert.Error(dbfs.ErrFileFolderExists, err)

		err = dbfs.DeleteFileById(db, file2.FileId, &superUser, "deleteServer")
		Assert.Nil(err)

		fragments, err = file2.GetFileFragments(db, &superUser)
		Assert.Error(dbfs.ErrFileNotFound, err)
		Assert.Equal(0, len(fragments))

		if debugPrint {
			fmt.Println(fragments)
		}

		// Test GetFileMeta
		// To test for GetFileMeta, need to add some permissions to a file first

		// Create a user
		userForGetFileMeta, err := dbfs.CreateNewUser(db, "getFileMetaUser", "user1Name", dbfs.AccountTypeEndUser, "getFileMetaUser",
			"refreshToken", "accessToken", "idToken", "testServer")
		uselessUser, err := dbfs.CreateNewUser(db, "uselessUser", "user1Name", dbfs.AccountTypeEndUser, "uselessUser",
			"refreshToken", "accessToken", "idToken", "testServer")
		Assert.Nil(err)

		newFile, err := dbfs.GetFileByPath(db, "/TestFakeFiles/somefile.txt", &superUser, false)

		err = newFile.AddPermissionUsers(db, &dbfs.PermissionNeeded{Read: true, Write: true}, &superUser, *userForGetFileMeta)
		Assert.Nil(err)

		newFile, err = dbfs.GetFileByPath(db, "/TestFakeFiles/somefile.txt", &superUser, false)
		Assert.Nil(err)
		err = newFile.GetFileMeta(db, &superUser)
		Assert.Nil(err)
		Assert.Equal(newFile.ModifiedUser.UserId, superUser.UserId)
		Assert.Equal(newFile.VersionNo, 1)

		// See if the new user can get the file
		newFile2, err := dbfs.GetFileById(db, newFile.FileId, userForGetFileMeta)
		Assert.Nil(err)
		Assert.Equal(newFile.FileId, newFile2.FileId)
		_, err = dbfs.GetFileById(db, newFile.FileId, uselessUser)
		Assert.Error(dbfs.ErrFileNotFound, err)

		// Test GetOldVersion

		newFileV0, err := newFile.GetOldVersion(db, userForGetFileMeta, 0)
		Assert.Nil(err)
		Assert.Equal(newFileV0.VersionNo, 0)

		newFileV0, err = newFile.GetOldVersion(db, userForGetFileMeta, 20)
		Assert.Error(dbfs.ErrVersionNotFound, err)

		newFileV0, err = newFile.GetOldVersion(db, userForGetFileMeta, -10)
		Assert.Error(dbfs.ErrVersionNotFound, err)

		// Rename a file

		err = newFile.UpdateMetaData(db, dbfs.FileMetadataModification{FileName: "pogfile.txt",
			MIMEType: "text", VersioningMode: dbfs.VersioningOff}, &superUser)
		Assert.Nil(err)
		newFile, err = dbfs.GetFileByPath(db, "/TestFakeFiles/pogfile.txt", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(newFile.FileName, "pogfile.txt")
		Assert.Equal(newFile.MIMEType, "text")
		Assert.Equal(newFile.VersioningMode, dbfs.VersioningOff)
		Assert.Equal(newFile.VersionNo, 2)

		// Move. Attempting to move the file to the root folder.

		// Getting root folder
		rootFolder, err := dbfs.GetRootFolder(db)
		Assert.Nil(err)

		// Moving the file
		err = newFile.Move(db, rootFolder, &superUser)
		Assert.Nil(err)
		files, err := rootFolder.ListContents(db, &superUser)
		Assert.Nil(err)
		Assert.Equal(8, len(files))

		// Assuming updating pogfile.txt
		err = EXAMPLEUpdateFile(db, newFile, "", &superUser)
		Assert.Nil(err)
		newFile, err = dbfs.GetFileByPath(db, "/pogfile.txt", &superUser, false)
		Assert.Nil(err)
		Assert.Equal(newFile.VersionNo, 3)
		fmt.Println(newFile.DataIdVersion)
		Assert.Equal(newFile.DataIdVersion, 1)

		// Trying to RemovePermission user
		// Getting all permissions

		permissions, err := newFile.GetPermissions(db, &superUser)
		Assert.Nil(err)
		for i, permission := range permissions {
			fmt.Println(i)
			if *permission.UserId == userForGetFileMeta.UserId {
				err = newFile.RemovePermission(db, &permission, &superUser)
				Assert.Nil(err)
			}
		}

		Assert.Equal(newFile.VersionNo, 4)

		// Trying to get the file again
		_, err = dbfs.GetFileById(db, newFile.FileId, userForGetFileMeta)
		Assert.Error(dbfs.ErrFileNotFound, err)
		_, err = dbfs.GetFileById(db, newFile.FileId, &superUser)
		Assert.Nil(err)

		// Trying encryption
		keyphrase := "HelloPassword"

		oldKey, oldIv, err := newFile.GetDecryptionKey(db, &superUser, "")
		Assert.Nil(err)

		err = newFile.PasswordProtect(db, "", keyphrase, "hinty", &superUser)
		Assert.Nil(err)
		sameKey, sameIv, err := newFile.GetDecryptionKey(db, &superUser, keyphrase)
		Assert.Equal(oldKey, sameKey)
		Assert.Equal(oldIv, sameIv)

		Assert.Nil(err)

	})

	// Encryption stuff testing
	t.Run("Encryption", func(t *testing.T) {
		//dbfs.GetAES("HelloPassword")

		keyToEncrypt, ivToEncrypt, err := dbfs.GenerateKeyIV()
		assert.Nil(t, err)

		key, iv, err := dbfs.GenerateKeyIV()
		assert.Nil(t, err)

		encrypted, err := dbfs.EncryptWithKeyIV(keyToEncrypt, key, iv)
		assert.Nil(t, err)
		plaintext, err := dbfs.DecryptWithKeyIV(encrypted, key, iv)
		assert.Nil(t, err)
		assert.Equal(t, keyToEncrypt, plaintext)

		encrypted, err = dbfs.EncryptWithKeyIV(ivToEncrypt, key, iv)
		assert.Nil(t, err)
		plaintext, err = dbfs.DecryptWithKeyIV(encrypted, key, iv)
		assert.Nil(t, err)
		assert.Equal(t, ivToEncrypt, plaintext)

	})

}

// EXAMPLECreateFile is an example driver for creating a File.
func EXAMPLECreateFile(tx *gorm.DB, user *dbfs.User, filename string, parentFolderId string) (*dbfs.File, error) {

	// This is an example script to show how the process should work.

	// First, the system receives the whole file from the user
	// Then, the system creates the record in the system

	file := dbfs.File{
		FileId:             uuid.New().String(),
		FileName:           filename,
		MIMEType:           "",
		ParentFolderFileId: &parentFolderId, // root folder for now
		Size:               512,
		VersioningMode:     dbfs.VersioningOff,
		Checksum:           "CHECKSUM",
		TotalShards:        ExampleTotalShards,
		DataShards:         ExampleDataShards,
		ParityShards:       ExampleParityShards,
		KeyThreshold:       ExampleKeyThreshold,
		PasswordProtected:  false,
		HandledServer:      "ThisServer",
	}

	fileKey, fileIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return nil, err
	}

	passwordProtect := dbfs.PasswordProtect{
		FileId:         file.FileId,
		FileKey:        fileKey,
		FileIv:         fileIv,
		PasswordActive: false,
	}

	// This is the key and IV from the pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return nil, err
	}

	// Encrypt dataKey and dataIv with the fileKey and fileIv
	dataKey, err = dbfs.EncryptWithKeyIV(dataKey, fileKey, fileIv)
	if err != nil {
		return nil, err
	}
	dataIv, err = dbfs.EncryptWithKeyIV(dataIv, fileKey, fileIv)
	if err != nil {
		return nil, err
	}

	err = tx.Transaction(func(tx *gorm.DB) error {

		err := dbfs.CreateInitialFile(tx, &file, fileKey, fileIv, dataKey, dataIv, user)
		if err != nil {
			return err
		}

		err = tx.Create(&passwordProtect).Error
		if err != nil {
			return err
		}

		err = dbfs.CreatePermissions(tx, &file)
		if err != nil {
			// By right, there should be no error possible? If any error happens, it's likely a system error.
			// However, in the case there is an error, we will revert the transaction (thus deleting the file entry)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// If no error, file is created. Now the system need to process the file in the pipeline and send it to each server
	// Pipeline can get the amount of parity bits based on the Parity Count, and get the amount of shards based on the amount of servers

	// Then, the system splits the files accordingly

	// Then, the system send each shard to each server, and once it's sent successfully

	for i := 1; i <= file.TotalShards; i++ {
		fragId := i
		fragmentPath := uuid.New().String()
		serverId := "Server" + strconv.Itoa(i)

		err = dbfs.CreateFragment(tx, file.FileId, file.DataId, file.VersionNo, fragId, serverId, fragmentPath)
		if err != nil {
			// Not sure how to handle this multiple error situation that is possible.
			// Don't necessarily want to put it in a transaction because I'm worried it'll be too long?
			// or does that make no sense?
			err2 := file.Delete(tx, user, "deleteServer")
			if err2 != nil {
				return nil, err2
			}

			return nil, err
			// If fails, delete File record and return error.
		}
	}

	// Now we can update it to indicate that it has been saved successfully.

	err = dbfs.FinishFile(tx, &file, user, 412, "checksum")
	if err != nil {
		// If fails, delete File record and return error.
	}

	err = dbfs.CreateFileVersionFromFile(tx, &file, user)
	if err != nil {
		// If fails, delete File record and return error.
	}

	return &file, nil

}

func EXAMPLEUpdateFile(tx *gorm.DB, file *dbfs.File, password string, user *dbfs.User) error {

	// Key and IV from pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return err
	}

	err = file.UpdateFile(tx, 2020, 2020, "wowKey", "wowServer", dataKey, dataIv, password, user)
	if err != nil {
		return err
	}

	// As each fragment is uploaded, each fragment is added to the database.
	for i := 1; i <= file.TotalShards; i++ {
		err = file.UpdateFragment(tx, i, "wowFragmentPath", "wowChecksum", "wowServerId")
		if err != nil {
			return err
		}
	}

	// Once all fragments are uploaded, the file is marked as finished.
	return file.FinishUpdateFile(tx, "")

}
