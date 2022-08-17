package controller

import (
	"errors"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetSharedWithUser is a route to get the files shared with a user
func (bc *BackendController) GetSharedWithUser(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	files, err := user.GetSharedWithUser(bc.Db)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, files)

}

// GetFavorites is a route to get the files favorited by a user
func (bc *BackendController) GetFavorites(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	startNumString := r.Header.Get("start_num")
	startNum := 0
	if startNumString != "" {
		startNum, err = strconv.Atoi(startNumString)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	files, err := user.GetFavoriteFiles(bc.Db, uint(startNum))
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, files)
}

// GetFavoriteItem returns the favorite item for the given file and user if exists
func (bc *BackendController) GetFavoriteItem(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	fileId, ok := mux.Vars(r)["fileID"]
	if !ok {
		util.HttpError(w, http.StatusBadRequest, "No fileID specified")
		return
	} else if fileId == "" {
		util.HttpError(w, http.StatusBadRequest, "fileID is required")
		return
	}

	file, err := user.GetFavoriteFileByFileId(bc.Db, fileId)
	if err != nil {
		if errors.Is(err, dbfs.ErrFileNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
	}

	util.HttpJson(w, http.StatusOK, file)

}

// AddFavorite is a route to add a file to the favorites of a user
func (bc *BackendController) AddFavorite(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	fileId, ok := mux.Vars(r)["fileID"]
	if !ok {
		util.HttpError(w, http.StatusBadRequest, "No file_id specified")
		return
	} else if fileId == "" {
		util.HttpError(w, http.StatusBadRequest, "file_id is required")
		return
	}

	// get file

	file, err := dbfs.GetFileById(bc.Db, fileId, user)
	if err != nil {
		if errors.Is(err, dbfs.ErrFileNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
	}

	// add to favorites

	err = file.AddToFavorites(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, true)
}

// RemoveFavorite is a route to remove a file from the favorites of a user
func (bc *BackendController) RemoveFavorite(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	fileId, ok := mux.Vars(r)["fileID"]
	if !ok {
		util.HttpError(w, http.StatusBadRequest, "No file_id specified")
		return
	} else if fileId == "" {
		util.HttpError(w, http.StatusBadRequest, "file_id is required")
		return
	}

	// get file

	file, err := dbfs.GetFileById(bc.Db, fileId, user)
	if err != nil {
		if errors.Is(err, dbfs.ErrFileNotFound) {
			util.HttpError(w, http.StatusNotFound, err.Error())
		} else {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
		}
	}

	// remove from favorites
	err = file.RemoveFromFavorites(bc.Db, user)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, true)

}
