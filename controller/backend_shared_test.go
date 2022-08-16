package controller_test

import (
	"context"
	"encoding/json"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	dbfstestutils "github.com/OhanaFS/ohana/dbfs/test_utils"
	selfsigntestutils "github.com/OhanaFS/ohana/selfsign/test_utils"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

func TestBackendController_SharedLinks(t *testing.T) {

	Assert := assert.New(t)

	// Generate dummy certificates for Inc
	tmpDir, err := os.MkdirTemp("", "ohana-test-")
	Assert.NoError(err)
	defer os.RemoveAll(tmpDir)
	certs, err := selfsigntestutils.GenCertsTest(tmpDir)
	Assert.NoError(err)
	shardsLocation := path.Join(tmpDir, "shards")
	Assert.NoError(os.MkdirAll(shardsLocation, 0755))

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
			Port:       "5558",
		},
	}
	logger := config.NewLogger(configFile)
	db := testutil.NewMockDB(t)
	session := testutil.NewMockSession(t)
	sessionId, err := session.Create(nil, "superuser", time.Hour)
	Assert.NoError(err)

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
	Assert.NotNil(user)

	// Get root folder

	rootFolder, err := dbfs.GetRootFolder(db)

	// Create file

	file, err := dbfstestutils.EXAMPLECreateFile(db, user, dbfstestutils.ExampleFile{
		FileName:       "blah",
		ParentFolderId: rootFolder.FileId,
		Server:         "localhost",
		FragmentPath:   bc.Inc.ShardsLocation,
		FileData:       "abc123",
		Size:           50,
		ActualSize:     80,
	})
	Assert.NoError(err)
	Assert.NotNil(file)

	t.Run("CreateSharedLink", func(t *testing.T) {

		// Creating a new link (Providing own url) No error.

		req := httptest.NewRequest("POST", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": file.FileId,
			"link":   "hello123"})

		w := httptest.NewRecorder()
		bc.CreateFileSharedLink(w, req)

		Assert.Equal(http.StatusOK, w.Code)

		var sharedLink dbfs.SharedLink
		err = json.Unmarshal(w.Body.Bytes(), &sharedLink)
		Assert.NoError(err)
		Assert.Equal("hello123", sharedLink.ShortenedLink)
		Assert.Equal(file.FileId, sharedLink.FileId)

		// Creating a link that already exists.

		req = httptest.NewRequest("POST", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": file.FileId,
			"link":   "hello123"})
		w = httptest.NewRecorder()
		bc.CreateFileSharedLink(w, req)

		Assert.Equal(http.StatusConflict, w.Code)

		// Creating a link without providing a shortname. Should randomly generate one.
		req = httptest.NewRequest("POST", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": file.FileId})
		req.Header.Add("link", "")

		w = httptest.NewRecorder()
		bc.CreateFileSharedLink(w, req)

		Assert.Equal(http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &sharedLink)
		Assert.NoError(err)
		t.Log(sharedLink.ShortenedLink)
		Assert.Equal(file.FileId, sharedLink.FileId)

	})

	t.Run("GetSharedLink", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": file.FileId})

		w := httptest.NewRecorder()
		bc.GetFileSharedLinks(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		var sharedLinks []dbfs.SharedLink
		err = json.Unmarshal(w.Body.Bytes(), &sharedLinks)
		Assert.NoError(err)
		Assert.Equal(2, len(sharedLinks)) // Dependent on the number of links created above.
		Assert.Equal(file.FileId, sharedLinks[0].FileId)
		Assert.Equal(file.FileId, sharedLinks[1].FileId)
		Assert.NotEqual(sharedLinks[0].ShortenedLink, sharedLinks[1].ShortenedLink)

		// Making sure they actually link to something
		Assert.NotEqual(sharedLinks[0].ShortenedLink, "")
		Assert.NotEqual(sharedLinks[1].ShortenedLink, "")

		// Checking that a random file id doesn't work

		req = httptest.NewRequest("GET", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": "ajldsjflkalkdsfjl"})

		w = httptest.NewRecorder()
		bc.GetFileSharedLinks(w, req)

		Assert.Equal(http.StatusNotFound, w.Code)

	})

	t.Run("Attempt to download file", func(t *testing.T) {

		// Bad attempt to download a file that doesn't exist.

		req := httptest.NewRequest("GET", "/api/v1/shared/{shortenedLink}/metadata", nil)
		req = mux.SetURLVars(req, map[string]string{"shortenedLink": "ajsdlfj"})

		w := httptest.NewRecorder()
		bc.GetMetadataSharedLink(w, req)

		Assert.Equal(http.StatusNotFound, w.Code)

		// Good attempt to get metadata of the real thing

		req = httptest.NewRequest("GET", "/api/v1/shared/{shortenedLink}/metadata", nil)
		req = mux.SetURLVars(req, map[string]string{"shortenedLink": "hello123"})

		w = httptest.NewRecorder()
		bc.GetMetadataSharedLink(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		var incomingFile dbfs.File
		err = json.Unmarshal(w.Body.Bytes(), &incomingFile)
		Assert.NoError(err)
		Assert.Equal(file.FileId, incomingFile.FileId)

		// Bad attempt to download a fake file

		req = httptest.NewRequest("GET", "/api/v1/shared/{shortenedLink}", nil)
		req = mux.SetURLVars(req, map[string]string{"shortenedLink": "kajsdklfjaklsdjkl"})

		w = httptest.NewRecorder()
		bc.DownloadSharedLink(w, req)

		Assert.Equal(http.StatusNotFound, w.Code, w.Body.String())

		// Good attempt to a real thing

		req = httptest.NewRequest("GET", "/api/v1/shared/{shortenedLink}", nil)
		req = mux.SetURLVars(req, map[string]string{"shortenedLink": "hello123"})

		w = httptest.NewRecorder()
		bc.DownloadSharedLink(w, req)

		Assert.Equal(http.StatusOK, w.Code, w.Body.String())

	})

	t.Run("PatchSharedLink", func(t *testing.T) {

		// amending hello123 to hello456

		req := httptest.NewRequest("PATCH", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": file.FileId,
			"link":   "hello123"})
		req.Header.Add("new_link", "hello456")

		w := httptest.NewRecorder()
		bc.PatchFileSharedLink(w, req)

		Assert.Equal(http.StatusOK, w.Code)

		// Getting to see if it was amended

		req = httptest.NewRequest("GET", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": file.FileId})

		w = httptest.NewRecorder()
		bc.GetFileSharedLinks(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		var sharedLinks []dbfs.SharedLink
		err = json.Unmarshal(w.Body.Bytes(), &sharedLinks)
		Assert.NoError(err)
		Assert.Equal(2, len(sharedLinks)) // Dependent on the number of links created above.

		// too lazy to do counts, this works.
		Assert.Contains(w.Body.String(), "hello456")

	})

	t.Run("DeleteSharedLink", func(t *testing.T) {

		// deleting hello456

		req := httptest.NewRequest("DELETE", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": file.FileId,
			"link":   "hello456"})

		w := httptest.NewRecorder()
		bc.DeleteFileSharedLink(w, req)

		Assert.Equal(http.StatusOK, w.Code)

		// Getting to see if it was amended

		req = httptest.NewRequest("GET", "/api/v1/file/{fileID}/share", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": file.FileId})

		w = httptest.NewRecorder()
		bc.GetFileSharedLinks(w, req)

		Assert.Equal(http.StatusOK, w.Code)
		var sharedLinks []dbfs.SharedLink
		err = json.Unmarshal(w.Body.Bytes(), &sharedLinks)
		Assert.NoError(err)
		Assert.Equal(1, len(sharedLinks)) // Dependent on the number of links created above.
		Assert.NotEqual(sharedLinks[0].ShortenedLink, "hello456")

	})

	t.Run("Favorites test", func(t *testing.T) {

		// Trying to get favorites (should be 0)
		req := httptest.NewRequest("GET", "/api/v1/favorites", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetFavorites(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		var favorites []dbfs.File
		err = json.Unmarshal(w.Body.Bytes(), &favorites)
		Assert.NoError(err)
		Assert.Equal(0, len(favorites))

		// Adding a favorite
		req = httptest.NewRequest("PUT", "/api/v1/favorites/{fileID}", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": file.FileId})
		w = httptest.NewRecorder()
		bc.AddFavorite(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		// Trying to get favorites (should be 1)
		req = httptest.NewRequest("GET", "/api/v1/favorites", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetFavorites(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &favorites)
		Assert.NoError(err)
		Assert.Equal(1, len(favorites))

		// Removing said favorite

		req = httptest.NewRequest("DELETE", "/api/v1/favorites/{fileID}", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{"fileID": file.FileId})
		w = httptest.NewRecorder()
		bc.RemoveFavorite(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		// Trying to get favorites (should be 0)
		req = httptest.NewRequest("GET", "/api/v1/favorites", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetFavorites(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &favorites)
		Assert.NoError(err)
		Assert.Equal(0, len(favorites))

	})

	t.Run("Shared To User", func(t *testing.T) {

		Assert := assert.New(t)

		// create a fake user
		user1, err := dbfs.CreateNewUser(db, "permissionCheckUser", "user1Name", 1,
			"permissionCheckUser", "refreshToken", "accessToken", "idToken", "testServer")
		Assert.NoError(err)

		// Checking files shared to user should be 0
		req := httptest.NewRequest("GET", "/api/v1/file/sharedWith", nil).
			WithContext(ctxutil.WithUser(context.Background(), user1))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetSharedWithUser(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		var sharedFiles []dbfs.File
		err = json.Unmarshal(w.Body.Bytes(), &sharedFiles)
		Assert.NoError(err)
		Assert.Equal(0, len(sharedFiles))

		// adding user 1 file

		req = httptest.NewRequest("POST", "/api/v1/file/"+file.FileId+"/permissions/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})

		req = mux.SetURLVars(req, map[string]string{
			"fileID": file.FileId,
		})
		w = httptest.NewRecorder()
		req.Header.Add("can_read", "true")
		req.Header.Add("can_write", "true")
		req.Header.Add("can_execute", "true")
		req.Header.Add("can_share", "true")
		req.Header.Add("users", "[\""+user1.UserId+"\"]")

		bc.AddPermissionsFile(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		// Checking files shared to user should be 1
		req = httptest.NewRequest("GET", "/api/v1/file/sharedWith", nil).
			WithContext(ctxutil.WithUser(context.Background(), user1))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetSharedWithUser(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &sharedFiles)
		Assert.NoError(err)
		Assert.Equal(1, len(sharedFiles))

		// Deleting the file
		req = httptest.NewRequest("DELETE", "/api/v1/file/"+file.FileId, nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		req = mux.SetURLVars(req, map[string]string{
			"fileID": file.FileId,
		})
		w = httptest.NewRecorder()
		bc.DeleteFile(w, req)
		Assert.Equal(http.StatusOK, w.Code)

		// Checking files shared to user should be 0
		req = httptest.NewRequest("GET", "/api/v1/file/sharedWith", nil).
			WithContext(ctxutil.WithUser(context.Background(), user1))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w = httptest.NewRecorder()
		bc.GetSharedWithUser(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &sharedFiles)
		Assert.NoError(err)
		Assert.Equal(0, len(sharedFiles))

	})

	t.Run("Favorites should be 0", func(t *testing.T) {

		// Trying to get favorites (should be 0)
		req := httptest.NewRequest("GET", "/api/v1/favorites", nil).
			WithContext(ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		bc.GetFavorites(w, req)
		Assert.Equal(http.StatusOK, w.Code)
		var favorites []dbfs.File
		err = json.Unmarshal(w.Body.Bytes(), &favorites)
		Assert.NoError(err)
		Assert.Equal(0, len(favorites))

	})
}
