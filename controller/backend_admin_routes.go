package controller

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"net/http"
)

// CronDeleteFragments clears fragments that are older than the configured amount of time and deleted fragments
func (bc *BackendController) CronDeleteFragments(w http.ResponseWriter, r *http.Request) {

	user, err := ctxutil.GetUser(r.Context())
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if user is admin
	if user.AccountType != dbfs.AccountTypeAdmin {
		util.HttpError(w, http.StatusForbidden, "You are not an admin")
		return
	}

	err, s := bc.Inc.CronJobDeleteFragments(true)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, s)
}
