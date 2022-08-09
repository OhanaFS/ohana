package inc

import (
	"github.com/OhanaFS/ohana/util"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
)

const (
	FragmentHealthCheckPath = "/api/v1/node/fragment/{fragmentPath}/health"
	FragmentPath            = "/api/v1/node/fragment/{fragmentPath}"
	FragmentOrphanedPath    = "/api/v1/node/orphaned"
	FragmentMissingPath     = "/api/v1/node/missing"
	ShutdownPath            = "/api/v1/node/shutdown"
)

func (i Inc) FragmentHealthCheckRoute(w http.ResponseWriter, r *http.Request) {

	// check if fragment exists on the server

	vars := mux.Vars(r)
	fragmentPath := vars["fragmentPath"]

	if fragmentPath == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing fragment path")
		return
	}

	// check if fragment exists

	check, err := i.LocalIndividualFragHealthCheck(fragmentPath)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, *check)

}

func (i Inc) DeleteFragmentRoute(w http.ResponseWriter, r *http.Request) {

	// check if fragment exists on the server

	vars := mux.Vars(r)
	fragmentPath := vars["fragmentPath"]

	// Delete

	filePath := path.Join(i.ShardsLocation, fragmentPath)
	err := os.Remove(filePath)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, true)

}

func (i Inc) OrphanedFragmentsRoute(w http.ResponseWriter, r *http.Request) {

	result, err := i.LocalOrphanedShardsCheck()
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, result)

}

func (i Inc) MissingFragmentsRoute(w http.ResponseWriter, r *http.Request) {

	result, err := i.LocalMissingShardsCheck()
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, result)
}

func (i Inc) ShutdownRoute(w http.ResponseWriter, r *http.Request) {
	util.HttpJson(w, http.StatusOK, true)
	i.Shutdown <- true
}
