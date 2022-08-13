package controller

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/gorilla/mux"
	"net/http"
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
		if filter < 0 || filter > dbfs.JobStatusCompleteNoErrors {
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

// RebuildFileVersion If a shard is broken from a file, the system will attempt to rebuild
// it by creating a new shard with the same fileId and shardId.
func (bc *BackendController) RebuildFileVersion(fileId string, versionId int, password string) error {

	// Get file version, get IVs, download whatever is good and upload it again

	var fv dbfs.FileVersion

	err := bc.Db.Model(&dbfs.FileVersion{}).
		Where("file_id = ? AND version_id = ?", fileId, versionId).First(&fv).Error

	if err != nil {
		return dbfs.ErrFileNotFound
	}

	if fv.EntryType != dbfs.IsFile {
		return dbfs.ErrNotFile
	}

	var fragments []dbfs.Fragment
	err = bc.Db.Model(&fragments).
		Where("file_version_data_id = ?", fv.DataId).
		Find(&fragments).Error
	if err != nil {
		return err
	}

	// Get FileKey, FileIv from PasswordProtect

	var passwordProtect dbfs.PasswordProtect
	err = bc.Db.Model(&dbfs.PasswordProtect{}).Where("file_id = ?", fv.FileId).First(&passwordProtect).Error
	if err != nil {
		return fmt.Errorf("could not find password protect for file %s", fv.FileId)
	}

	var fileKey, fileIv string

	// if password nil
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

	key, err := hex.DecodeString(fileKey)
	if err != nil {
		return err
	}
	iv, err := hex.DecodeString(fileIv)
	if err != nil {
		return err
	}

	return nil
}
