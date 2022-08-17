package controller

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/stitch"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BackendController struct {
	Db         *gorm.DB
	Logger     *zap.Logger
	Path       string
	ServerName string
	Inc        *inc.Inc
}

// NewBackend takes in config, dbfs, loggers, and middleware and registers the backend
// routes.
func NewBackend(
	router *mux.Router,
	logger *zap.Logger,
	db *gorm.DB,
	mw *middleware.Middlewares,
	config *config.Config,
	inc *inc.Inc,
) error {

	bc := &BackendController{
		Db:         db,
		Logger:     logger,
		Path:       config.Stitch.ShardsLocation,
		ServerName: config.Inc.ServerName,
		Inc:        inc,
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
	r.HandleFunc("/api/v1/file/{fileID}/path", bc.GetPath).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}", bc.DownloadFileVersion).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}", bc.DeleteFile).Methods("DELETE")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.GetFolderPermissions).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.AddPermissionsFolder).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/permissions", bc.UpdateFolderMetadata).Methods("PATCH")
	r.HandleFunc("/api/v1/file/{fileID}/share", bc.GetFileSharedLinks).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/share", bc.CreateFileSharedLink).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/share/{link}", bc.PatchFileSharedLink).Methods("PATCH")
	r.HandleFunc("/api/v1/file/{fileID}/share/{link}", bc.DeleteFileSharedLink).Methods("DELETE")
	r.HandleFunc("/api/v1/file/{fileID}/share/{link}", bc.CreateFileSharedLink).Methods("POST")
	r.HandleFunc("/api/v1/file/{fileID}/versions/{versionsID}/metadata", bc.GetFileVersionMetadata).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/versions/{versionsID}", bc.DownloadFileVersion).Methods("GET")
	r.HandleFunc("/api/v1/file/{fileID}/versions/{versionsID}", bc.DeleteFileVersion).Methods("DELETE")
	r.HandleFunc("/api/v1/file/{fileID}/versions", bc.GetFileVersionHistory).Methods("GET")

	// Folder
	r.HandleFunc("/api/v1/folder/{folderID}", bc.LsFolderID).Methods("GET")
	r.HandleFunc("/api/v1/folder/{folderID}", bc.UpdateFolderMetadata).Methods("PATCH")
	r.HandleFunc("/api/v1/folder/{folderID}", bc.DeleteFolder).Methods("DELETE")
	r.HandleFunc("/api/v1/folder/{folderID}/path", bc.GetPath).Methods("GET")
	r.HandleFunc("/api/v1/folder", bc.GetFolderIDFromPath).Methods("GET")
	r.HandleFunc("/api/v1/folder", bc.CreateFolder).Methods("POST")
	r.HandleFunc("/api/v1/folder/{folderID}/permissions", bc.GetPermissionsFile).Methods("GET")
	r.HandleFunc("/api/v1/folder/{folderID}/permissions", bc.UpdateMetadataFile).Methods("PATCH")
	r.HandleFunc("/api/v1/folder/{folderID}/permissions", bc.AddPermissionsFolder).Methods("POST")
	r.HandleFunc("/api/v1/folder/{folderID}/move", bc.MoveFolder).Methods("POST")
	r.HandleFunc("/api/v1/file/{folderID}/copy", bc.CopyFile).Methods("POST")
	r.HandleFunc("/api/v1/folder/{folderID}/details", bc.GetMetadataFile).Methods("GET")

	// Get Favorites, Get Shared
	r.HandleFunc("/api/v1/favorites", bc.GetFavorites).Methods("GET")
	r.HandleFunc("/api/v1/favorites/{fileID}", bc.GetFavoriteItem).Methods("GET")
	r.HandleFunc("/api/v1/favorites/{fileID}", bc.AddFavorite).Methods("PUT")
	r.HandleFunc("/api/v1/favorites/{fileID}", bc.RemoveFavorite).Methods("DELETE")
	r.HandleFunc("/api/v1/sharedWith", bc.GetSharedWithUser).Methods("GET")

	// Shared Routes
	// Use a fresh subrouter to skip auth
	rPub := router.NewRoute().Subrouter()
	rPub.HandleFunc("/api/v1/shared/{shortenedLink}/metadata", bc.GetMetadataSharedLink).Methods("GET")
	rPub.HandleFunc("/api/v1/shared/{shortenedLink}", bc.DownloadSharedLink).Methods("GET")

	// Cluster Routes
	r.HandleFunc("/api/v1/cluster/stats/num_of_files", bc.GetNumOfFiles).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/num_of_files_historical", bc.GetNumOfFilesHistorical).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/non_replica_used", bc.GetStorageUsed).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/non_replica_used_historical", bc.GetStorageUsedHistorical).
		Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/replica_used", bc.GetStorageUsedReplica).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/replica_used_historical", bc.GetStorageUsedReplicaHistorical).
		Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/alerts", bc.GetAllAlerts).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/alerts", bc.ClearAllAlerts).Methods("DELETE")
	r.HandleFunc("/api/v1/cluster/stats/alerts/{id}", bc.GetAlert).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/alerts/{id}", bc.ClearAlert).Methods("DELETE")
	r.HandleFunc("/api/v1/cluster/stats/logs", bc.GetAllLogs).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/logs", bc.ClearAllLogs).Methods("DELETE")
	r.HandleFunc("/api/v1/cluster/stats/logs/{id}", bc.GetLog).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/logs/{id}", bc.ClearLog).Methods("DELETE")
	r.HandleFunc("/api/v1/cluster/stats/servers", bc.GetServerStatuses).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/servers/{serverName}", bc.GetSpecificServerStatus).Methods("GET")
	r.HandleFunc("/api/v1/cluster/stats/servers/{serverName}", bc.DeleteServer).Methods("DELETE")

	// Maintenance Routes:
	r.HandleFunc("/api/v1/maintenance/all", bc.GetAllJobs).Methods("GET")
	r.HandleFunc("/api/v1/maintenance/start", bc.StartJob).Methods("POST")
	r.HandleFunc("/api/v1/maintenance/job/{id}", bc.GetJob).Methods("GET")
	r.HandleFunc("/api/v1/maintenance/job/{id}", bc.DeleteJob).Methods("DELETE")
	r.HandleFunc("/api/v1/maintenance/job/{id}/full_shards", bc.GetFullShardsResult).Methods("GET")
	r.HandleFunc("/api/v1/maintenance/job/{id}/full_shards", bc.FixFullShardsResult).Methods("POST")
	r.HandleFunc("/api/v1/maintenance/job/{id}/quick_shards", bc.GetQuickShardsResult).Methods("GET")
	r.HandleFunc("/api/v1/maintenance/job/{id}/quick_shards", bc.FixQuickShardsResult).Methods("POST")
	r.HandleFunc("/api/v1/maintenance/job/{id}/orphaned_shards", bc.GetOrphanedShardsResult).Methods("GET")
	r.HandleFunc("/api/v1/maintenance/job/{id}/orphaned_shards", bc.FixOrphanedShardsResult).Methods("POST")
	r.HandleFunc("/api/v1/maintenance/job/{id}/orphaned_files", bc.GetOrphanedFilesResult).Methods("GET")
	r.HandleFunc("/api/v1/maintenance/job/{id}/orphaned_files", bc.FixOrphanedFilesResult).Methods("POST")

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

