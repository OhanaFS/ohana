package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/service"
	"github.com/OhanaFS/ohana/util"
	"github.com/gorilla/mux"
)

type AuthController struct {
	service service.Auth
}

func RegisterAuth(r *mux.Router, service service.Auth) {
	s := &AuthController{service}
	r.HandleFunc("/auth/login", s.GetAuth).Methods("GET")
	r.HandleFunc("/auth/callback", s.HandCallback).Methods("GET")

	r.HandleFunc("/auth/logout", s.Logout).Methods("GET")
}

func (s *AuthController) GetAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawAccessToken := r.Header.Get("Authorization")

	ourl, err := s.service.SendRequest(ctx, rawAccessToken)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error sending request: %s", err))
		return
	}

	http.Redirect(w, r, ourl, http.StatusFound)
}

func (s *AuthController) HandCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	checkState := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	result, err := s.service.Callback(ctx, code, checkState)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting roles: %s", err))
		return
	}

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    middleware.SessionCookieName,
		Value:   result.SessionId,
		Expires: time.Now().Add(result.TTL),
		Path:    "/",
	})

	util.HttpJson(w, http.StatusOK, result)
}

func (s *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := r.Cookie(middleware.SessionCookieName)
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, "No session cookie")
	}
	sessionId := c.Value

	// remove the users session from the session map
	err = s.service.InvalidateSession(ctx, sessionId)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error invalidating user: %s", err))
		return
	}

	// Expiring session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    middleware.SessionCookieName,
		Value:   "",
		Expires: time.Now(),
	})
}
