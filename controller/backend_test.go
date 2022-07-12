package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

func TestBackendController(t *testing.T) {
	assert := assert.New(t)

	//Set up mock Db and session store
	stichConfig := config.StitchConfig{
		ShardsLocation: "shards/",
	}
	configFile := &config.Config{Stitch: stichConfig}
	logger := config.NewLogger(configFile)
	db := testutil.NewMockDB(t)
	err := dbfs.InitDB(db)
	session := testutil.NewMockSession(t)
	assert.NoError(err)

	// Setting up controller
	bc := &controller.BackendController{
		Db:     db,
		Logger: logger,
		Path:   configFile.Stitch.ShardsLocation,
	}

	// Getting Superuser to use with testing
	user, err := dbfs.GetUser(db, "superuser")
	if err != nil {
		return
	}

	// Creating a user with no permisisons
	noPermUser, err := dbfs.CreateNewUser(db, "noPermUser", "noPermUser", dbfs.AccountTypeEndUser, "noPermUser",
		"noPermUser", "noPermUser", "noPermUser")

	// Create new session
	sessionId, err := session.Create(nil, "superuser", time.Hour)
	assert.NoError(err)

	// Get root folder
	req := httptest.NewRequest("GET", "/api/v1/file/00000000-0000-0000-0000-000000000000/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w := httptest.NewRecorder()
	bc.GetMetadataFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Convert body to json
	body := w.Body.String()
	file := &dbfs.File{}
	err = json.Unmarshal([]byte(body), file)
	assert.NoError(err)
	assert.Equal("00000000-0000-0000-0000-000000000000", file.FileId)
	assert.Equal("root", file.FileName)

	// Create folder at root directory

	req = httptest.NewRequest("POST", "/api/v1/folder", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req.Header.Set("folder_name", "new_folder_name")
	req.Header.Set("parent_folder_id", "00000000-0000-0000-0000-000000000000")
	w = httptest.NewRecorder()
	bc.CreateFolder(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Convert body to json
	body = w.Body.String()
	file = &dbfs.File{}
	err = json.Unmarshal([]byte(body), file)
	assert.NoError(err)
	assert.Equal("00000000-0000-0000-0000-000000000000", *file.ParentFolderFileId)
	assert.Equal("new_folder_name", file.FileName)

	// Upload file to root directory
	fileDir, err := os.Getwd()
	assert.NoError(err)
	filePath := path.Join(fileDir, "testdata/test.txt")

	testFile, err := os.Open(filePath)
	assert.NoError(err)
	defer testFile.Close()

	sendBody := &bytes.Buffer{}
	writer := multipart.NewWriter(sendBody)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(err)
	_, err = io.Copy(part, testFile)
	assert.NoError(err)
	err = writer.Close()
	assert.NoError(err)

	req = httptest.NewRequest("POST", "/api/v1/file", sendBody).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("file_name", "test.txt")
	req.Header.Set("folder_id", "00000000-0000-0000-0000-000000000000")

	bc.UploadFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// After uploading file, check that we can ls and find the file

	req = httptest.NewRequest("GET", "/api/v1/folder/00000000-0000-0000-0000-000000000000/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req = mux.SetURLVars(req, map[string]string{
		"folderID": "00000000-0000-0000-0000-000000000000",
	})
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// body

	body = w.Body.String()
	var files []dbfs.File
	err = json.Unmarshal([]byte(body), &files)

	assert.Equal(2, len(files))
	assert.Equal("test.txt", files[1].FileName)

	// Getting file ID

	newFolderID := files[0].FileId
	newFileID := files[1].FileId

	// Download File
	req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.DownloadFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	fmt.Println(w.Body.String())

	// Copy File

	// Check ls folder before copying
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// body
	body = w.Body.String()
	var files2 []dbfs.File
	err = json.Unmarshal([]byte(body), &files2)
	assert.Equal(0, len(files2))

	req = httptest.NewRequest("POST", "/api/v1/file/"+newFileID+"/copy", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	req.Header.Add("folder_id", newFolderID)
	bc.CopyFile(w, req)
	assert.Equal(http.StatusOK, w.Code)
	body = w.Body.String()

	// ls and check again
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// body
	body = w.Body.String()
	err = json.Unmarshal([]byte(body), &files2)
	assert.Equal(1, len(files2))

	// Modify File Metadata (rename, delta)
	// dbfs.UpdateFileMetadata

	newBody, err := json.Marshal(dbfs.FileMetadataModification{
		FileName:       "thisfileisgreat.txt",
		VersioningMode: dbfs.VersioningOnVersions,
	})

	req = httptest.NewRequest("PATCH", "/api/v1/file/"+newFileID+"/metadata", bytes.NewReader(newBody)).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.UpdateMetadataFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// check if the file is renamed through
	// GetMetadataFile

	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/metadata", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	bc.GetMetadataFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// body
	body = w.Body.String()
	var anotherFile dbfs.File
	err = json.Unmarshal([]byte(body), &anotherFile)
	assert.Equal("thisfileisgreat.txt", anotherFile.FileName)

	// Seeing if we can GetFolderIDFromPath
	req = httptest.NewRequest("GET", "/api/v1/folder/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req.Header.Add("folder_path", "/new_folder_name")
	w = httptest.NewRecorder()
	bc.GetFolderIDFromPath(w, req)
	assert.Equal(http.StatusOK, w.Code)
	fmt.Println(w.Body.String())

	var returnFile dbfs.File
	err = json.Unmarshal([]byte(w.Body.String()), &returnFile)
	assert.Nil(err)
	assert.Equal(newFolderID, returnFile.FileId)

	// Testing GetFolderIDFromPath with a user that has no permissions
	req = httptest.NewRequest("GET", "/api/v1/folder/", nil).WithContext(
		ctxutil.WithUser(context.Background(), noPermUser))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req.Header.Add("folder_path", "/new_folder_name")
	w = httptest.NewRecorder()
	bc.GetFolderIDFromPath(w, req)
	assert.Equal(http.StatusNotFound, w.Code)

	// Trying to get a folder that doesn't exist
	req = httptest.NewRequest("GET", "/api/v1/folder/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req.Header.Add("folder_path", "/adsfas")
	w = httptest.NewRecorder()
	bc.GetFolderIDFromPath(w, req)
	assert.Equal(http.StatusNotFound, w.Code)

	// GetFolderDetails

	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.GetFolderDetails(w, req)
	assert.Equal(http.StatusOK, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &returnFile)
	assert.Nil(err)
	assert.Equal(newFolderID, returnFile.FileId)

	// Rename the inner folder
	// dbfs.UpdateFolderMetadata
	// ls root first
	req = httptest.NewRequest("GET", "/api/v1/folder/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": "00000000-0000-0000-0000-000000000000",
	})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)
	body = w.Body.String()

	var lsFiles []dbfs.File
	err = json.Unmarshal([]byte(body), &lsFiles)
	assert.Equal(2, len(lsFiles))

	var innerFolderID string
	var innerFolderName string

	// ensure the inner folder is there
	for _, file := range lsFiles {
		if file.EntryType == dbfs.IsFolder {
			innerFolderID = file.FileId
			innerFolderName = file.FileName
		}
	}

	assert.NotEmpty(innerFolderID)

	// Renaming the folder

	req = httptest.NewRequest("PATCH", "/api/v1/folder/"+innerFolderID+"/metadata", bytes.NewReader(newBody)).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": innerFolderID,
	})
	req.Header.Add("new_name", innerFolderName+"NEW")
	w = httptest.NewRecorder()
	bc.UpdateFolderMetadata(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// ls root again
	req = httptest.NewRequest("GET", "/api/v1/folder/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": "00000000-0000-0000-0000-000000000000",
	})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)
	body = w.Body.String()

	err = json.Unmarshal([]byte(body), &lsFiles)
	assert.Equal(2, len(lsFiles))

	foundFile := false
	// ensure the inner folder is there
	for _, file := range lsFiles {
		if file.EntryType == dbfs.IsFolder && file.FileName == innerFolderName+"NEW" && file.FileId == innerFolderID {
			foundFile = true
		}
	}

	assert.True(foundFile)

	// Test Update File
	// dbfs.UpdateFile

	// Upload file to root directory
	updateFilePath := path.Join(fileDir, "testdata/test2.txt")

	testFile, err = os.Open(updateFilePath)
	assert.NoError(err)
	defer testFile.Close()

	sendBody = &bytes.Buffer{}
	writer = multipart.NewWriter(sendBody)
	part, err = writer.CreateFormFile("file", "test.txt")
	assert.NoError(err)
	_, err = io.Copy(part, testFile)
	assert.NoError(err)
	err = writer.Close()
	assert.NoError(err)

	req = httptest.NewRequest("POST", "/api/v1/file/"+newFileID+"/update", sendBody).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	req.Header.Add("Content-Type", writer.FormDataContentType())

	bc.UpdateFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// After uploading file, check that we can download the file

	req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.DownloadFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	fmt.Println(w.Body.String())

	// Finding all versions of the file
	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/versions", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.GetFileVersionHistory(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// List of File version
	var fileVersions []dbfs.FileVersion
	err = json.Unmarshal([]byte(w.Body.String()), &fileVersions)
	assert.NoError(err)
	assert.Equal(3, len(fileVersions))

	// Trying to get an old version metadata

	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/versions/"+"0/metadata", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID":    newFileID,
		"versionID": "0",
	})
	w = httptest.NewRecorder()
	bc.GetFileVersionMetadata(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// FileVersion check

	var fileVersion dbfs.FileVersion
	err = json.Unmarshal([]byte(w.Body.String()), &fileVersion)
	assert.NoError(err)
	assert.Equal(fileVersion.VersioningMode, dbfs.VersioningOff)
	fmt.Println(w.Body.String())

	// Trying to download the old version of the file

	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/versions/"+"1", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID":    newFileID,
		"versionID": "1",
	})
	w = httptest.NewRecorder()
	bc.DownloadFileVersion(w, req)
	assert.Equal(http.StatusOK, w.Code)

	fmt.Println(w.Body.String())

	// Deleting the original file version
	req = httptest.NewRequest("DELETE", "/api/v1/file/"+newFileID+"/versions/"+"0", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID":    newFileID,
		"versionID": "0",
	})
	w = httptest.NewRecorder()
	bc.DeleteFileVersion(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Finding all versions of the file
	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/versions", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.GetFileVersionHistory(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// List of File version
	err = json.Unmarshal([]byte(w.Body.String()), &fileVersions)
	assert.NoError(err)
	assert.Equal(2, len(fileVersions))

	// Creating multiple users to test stuff with

	user1, err := dbfs.CreateNewUser(bc.Db, "testuser1", "testuser1", dbfs.AccountTypeEndUser,
		"testuser1", "testuser1", "testuser1", "testuseer1")
	assert.NoError(err)

	//group2, err := dbfs.CreateNewGroup(bc.Db, "testGroup2", "testGroup2")
	//assert.NoError(err)

	// Testing if user 1 has permission to newFolderID
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user1))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusNotFound, w.Code)

	// adding user 1 newFolderID

	req = httptest.NewRequest("POST", "/api/v1/folder/"+newFolderID+"/permissions/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})

	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	req.Header.Add("can_read", "true")
	req.Header.Add("can_write", "true")
	req.Header.Add("can_execute", "true")
	req.Header.Add("can_share", "true")
	req.Header.Add("users", "[\""+user1.UserId+"\"]")

	bc.AddPermissionsFolder(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Trying again
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user1))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Create Group and add to it

	user2, err := dbfs.CreateNewUser(bc.Db, "testuser2", "testuser2", dbfs.AccountTypeEndUser,
		"testuser2", "testuser2", "testuser2", "testuseer2")
	assert.NoError(err)
	user3, err := dbfs.CreateNewUser(bc.Db, "testuser3", "testuser3", dbfs.AccountTypeEndUser,
		"testuser3", "testuser3", "testuser3", "testuseer3")
	assert.NoError(err)

	group1, err := dbfs.CreateNewGroup(bc.Db, "testGroup1", "testGroup1")
	assert.NoError(err)

	err = user2.AddToGroup(bc.Db, group1)
	assert.NoError(err)
	err = user3.AddToGroup(bc.Db, group1)
	assert.NoError(err)

	// Testing if user2 can access newFolderID
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user2))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusNotFound, w.Code)

	// Adding group 1 to newFolderID
	req = httptest.NewRequest("POST", "/api/v1/folder/"+newFolderID+"/permissions/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})

	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	req.Header.Add("can_read", "true")
	req.Header.Add("can_write", "true")
	req.Header.Add("can_execute", "true")
	req.Header.Add("can_share", "true")
	req.Header.Add("groups", "[\""+group1.GroupId+"\"]")
	w = httptest.NewRecorder()
	bc.AddPermissionsFolder(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Trying again
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user2))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusNotFound, w.Code)

	// Get Permissions Folder / File
	req = httptest.NewRequest("GET", "/api/v1/folder/"+newFolderID+"/permissions/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.GetFolderPermissions(w, req)
	assert.Equal(http.StatusOK, w.Code)
	fmt.Println(w.Body.String())

	var incomingPermissions []dbfs.Permission
	err = json.Unmarshal(w.Body.Bytes(), &incomingPermissions)
	assert.NoError(err)

	// TODO: The amount of items being returned right now is not right. Need to fix, but for now it works.

	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/permissions/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.GetPermissionsFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &incomingPermissions)
	assert.NoError(err)
	assert.Equal(len(incomingPermissions), 1) // Only Superuser should see it

	// Adding group 1 and user1 to newFileID
	req = httptest.NewRequest("POST", "/api/v1/file/"+newFileID+"/permissions/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	req.Header.Add("can_read", "true")
	req.Header.Add("can_write", "true")
	req.Header.Add("can_execute", "true")
	req.Header.Add("can_share", "true")
	req.Header.Add("users", "[\""+user1.UserId+"\"]")
	req.Header.Add("groups", "[\""+group1.GroupId+"\"]")
	w = httptest.NewRecorder()
	bc.AddPermissionsFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Getting Permissions File again

	req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/permissions/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.GetPermissionsFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &incomingPermissions)
	assert.NoError(err)
	// TODO: Again, length is wrong, likely when addiing permissions, it's doing some weird recursive thing.
	// Will look into it

	// TODO: Modify Permissions not working... ID isn't auto incrementing

	// Move Folder

	// First, we'll create a new folder to move to
	req = httptest.NewRequest("POST", "/api/v1/folder/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req.Header.Add("folder_name", "moveFolder")
	req.Header.Add("parent_folder_id", "00000000-0000-0000-0000-000000000000")
	w = httptest.NewRecorder()
	bc.CreateFolder(w, req)
	assert.Equal(http.StatusOK, w.Code)
	var moveFolder dbfs.File
	err = json.Unmarshal(w.Body.Bytes(), &moveFolder)

	// Move newFileID and newFolderID to moveFolder

	req = httptest.NewRequest("PUT", "/api/v1/file/"+newFileID+"/move/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	req.Header.Add("folder_id", moveFolder.FileId)
	w = httptest.NewRecorder()
	bc.MoveFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	req = httptest.NewRequest("PUT", "/api/v1/folder/"+newFolderID+"/move/", nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	req.Header.Add("folder_id", moveFolder.FileId)
	w = httptest.NewRecorder()
	bc.MoveFolder(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// ls and see if newFileID and newFolderID are in the new folder
	req = httptest.NewRequest("GET", "/api/v1/folder/"+moveFolder.FileId, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req = mux.SetURLVars(req, map[string]string{
		"folderID": moveFolder.FileId,
	})
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// body
	body = w.Body.String()
	err = json.Unmarshal([]byte(body), &files2)
	assert.NoError(err)
	assert.Equal(2, len(files2))

	// Delete stuff in the folder
	// Delete newFileID
	req = httptest.NewRequest("DELETE", "/api/v1/file/"+newFileID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"fileID": newFileID,
	})
	w = httptest.NewRecorder()
	bc.DeleteFile(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Delete newFolderID
	req = httptest.NewRequest("DELETE", "/api/v1/folder/"+newFolderID, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	req = mux.SetURLVars(req, map[string]string{
		"folderID": newFolderID,
	})
	w = httptest.NewRecorder()
	bc.DeleteFolder(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// ls and see if newFileID and newFolderID are in the new folder
	req = httptest.NewRequest("GET", "/api/v1/folder/"+moveFolder.FileId, nil).WithContext(
		ctxutil.WithUser(context.Background(), user))
	req = mux.SetURLVars(req, map[string]string{
		"folderID": moveFolder.FileId,
	})
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	bc.LsFolderID(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// body
	body = w.Body.String()
	err = json.Unmarshal([]byte(body), &files2)
	assert.NoError(err)
	assert.Equal(0, len(files2))

}
