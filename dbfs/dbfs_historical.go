package dbfs

import (
	"gorm.io/gorm"
	"time"
)

const (
	HistoricalRangeDay   = 1
	HistoricalRangeWeek  = 2
	HistoricalRangeMonth = 3
)

type HistoricalStats struct {
	Day            int
	Month          int
	Year           int
	NonReplicaUsed int64
	ReplicaUsed    int64
	NumOfFiles     int64
}

type DateInt64Value struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// DumpDailyStats calculates statistics such as NonReplicaUsed Data
// and more in HistoricalStats and dump it into it.
func DumpDailyStats(db *gorm.DB) error {

	// Get stats

	// NonReplicaUsed

	var storageUsedReplica int64

	if db.Model(FileVersion{}).Select("sum(actual_size)").
		Where("entry_type = ? AND status <> ?", IsFile, FileStatusDeleted).
		Take(&storageUsedReplica).Error != nil {
		storageUsedReplica = 0
	}

	// ReplicaUsed

	var storageUsedNonReplica int64

	if db.Model(File{}).Select("sum(size)").Where("entry_type = ?", IsFile).
		Take(&storageUsedNonReplica).Error != nil {
		storageUsedNonReplica = 0
	}

	// NumOfFiles

	var numOfFiles int64

	_ = db.Model(File{}).Where("entry_type = ?", IsFile).Count(&numOfFiles).Error

	var existingHistoricalStats HistoricalStats

	if err := db.Where("day = ? AND month = ? and year = ?",
		time.Now().Day(), int(time.Now().Month()), time.Now().Year()).
		First(&existingHistoricalStats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new stat
			existingHistoricalStats = HistoricalStats{
				Day:            time.Now().Day(),
				Month:          int(time.Now().Month()),
				Year:           time.Now().Year(),
				NonReplicaUsed: storageUsedNonReplica,
				ReplicaUsed:    storageUsedReplica,
				NumOfFiles:     numOfFiles,
			}
			return db.Create(&existingHistoricalStats).Error
		}
	} else {
		return err
	}

	// already exists, do an update instead

	return db.Model(&existingHistoricalStats).Updates(map[string]interface{}{
		"NonReplicaUsed": storageUsedNonReplica,
		"ReplicaUsed":    storageUsedReplica,
		"NumOfFiles":     numOfFiles,
	}).Error

}

func GetDayStat(db *gorm.DB, day int, month int, year int) (HistoricalStats, error) {
	var stats HistoricalStats
	if err := db.Where("day = ? AND month = ? and year = ?", day, month, year).First(&stats).Error; err != nil {
		return stats, err
	}
	return stats, nil
}

func GetTodayStat(db *gorm.DB) (HistoricalStats, error) {
	return GetDayStat(db, time.Now().Day(), int(time.Now().Month()), time.Now().Year())
}