type MultipartFileStream struct {
	reader      *multipart.Reader
	part        *multipart.Part
	formName    string
	isDone      bool
	ContentType string
	Filename    string
}

var _ io.ReadCloser = &MultipartFileStream{}

// Read implements io.ReadCloser
func (mf *MultipartFileStream) Read(p []byte) (int, error) {
	if mf.isDone {
		return 0, io.EOF
	}

	n, err := io.ReadAtLeast(mf.part, p, len(p))
	if n < len(p) {
		mf.isDone = true
		return n, nil
	}

	return n, err
}

// Close implements io.ReadCloser
func (mf *MultipartFileStream) Close() error {
	return mf.part.Close()
}

func NewMultipartFileReader(req *http.Request, fieldName string) (*MultipartFileStream, error) {
	_, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Content-Type header: %w", err)
	}
	boundary, ok := params["boundary"]
	if !ok {
		return nil, fmt.Errorf("Failed to get multipart boundary")
	}
	reader := multipart.NewReader(req.Body, boundary)

	mfs := &MultipartFileStream{
		reader:      reader,
		part:        nil,
		formName:    fieldName,
		ContentType: "",
		Filename:    "",
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		for part.FormName() != fieldName {
			io.Copy(io.Discard, part)
			part.Close()

			if part, err = reader.NextPart(); err == io.EOF {
				break
			}
		}

		mfs.part = part
		mfs.ContentType = part.Header.Get("Content-Type")
		mfs.Filename = part.FileName()

		return mfs, nil
	}

	return nil, fmt.Errorf("No matching field name")
}

