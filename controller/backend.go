package controller

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/stitch"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
)

type BackendController struct {
	Db         *gorm.DB
	Logger     *zap.Logger
	Path       string
	ServerName string
}

// NewBackend takes in config, dbfs, loggers, and middleware and registers the backend
// routes.
func NewBackend(
	router *mux.Router,
	logger *zap.Logger,
	db *gorm.DB,
	mw *middleware.Middlewares,
	config *config.Config,
) error {

	bc := &BackendController{
		Db:         db,
		Logger:     logger,
		Path:       config.Stitch.ShardsLocation,
		ServerName: config.Database.ServerName,
	}

	bc.InitialiseShardsFolder()

	// Register routes
	r := router.NewRoute().Subrouter()

	// File
	r.HandleFunc("/api/v1/file", bc.UploadFile).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/update", bc.UpdateFile).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/metadata", bc.GetMetadataFile).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/metadata", bc.UpdateMetadataFile).Methods("PATCH")
	r.HandleFunc("/api/v1/file/{fileID}/move", bc.MoveFile).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/copy", bc.CopyFile).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}", bc.DownloadFileVersion).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}", bc.DeleteFile).Methods("DELETE")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.GetFolderPermissions).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.AddPermissionsFolder).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.UpdateFolderMetadata).Methods("PATCH")
	r.HandleFunc("/api/v1/file/{fileID}/versions/{versionsID}/metadata", bc.GetFileVersionMetadata).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/versions/{versionsID}", bc.DownloadFileVersion).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/versions/{versionsID}", bc.DeleteFileVersion).Methods("DELETE")
	r.HandleFunc("/api/v1/file/{fileID}/versions/", bc.GetFileVersionHistory).Methods("GET")

	// Folder
	r.HandleFunc("/api/v1/folder/{folderID}", bc.LsFolderID).Methods("GET")
	r.HandleFunc("/api/v1/folder/{folderID}", bc.UpdateFolderMetadata).Methods("PATCH")
	r.HandleFunc("/api/v1/folder/{folderID}", bc.DeleteFolder).Methods("DELETE")
	r.HandleFunc("/api/v1/folder", bc.GetFolderIDFromPath).Methods("GET")
	r.HandleFunc("/api/v1/folder", bc.CreateFolder).Methods("POST")
	r.HandleFunc("/api/v1/folder/{folderID}/permissions", bc.GetPermissionsFile).Methods("GET")
	r.HandleFunc("/api/v1/folder/{folderID}/permissions", bc.UpdateMetadataFile).Methods("PATCH")
	r.HandleFunc("/api/v1/folder/{folderID}/permissions", bc.AddPermissionsFolder).Methods("POST")
	r.HandleFunc("/api/v1/folder/{folderID}/move", bc.MoveFolder).Methods("POST")
	r.HandleFunc("/api/v1/file/{folderID}/copy", bc.CopyFile).Methods("POST")
	r.HandleFunc("/api/v1/folder/{folderID}/details", bc.GetMetadataFile).Methods("GET")

	r.Use(mw.UserAuth)

	return nil
}

func (bc *BackendController) InitialiseShardsFolder() {
	if w, err := os.Stat(bc.Path); os.IsNotExist(err) {
		err := os.MkdirAll(bc.Path, 0755)
		if err != nil {
			panic("ERROR. CANNOT CREATE SHARDS FOLDER.")
		}
	} else if !w.IsDir() {
		panic("ERROR. SHARDS FOLDER IS NOT A DIRECTORY.")
	}
}

