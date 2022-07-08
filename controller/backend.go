package controller

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/stitch"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
)

type BackendController struct {
	db      *gorm.DB
	encoder *stitch.Encoder
	logger  *zap.Logger
}

// NewBackend takes in config, dbfs, loggers, and middleware and registers the backend
// routes.
func NewBackend(
	router *mux.Router,
	logger *zap.Logger,
	db *gorm.DB,
) error {

	bc := &BackendController{
		db: db,
		encoder: stitch.NewEncoder(&stitch.EncoderOptions{
			DataShards:   2,
			ParityShards: 1,
			KeyThreshold: 2,
		}),
		logger: logger,
	}

	r := router.NewRoute().Subrouter()
	r.HandleFunc("/api/v1/file/{fileID}", bc.GetMetadataFile).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}", bc.UpdateMetadataFile).Methods("PUT")
	r.HandleFunc("/api/v1/file/{fileID}/move", bc.MoveFile).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/copy", bc.CopyFile).Methods("POST")
	r.HandleFunc("/api/v1/file", bc.DownloadFile).Methods("GET")
	r.HandleFunc("/api/v1/file", bc.DeleteFile).Methods("DELETE")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.GetPermissionsFile).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.AddPermissionsFile).Methods("POST")
	// Register routes

	return nil
}

// GetMetadataFile returns the metadata for the requested file based on ID
func (bc *BackendController) GetMetadataFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// json encode file
	util.HttpJson(w, http.StatusOK, file)
}

//UpdateMetadataFile updates the metadata for the requested file based on ID
func (bc *BackendController) UpdateMetadataFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get Request Body and FileMetadataModification
	var fmm dbfs.FileMetadataModification
	err = json.NewDecoder(r.Body).Decode(&fmm)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update file
	err = file.UpdateMetaData(bc.db, fmm, user)
	if errors.Is(err, dbfs.ErrNoPermission) {
		util.HttpError(w, http.StatusForbidden, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, file)

}

//MoveFile moves the file to the new location
func (bc *BackendController) MoveFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Getting params
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}
	folderID := r.Header.Get("folder_id")
	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// Getting file and dest folder from db

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	destFolder, err := dbfs.GetFileByPath(bc.db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Move file (Permission check will be done by dbfs)
	err = file.Move(bc.db, destFolder, user)
	if errors.Is(err, dbfs.ErrNoPermission) {
		util.HttpError(w, http.StatusForbidden, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, nil)

}

//CopyFile copies the file to the new location
func (bc *BackendController) CopyFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Getting params
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}
	folderID := r.Header.Get("folder_id")
	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// Getting file and dest folder from db
	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, "File not found")
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	destFolder, err := dbfs.GetFileById(bc.db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, "Destination folder not found")
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Copy file (Permission check will be done by dbfs)
	err = file.Copy(bc.db, destFolder, user)
	if errors.Is(err, dbfs.ErrNoPermission) {
		util.HttpError(w, http.StatusForbidden, "No write permisison on destination folder")
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Success
	util.HttpJson(w, http.StatusOK, nil)

}

