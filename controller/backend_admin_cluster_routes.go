package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CronDeleteFragments clears fragments that are older than the configured amount of time and deleted fragments
func (bc *BackendController) CronDeleteFragments(w http.ResponseWriter, r *http.Request) {

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

	s, err := bc.Inc.CronJobDeleteShards(true)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, s)
}

// GetNumOfFiles returns the number of files currently in the database
func (bc *BackendController) GetNumOfFiles(w http.ResponseWriter, r *http.Request) {

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

	var numOfFiles int64

	bc.Db.Model(&dbfs.File{}).Where("entry_type = ?", dbfs.IsFile).Count(&numOfFiles)

	// success
	util.HttpJson(w, http.StatusOK, numOfFiles)
}

// GetNumOfFilesHistorical returns the historical history of the number of files in the database
// based on the time period specified
func (bc *BackendController) GetNumOfFilesHistorical(w http.ResponseWriter, r *http.Request) {

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

	timePeriod := r.Header.Get("range_type")
	startDate := r.Header.Get("start_date")
	endDate := r.Header.Get("end_date")

	// Process startDate and endDate

	if timePeriod == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing time period")
		return
	}

	// Converting to int the type timePeriod
	timePeriodInt, err := strconv.Atoi(timePeriod)
	if err != nil {
		util.HttpError(w, http.StatusBadRequest, "Invalid time period. "+
			"Expect 1 for day, 2 for week, 3 for month")
		return
	}

	// Check if the time period is valid
	if timePeriodInt < dbfs.HistoricalRangeDay || timePeriodInt > dbfs.HistoricalRangeMonth {
		util.HttpError(w, http.StatusBadRequest, "Invalid time period"+
			"Expect 1 for day, 2 for week, 3 for month")
		return
	}

	var NumOfFiles []dbfs.HistoricalStats
	var startDateTime, endDateTime time.Time

	switch timePeriodInt {
	case dbfs.HistoricalRangeDay:
		{
			// If there is no time period specified, return the last 10 days
			if startDate == "" {
				startDateTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
					0, 0, 0, 0, time.UTC)
				startDateTime = startDateTime.AddDate(0, 0, -9)
			} else {
				// truncate to just get date part
				startDateTime, err = time.Parse("2006-01-02", startDate[:10])
			}
			if endDate == "" {
				endDateTime = time.Now()
			} else {
				endDateTime, err = time.Parse("2006-01-02", endDate[:10])
			}

			// Get the number of files for each day in the week
			for startDateTime.Before(endDateTime) || startDateTime.Equal(endDateTime) {
				// manually loop
				var tempFile dbfs.HistoricalStats
				err := bc.Db.Model(&dbfs.HistoricalStats{}).Where("day = ? AND month = ? and year = ?",
					startDateTime.Day(), int(startDateTime.Month()), startDateTime.Year()).Find(&tempFile).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// the record isn't found so continue
						continue
					}
					util.HttpError(w, http.StatusInternalServerError, "Couldn't extract data: "+err.Error())
					return
				}

				NumOfFiles = append(NumOfFiles, tempFile)

				startDateTime = startDateTime.AddDate(0, 0, 1)
			}

		}
	case dbfs.HistoricalRangeWeek:
		{
			// for week, we'll just grab each sunday's data, else we'll grab the last day
			if endDate == "" {
				endDateTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
					0, 0, 0, 0, time.UTC)
				if endDateTime.Weekday() != time.Sunday {
					endDateTime = endDateTime.AddDate(0, 0, -int(endDateTime.Weekday()))
				}
			} else {
				endDateTime, err = time.Parse("2006-01-02", endDate[:10])

				if endDateTime.Weekday() != time.Sunday {
					endDateTime = endDateTime.AddDate(0, 0, -int(endDateTime.Weekday()))
				}
			}
			if startDate == "" {
				startDateTime = time.Date(endDateTime.Year(), endDateTime.Month(), endDateTime.Day(),
					0, 0, 0, 0, time.UTC)
				startDateTime = startDateTime.AddDate(0, 0, -(7 * 9))
			} else {
				startDateTime, err = time.Parse("2006-01-02", startDate[:10])
				if startDateTime.Weekday() != time.Sunday {
					startDateTime = startDateTime.AddDate(0, 0, -int(startDateTime.Weekday()))
				}
			}

			// Get the number of files for each day in the week
			for startDateTime.Before(endDateTime) || startDateTime.Equal(endDateTime) {
				// manually loop
				var tempFile dbfs.HistoricalStats
				err := bc.Db.Model(&dbfs.HistoricalStats{}).Where("day = ? AND month = ? and year = ?",
					startDateTime.Day(), int(startDateTime.Month()), startDateTime.Year()).Find(&tempFile).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// the record isn't found so continue
						continue
					}
					util.HttpError(w, http.StatusInternalServerError, "Couldn't extract data: "+err.Error())
					return
				}

				NumOfFiles = append(NumOfFiles, tempFile)

				startDateTime = startDateTime.AddDate(0, 0, 7)
			}

		}
	case dbfs.HistoricalRangeMonth:
		{
			// for month, we'll just grab each first day of the month
			if endDate == "" {
				endDateTime = time.Date(time.Now().Year(), time.Now().Month(), 1,
					0, 0, 0, 0, time.UTC)
			} else {
				endDateTime, err = time.Parse("2006-01-02", endDate[:10])
				endDateTime = time.Date(endDateTime.Year(), endDateTime.Month(), 1,
					0, 0, 0, 0, time.UTC)
			}
			if startDate == "" {
				startDateTime = time.Date(endDateTime.Year(), endDateTime.Month(), 1,
					0, 0, 0, 0, time.UTC)
				startDateTime = startDateTime.AddDate(0, -9, 0)
			} else {
				startDateTime, err = time.Parse("2006-01-02", startDate[:10])
				startDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), 1,
					0, 0, 0, 0, time.UTC)
			}

			// Get the number of files for each day in the week
			for startDateTime.Before(endDateTime) || startDateTime.Equal(endDateTime) {
				// manually loop
				var tempFile dbfs.HistoricalStats
				err := bc.Db.Model(&dbfs.HistoricalStats{}).Where("day = ? AND month = ? and year = ?",
					startDateTime.Day(), int(startDateTime.Month()), startDateTime.Year()).Find(&tempFile).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// the record isn't found so continue
						continue
					}
					util.HttpError(w, http.StatusInternalServerError, "Couldn't extract data: "+err.Error())
					return
				}

				NumOfFiles = append(NumOfFiles, tempFile)

				startDateTime = startDateTime.AddDate(0, 1, 0)
			}

		}
	}

	util.HttpJson(w, http.StatusOK, NumOfFiles)
}

