package controller

import (
	"fmt"
	"net/http"

	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/service"
	"github.com/OhanaFS/ohana/util"
	"github.com/gorilla/mux"
)

type Authentication struct {
	service service.Auth
}

func RegisterAuth(r *mux.Router, service service.Auth) {
	s := &Authentication{service}
	r.HandleFunc("/auth/login", s.GetAuth).Methods("GET")
	r.HandleFunc("/auth/callback", s.HandCallback).Methods("GET")
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

	util.HttpJson(w, http.StatusOK, result)
}
