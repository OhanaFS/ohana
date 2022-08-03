package dbfs

import (
	"errors"
	"github.com/OhanaFS/ohana/util"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Key: CronJobDeleteFragmentsInProgress
// Value: 1 or 0
// If job is 1 and CronJobDeleteFragmentsLastStart is longer than 1 hour ago, then send a warning to the sysadmin.
// as it seems like the job is stuck.

// Key: CronJobDeleteFragmentsHandledServer
// Value (string): The server that handled the job last time.

// Key: CronJobDeleteFragmentsLastStart
// Value (string): Unix timestamp

// Key: CronJobDeleteFragmentsLastEnd
// Value (string): Unix timestamp
// If the timestamp is older than 1 hour, the job is ran. Unless manually ran.
const (
	CronJobDeleteFragmentsInProgress    = "CronJobDeleteFragmentsInProgress"
	CronJobDeleteFragmentsHandledServer = "CronJobDeleteFragmentsHandledServer"
	CronJobDeleteFragmentsLastStart     = "CronJobDeleteFragmentsLastStart"
	CronJobDeleteFragmentsLastEnd       = "CronJobDeleteFragmentsLastEnd"
	CronJobDeleteKeepVersionsFor        = "CronJobDeleteKeepVersionsFor"
)

var (
	ErrCronJobPropertyNotSet  = errors.New("cron job property not set")
	ErrInvalidCronJobProperty = errors.New("invalid cron job property")
)

func GetToBeDeletedFragments(tx *gorm.DB) ([]Fragment, error) {

	var fragments []Fragment
	var fileVersions []FileVersion

	// Get FileVersions that are marked as to be deleted
	err := tx.Where("status = ?", FileStatusToBeDeleted).Find(&fileVersions).Error
	if err != nil {
		return nil, err
	}

	if len(fileVersions) == 0 {
		return fragments, nil
	}

	// Get Fragments that are marked as to be deleted

	dataIdSeen := util.NewSet[string]()

	for _, fileVersion := range fileVersions {

		if dataIdSeen.Has(fileVersion.DataId) {
			continue
		}

		// Check if any of the file version's data fragments are still being used

		var copiesOfData int64
		err2 := tx.Model(&DataCopies{}).Where("data_id = ?", fileVersion.DataId).Count(&copiesOfData).Error
		if err2 != nil {
			return nil, err2
		}
		if copiesOfData >= 1 {
			// There are still data copies of this file
			// Thus, search in the DB to verify that another file is still using it

			var copiesOfDataUsing int64
			err2 = tx.Model(&FileVersion{}).Where("data_id = ? AND status = ?", fileVersion.DataId, FileStatusGood).
				Count(&copiesOfDataUsing).Error

			if copiesOfDataUsing >= 1 {
				continue
			} else {
				// There are no other files using this data fragment
				// Thus, we can delete it.
				err3 := tx.Delete(&DataCopies{}, "data_id = ?", fileVersion.DataId).Error
				if err3 != nil {
					return nil, err3
				}
			}

		}

		// There are no data copies of this file
		// Thus, we can get the fragments and append to fragment.
		var tempFragments []Fragment

		err3 := tx.Where("file_version_data_id = ?", fileVersion.DataId).Find(&tempFragments).Error
		if err3 != nil {
			if errors.Is(err3, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, err3
		}

		fragments = append(fragments, tempFragments...)

		dataIdSeen.Add(fileVersion.DataId)
	}

	return fragments, err
}

func FinishDeleteDataId(tx *gorm.DB, dataId string) error {

	// Delete all fragments with this dataId
	err := tx.Where("file_version_data_id = ?", dataId).Delete(&Fragment{}).Error
	if err != nil {
		return err
	}

	// In that case, deleting a folder should just go ahead and delete everything else as well.
	return tx.Model(&FileVersion{}).Where("data_id = ?", dataId).Update("status", FileStatusDeleted).Error
}

// MarkOldFileVersions goes through the database and see what file versions are old and can be deleted.
func MarkOldFileVersions(tx *gorm.DB) (int64, error) {

	// Get the timeframe to consider "old"

	var days KeyValueDBPair

	err := tx.First(&days, "key = ?", CronJobDeleteKeepVersionsFor).Error
	if err != nil {
		return 0, err
	}

	beforeDate := time.Now().AddDate(0, 0, -1*days.ValueInt)

	// Getting files that are currently in use
	filesCurrentlyUsed := tx.Select("data_id").Where("entry_type = ?", IsFile).Table("files")

	// Set date
	result := tx.Model(&FileVersion{}).Where(
		"status = ? AND modified_time < ? AND entry_type = ? AND data_id NOT IN (?)",
		FileStatusGood, beforeDate, IsFile, filesCurrentlyUsed).Update("status", FileStatusToBeDeleted)

	// This will only delete files. Will not delete folders since they don't take up disk space anyway so we don't care.

	if result.Error != nil {
		return 0, result.Error
	} else {
		return result.RowsAffected, nil
	}

}

func ClearFileStatusDeletedEntries(tx *gorm.DB) error {
	return tx.Model(&FileVersion{}).Where("status = ? AND modified_time < ?", FileStatusDeleted,
		time.Now()).Delete(&FileVersion{}).Error
}

// createCronJobKeyValues creates the parameters for the cron job.
func createCronJobKeyValues(tx *gorm.DB) error {
	keyValueCronJobs := []string{CronJobDeleteFragmentsInProgress,
		CronJobDeleteFragmentsHandledServer,
		CronJobDeleteFragmentsLastStart,
		CronJobDeleteFragmentsLastEnd,
		CronJobDeleteKeepVersionsFor,
	}

	return tx.Transaction(func(tx *gorm.DB) error {
		for _, keyValueCronJob := range keyValueCronJobs {
			var keyValue KeyValueDBPair
			err := tx.Where("key = ?", keyValueCronJob).First(&keyValue).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					keyValue = KeyValueDBPair{Key: keyValueCronJob}
					err = tx.Create(&keyValue).Error
					if err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}

		return nil
	})
}

// IsCronDeleteRunning checks if the cron job is running.
// Returns a string of the server currently running, the time it started (or last ended),
// and an error if DB encounters an issue.
func IsCronDeleteRunning(tx *gorm.DB) (string, time.Time, error) {
	var cronDeleteRunning KeyValueDBPair
	err := tx.Model(&KeyValueDBPair{}).Where("key = ?", CronJobDeleteFragmentsInProgress).
		First(&cronDeleteRunning).Error
	if err != nil {
		return "", time.Unix(0, 0), err
	}

	if cronDeleteRunning.ValueInt == 0 {

		// get the last time it finished
		var cronDeleteLastEnd KeyValueDBPair
		err = tx.Model(&KeyValueDBPair{}).Where("key = ?", CronJobDeleteFragmentsLastEnd).
			First(&cronDeleteLastEnd).Error
		if err != nil {
			return "", time.Unix(0, 0), err
		}

		// string unix timestamp to time.Time
		if cronDeleteLastEnd.ValueString == "" {
			cronDeleteLastEnd.ValueString = "0"
		}
		int64Timestamp, err := strconv.ParseInt(cronDeleteLastEnd.ValueString, 10, 64)
		if err != nil {
			return "", time.Unix(0, 0), err
		}
		return "", time.Unix(int64Timestamp, 0), nil

	} else {
		// get the server that is running it

		var cronDeleteRunningServer KeyValueDBPair
		err = tx.Model(&KeyValueDBPair{}).Where("key = ?", CronJobDeleteFragmentsHandledServer).
			First(&cronDeleteRunningServer).Error

		// get the last time it started
		var cronDeleteLastStart KeyValueDBPair
		err = tx.Model(&KeyValueDBPair{}).Where("key = ?", CronJobDeleteFragmentsLastStart).
			First(&cronDeleteLastStart).Error
		if err != nil {
			return "", time.Unix(0, 0), err
		}

		// string unix timestamp to time.Time
		if cronDeleteLastStart.ValueString == "" {
			cronDeleteLastStart.ValueString = "0"
		}
		int64Timestamp, err := strconv.ParseInt(cronDeleteLastStart.ValueString, 10, 64)
		if err != nil {
			return "", time.Unix(0, 0), err
		}

		return cronDeleteRunningServer.ValueString, time.Unix(int64Timestamp, 0), nil
	}

}

// SetHowLongToKeepFileVersions sets the number of days to keep file versions.
func SetHowLongToKeepFileVersions(tx *gorm.DB, days int) error {

	if days < 0 {
		return ErrInvalidCronJobProperty
	}

	return tx.Model(&KeyValueDBPair{}).Where("key = ?", CronJobDeleteKeepVersionsFor).
		Update("value_int", days).Error
}