// GetStorageUsed returns the amount of storage used (not including replica, versioning) in bytes
func (bc *BackendController) GetStorageUsed(w http.ResponseWriter, r *http.Request) {

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

	var storageUsed int64

	err = bc.Db.Model(&dbfs.File{}).Select("sum(size)").Where("entry_type = ?", dbfs.IsFile).
		Take(&storageUsed).Error
	if err != nil {
		storageUsed = 0
	}

	// success
	util.HttpJson(w, http.StatusOK, storageUsed)

}

// GetStorageUsedReplica returns the amount of storage used (including replica, versioning) in bytes
func (bc *BackendController) GetStorageUsedReplica(w http.ResponseWriter, r *http.Request) {

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

	var storageUsedReplica int64

	err = bc.Db.Model(&dbfs.FileVersion{}).Select("sum(actual_size)").
		Where("entry_type = ? AND status <> ?", dbfs.IsFile, dbfs.FileStatusDeleted).Take(&storageUsedReplica).Error
	if err != nil {
		storageUsedReplica = 0
	}
	// success
	util.HttpJson(w, http.StatusOK, storageUsedReplica)

}

// GetAllAlerts returns all alerts
func (bc *BackendController) GetAllAlerts(w http.ResponseWriter, r *http.Request) {

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

	var alerts []dbfs.Alert
	err = bc.Db.Find(&alerts).Error

	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, alerts)
}

// ClearAllAlerts clears all alerts
func (bc *BackendController) ClearAllAlerts(w http.ResponseWriter, r *http.Request) {

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

	err = bc.Db.Where("1=1").Delete(&dbfs.Alert{}).Error
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, true)
}

// GetAlert returns the alert based on the ID provided.
func (bc *BackendController) GetAlert(w http.ResponseWriter, r *http.Request) {

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
	alertId := vars["id"]

	if alertId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing alert id")
		return
	}

	var alert dbfs.Alert
	err = bc.Db.Where("log_id = ?", alertId).First(&alert).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			util.HttpError(w, http.StatusNotFound, "Alert not found")
			return
		}
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, alert)
}

// ClearAlert deletes the alert based on the ID provided.
func (bc *BackendController) ClearAlert(w http.ResponseWriter, r *http.Request) {

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
	alertId := vars["id"]

	if alertId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing alert id")
		return
	}

	err = bc.Db.Where("log_id = ?", alertId).Delete(&dbfs.Alert{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			util.HttpError(w, http.StatusNotFound, "Alert not found")
			return
		}
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, true)
}

