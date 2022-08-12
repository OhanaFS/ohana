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
	FragmentHealthCheckPath = "/api/v1/node/shard/{shardPath}/health"
	FragmentPath            = "/api/v1/node/shard/{shardPath}"
	FragmentOrphanedPath    = "/api/v1/node/orphaned"
	FragmentMissingPath     = "/api/v1/node/missing"
	ShutdownPath            = "/api/v1/node/shutdown"
	CurrentFilesHealthPath  = "/api/v1/node/health_current_files"
	AllFilesHealthPath      = "/api/v1/node/health_all_files"
	PathReplaceShardString  = "{shardPath}"
	PathFindString          = "shardPath"
)

// ShardHealthCheckRoute calls FragmentHealthCheck on the local server
// to ensure that the fragment is healthy.
// Expects shardPath in the URL
// Returns a JSON report based on the marshalling of stitch.ShardVerificationResult
func (i Inc) ShardHealthCheckRoute(w http.ResponseWriter, r *http.Request) {

	// check if fragment exists on the server

	vars := mux.Vars(r)
	shardPath := vars[PathFindString]

	if shardPath == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing fragment path")
		return
	}

	// check if fragment exists

	check, err := i.LocalIndividualShardHealthCheck(shardPath)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			"failed to get fragment "+err.Error())
		return
	}

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, *check)

}

// DeleteShardRoute deletes a shard from the local server.
// Expects shardPath in the URL
// Returns a http.StatusOK if successful
// Returns a http.StatusBadRequest if the shardPath is missing
func (i Inc) DeleteShardRoute(w http.ResponseWriter, r *http.Request) {

	// check if fragment exists on the server

	vars := mux.Vars(r)
	shardPath := vars[PathFindString]

	// Delete

	filePath := path.Join(i.ShardsLocation, shardPath)
	err := os.Remove(filePath)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			"failed to delete shard"+err.Error())
		return
	}

	// marshal check to json and return
	util.HttpJson(w, http.StatusOK, true)

}

// OrphanedShardsRoute starts a job to checked for orphaned shards.
// Expects job_id in the header
// Returns a http.StatusOK once job starts
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

// MissingShardsRoute starts a job to checked for missing shards.
// Expects job_id in the header
// Returns a http.StatusOK once job starts
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

// CurrentFilesFragmentsHealthCheckRoute starts a job to check the health of
// current file shards
// Expects job_id in the header
// Returns a http.StatusOK once job starts
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
		err := i.LocalCurrentFilesShardsHealthCheck(jobIdInt)
		if err != nil {
			// close the JobProgressMissingShard
			i.Db.Model(&dbfs.JobProgressCFSHC{}).
				Where("job_id = ? and server_id = ?", jobId, i.ServerName).
				Updates(map[string]interface{}{"in_progress": false, "msg": err.Error()})

		}
	}()

	util.HttpJson(w, http.StatusOK, true)
}

// AllFilesFragmentsHealthCheckRoute starts a job to check the health of
// all file shards
// Expects job_id in the header
// Returns a http.StatusOK once job starts
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
		err := i.LocalAllFilesShardsHealthCheck(jobIdInt)
		if err != nil {
			// close the JobProgressMissingShard
			i.Db.Model(&dbfs.JobProgressAFSHC{}).
				Where("job_id = ? and server_id = ?", jobId, i.ServerName).
				Updates(map[string]interface{}{"in_progress": false, "msg": err.Error()})

		}
	}()

	util.HttpJson(w, http.StatusOK, true)
}