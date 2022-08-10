package inc

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	FragmentHealthCheckPath = "/api/v1/node/fragment/{fragmentPath}/health"
	FragmentPath            = "/api/v1/node/fragment/{fragmentPath}"
	FragmentOrphanedPath    = "/api/v1/node/orphaned"
	FragmentMissingPath     = "/api/v1/node/missing"
	ShutdownPath            = "/api/v1/node/shutdown"
	CurrentFilesHealthPath  = "/api/v1/node/health_current_files"
	AllFilesHealthPath      = "/api/v1/node/health_all_files"
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

func (i Inc) OrphanedShardsRoute(w http.ResponseWriter, r *http.Request) {

	// get job_id from header
	jobId := r.Header.Get("job_id")
	if jobId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing job_id")
		return
	}

	// convert job_id to int
	jobIdInt, err := strconv.Atoi(jobId)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid job_id")
		return
	}

	go func() {
		_, err := i.LocalOrphanedShardsCheck(jobIdInt, true)
		if err != nil {
			// close the JobProgressMissingShard
			i.Db.Model(&dbfs.JobProgressOrphanedShard{}).
				Where("job_id = ? and server_id = ?", jobId, i.ServerName).
				Updates(map[string]interface{}{"in_progress": false, "msg": err.Error()})

		}
	}()

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, true)

}

func (i Inc) MissingShardsRoute(w http.ResponseWriter, r *http.Request) {

	// get job_id from header
	jobId := r.Header.Get("job_id")
	if jobId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing job_id")
		return
	}

	// convert job_id to int
	jobIdInt, err := strconv.Atoi(jobId)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid job_id")
		return
	}

	go func() {
		_, err := i.LocalMissingShardsCheck(jobIdInt, true)
		if err != nil {
			// close the JobProgressMissingShard
			i.Db.Model(&dbfs.JobProgressMissingShard{}).
				Where("job_id = ? and server_id = ?", jobId, i.ServerName).
				Updates(map[string]interface{}{"in_progress": false, "msg": err.Error()})

		}
	}()

	util.HttpJson(w, http.StatusOK, true)
}

func (i Inc) CurrentFilesFragmentsHealthCheckRoute(w http.ResponseWriter, r *http.Request) {

	// get job_id from header
	jobId := r.Header.Get("job_id")
	if jobId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing job_id")
		return
	}

	// convert job_id to int
	jobIdInt, err := strconv.Atoi(jobId)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid job_id")
		return
	}

	// run the job
	go func() {
		err := i.LocalCurrentFilesFragmentsHealthCheck(jobIdInt)
		if err != nil {
			// close the JobProgressMissingShard
			i.Db.Model(&dbfs.JobprogressCffhc{}).
				Where("job_id = ? and server_id = ?", jobId, i.ServerName).
				Updates(map[string]interface{}{"in_progress": false, "msg": err.Error()})

		}
	}()

	util.HttpJson(w, http.StatusOK, true)
}

func (i Inc) AllFilesFragmentsHealthCheckRoute(w http.ResponseWriter, r *http.Request) {

	// get job_id from header
	jobId := r.Header.Get("job_id")
	if jobId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing job_id")
		return
	}

	// convert job_id to int
	jobIdInt, err := strconv.Atoi(jobId)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid job_id")
		return
	}

	// run the job
	go func() {
		err := i.LocalAllFilesFragmentsHealthCheck(jobIdInt)
		if err != nil {
			// close the JobProgressMissingShard
			i.Db.Model(&dbfs.JobprogressAffhc{}).
				Where("job_id = ? and server_id = ?", jobId, i.ServerName).
				Updates(map[string]interface{}{"in_progress": false, "msg": err.Error()})

		}
	}()

	util.HttpJson(w, http.StatusOK, true)
}