// GetAllLogs returns the 50 newest logs based on the parameters provided in the header
// start_num (if not will get latest)
// start_date
// end_date
// server_filter
// type_filter
func (bc *BackendController) GetAllLogs(w http.ResponseWriter, r *http.Request) {

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

	// header vars
	startNumString := r.Header.Get("start_num")
	startDateString := r.Header.Get("start_date")
	endDateString := r.Header.Get("end_date")
	serverFilterString := r.Header.Get("server_filter")
	typeFilterString := r.Header.Get("type_filter")

	// parse vars
	// We will build upon a string and an array of interfaces and append to it as we parse the vars

	parseStringArray := make([]string, 0)
	parseObjectsArray := make([]interface{}, 0)

	// convert to int
	startNum, err := strconv.ParseInt(startNumString, 10, 64)
	if err != nil {
		startNum = 0
	}

	// check if start date is valid
	if startDateString != "" {
		startDate, err := time.Parse(time.RFC3339, startDateString)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, "Invalid startDateString. Follow RFC3339 2006-01-02T15:04:05Z07:00")
			return
		}

		parseStringArray = append(parseStringArray, "time_stamp >= ?")
		parseObjectsArray = append(parseObjectsArray, startDate)

	}
	if endDateString != "" {
		endDate, err := time.Parse(time.RFC3339, endDateString)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, "Invalid endDateString. Follow YYYY-MM-DD")
			return
		}

		parseStringArray = append(parseStringArray, "time_stamp <= ?")
		parseObjectsArray = append(parseObjectsArray, endDate)
	}

	if serverFilterString != "" {
		parseStringArray = append(parseStringArray, "server_name LIKE ?")
		parseObjectsArray = append(parseObjectsArray, "%"+serverFilterString+"%")
	}
	if typeFilterString != "" {
		typeFilter, err := strconv.ParseInt(typeFilterString, 10, 64)
		if err != nil {
			util.HttpError(w, http.StatusBadRequest, "Invalid typeFilterString")
			return
		}

		parseStringArray = append(parseStringArray, "log_type = ?")
		parseObjectsArray = append(parseObjectsArray, typeFilter)
	}

	var logs []dbfs.Log

	// build query

	err = bc.Db.Where(strings.Join(parseStringArray[:], " AND "),
		parseObjectsArray...).Order("log_id desc").Offset(int(startNum * 50)).Limit(50).Find(&logs).Error

	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, logs)
}

// GetLog returns the log based on the ID provided.
func (bc *BackendController) GetLog(w http.ResponseWriter, r *http.Request) {

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
	logId := vars["id"]

	if logId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing log id")
		return
	}

	var log dbfs.Log
	err = bc.Db.Where("log_id = ?", logId).First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			util.HttpError(w, http.StatusNotFound, "Log not found")
			return
		}
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, log)
}

// ClearLog clears the log based on the ID provided.
func (bc *BackendController) ClearLog(w http.ResponseWriter, r *http.Request) {

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
	logId := vars["id"]

	if logId == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing log id")
		return
	}

	err = bc.Db.Where("log_id = ?", logId).Delete(&dbfs.Log{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			util.HttpError(w, http.StatusNotFound, "Log not found")
			return
		}
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, true)
}

// ClearAllLogs clears all logs
func (bc *BackendController) ClearAllLogs(w http.ResponseWriter, r *http.Request) {

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

	err = bc.Db.Where("1=1").Delete(&dbfs.Log{}).Error
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// success
	util.HttpJson(w, http.StatusOK, true)
}

// GetServerStatuses Get All Server Status. May not be the most recent as it's from the cache.
func (bc *BackendController) GetServerStatuses(w http.ResponseWriter, r *http.Request) {

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

	servers, err := dbfs.GetServers(bc.Db)

	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, servers)
}

// GetSpecificServerStatus Get a specific server status. Will ping the server for the latest info
func (bc *BackendController) GetSpecificServerStatus(w http.ResponseWriter, r *http.Request) {

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
	serverName := vars["serverName"]

	if serverName == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing server name")
		return
	}

	serverReport, err := bc.Inc.GetServerStatusReport(serverName)
	if err != nil {
		util.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpJson(w, http.StatusOK, serverReport)
}

// DeleteServer de-registers a server, deleting it from the database. Will cause shards to be marked as offline.
// Should only be used when a server is stuck registering. This will not delete the data on the server.
func (bc *BackendController) DeleteServer(w http.ResponseWriter, r *http.Request) {

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
	serverName := vars["serverName"]

	if serverName == "" {
		util.HttpError(w, http.StatusBadRequest, "Missing server name")
		return
	} else if serverName == bc.ServerName {
		// de-register server, and shutdown gracefully
		err := dbfs.MarkServerOffline(bc.Db, serverName)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		bc.Inc.Shutdown <- true

	} else {
		// check if server exists, then send shutdown signal. If server doesn't exist, return error.

		hostname, err := dbfs.GetServerAddress(bc.Db, serverName)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("Marking server offline: " + serverName + " at " + hostname)

		//// timeout of 10 seconds. If it doesn't respond, then it's dead.

		err = dbfs.MarkServerOffline(bc.Db, serverName)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}
	// success
	util.HttpJson(w, http.StatusOK, true)
}
