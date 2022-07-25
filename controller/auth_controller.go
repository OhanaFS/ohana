package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/service"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/gorilla/mux"
)

type Authentication struct {
	service service.Auth
}

func RegisterAuth(r *mux.Router, service service.Auth, mw *middleware.Middlewares) {
	s := &Authentication{service}
	r.HandleFunc("/auth/login", s.GetAuth).Methods("GET")
	r.HandleFunc("/auth/callback", s.HandCallback).Methods("GET")
	r.HandleFunc("/auth/logout", s.Logout).Methods("GET")

	r2 := r.NewRoute().Subrouter()
	r2.HandleFunc("/auth/whoami", s.Whoami).Methods("GET")
	r2.Use(mw.UserAuth)
}

func (s *Authentication) Whoami(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, err.Error())
	}
	util.HttpJson(w, http.StatusOK, user)
}

func (s *Authentication) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := r.Cookie(middleware.SessionCookieName)
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, "No session cookie")
		return
	}

	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cookie: %s", err))
		return
	}

	// remove the users session from the session map
	err = s.service.InvalidateSession(ctx, c.Value)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error invalidating user: %s", err))
		return
	}

	// Expiring session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    middleware.SessionCookieName,
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Authentication) GetAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawAccessToken := r.Header.Get("Authorization")

	ourl, err := s.service.SendRequest(ctx, rawAccessToken)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error sending request: %s", err))
		return
	}

	http.Redirect(w, r, ourl, http.StatusFound)
}

func (s *Authentication) HandCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	checkState := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	result, err := s.service.Callback(ctx, code, checkState)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting roles: %s", err))
		return
	}

	// Assign session to user with cookie
	http.SetCookie(w, &http.Cookie{
		Name:  middleware.SessionCookieName,
		Value: result.SessionID,
		Path:  "/",
	})

	http.Redirect(w, r, "/home", http.StatusFound)
}
