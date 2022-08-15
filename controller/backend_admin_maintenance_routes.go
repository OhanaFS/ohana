package controller

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/stitch"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// GetAllJobs returns all jobs and their progress (in blocks of 25.)
// Requires a startNum, startDate, endDaate, and filter
func (bc *BackendController) GetAllJobs(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	startNumString := r.Header.Get("start_num")
	startNum, err := strconv.Atoi(startNumString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid start_num")
		return
	}

	startDateString := r.Header.Get("start_date")
	var startDate time.Time
	if startDateString == "" {
		startDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		startDate, err = util.Rfc3339ToDateOnly(startDateString)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, "Invalid start_date")
			return
		}

	}

	endDateString := r.Header.Get("end_date")
	var endDate time.Time
	if endDateString == "" {
		endDate = time.Now()
	} else {
		endDate, err = util.Rfc3339ToDateOnly(endDateString)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, "Invalid end_date")
			return
		}
	}

	filterString := r.Header.Get("filter")
	var filter int
	if filterString == "" {
		filter = 0 // assume 0 is all jobs
	} else {
		filter, err = strconv.Atoi(filterString)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, "Invalid filter")
			return
		}
		if filter < 0 || filter > dbfs.JobHasErrors {
			util.HttpError(w, http.StatusBadRequest, "Invalid filter")
			return
		}
	}

	jobs, err := dbfs.GetAllJobs(bc.Db, startNum, startDate, endDate, filter)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.HttpJson(w, http.StatusOK, jobs)

}

// GetJob returns the job with the given jobId
func (bc *BackendController) GetJob(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}
	job, err := dbfs.GetJob(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.HttpJson(w, http.StatusOK, job)
}

// StartJob starts the job with the given jobId
// and calls Inc for initalization of all the jobs
// returns a job which can be used to query the progress
func (bc *BackendController) StartJob(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	// getting the different options from params

	parameters := dbfs.JobParameters{
		MissingShardsCheck:  strings.TrimSpace(r.Header.Get("missing_shards_check")) == "true",
		OrphanedShardsCheck: strings.TrimSpace(r.Header.Get("orphaned_shards_check")) == "true",
		QuickShardsCheck:    strings.TrimSpace(r.Header.Get("quick_shards_check")) == "true",
		AllFilesShardsCheck: strings.TrimSpace(r.Header.Get("full_shards_check")) == "true",
		PermissionCheck:     strings.TrimSpace(r.Header.Get("permission_check")) == "true",
		DeleteFragments:     strings.TrimSpace(r.Header.Get("delete_fragments")) == "true",
	}

	// Starting the job

	job, err := dbfs.InitializeJob(bc.Db, parameters)

	fmt.Println(job.JobId)

	// StartingJobs, all these should return asap once all servers receive the request

	// Calling inc
	allErrors, err := bc.Inc.StartJob(job)
	if err != nil {
		// marshal allErrors to json
		allErrorsString, _ := json.Marshal(allErrors)
		util.HttpError(w, http.StatusInternalServerError, string(allErrorsString))
		return
	}
	// Calling dbfs jobs
	err = dbfs.StartJob(bc.Db, job)

	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.HttpJson(w, http.StatusOK, job)
}

// DeleteJob deletes the job with the given jobId
func (bc *BackendController) DeleteJob(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}
	err = dbfs.DeleteJob(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.HttpJson(w, http.StatusOK, true)
}

// GetFullShardsResult returns the result of the full shards check for the given jobId
func (bc *BackendController) GetFullShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}
	result, err := dbfs.GetResultsAFSHC(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.HttpJson(w, http.StatusOK, result)
}

