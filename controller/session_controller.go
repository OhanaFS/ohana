package controller

import (
	"net/http"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/service"
	"github.com/gorilla/mux"
)

type SessionService struct {
	session service.Session
}

// create session controller for session service
func RegisterSession(r *mux.Router, session service.Session) {
	s := &SessionService{session}
	// session route - I guess need to be changed
	r.HandleFunc("/api/v1/session/create", s.CreateSession).Methods("POST")
	r.HandleFunc("/api/v1/session/get", s.GetSession).Methods("GET")
	r.HandleFunc("/api/v1/session/invalidate", s.InvalidateSession).Methods("POST")
}

// create session from cookie
func (s *SessionService) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user id from request
	userId := r.FormValue("userId")
	// get ttl from request
	ttl := time.Duration(time.Hour * 24 * 7)
	// create session
	sessionId, err := s.session.Create(ctx, userId, ttl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set a cookie with the token
	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieSessionName,
		Value:    sessionId,
		Expires:  time.Now().Add(ttl),
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *SessionService) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get session cookie
	cookie, err := r.Cookie("ohana_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session id from cookie
	sessionId := cookie.Value

	// get session
	userId, err := s.session.Get(ctx, sessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userId = userId // to avoid warning
}

// invalidate session from cookie if more than one week
func (s *SessionService) InvalidateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get session cookie
	cookie, err := r.Cookie(config.CookieSessionName)
	if err != nil {
		return
	}

	// get session id from cookie
	sessionId := cookie.Value

	// invalidate session
	err = s.session.Invalidate(ctx, sessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
