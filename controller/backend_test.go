package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestBackendController(t *testing.T) {
	assert := assert.New(t)

	// Generate dummy certificates for Inc
	tmpDir, err := os.MkdirTemp("", "ohana-test-")
	assert.NoError(err)
	defer os.RemoveAll(tmpDir)
	certs, err := selfsigntestutils.GenCertsTest(tmpDir)
	assert.NoError(err)
	shardsLocation := path.Join(tmpDir, "shards")
	assert.NoError(os.MkdirAll(shardsLocation, 0755))

	//Set up mock Db and session store
	configFile := &config.Config{
		Stitch: config.StitchConfig{
			ShardsLocation: shardsLocation,
		},
		Inc: config.IncConfig{
			CaCert:     certs.CaCertPath,
			PublicCert: certs.PublicCertPath,
			PrivateKey: certs.PrivateKeyPath,
			ServerName: "localhost",
			HostName:   "localhost",
			Port:       "65432",
		},
	}
	logger := config.NewLogger(configFile)
	db := testutil.NewMockDB(t)
	session := testutil.NewMockSession(t)

	// Setting up controller
	bc := &controller.BackendController{
		Db:         db,
		Logger:     logger,
		Path:       configFile.Stitch.ShardsLocation,
		ServerName: "localhost",
		Inc:        inc.NewInc(configFile, db),
	}

	// Register inc services
	inc.RegisterIncServices(bc.Inc)
	time.Sleep(time.Second * 3)

	bc.InitialiseShardsFolder()

	// Getting Superuser to use with testing
	user, err := dbfs.GetUser(db, "superuser")
	if err != nil {
		return
	}

	// Creating a user with no permisisons
	noPermUser, err := dbfs.CreateNewUser(db, "noPermUser", "noPermUser", dbfs.AccountTypeEndUser, "noPermUser",
		"noPermUser", "noPermUser", "noPermUser", "server")

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

	var returnFile, moveFolder dbfs.File
	var files, files2, lsFiles []dbfs.File
	var newFolderID, newFileID, fileDir, innerFolderID, innerFolderName string
	var newBody []byte
	var testFile *os.File
	sendBody := &bytes.Buffer{}
	var writer *multipart.Writer
	var part io.Writer

	newPassword := "newpassword"

	t.Run("Create Folder at root", func(t *testing.T) {
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

	})

	t.Run("Upload file to root directory", func(t *testing.T) {

		// Upload file to root directory
		fileDir, err = os.Getwd()
		assert.NoError(err)
		filePath := path.Join(fileDir, "testdata/test.txt")

		testFile, err = os.Open(filePath)
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
		err = json.Unmarshal([]byte(body), &files)

		assert.Equal(4, len(files), body)

		// Getting file ID
		for _, file := range files {
			if file.FileName == "test.txt" {
				newFileID = file.FileId
			}
			if file.FileName == "new_folder_name" {
				newFolderID = file.FileId
			}
		}

		assert.NotEqual("", newFileID)

		// Download File
		req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": newFileID,
		})
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		assert.Equal(http.StatusOK, w.Code, w.Body.String())
	})

	t.Run("Copy File", func(t *testing.T) {

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

		// Try to download the file
		// Download File
		req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": files2[0].FileId,
		})
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		assert.Equal(http.StatusOK, w.Code, w.Body.String())
	})

	t.Run("Copy Folder", func(t *testing.T) {

		// Create copyFolder at root directory
		req = httptest.NewRequest("POST", "/api/v1/folder", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Set("folder_name", "copyFolder")
		req.Header.Set("parent_folder_id", "00000000-0000-0000-0000-000000000000")
		w = httptest.NewRecorder()
		bc.CreateFolder(w, req)
		assert.Equal(http.StatusOK, w.Code)

		// Convert body to json
		body = w.Body.String()
		copyFolder := &dbfs.File{}
		err = json.Unmarshal([]byte(body), copyFolder)
		assert.NoError(err)
		assert.Equal("00000000-0000-0000-0000-000000000000", *copyFolder.ParentFolderFileId)
		assert.Equal("copyFolder", copyFolder.FileName)

		// Copy folder into it

		req = httptest.NewRequest("POST", "/api/v1/folder/"+copyFolder.FileId+"/copy", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": newFolderID,
		})
		req.Header.Add("folder_id", copyFolder.FileId)
		bc.CopyFile(w, req)
		assert.Equal(http.StatusOK, w.Code)

		// ls and check again
		req = httptest.NewRequest("GET", "/api/v1/folder/"+copyFolder.FileId, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": copyFolder.FileId,
		})
		w = httptest.NewRecorder()
		bc.LsFolderID(w, req)
		assert.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal([]byte(w.Body.String()), &files2)
		assert.NoError(err)
		assert.Equal(1, len(files2))

		// ls inner
		req = httptest.NewRequest("GET", "/api/v1/folder/"+files2[0].FileId, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": files2[0].FileId,
		})
		w = httptest.NewRecorder()
		bc.LsFolderID(w, req)
		assert.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal([]byte(w.Body.String()), &files2)
		assert.NoError(err)
		assert.Equal(1, len(files2))

		// Delete folder

		req = httptest.NewRequest("DELETE", "/api/v1/folder/"+copyFolder.FileId, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": copyFolder.FileId,
		})

		w = httptest.NewRecorder()
		bc.DeleteFolder(w, req)

		assert.Equal(http.StatusOK, w.Code)

	})

	t.Run("Modify File Metadata (rename, delta)", func(t *testing.T) {

		newBody, err = json.Marshal(dbfs.FileMetadataModification{
			FileName:          "thisfileisgreat.txt",
			VersioningMode:    dbfs.VersioningOnVersions,
			PasswordProtected: true,
			NewPassword:       newPassword,
		})

		assert.NoError(err)

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
	})

	t.Run("GetFolderIDFromPath", func(t *testing.T) {

		// Seeing if we can GetFolderIDFromPath
		req = httptest.NewRequest("GET", "/api/v1/folder/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("folder_path", "/new_folder_name")
		w = httptest.NewRecorder()
		bc.GetFolderIDFromPath(w, req)
		assert.Equal(http.StatusOK, w.Code)
		fmt.Println(w.Body.String())

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
	})

	t.Run("GetFolderDetails", func(t *testing.T) {
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
	})

	t.Run("Rename inner folder", func(t *testing.T) {

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

		err = json.Unmarshal([]byte(body), &lsFiles)
		assert.Equal(4, len(lsFiles))

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
		assert.Equal(4, len(lsFiles))

		foundFile := false
		// ensure the inner folder is there
		for _, file := range lsFiles {
			if file.EntryType == dbfs.IsFolder && file.FileName == innerFolderName+"NEW" && file.FileId == innerFolderID {
				foundFile = true
			}
		}

		assert.True(foundFile)
	})

	t.Run("Update File", func(t *testing.T) {

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
		req.Header.Add("password", newPassword)

		bc.UpdateFile(w, req)
		assert.Equal(http.StatusOK, w.Code)

		// After uploading file, check that we can download the file

		req = httptest.NewRequest("GET", "/api/v1/file/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": newFileID,
		})
		req.Header.Add("password", newPassword)
		w = httptest.NewRecorder()
		bc.DownloadFileVersion(w, req)
		assert.Equal(http.StatusOK, w.Code, w.Body.String())
	})

	t.Run("Versioning File", func(t *testing.T) {

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
		//fmt.Println(w.Body.String())

		// Trying to download the old version of the file

		req = httptest.NewRequest("GET", "/api/v1/file/"+newFileID+"/versions/"+"1", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID":    newFileID,
			"versionID": "1",
		})
		w = httptest.NewRecorder()
		req.Header.Add("password", newPassword)
		bc.DownloadFileVersion(w, req)
		assert.Equal(http.StatusOK, w.Code, w.Body.String())

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

	})

	user1, err := dbfs.CreateNewUser(bc.Db, "testuser1", "testuser1", dbfs.AccountTypeEndUser,
		"testuser1", "testuser1", "testuser1", "testuser1", "server")
	assert.NoError(err)

	t.Run("Permissions Check", func(t *testing.T) {

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
	})

	// Create Group and add to it
	t.Run("Group Check", func(t *testing.T) {

		user2, err := dbfs.CreateNewUser(bc.Db, "testuser2", "testuser2", dbfs.AccountTypeEndUser,
			"testuser2", "testuser2", "testuser2", "testuseer2", "server")
		assert.NoError(err)
		user3, err := dbfs.CreateNewUser(bc.Db, "testuser3", "testuser3", dbfs.AccountTypeEndUser,
			"testuser3", "testuser3", "testuser3", "testuseer3", "server")
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
		//fmt.Println(w.Body.String())

		var incomingPermissions []dbfs.Permission
		err = json.Unmarshal(w.Body.Bytes(), &incomingPermissions)
		assert.NoError(err)
		assert.Equal(3, len(incomingPermissions)) // Should be superuser, user1, and group1

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
		assert.Equal(3, len(incomingPermissions))

		// Remove can_share from user1
		// Finding the permission ID is

		var permissionID string

		for _, permission := range incomingPermissions {
			if permission.UserId != nil {
				if *permission.UserId == user1.UserId {
					permissionID = permission.PermissionId
				}
			}
		}

		req = httptest.NewRequest("PUT", "/api/v1/file/"+newFileID+"/permissions/"+permissionID, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": newFileID,
		})
		req.Header.Add("permission_id", permissionID)
		req.Header.Add("can_read", "true")
		req.Header.Add("can_write", "true")
		req.Header.Add("can_execute", "false")
		req.Header.Add("can_share", "false")
		req.Header.Add("can_audit", "false")

		w = httptest.NewRecorder()
		bc.ModifyPermissionsFile(w, req)
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
		assert.Equal(3, len(incomingPermissions))

		for _, permission := range incomingPermissions {
			if permission.UserId != nil {
				if *permission.UserId == user1.UserId {
					assert.Equal(false, permission.CanShare)
				}
			}
			fmt.Println(permission)
		}

	})

	t.Run("Move", func(t *testing.T) {

		// First, we'll create a new folder to move to
		req = httptest.NewRequest("POST", "/api/v1/folder/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req.Header.Add("folder_name", "moveFolder")
		req.Header.Add("parent_folder_id", "00000000-0000-0000-0000-000000000000")
		w = httptest.NewRecorder()
		bc.CreateFolder(w, req)
		assert.Equal(http.StatusOK, w.Code)
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
	})

	t.Run("Delete", func(t *testing.T) {

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

	})

	t.Run("GetPath", func(t *testing.T) {

		// create path1, path2, path3, path4

		rootFolder, err := dbfs.GetRootFolder(db)
		assert.NoError(err)

		path1, err := rootFolder.CreateSubFolder(db, "path1", user, "")
		assert.NoError(err)

		path2, err := path1.CreateSubFolder(db, "path2", user, "")
		assert.NoError(err)

		path3, err := path2.CreateSubFolder(db, "path3", user, "")
		assert.NoError(err)

		path4, err := path3.CreateSubFolder(db, "path4", user, "")
		assert.NoError(err)

		// Create a user
		userForGetFileMeta, err := dbfs.CreateNewUser(db, "GettingPath", "GettingPath", dbfs.AccountTypeEndUser, "GettingPath",
			"GettingPath", "GettingPath", "GettingPath", "testServer")
		assert.NoError(err)

		err = path2.AddPermissionUsers(db, &dbfs.PermissionNeeded{Read: true, Write: true}, user, *userForGetFileMeta)
		assert.NoError(err)

		var files []dbfs.File

		fakeSessionId, err := session.Create(nil, "GettingPath", time.Hour)

		req = httptest.NewRequest("GET", "/api/v1/folder/"+path2.FileId+"/path", nil).WithContext(
			ctxutil.WithUser(context.Background(), userForGetFileMeta))

		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: fakeSessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": path4.FileId,
		})
		w = httptest.NewRecorder()
		bc.GetPath(w, req)
		assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		err = json.Unmarshal([]byte(body), &files)
		assert.Equal(3, len(files))

		// folder 2
		req = httptest.NewRequest("GET", "/api/v1/folder/"+path2.FileId+"/path", nil).WithContext(
			ctxutil.WithUser(context.Background(), userForGetFileMeta))

		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: fakeSessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": path2.FileId,
		})
		w = httptest.NewRecorder()
		bc.GetPath(w, req)
		assert.Equal(http.StatusOK, w.Code)
		body = w.Body.String()

		err = json.Unmarshal([]byte(body), &files)
		assert.Equal(1, len(files))

		// folder 1 error since no permission
		req = httptest.NewRequest("GET", "/api/v1/folder/"+path2.FileId+"/path", nil).WithContext(
			ctxutil.WithUser(context.Background(), userForGetFileMeta))

		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: fakeSessionId})
		req = mux.SetURLVars(req, map[string]string{
			"folderID": path1.FileId,
		})
		w = httptest.NewRecorder()
		bc.GetPath(w, req)
		assert.Equal(http.StatusNotFound, w.Code)

	})

}
