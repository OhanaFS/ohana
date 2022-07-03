package controller

import (
	"net/http"

	"github.com/OhanaFS/ohana/service"
	"github.com/OhanaFS/ohana/util"
	"github.com/gorilla/mux"
)

// Health is the controller for the health service.
type Health struct {
	service service.Health
}

// RegisterHealth registers routes to the supplied router.
func RegisterHealth(r *mux.Router, service service.Health) {
	c := &Health{service}
	r.HandleFunc("/v1/_health", c.GetHealth).Methods("GET")
}

func (c *Health) GetHealth(w http.ResponseWriter, r *http.Request) {
	util.HttpJson(w, http.StatusOK, c.service.Status())
}