// UploadFile allows a user to upload a file to the system
func (bc *BackendController) UploadFile(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Get file
	file, err := NewMultipartFileReader(r, "file")
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "file not found: "+err.Error())
		return
	}
	defer file.Close()

	// Get parameters
	fileName := file.Filename
	folderId := r.Header.Get("folder_id")

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
	dataShards, parityShards, keyThreshold :=
		stitchParams.DataShards, stitchParams.ParityShards, stitchParams.KeyThreshold
	totalShards := dataShards + parityShards

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   uint8(dataShards),
		ParityShards: uint8(parityShards),
		KeyThreshold: uint8(keyThreshold),
	})

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

	// File, PasswordProtect entries for dbfs.
	dbfsFile := dbfs.File{
		FileId:             uuid.New().String(),
		FileName:           fileName,
		MIMEType:           file.ContentType,
		ParentFolderFileId: &folderId, // root folder for now
		Size:               1024,      // placeholder size
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
		if err := dbfs.CreateInitialFile(tx, &dbfsFile,
			fileKey, fileIv, dataKey, dataIv, user); err != nil {
			return fmt.Errorf("failed to create initial file: %w", err)
		}

		if err := tx.Create(&passwordProtect).Error; err != nil {
			return fmt.Errorf("failed to create PasswordProtect row: %w", err)
		}

		if err := dbfs.CreatePermissions(tx, &dbfsFile); err != nil {
			// By right, there should be no error possible? If any error happens, it's
			// likely a system error. However, in the case there is an error, we will
			// revert the transaction (thus deleting the file entry).
			return fmt.Errorf("failed to create permissions: %w", err)
		}
		return nil
	})
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create file metadata: %s", err.Error()))
		return
	}

	// Fetch a list of servers
	servers, err := bc.Inc.AssignShardServer(r.Context(), totalShards)
	if err != nil {
		util.HttpError(w,
			http.StatusInternalServerError,
			fmt.Sprintf("failed to assign servers: %s", err.Error()),
		)
		return
	}

	// Prepare the writers
	shardWriters := make([]io.Writer, totalShards)
	shardNames := make([]string, totalShards)

	// Generate names for the shards
	for i := 0; i < totalShards; i++ {
		shardNames[i] = dbfsFile.DataId + ".shard" + strconv.Itoa(i)
	}

	// Open the output writers
	for i := 0; i < totalShards; i++ {
		shardWriter, err := bc.Inc.NewShardWriter(
			r.Context(), servers[i].Name, shardNames[i])
		if err != nil {
			util.HttpError(w,
				http.StatusInternalServerError,
				fmt.Sprintf("failed to initialize shard writer: %s", err.Error()),
			)
			return
		}
		shardWriters[i] = shardWriter
	}

	// Decode the data key and iv
	dataKeyBytes, err := hex.DecodeString(dataKey)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to decode data key: %s", err.Error()))
		return
	}
	dataIvBytes, err := hex.DecodeString(dataIv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to decode data iv: %s", err.Error()))
		return
	}

	// Encode the file
	result, err := encoder.Encode(file, shardWriters, dataKeyBytes, dataIvBytes)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to encode file: %s", err.Error()))
		return
	}

	// Close the writers
	for _, writer := range shardWriters {
		if err := writer.(io.WriteCloser).Close(); err != nil {
			util.HttpError(w, http.StatusInternalServerError,
				fmt.Sprintf("failed to close shard writer: %s", err.Error()))
			return
		}
	}

	// Insert fragments into the database
	err = bc.Db.Transaction(func(tx *gorm.DB) error {
		for i := 1; i <= int(dbfsFile.TotalShards); i++ {
			fragId := int(i)
			fragmentPath := shardNames[i-1]

			if err := dbfs.CreateFragment(tx,
				dbfsFile.FileId, dbfsFile.DataId, dbfsFile.VersionNo,
				fragId, servers[i-1].Name, fragmentPath); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if err := dbfsFile.Delete(bc.Db, user, bc.ServerName); err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checksum
	checksum := hex.EncodeToString(result.FileHash)

	dbfsFile.Size = int(result.FileSize)
	err = dbfs.FinishFile(bc.Db, &dbfsFile, user, int(result.FileSize), checksum)
	if err != nil {
		err2 := dbfsFile.Delete(bc.Db, user, bc.ServerName)
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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

	// Get file
	file, err := NewMultipartFileReader(r, "file")
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "file not found: "+err.Error())
		return
	}
	defer file.Close()

	// Use placeholder size values as it is not yet known at this point
	err = dbfsFile.UpdateFile(bc.Db, 1024, 1024, "TODO:CHECKSUM", "", dataKey, dataIv, password, user, file.Filename)
	if err != nil {
		if errors.Is(err, dbfs.ErrIncorrectPassword) {
			util.HttpError(w, http.StatusForbidden, err.Error())
			return
		}
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch a list of servers
	servers, err := bc.Inc.AssignShardServer(r.Context(), totalShards)
	if err != nil {
		util.HttpError(w,
			http.StatusInternalServerError,
			fmt.Sprintf("failed to assign servers: %s", err.Error()),
		)
		return
	}

	// Prepare the writers
	shardWriters := make([]io.Writer, totalShards)
	shardNames := make([]string, totalShards)

	// Generate names for the shards
	for i := 0; i < totalShards; i++ {
		shardNames[i] = dbfsFile.DataId + ".shard" + strconv.Itoa(i)
	}

	// Open the output writers
	for i := 0; i < totalShards; i++ {
		shardWriter, err := bc.Inc.NewShardWriter(
			r.Context(), servers[i].Name, shardNames[i])
		if err != nil {
			util.HttpError(w,
				http.StatusInternalServerError,
				fmt.Sprintf("failed to initialize shard writer: %s", err.Error()),
			)
			return
		}
		shardWriters[i] = shardWriter
	}

	// Decode the data key and iv
	dataKeyBytes, err := hex.DecodeString(dataKey)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to decode data key: %s", err.Error()))
		return
	}
	dataIvBytes, err := hex.DecodeString(dataIv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to decode data iv: %s", err.Error()))
		return
	}

	// Encode the file
	result, err := encoder.Encode(file, shardWriters, dataKeyBytes, dataIvBytes)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to encode file: %s", err.Error()))
		return
	}

	// Close the writers
	for _, writer := range shardWriters {
		if err := writer.(io.WriteCloser).Close(); err != nil {
			util.HttpError(w, http.StatusInternalServerError,
				fmt.Sprintf("failed to close shard writer: %s", err.Error()))
			return
		}
	}

	// Insert fragments into the database
	for i := 1; i <= int(dbfsFile.TotalShards); i++ {
		fragId := int(i)
		fragmentPath := shardNames[i-1]

		// TODO: Figure out checksum
		fragChecksum := "CHECKSUM" + strconv.Itoa(i)

		err = dbfsFile.UpdateFragment(bc.Db, fragId, fragmentPath, fragChecksum, bc.ServerName)
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

	dbfsFile.Size = int(result.FileSize)
	dbfsFile.ActualSize = int(result.FileSize)
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
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

	// Delete file
	err = file.Delete(bc.Db, user, bc.ServerName)
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	versionID := vars["versionID"]
	if !ok || fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}
	versionID := vars["versionID"]

	queries := r.URL.Query()
	inline := queries.Get("inline")

	isDownload := true

	// get password
	password := r.Header.Get("password")

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}

	if inline == "true" {
		isDownload = false
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

	// Opening input files
	var shards []io.ReadSeeker
	for _, shardMeta := range shardsMeta {
		shardReader, err := bc.Inc.NewShardReader(
			r.Context(), shardMeta.ServerName, shardMeta.FileFragmentPath)
		if err == nil {
			shards = append(shards, shardReader)
		}
	}

	hexKey, hexIv, err := version.GetDecryptionKey(bc.Db, user, password)
	if err != nil {
		if errors.Is(err, dbfs.ErrIncorrectPassword) {
			util.HttpError(w, http.StatusForbidden, err.Error())
			return
		}
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

	w.Header().Set("Content-Type", file.MIMEType)
	if isDownload {
		w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	} else {
		w.Header().Set("Content-Disposition", "inline; filename="+file.FileName)
	}

	http.ServeContent(w, r, file.FileName, file.ModifiedTime, reader)
}

// DeleteFileVersion deletes a file version
func (bc *BackendController) DeleteFileVersion(w http.ResponseWriter, r *http.Request) {

	// somehow get user idk
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID and versionID
	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID provided")
		return
	}
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID, ok := vars["folderID"]
	if !ok || folderID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folderID provided")
		return
	}

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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
	err = folder.Delete(bc.Db, user, bc.ServerName)
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// folderPath from header
	folderPath := r.Header.Get("folder_path")

	if folderPath == "" {
		util.HttpError(w, http.StatusBadRequest, "No folder_path provided")
		return
	}

	// get folder
	folder, err := dbfs.GetFileByPath(bc.Db, folderPath, user, true)
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusBadRequest, "FileName or Folder ID not provided")
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
		parentFolder, err = dbfs.GetFileByPath(bc.Db, parentFolderPath, user, true)
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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
		util.HttpError(w, http.StatusUnauthorized, err.Error())
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