//DownloadFile gives the user the file back with a application/octet-stream header
func (bc *BackendController) DownloadFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Getting params
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	password := r.Header.Get("password")

	// Getting file from db
	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shardsMeta, err := file.GetFileFragments(bc.db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: Now assuming that all fragments are on the same server, else you need to query multiple
	// servers for the file fragments

	// Opening input files
	var shards []io.ReadSeeker
	var shardFiles []*os.File
	for _, shardMeta := range shardsMeta {
		// TODO: Should pick up the inital part of the path from config file.
		shardFile, err := os.Open(shardMeta.FileFragmentPath)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		shards = append(shards, shardFile)
		shardFiles = append(shardFiles, shardFile)
	}

	hexKey, hexIv, err := file.GetDecryptionKey(bc.db, user, password)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Getting key and iv
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	iv, err := hex.DecodeString(hexIv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Decode file
	reader, err := bc.encoder.NewReadSeeker(shards, key, iv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, reader)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, shardFile := range shardFiles {
		err := shardFile.Close()
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

}

//DeleteFile deletes a file from the db
func (bc *BackendController) DeleteFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete file
	err = file.Delete(bc.db, user)
	if errors.Is(err, dbfs.ErrNoPermission) {
		util.HttpError(w, http.StatusForbidden, "No write permisison on destination folder")
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Success
	util.HttpJson(w, http.StatusOK, nil)
}

//GetPermissionsFile returns the permissions of a file
func (bc *BackendController) GetPermissionsFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get permissions
	permissions, err := file.GetPermissions(bc.db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Success
	util.HttpJson(w, http.StatusOK, permissions)
}

//ModifyPermissionsFile modifies the permissions of a file based on the id obtained from GetPermissionsFile
func (bc *BackendController) ModifyPermissionsFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	// get other headers
	permissionID := r.Header.Get("permission_id")
	canRead := r.Header.Get("can_read")
	canWrite := r.Header.Get("can_write")
	canExecute := r.Header.Get("can_execute")
	canShare := r.Header.Get("can_share")
	canAudit := r.Header.Get("can_audit")

	// if any empty return error
	if fileID == "" || permissionID == "" || canRead == "" || canWrite == "" || canExecute == "" || canShare == "" || canAudit == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing params")
		return
	}

	// convert permissionID to int
	permissionIDInt, err := strconv.Atoi(permissionID)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid permissionID")
		return
	}

	// convert to bool
	canReadBool, err := strconv.ParseBool(canRead)
	if err != nil {
		canReadBool = false
	}
	canWriteBool, err := strconv.ParseBool(canWrite)
	if err != nil {
		canWriteBool = false
	}
	canExecuteBool, err := strconv.ParseBool(canExecute)
	if err != nil {
		canExecuteBool = false
	}
	canShareBool, err := strconv.ParseBool(canShare)
	if err != nil {
		canShareBool = false
	}
	canAuditBool, err := strconv.ParseBool(canAudit)
	if err != nil {
		canAuditBool = false
	}

	// get file
	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get old permission struct entry
	oldPermission, err := file.GetPermissionById(bc.db, permissionIDInt, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// create new permission struct
	newPermission := &dbfs.Permission{
		FileId:       oldPermission.FileId,
		PermissionId: oldPermission.PermissionId,
		User:         oldPermission.User,
		UserId:       oldPermission.UserId,
		Group:        oldPermission.Group,
		GroupId:      oldPermission.GroupId,
		CanRead:      canReadBool,
		CanWrite:     canWriteBool,
		CanExecute:   canExecuteBool,
		CanShare:     canShareBool,
		VersionNo:    oldPermission.VersionNo,
		Audit:        canAuditBool,
		CreatedAt:    oldPermission.CreatedAt,
		UpdatedAt:    oldPermission.UpdatedAt,
		DeletedAt:    oldPermission.DeletedAt,
		Status:       oldPermission.Status,
	}

	// call update

	err = file.UpdatePermission(bc.db, oldPermission, newPermission, user)

	if errors.Is(err, dbfs.ErrNoPermission) {
		util.HttpError(w, http.StatusForbidden, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Success
	util.HttpJson(w, http.StatusOK, nil)

}

//AddPermissionsFile adds permissions to a file
func (bc *BackendController) AddPermissionsFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	// get other headers
	canRead := r.Header.Get("can_read")
	canWrite := r.Header.Get("can_write")
	canExecute := r.Header.Get("can_execute")
	canShare := r.Header.Get("can_share")
	canAudit := r.Header.Get("can_audit")
	// users := r.Header.Get("users")
	// groups := r.Header.Get("groups")

	// need to process 'em

	// convert to bool
	canReadBool, err := strconv.ParseBool(canRead)
	if err != nil {
		canReadBool = false
	}
	canWriteBool, err := strconv.ParseBool(canWrite)
	if err != nil {
		canWriteBool = false
	}
	canExecuteBool, err := strconv.ParseBool(canExecute)
	if err != nil {
		canExecuteBool = false
	}
	canShareBool, err := strconv.ParseBool(canShare)
	if err != nil {
		canShareBool = false
	}
	canAuditBool, err := strconv.ParseBool(canAudit)
	if err != nil {
		canAuditBool = false
	}

	permissionNeeded := dbfs.PermissionNeeded{
		Read:    canReadBool,
		Write:   canWriteBool,
		Execute: canExecuteBool,
		Share:   canShareBool,
		Audit:   canAuditBool,
	}

	//TODO: convert to array. Need to write a function to handle this

	var userObjects []dbfs.User
	var groupObjects []dbfs.Group

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// add permissions
	err = file.AddPermissionUsers(bc.db, &permissionNeeded, user, userObjects...)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = file.AddPermissionGroups(bc.db, &permissionNeeded, user, groupObjects...)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Success
	util.HttpJson(w, http.StatusOK, nil)
}

//GetFileVersionMetadata gets the metadata of a file version
func (bc *BackendController) GetFileVersionMetadata(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	versionID := vars["versionID"]

	// convert versionID into int
	versionIDInt, err := strconv.Atoi(versionID)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	if fileID == "" || versionID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get version
	version, err := file.GetOldVersion(bc.db, user, versionIDInt)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, version)

}

//GetFileVersionHistory returns all the historical versions of a file
func (bc *BackendController) GetFileVersionHistory(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get version
	versions, err := file.GetAllVersions(bc.db, user)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, versions)

}

//DownloadFileVersion downloads a file version
func (bc *BackendController) DownloadFileVersion(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	versionID := vars["versionID"]

	// get password
	password := r.Header.Get("password")

	// convert versionID into int
	versionIDInt, err := strconv.Atoi(versionID)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	if fileID == "" || versionID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	// get file
	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get version
	version, err := file.GetOldVersion(bc.db, user, versionIDInt)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	// get file
	shardsMeta, err := version.GetFragments(bc.db, user)
	if err != nil {
		return
	}

	// TODO: Now assuming that all fragments are on the same server, else you need to query multiple
	// servers for the file fragments

	// Opening input files
	var shards []io.ReadSeeker
	var shardFiles []*os.File
	for _, shardMeta := range shardsMeta {
		// TODO: Should pick up the inital part of the path from config file.
		shardFile, err := os.Open(shardMeta.FileFragmentPath)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		shards = append(shards, shardFile)
		shardFiles = append(shardFiles, shardFile)
	}

	hexKey, hexIv, err := file.GetDecryptionKey(bc.db, user, password)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Getting key and iv
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	iv, err := hex.DecodeString(hexIv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Decode file
	reader, err := bc.encoder.NewReadSeeker(shards, key, iv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, reader)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, shardFile := range shardFiles {
		err := shardFile.Close()
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// DeleteFileVersion deletes a file version
func (bc *BackendController) DeleteFileVersion(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := dbfs.GetUser(bc.db, "1")
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	versionID := vars["versionID"]

	// convert versionID into int
	versionIDInt, err := strconv.Atoi(versionID)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	if fileID == "" || versionID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	// get file
	file, err := dbfs.GetFileById(bc.db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// delete version
	err = file.DeleteFileVersion(bc.db, user, versionIDInt)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, nil)

}