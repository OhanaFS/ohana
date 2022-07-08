package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	fmt.Print(ourl)
	if err == nil {
		http.Redirect(w, r, ourl, http.StatusFound)
	}
}

func (s *Authentication) HandCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	checkState := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	roles, _ := s.service.Callback(ctx, code, checkState)
	if roles.Fetched == false {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(roles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.HttpJson(w, http.StatusOK, string(j))
}