func (bc *BackendController) GetPath(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// folderID
	vars := mux.Vars(r)
	folderID := vars["folderID"]
	fileID := vars["fileID"]
	if folderID == "" && fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No folder_id or file_id provided")
		return
	} else if folderID != "" && fileID != "" {
		util.HttpError(w, http.StatusBadRequest, "Both folder_id and file_id provided")
		return
	} else if folderID == "" && fileID != "" {
		folderID = fileID
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

	// success

	folders, err := folder.GetPath(bc.Db, user)

	util.HttpJson(w, http.StatusOK, folders)
}

// GetFileSharedLinks returns all shared links for a file
func (bc *BackendController) GetFileSharedLinks(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No FileID provided")
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
	// Check if file is actually a file
	if file.EntryType != dbfs.IsFile {
		util.HttpError(w, http.StatusBadRequest, "File is not a file")
		return
	}

	// Get File Shared Links
	sharedLinks, err := file.GetSharedLinks(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("Error getting shared links: %s", err.Error()))
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, sharedLinks)
}

// CreateFileSharedLink creates a shared link for a file
func (bc *BackendController) CreateFileSharedLink(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	link := vars["link"]

	if fileID == "" {
		util.HttpError(w, http.StatusBadRequest, "No FileID provided")
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
	// Check if file is actually a file
	if file.EntryType != dbfs.IsFile {
		util.HttpError(w, http.StatusBadRequest, "File is not a file")
		return
	}

	// Create Shared Link
	sharedLink, err := file.CreateSharedLink(bc.Db, user, link)
	if err != nil {
		if errors.Is(err, dbfs.ErrLinkExists) {
			util.HttpError(w, http.StatusConflict, err.Error())
			return
		}
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("Error creating shared link: %s", err.Error()))
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, sharedLink)
}

