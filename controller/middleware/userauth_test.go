package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	assert := assert.New(t)

	// Set up mock db and session store
	session := testutil.NewMockSession(t)
	db := testutil.NewMockDB(t)

	// Set up next handler
	nextCalledTimes := 0
	nextLastContext := context.Background()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalledTimes++
		nextLastContext = r.Context()
		w.WriteHeader(http.StatusOK)
	})

	// Set up middleware
	mw := middleware.NewUserAuthMW(session, db)(next)

	// Test with blank request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	assert.Equal(http.StatusUnauthorized, w.Code)

	// Test with invalid session
	req = httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: "invalid"})
	w = httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	assert.Equal(http.StatusUnauthorized, w.Code)

	// Create session for nonexistent user
	sessionId, err := session.Create(nil, "1234", time.Hour)
	assert.Nil(err)

	// Test session
	req = httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	assert.Equal(http.StatusUnauthorized, w.Code)

	// Create user
	user, err := dbfs.CreateNewUser(
		db, "test@email.co", "John",
		dbfs.AccountTypeEndUser,
		"some-id", "refresh token",
		"access token", "id token", "server",
	)

	// Create new session
	sessionId, err = session.Create(nil, user.UserId, time.Hour)
	assert.NoError(err)

	// Test session
	req = httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
	w = httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal(1, nextCalledTimes)
	ctxuser, err := ctxutil.GetUser(nextLastContext)
	assert.Equal(user.Email, ctxuser.Email)
}