// FixFullShardsResult takes the request body,
// decodes the results, and fixes based on the user input
// and fixes the full shards check for the given jobId
func (bc *BackendController) FixFullShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}

	// Check that jobId exists
	originalResults, err := dbfs.GetResultsAFSHC(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// make it a map so I can easily check if a file is in the map
	originalResultsMap := make(map[string]bool)
	for _, file := range originalResults {
		originalResultsMap[file.DataId] = true
	}

	// Decoding the request body
	var results []dbfs.ShardActions
	err = json.NewDecoder(r.Body).Decode(&results)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for _, result := range results {
		_, ok := originalResultsMap[result.DataId]
		if !ok {
			continue
		}
		// TODO: Should be a thread pool
		if result.Fix {
			err = bc.RebuildShard(result.DataId, result.Password)
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Mark as good
			err := bc.Db.Model(&dbfs.ResultsAFSHC{}).Where("data_id = ? and job_id = ?",
				result.DataId, jobId).
				Updates(map[string]interface{}{
					"error_type": dbfs.CronErrorTypeSolved,
					"error":      "Reconstructed",
				}).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else if result.Delete {
			// get/delete file
			var files []dbfs.File
			bc.Db.Model(&dbfs.File{}).Where("data_id = ?", result.DataId).Find(&files)
			for _, file := range files {
				err := file.Delete(bc.Db, user, "")
				if err != nil {
					util.HttpError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
			// get/delete file version
			err = bc.Db.Model(&dbfs.FileVersion{}).Where("data_id = ?", result.DataId).
				Update("status", dbfs.FileStatusToBeDeleted).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}

			// Mark as good
			err := bc.Db.Model(&dbfs.ResultsAFSHC{}).Where("data_id = ? and job_id = ?",
				result.DataId, jobId).
				Updates(map[string]interface{}{
					"error_type": dbfs.CronErrorTypeSolved,
					"error":      "Deleted",
				}).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

	}

	util.HttpJson(w, http.StatusOK, true)
}

// GetQuickShardsResult returns the result of the quick shards check for the given jobId
func (bc *BackendController) GetQuickShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}
	result, err := dbfs.GetResultsCFSHC(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.HttpJson(w, http.StatusOK, result)
}

// FixQuickShardsResult takes the request body,
// decodes the results, and fixes based on the user input
// and fixes the quick shards check for the given jobId
func (bc *BackendController) FixQuickShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}

	// Check that jobId exists
	originalResults, err := dbfs.GetResultsCFSHC(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// make it a map so I can easily check if a file is in the map
	originalResultsMap := make(map[string]bool)
	for _, file := range originalResults {
		originalResultsMap[file.DataId] = true
	}

	// Decoding the request body
	var results []dbfs.ShardActions
	err = json.NewDecoder(r.Body).Decode(&results)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for _, result := range results {
		_, ok := originalResultsMap[result.DataId]
		if !ok {
			continue
		}
		// TODO: Should be a thread pool
		if result.Fix {
			err = bc.RebuildShard(result.DataId, result.Password)
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Mark as good
			err := bc.Db.Model(&dbfs.ResultsCFSHC{}).Where("data_id = ? and job_id = ?",
				result.DataId, jobId).
				Updates(map[string]interface{}{
					"error_type": dbfs.CronErrorTypeSolved,
					"error":      "Reconstructed",
				}).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else if result.Delete {
			// get/delete file
			var files []dbfs.File
			bc.Db.Model(&dbfs.File{}).Where("data_id = ?", result.DataId).Find(&files)
			for _, file := range files {
				err := file.Delete(bc.Db, user, "")
				if err != nil {
					util.HttpError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
			// get/delete file version
			err = bc.Db.Model(&dbfs.FileVersion{}).Where("data_id = ?", result.DataId).
				Update("status", dbfs.FileStatusToBeDeleted).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}

			// Mark as good
			err := bc.Db.Model(&dbfs.ResultsCFSHC{}).Where("data_id = ? and job_id = ?",
				result.DataId, jobId).
				Updates(map[string]interface{}{
					"error_type": dbfs.CronErrorTypeSolved,
					"error":      "Deleted",
				}).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

	}

	util.HttpJson(w, http.StatusOK, true)
}

// GetMissingShardsResult returns the result of the missing shards check for the given jobId
func (bc *BackendController) GetMissingShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}
	result, err := dbfs.GetResultsMissingShard(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.HttpJson(w, http.StatusOK, result)
}

// FixMissingShardsResult takes the request body,
// decodes the results, and fixes based on the user input
// and fixes the missing shards check for the given jobId
func (bc *BackendController) FixMissingShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}

	// Check that jobId exists
	originalResults, err := dbfs.GetResultsMissingShard(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// make it a map so I can easily check if a file is in the map
	originalResultsMap := make(map[string]bool)
	for _, file := range originalResults {
		originalResultsMap[file.DataId] = true
	}

	// Decoding the request body
	var results []dbfs.ShardActions
	err = json.NewDecoder(r.Body).Decode(&results)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for _, result := range results {
		_, ok := originalResultsMap[result.DataId]
		if !ok {
			continue
		}
		// TODO: Should be a thread pool
		if result.Fix {
			err = bc.RebuildShard(result.DataId, result.Password)
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Mark as good
			err := bc.Db.Model(&dbfs.ResultsMissingShard{}).Where("data_id = ? and job_id = ?",
				result.DataId, jobId).
				Updates(map[string]interface{}{
					"error_type": dbfs.CronErrorTypeSolved,
					"error":      "Reconstructed",
				}).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else if result.Delete {
			// get/delete file
			var files []dbfs.File
			bc.Db.Model(&dbfs.File{}).Where("data_id = ?", result.DataId).Find(&files)
			for _, file := range files {
				err := file.Delete(bc.Db, user, "")
				if err != nil {
					util.HttpError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
			// get/delete file version
			err = bc.Db.Model(&dbfs.FileVersion{}).Where("data_id = ?", result.DataId).
				Update("status", dbfs.FileStatusToBeDeleted).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}

			// Mark as good
			err := bc.Db.Model(&dbfs.ResultsMissingShard{}).Where("data_id = ? and job_id = ?",
				result.DataId, jobId).
				Updates(map[string]interface{}{
					"error_type": dbfs.CronErrorTypeSolved,
					"error":      "Deleted",
				}).Error
			if err != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

	}

	util.HttpJson(w, http.StatusOK, true)
}

// GetOrphanedShardsResult returns the result of the full shards check for the given jobId
func (bc *BackendController) GetOrphanedShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}
	result, err := dbfs.GetResultsOrphanedShard(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.HttpJson(w, http.StatusOK, result)
}

// FixOrphanedShardsResult takes the request body,
// decodes the results, and fixes based on the user input
// (delete or leave alone)
func (bc *BackendController) FixOrphanedShardsResult(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	vars := mux.Vars(r)
	jobIdString := vars["id"]
	jobId, err := strconv.Atoi(jobIdString)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid jobId")
		return
	}

	// Check that jobId exists
	originalResults, err := dbfs.GetResultsOrphanedShard(bc.Db, jobId)
	if err != nil {
		if errors.Is(err, dbfs.ErrorCronJobDoesNotExist) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Decoding the request body
	var results []dbfs.OrphanedShardActions
	err = json.NewDecoder(r.Body).Decode(&results)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check that all results are valid
	type ServerPathKey struct {
		ServerId string
		Path     string
	}

	serverPaths := make(map[ServerPathKey]bool)

	for _, result := range originalResults {
		serverPaths[ServerPathKey{result.ServerId, result.FileName}] = true
	}

	for _, result := range results {

		if !serverPaths[ServerPathKey{result.ServerId, result.Path}] {
			continue
		}

		if result.Delete {

			if result.ServerId == bc.ServerName {
				delPath := path.Join(bc.Inc.ShardsLocation, result.Path)
				err = os.Remove(delPath)
				if err != nil {
					bc.Db.Model(&dbfs.ResultsOrphanedShard{}).
						Where("job_id = ? AND server_id = ? AND file_name = ?",
							jobId, result.ServerId, result.Path).Updates(map[string]interface{}{
						"error_type": dbfs.CronErrorTypeInternalError,
						"error":      err.Error(),
					})
					continue
				}
			} else {
				// Build a POST
				url := fmt.Sprintf("https://%s%s", result.ServerId,
					strings.Replace(inc.FragmentPath, "{shardPath}", result.Path, 1))
				req, err := http.NewRequest("DELETE", url, nil)
				if err != nil {
					bc.Db.Model(&dbfs.ResultsOrphanedShard{}).
						Where("job_id = ? AND server_id = ? AND file_name = ?",
							jobId, result.ServerId, result.Path).Updates(map[string]interface{}{
						"error_type": dbfs.CronErrorTypeInternalError,
						"error":      err.Error(),
					})
					continue
				}
				resp, err := bc.Inc.HttpClient.Do(req)
				if err != nil {
					bc.Db.Model(&dbfs.ResultsOrphanedShard{}).
						Where("job_id = ? AND server_id = ? AND file_name = ?",
							jobId, result.ServerId, result.Path).Updates(map[string]interface{}{
						"error_type": dbfs.CronErrorTypeInternalError,
						"error":      err.Error(),
					})
					continue
				}
				if resp.StatusCode != http.StatusOK {
					bc.Db.Model(&dbfs.ResultsOrphanedShard{}).
						Where("job_id = ? AND server_id = ? AND file_name = ?",
							jobId, result.ServerId, result.Path).Updates(map[string]interface{}{
						"error_type": dbfs.CronErrorTypeInternalError,
						"error":      resp.Body,
					})
					continue
				}
				defer resp.Body.Close()

			}

			bc.Db.Model(&dbfs.ResultsOrphanedShard{}).
				Where("job_id = ? AND server_id = ? AND file_name = ?",
					jobId, result.ServerId, result.Path).Updates(map[string]interface{}{
				"error_type": dbfs.CronErrorTypeSolved,
				"error":      "Deleted",
			})

		} else {
			bc.Db.Model(&dbfs.ResultsOrphanedShard{}).
				Where("job_id = ? AND server_id = ? AND file_name = ?",
					jobId, result.ServerId, result.Path).Updates(map[string]interface{}{
				"error_type": dbfs.CronErrorTypeSolved,
				"error":      "Deferred/Ignored",
			})
		}

	}

	util.HttpJson(w, http.StatusOK, true)
}

// RebuildShard If a shard is broken from a file, the system will attempt to rebuild
// it by creating a new shard with the same fileId and shardId.
func (bc *BackendController) RebuildShard(dataId, password string) error {

	// Get file version, get IVs, download whatever is good and upload it again

	var fv dbfs.FileVersion

	err := bc.Db.Model(&dbfs.FileVersion{}).
		Where("data_id = ?", dataId).First(&fv).Error

	if err != nil {
		return dbfs.ErrFileNotFound
	}

	if fv.EntryType != dbfs.IsFile {
		return dbfs.ErrNotFile
	}

	var shardsMeta []dbfs.Fragment
	err = bc.Db.Model(&shardsMeta).
		Where("file_version_data_id = ?", fv.DataId).
		Find(&shardsMeta).Error
	if err != nil {
		return err
	}

	// Get FileKey, FileIv from PasswordProtect
	// We'll use the exact same key, iv for the rebuild

	// We don't call dbfs here as we don't want to create a new file version
	// and dbfs requires a user to be passed in

	var passwordProtect dbfs.PasswordProtect
	err = bc.Db.Model(&dbfs.PasswordProtect{}).Where("file_id = ?", fv.FileId).First(&passwordProtect).Error
	if err != nil {
		return fmt.Errorf("could not find password protect for file %s", fv.FileId)
	}

	var fileKey, fileIv string

	// Password validation and extraction of key and iv

	if password == "" && !passwordProtect.PasswordActive {

		// decrypt file key with PasswordProtect
		fileKey, err = dbfs.DecryptWithKeyIV(fv.EncryptionKey, passwordProtect.FileKey, passwordProtect.FileIv)
		if err != nil {
			return fmt.Errorf("could not find decrypt password file %s", fv.FileId)
		}
		fileIv, err = dbfs.DecryptWithKeyIV(fv.EncryptionIv, passwordProtect.FileKey, passwordProtect.FileIv)
		if err != nil {
			return fmt.Errorf("could not find decrypt password file %s", fv.FileId)
		}
	} else if password == "" && passwordProtect.PasswordActive {
		return dbfs.ErrPasswordRequired
	} else {
		decryptedFileKey, decryptedFileIv, err := passwordProtect.DecryptWithPassword(password)
		if err != nil {
			return fmt.Errorf("invalid password for %s", fv.FileId)
		}
		fileKey, err = dbfs.DecryptWithKeyIV(fv.EncryptionKey, decryptedFileKey, decryptedFileIv)
		if err != nil {
			return fmt.Errorf("could not find decrypt password file %s", fv.FileId)
		}
		fileIv, err = dbfs.DecryptWithKeyIV(fv.EncryptionIv, decryptedFileKey, decryptedFileIv)
		if err != nil {
			return fmt.Errorf("could not find decrypt password file %s", fv.FileId)
		}
	}

	// We'll use the exact same key and iv for the rebuild
	key, err := hex.DecodeString(fileKey)
	if err != nil {
		return err
	}
	iv, err := hex.DecodeString(fileIv)
	if err != nil {
		return err
	}

	// Setting up the reader

	var shards []io.ReadSeeker

	for _, shardMeta := range shardsMeta {
		shardReader, err := bc.Inc.NewShardReader(
			context.Background(), shardMeta.ServerName, shardMeta.FileFragmentPath)
		if err == nil {
			shards = append(shards, shardReader)
		}
	}

	decoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   uint8(fv.DataShards),
		ParityShards: uint8(fv.ParityShards),
		KeyThreshold: uint8(fv.KeyThreshold)},
	)

	reader, err := decoder.NewReadSeeker(shards, key, iv)

	// read to the disk first, then upload (because weird bug)

	file, err := ioutil.TempFile("", "temp-file-reconstruct")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	file.Seek(0, 0)

	// Setting up the writer, write with .shardNEW extension

	stitchParams, err := dbfs.GetStitchParams(bc.Db, bc.Logger)
	dataShards, parityShards, keyThreshold := stitchParams.DataShards, stitchParams.ParityShards, stitchParams.KeyThreshold
	totalShards := dataShards + parityShards

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   uint8(dataShards),
		ParityShards: uint8(parityShards),
		KeyThreshold: uint8(keyThreshold),
	})

	beforeShardName := make([]string, totalShards)
	finalShardName := make([]string, totalShards)

	for i := 0; i < totalShards; i++ {
		beforeShardName[i] = fv.DataId + ".shardnew" + strconv.Itoa(i+1)
		finalShardName[i] = fv.DataId + ".shard" + strconv.Itoa(i+1)
	}

	// Open the output writers

	servers, err := bc.Inc.AssignShardServer(context.Background(), totalShards)

	// Prepare the writers
	shardWriters := make([]io.Writer, totalShards)

	ctx := context.Background()

	// Open the output writers
	for i := 0; i < totalShards; i++ {
		shardWriter, err := bc.Inc.NewShardWriter(
			ctx, servers[i].Name, beforeShardName[i])
		if err != nil {
			return err
		}
		shardWriters[i] = shardWriter
	}

	// Let the copying begin
	_, err = encoder.Encode(file, shardWriters, key, iv)
	if err != nil {
		return err
	}

	file.Close()

	os.Remove(file.Name())

	// Close the output writers
	for _, writer := range shardWriters {
		if err := writer.(io.WriteCloser).Close(); err != nil {
			return err
		}
	}

	// Updating fragments on the server to be named correctly (.shardNEW -> .shard)
	for i := 0; i < totalShards; i++ {
		err = bc.Inc.ReplaceShard(servers[i].Name, beforeShardName[i], finalShardName[i])
		if err != nil {
			return err
		}
	}

	// Update all file version/files that have that dataID to have
	// the right count for dataShards, parityShards, and keyThreshold
	err = bc.Db.Transaction(func(tx *gorm.DB) error {

		// ERRORING
		err := tx.Model(&dbfs.FileVersion{}).Where("data_id = ?", fv.DataId).
			Updates(map[string]interface{}{
				"data_shards":   dataShards,
				"parity_shards": parityShards,
				"key_threshold": keyThreshold,
				"total_shards":  totalShards,
				"status":        dbfs.FileStatusGood},
			).Error
		if err != nil {
			return fmt.Errorf("could not update file version %s", fv.DataId)
		}

		// Update file as well if exists
		err = tx.Model(&dbfs.File{}).
			Where("data_id = ?", fv.DataId).Updates(map[string]interface{}{
			"data_shards":   dataShards,
			"parity_shards": parityShards,
			"key_threshold": keyThreshold,
			"total_shards":  totalShards,
			"status":        dbfs.FileStatusGood},
		).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("could not update file %s", fv.DataId)
			}
		}

		// Update shards

		fvFId := shardsMeta[0].FileVersionFileId
		fvDId := shardsMeta[0].FileVersionDataId
		fvVn := shardsMeta[0].FileVersionVersionNo
		fvDidV := shardsMeta[0].FileVersionDataIdVersion

		// delete old shards.

		err = tx.Where("file_version_data_id = ?", fvDId).Delete(&dbfs.Fragment{}).Error

		// create new ones
		for i := 1; i <= totalShards; i++ {
			shard := dbfs.Fragment{
				FileVersionFileId:        fvFId,
				FileVersionDataId:        fvDId,
				FileVersionVersionNo:     fvVn,
				FileVersionDataIdVersion: fvDidV,
				FileFragmentPath:         finalShardName[i-1],
				FragId:                   i,
				ServerName:               servers[i-1].Name,
				LastChecked:              time.Now(),
				TotalShards:              totalShards,
				Status:                   dbfs.FragmentStatusGood,
			}
			if tx.Create(&shard).Error != nil {
				return fmt.Errorf("could not create fragment %s", finalShardName[i])
			}
		}

		return nil

	})

	return nil
}