// DeleteFileSharedLink deletes a shared link for a file
func (bc *BackendController) DeleteFileSharedLink(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	link := vars["link"]
	if fileID == "" || link == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID or link provided")
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
	// Check if file is actually a file
	if file.EntryType != dbfs.IsFile {
		util.HttpError(w, http.StatusBadRequest, "File is not a file")
		return
	}

	// Delete Shared Link
	err = file.DeleteSharedLink(bc.Db, user, link)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("Error deleting shared link: %s", err.Error()))
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, true)
}

// PatchFileSharedLink updates a shared link for a file
func (bc *BackendController) PatchFileSharedLink(w http.ResponseWriter, r *http.Request) {
	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Get link from headers
	newLink := r.Header.Get("new_link")

	// fileID
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	link := vars["link"]
	if fileID == "" || link == "" || newLink == "" {
		util.HttpError(w, http.StatusBadRequest, "No fileID, link, new_link provided")
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
	// Check if file is actually a file
	if file.EntryType != dbfs.IsFile {
		util.HttpError(w, http.StatusBadRequest, "File is not a file")
		return
	}

	// Update Shared Link
	err = file.UpdateSharedLink(bc.Db, user, link, newLink)
	if err != nil {
		if errors.Is(err, dbfs.ErrLinkExists) {
			util.HttpError(w, http.StatusConflict, err.Error())
			return
		}
		util.HttpError(w, http.StatusInternalServerError,
			fmt.Sprintf("Error updating shared link: %s", err.Error()))
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, true)
}

// GetMetadataSharedLink returns the metadata for the requested file based on ID
func (bc *BackendController) GetMetadataSharedLink(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	shortenedLink := vars["shortenedLink"]

	if shortenedLink == "" {
		util.HttpError(w, http.StatusBadRequest, "No shortenedLink provided")
		return
	}

	// get file
	file, err := dbfs.GetFileFromShortenedLink(bc.Db, shortenedLink)
	if err != nil {
		if errors.Is(err, dbfs.ErrSharedLinkNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
			return
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	// json encode file
	util.HttpJson(w, http.StatusOK, file)
}

//DownloadSharedLink downloads a file version
func (bc *BackendController) DownloadSharedLink(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	shortenedLink := vars["shortenedLink"]

	queries := r.URL.Query()
	inline := queries.Get("inline")

	isDownload := true

	// get password
	password := r.Header.Get("password")

	if shortenedLink == "" {
		util.HttpError(w, http.StatusBadRequest, "No shortenedLink provided")
		return
	}

	if inline == "true" {
		isDownload = false
	}

	// get file
	file, err := dbfs.GetFileFromShortenedLink(bc.Db, shortenedLink)
	if err != nil {
		if errors.Is(err, dbfs.ErrSharedLinkNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
			return
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// get file
	shardsMeta, err := file.GetFileFragments(bc.Db, nil) // nil will check if the file is public
	if err != nil {
		util.HttpError(w, http.StatusNotFound, err.Error())
		return
	}

	// Opening input files
	var shards []io.ReadSeeker
	for _, shardMeta := range shardsMeta {
		shardReader, err := bc.Inc.NewShardReader(
			r.Context(), shardMeta.ServerName, shardMeta.FileFragmentPath)
		if err == nil {
			shards = append(shards, shardReader)
		}
	}

	hexKey, hexIv, err := file.GetDecryptionKey(bc.Db, nil, password)
	if err != nil {
		if errors.Is(err, dbfs.ErrIncorrectPassword) {
			util.HttpError(w, http.StatusForbidden, err.Error())
			return
		}
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
		DataShards:   uint8(file.DataShards),
		ParityShards: uint8(file.ParityShards),
		KeyThreshold: uint8(file.KeyThreshold)},
	)

	// Decode file
	reader, err := encoder.NewReadSeeker(shards, key, iv)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", file.MIMEType)
	if isDownload {
		w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	} else {
		w.Header().Set("Content-Disposition", "inline; filename="+file.FileName)
	}

	http.ServeContent(w, r, file.FileName, file.ModifiedTime, reader)
}