// UploadFile allows a user to upload a file to the system
func (bc *BackendController) UploadFile(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get header parameters
	folderId := r.Header.Get("folder_id")
	fileName := r.Header.Get("file_name")

	if folderId == "" {
		util.HttpError(w, http.StatusBadRequest, "folder_id is required")
		return
	}
	if fileName == "" {
		util.HttpError(w, http.StatusBadRequest, "file_name is required")
		return
	}

	// Check if the user has the permission to write to the folder
	folder, err := dbfs.GetFileById(bc.Db, folderId, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hasPermissions, err := user.HasPermission(bc.Db, folder, &dbfs.PermissionNeeded{Write: true})
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasPermissions {
		util.HttpError(w, http.StatusForbidden, "You do not have permission to write to this folder")
		return
	}

	// Get encoder params and create new encoder
	stitchParams, err := dbfs.GetStitchParams(bc.Db, bc.Logger)
	dataShards, parityShards, keyThreshold := stitchParams.DataShards, stitchParams.ParityShards, stitchParams.KeyThreshold
	totalShards := dataShards + parityShards

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   uint8(dataShards),
		ParityShards: uint8(parityShards),
		KeyThreshold: uint8(keyThreshold)},
	)

	// This is the fileKey and fileIV for the passwordProtect
	fileKey, fileIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// This is the key and IV for the pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Stream file
	file, header, err := r.FormFile("file")
	defer file.Close()
	fileSize := header.Size

	// File, PasswordProtect entries for dbfs.
	dbfsFile := dbfs.File{
		FileId:             uuid.New().String(),
		FileName:           fileName,
		MIMEType:           "",
		ParentFolderFileId: &folderId, // root folder for now
		Size:               int(fileSize),
		VersioningMode:     dbfs.VersioningOff,
		TotalShards:        totalShards,
		DataShards:         dataShards,
		ParityShards:       parityShards,
		KeyThreshold:       keyThreshold,
		PasswordProtected:  false,
		HandledServer:      bc.ServerName,
	}

	passwordProtect := dbfs.PasswordProtect{
		FileId:         dbfsFile.FileId,
		FileKey:        fileKey,
		FileIv:         fileIv,
		PasswordActive: false,
	}

	err = bc.Db.Transaction(func(tx *gorm.DB) error {

		err := dbfs.CreateInitialFile(tx, &dbfsFile, fileKey, fileIv, dataKey, dataIv, user)
		if err != nil {
			return err
		}

		err = tx.Create(&passwordProtect).Error
		if err != nil {
			return err
		}

		err = dbfs.CreatePermissions(tx, &dbfsFile)
		if err != nil {
			// By right, there should be no error possible? If any error happens, it's likely a system error.
			// However, in the case there is an error, we will revert the transaction (thus deleting the file entry)
			return err
		}
		return nil
	})
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Open the output files
	shardWriters := make([]io.Writer, totalShards)
	shardFiles := make([]*os.File, totalShards)
	shardNames := make([]string, totalShards)
	for i := 0; i < totalShards; i++ {
		shardNames[i] = bc.Path + dbfsFile.DataId + ".shard" + strconv.Itoa(i)
	}
	for i := 0; i < totalShards; i++ {
		shardFile, err := os.Create(shardNames[i])
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer shardFile.Close()
		shardWriters[i] = shardFile
		shardFiles[i] = shardFile
	}

	dataKeyBytes, err := hex.DecodeString(dataKey)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	dataIvBytes, err := hex.DecodeString(dataIv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := encoder.Encode(file, shardWriters, dataKeyBytes, dataIvBytes)

	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for i := 0; i < totalShards; i++ {
		if err = encoder.FinalizeHeader(shardFiles[i]); err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	for i := 1; i <= int(dbfsFile.TotalShards); i++ {
		fragId := int(i)
		fragmentPath := shardNames[i-1]
		serverId := "Server" + strconv.Itoa(i)

		err = dbfs.CreateFragment(bc.Db, dbfsFile.FileId, dbfsFile.DataId, dbfsFile.VersionNo, fragId, serverId, fragmentPath)
		if err != nil {
			err2 := dbfsFile.Delete(bc.Db, user)
			if err2 != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}

			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// checksum
	checksum := hex.EncodeToString(result.FileHash)

	err = dbfs.FinishFile(bc.Db, &dbfsFile, user, 412, checksum)
	if err != nil {
		err2 := dbfsFile.Delete(bc.Db, user)
		errorText := "Error finishing file: " + err.Error()
		if err2 != nil {
			errorText += " Error deleting file: " + err2.Error()
		}
		util.HttpError(w, http.StatusInternalServerError, errorText)
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, dbfsFile)

}

// UpdateFile allows a user to upload a file to the system
func (bc *BackendController) UpdateFile(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	// password from header
	password := r.Header.Get("password")

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "file_id is required")
		return
	}

	// Check if the user has the permission to write to the file
	dbfsFile, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hasPermissions, err := user.HasPermission(bc.Db, dbfsFile, &dbfs.PermissionNeeded{Write: true})
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasPermissions {
		util.HttpError(w, http.StatusForbidden, "You do not have permission to write to this folder")
		return
	}

	// Get encoder params and create new encoder
	stitchParams, err := dbfs.GetStitchParams(bc.Db, bc.Logger)
	dataShards, parityShards, keyThreshold := stitchParams.DataShards, stitchParams.ParityShards, stitchParams.KeyThreshold
	totalShards := dataShards + parityShards

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   uint8(dataShards),
		ParityShards: uint8(parityShards),
		KeyThreshold: uint8(keyThreshold)},
	)

	// This is the key and IV for the pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Stream file
	file, header, err := r.FormFile("file")
	defer file.Close()
	fileSize := header.Size

	err = dbfsFile.UpdateFile(bc.Db, int(fileSize), int(fileSize), "TODO:CHECKSUM",
		"", dataKey, dataIv, password, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Open the output files
	shardWriters := make([]io.Writer, totalShards)
	shardFiles := make([]*os.File, totalShards)
	shardNames := make([]string, totalShards)
	for i := 0; i < totalShards; i++ {
		shardNames[i] = bc.Path + dbfsFile.DataId + ".shard" + strconv.Itoa(i)
	}
	for i := 0; i < totalShards; i++ {
		shardFile, err := os.Create(shardNames[i])
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer shardFile.Close()
		shardWriters[i] = shardFile
		shardFiles[i] = shardFile
	}

	dataKeyBytes, err := hex.DecodeString(dataKey)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	dataIvBytes, err := hex.DecodeString(dataIv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := encoder.Encode(file, shardWriters, dataKeyBytes, dataIvBytes)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for i := 0; i < totalShards; i++ {
		if err = encoder.FinalizeHeader(shardFiles[i]); err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	for i := 1; i <= int(dbfsFile.TotalShards); i++ {
		fragId := int(i)
		fragmentPath := shardNames[i-1]
		serverId := "Server" + strconv.Itoa(i)
		fragChecksum := "CHECKSUM" + strconv.Itoa(i)

		err = dbfsFile.UpdateFragment(bc.Db, fragId, fragmentPath, fragChecksum, serverId)
		if err != nil {
			err2 := dbfsFile.RevertFileToVersion(bc.Db, dbfsFile.VersionNo-1, user)
			if err2 != nil {
				util.HttpError(w, http.StatusInternalServerError, err.Error())
				return
			}

			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// checksum
	checksum := hex.EncodeToString(result.FileHash)

	err = dbfsFile.FinishUpdateFile(bc.Db, checksum)
	if err != nil {
		errString := err.Error()
		err2 := dbfsFile.RevertFileToVersion(bc.Db, dbfsFile.VersionNo-1, user)
		if err2 != nil {
			errString += " " + err2.Error()
		}
		util.HttpError(w, http.StatusInternalServerError, errString)
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, dbfsFile)

}

// GetMetadataFile returns the metadata for the requested file based on ID
func (bc *BackendController) GetMetadataFile(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
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
	user, err := ctxutil.GetUser(r.Context())
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

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update file
	err = file.UpdateMetaData(bc.Db, fmm, user)
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
	user, err := ctxutil.GetUser(r.Context())
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

	// Getting file and dest folder from Db

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	destFolder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Move file (Permission check will be done by dbfs)
	err = file.Move(bc.Db, destFolder, user)
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
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Getting params
	vars := mux.Vars(r)
	ogFileID, fileOk := vars["fileID"]
	ogFolderID, folderOk := vars["folderID"]
	if !fileOk && !folderOk {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	var fileID string
	if !fileOk {
		fileID = ogFolderID
	} else {
		fileID = ogFileID
	}

	folderID := r.Header.Get("folder_id")
	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// Getting file and dest folder from Db
	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, "File not found")
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	destFolder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, "Destination folder not found")
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Copy file (Permission check will be done by dbfs)
	err = file.Copy(bc.Db, destFolder, user, bc.ServerName)
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

//DeleteFile deletes a file from the Db
func (bc *BackendController) DeleteFile(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete file
	err = file.Delete(bc.Db, user)
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
	user, err := ctxutil.GetUser(r.Context())
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

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get permissions
	permissions, err := file.GetPermissions(bc.Db, user)
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
	user, err := ctxutil.GetUser(r.Context())
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
	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if file is actually a file
	if file.EntryType != dbfs.IsFile {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	// get old permission struct entry
	oldPermission, err := file.GetPermissionById(bc.Db, permissionID, user)
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

	err = file.UpdatePermission(bc.Db, oldPermission, newPermission, user)

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
	user, err := ctxutil.GetUser(r.Context())
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
	usersString := r.Header.Get("users")
	groupsString := r.Header.Get("groups")

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

	var usersStringArray []string
	var groupsStringArray []string
	_ = json.Unmarshal([]byte(usersString), &usersStringArray)
	_ = json.Unmarshal([]byte(groupsString), &groupsStringArray)

	permissionNeeded := dbfs.PermissionNeeded{
		Read:    canReadBool,
		Write:   canWriteBool,
		Execute: canExecuteBool,
		Share:   canShareBool,
		Audit:   canAuditBool,
	}

	userObjects := make([]dbfs.User, len(usersStringArray))
	groupObjects := make([]dbfs.Group, len(groupsStringArray))

	for i, userString := range usersStringArray {
		tempUser, err := dbfs.GetUserById(bc.Db, userString)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		userObjects[i] = *tempUser
	}
	for i, groupString := range groupsStringArray {
		tempGroup, err := dbfs.GetGroupById(bc.Db, groupString)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		groupObjects[i] = *tempGroup
	}

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// add permissions
	err = file.AddPermissionUsers(bc.Db, &permissionNeeded, user, userObjects...)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = file.AddPermissionGroups(bc.Db, &permissionNeeded, user, groupObjects...)
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
	user, err := ctxutil.GetUser(r.Context())
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

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get version
	version, err := file.GetOldVersion(bc.Db, user, versionIDInt)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, version)

}

//GetFileVersionHistory returns all the historical versions of a file
func (bc *BackendController) GetFileVersionHistory(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
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

	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get version
	versions, err := file.GetAllVersions(bc.Db, user)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, versions)

}

//DownloadFileVersion downloads a file version
func (bc *BackendController) DownloadFileVersion(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
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

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID  provided")
		return
	}

	// convert versionID into int if it is not empty
	var versionIDInt int
	if versionID != "" {
		versionIDInt, err = strconv.Atoi(versionID)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// get file
	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// meaning none provided
	if versionID == "" {
		versionIDInt = file.VersionNo
	}

	// get version
	version, err := file.GetOldVersion(bc.Db, user, versionIDInt)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	// get file
	shardsMeta, err := version.GetFragments(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	// TODO: Now assuming that all fragments are on the same server, else you need to query multiple
	// servers for the file fragments

	// Opening input files
	var shards []io.ReadSeeker
	var shardFiles []*os.File
	for _, shardMeta := range shardsMeta {
		shardFile, err := os.Open(shardMeta.FileFragmentPath)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		shards = append(shards, shardFile)
		shardFiles = append(shardFiles, shardFile)
	}

	hexKey, hexIv, err := version.GetDecryptionKey(bc.Db, user, password)
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

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   uint8(version.DataShards),
		ParityShards: uint8(version.ParityShards),
		KeyThreshold: uint8(version.KeyThreshold)},
	)

	// Decode file
	reader, err := encoder.NewReadSeeker(shards, key, iv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", file.MIMEType)

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
	user, err := ctxutil.GetUser(r.Context())
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
	file, err := dbfs.GetFileById(bc.Db, fileID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// delete version
	err = file.DeleteFileVersion(bc.Db, user, versionIDInt)
	if errors.Is(err, dbfs.ErrVersionNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, nil)

}

// LsFolderID lists the contents of a folder based on the folderID
func (bc *BackendController) LsFolderID(w http.ResponseWriter, r *http.Request) {
	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// get folder
	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get folder contents
	files, err := folder.ListContents(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return contents
	util.HttpJson(w, http.StatusOK, files)
}

// UpdateFolderMetadata updates the filename or the versioningMode of a folder
func (bc *BackendController) UpdateFolderMetadata(w http.ResponseWriter, r *http.Request) {
	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// get headers
	newName := r.Header.Get("new_name")
	versioningMode := r.Header.Get("versioning_mode")

	versioningModeBool := (versioningMode != "")

	// get folder
	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var versioningModeInt int8

	if versioningModeBool {
		tempVersioningModeInt, err := strconv.Atoi(versioningMode)
		versioningModeInt = int8(tempVersioningModeInt)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		versioningModeInt = folder.VersioningMode
	}

	// attempt to amend folder

	fmm := dbfs.FileMetadataModification{
		FileName:       newName,
		VersioningMode: versioningModeInt,
	}

	err = folder.UpdateMetaData(bc.Db, fmm, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return contents
	util.HttpJson(w, http.StatusOK, nil)
}

// DeleteFolder deletes a folder
func (bc *BackendController) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// get folder
	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if folder is actually a folder
	if folder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	// attempt to delete folder
	err = folder.Delete(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return contents
	util.HttpJson(w, http.StatusOK, nil)
}

// GetFolderIDFromPath list the contents of a folder based on the folderPath
func (bc *BackendController) GetFolderIDFromPath(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderPath from header
	folderPath := r.Header.Get("folder_path")

	if folderPath == "" {
		util.HttpError(w, http.StatusBadRequest, "No folder_path provided")
		return
	}

	// get folder
	folder, err := dbfs.GetFileByPath(bc.Db, folderPath, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return contents
	util.HttpJson(w, http.StatusOK, folder)
}

// CreateFolder creates a folder given a folder_name, parent_folder_id or parent_folder_path
func (bc *BackendController) CreateFolder(w http.ResponseWriter, r *http.Request) {
	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get headers
	folderName := r.Header.Get("folder_name")
	parentFolderID := r.Header.Get("parent_folder_id")
	parentFolderPath := r.Header.Get("parent_folder_path")

	// check if folder_name is provided
	if folderName == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing new folder name")
		return
	}

	useID := true

	if parentFolderID == "" && parentFolderPath == "" {
		util.HttpError(w, http.StatusBadRequest, "Path or Folder ID not provided")
		return
	} else if parentFolderID == "" {
		useID = false
	}

	// get parent folder

	var parentFolder *dbfs.File

	if useID {
		parentFolder, err = dbfs.GetFileById(bc.Db, parentFolderID, user)
		if errors.Is(err, dbfs.ErrFileNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
			return
		}
	} else {
		parentFolder, err = dbfs.GetFileByPath(bc.Db, parentFolderPath, user)
		if errors.Is(err, dbfs.ErrFileNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
			return
		}
	}

	// Check that parent folder is actually a folder
	if parentFolder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "Parent folder is not a folder")
		return
	}

	// attempt to create folder (permisison check here as well)
	newFolder, err := parentFolder.CreateSubFolder(bc.Db, folderName, user, bc.ServerName)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return new folder
	util.HttpJson(w, http.StatusOK, newFolder)

}

// GetFolderPermissions returns the permissions of a folder
func (bc *BackendController) GetFolderPermissions(w http.ResponseWriter, r *http.Request) {
	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	// get folder
	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if folder is actually a folder
	if folder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	// attempt to get permissions
	permissions, err := folder.GetPermissions(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return permissions
	util.HttpJson(w, http.StatusOK, permissions)
}

// ModifyPermissionsFolder modifies the permissions of a folder based on the id obtained from GetFolderPermissions
func (bc *BackendController) ModifyPermissionsFolder(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	// get other headers
	permissionID := r.Header.Get("permission_id")
	canRead := r.Header.Get("can_read")
	canWrite := r.Header.Get("can_write")
	canExecute := r.Header.Get("can_execute")
	canShare := r.Header.Get("can_share")
	canAudit := r.Header.Get("can_audit")

	// if any empty return error
	if folderID == "" || permissionID == "" || canRead == "" || canWrite == "" || canExecute == "" || canShare == "" || canAudit == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing params")
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
	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if folder is actually a folder
	if folder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	// get old permission struct entry
	oldPermission, err := folder.GetPermissionById(bc.Db, permissionID, user)
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

	err = folder.UpdatePermission(bc.Db, oldPermission, newPermission, user)

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

// AddPermissionsFolder ads permission to a folder
func (bc *BackendController) AddPermissionsFolder(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	// get other headers
	canRead := r.Header.Get("can_read")
	canWrite := r.Header.Get("can_write")
	canExecute := r.Header.Get("can_execute")
	canShare := r.Header.Get("can_share")
	canAudit := r.Header.Get("can_audit")
	usersString := r.Header.Get("users")
	groupsString := r.Header.Get("groups")

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

	var usersStringArray []string
	var groupsStringArray []string
	_ = json.Unmarshal([]byte(usersString), &usersStringArray)
	_ = json.Unmarshal([]byte(groupsString), &groupsStringArray)

	userObjects := make([]dbfs.User, len(usersStringArray))
	groupObjects := make([]dbfs.Group, len(groupsStringArray))

	for i, userString := range usersStringArray {
		tempUser, err := dbfs.GetUserById(bc.Db, userString)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		userObjects[i] = *tempUser
	}
	for i, groupString := range groupsStringArray {
		tempGroup, err := dbfs.GetGroupById(bc.Db, groupString)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		groupObjects[i] = *tempGroup
	}

	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if folder is actually a folder
	if folder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	// add permissions
	err = folder.AddPermissionUsers(bc.Db, &permissionNeeded, user, userObjects...)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = folder.AddPermissionGroups(bc.Db, &permissionNeeded, user, groupObjects...)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Success
	util.HttpJson(w, http.StatusOK, nil)
}

// MoveFolder moves a folder to a new location
func (bc *BackendController) MoveFolder(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	// get other headers
	newFolderID := r.Header.Get("folder_id")

	if folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folder_id provided")
		return
	}

	// Getting og folder and dest folder from Db

	oldFolder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check that old folder is actually a folder
	if oldFolder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	destFolder, err := dbfs.GetFileById(bc.Db, newFolderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Move file (Permission check will be done by dbfs)
	err = oldFolder.Move(bc.Db, destFolder, user)
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

// GetFolderDetails returns details about a folder from folderID
func (bc *BackendController) GetFolderDetails(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]

	// get file
	folder, err := dbfs.GetFileById(bc.Db, folderID, user)
	if errors.Is(err, dbfs.ErrFileNotFound) {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if folder is actually a folder
	if folder.EntryType != dbfs.IsFolder {
		util.HttpError(w, http.StatusBadRequest, "File is not a folder")
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, folder)
}
