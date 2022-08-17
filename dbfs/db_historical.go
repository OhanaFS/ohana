package dbfs

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	HistoricalRangeDay       = 1
	HistoricalRangeWeek      = 2
	HistoricalRangeMonth     = 3
	HistoricalNumOfFiles     = 1
	HistoricalNonReplicaUsed = 2
	HistoricalReplicaUsed    = 3
)

var (
	ErrorMissingTimePeriod = errors.New("missing time period")
	ErrorInvalidTimePeriod = errors.New("invalid time period. Expect 1 for day, 2 for week, 3 for month")
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

func (di DateInt64Value) Year() int {
	decodedTime, err := time.Parse("2006-01-02", di.Date)
	if err != nil {
		return 0
	}
	return decodedTime.Year()
}

func (di DateInt64Value) Month() int {
	decodedTime, err := time.Parse("2006-01-02", di.Date)
	if err != nil {
		return 0
	}
	return int(decodedTime.Month())
}

func (di DateInt64Value) Day() int {
	decodedTime, err := time.Parse("2006-01-02", di.Date)
	if err != nil {
		return 0
	}
	return decodedTime.Day()
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

func GetHistoricalData(db *gorm.DB, timePeriod, startDate, endDate string, resultType int) ([]DateInt64Value, error) {

	if timePeriod == "" {
		return nil, ErrorMissingTimePeriod
	}

	// Converting to int the type timePeriod
	timePeriodInt, err := strconv.Atoi(timePeriod)
	if err != nil {
		return nil, ErrorInvalidTimePeriod

	}

	// Check if the time period is valid
	if timePeriodInt < HistoricalRangeDay || timePeriodInt > HistoricalRangeMonth {
		return nil, ErrorInvalidTimePeriod
	}

	var rawStats []HistoricalStats
	var startDateTime, endDateTime time.Time

	switch timePeriodInt {
	case HistoricalRangeDay:
		{
			// If there is no time period specified, return the last 10 days
			if startDate == "" {
				startDateTime = time.Now()
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
				var tempFile HistoricalStats

				err := db.Model(&HistoricalStats{}).Where("day = ? AND month = ? and year = ?",
					startDateTime.Day(), int(startDateTime.Month()), startDateTime.Year()).Find(&tempFile).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// the record isn't found so continue
						continue
					}

					return nil, fmt.Errorf("Couldn't extract data: " + err.Error())

				}

				rawStats = append(rawStats, tempFile)

				startDateTime = startDateTime.AddDate(0, 0, 1)
			}

		}
	case HistoricalRangeWeek:
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
				var tempFile HistoricalStats
				err := db.Model(&HistoricalStats{}).Where("day = ? AND month = ? and year = ?",
					startDateTime.Day(), int(startDateTime.Month()), startDateTime.Year()).Find(&tempFile).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// the record isn't found so continue
						continue
					}

					return nil, fmt.Errorf("Couldn't extract data: " + err.Error())
				}

				rawStats = append(rawStats, tempFile)

				startDateTime = startDateTime.AddDate(0, 0, 7)
			}

		}
	case HistoricalRangeMonth:
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
				var tempFile HistoricalStats
				err := db.Model(&HistoricalStats{}).Where("day = ? AND month = ? and year = ?",
					startDateTime.Day(), int(startDateTime.Month()), startDateTime.Year()).Find(&tempFile).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// the record isn't found so continue
						continue
					}
					return nil, fmt.Errorf("Couldn't extract data: " + err.Error())

				}

				rawStats = append(rawStats, tempFile)

				startDateTime = startDateTime.AddDate(0, 1, 0)
			}

		}
	}

	returnStats := make([]DateInt64Value, len(rawStats))

	for i, rawStat := range rawStats {

		switch resultType {
		case HistoricalNumOfFiles:
			{
				returnStats[i] = DateInt64Value{
					Date:  fmt.Sprintf("%04d-%02d-%02d", rawStat.Year, rawStat.Month, rawStat.Day),
					Value: rawStat.NumOfFiles,
				}
			}
		case HistoricalNonReplicaUsed:
			{
				returnStats[i] = DateInt64Value{
					Date:  fmt.Sprintf("%04d-%02d-%02d", rawStat.Year, rawStat.Month, rawStat.Day),
					Value: rawStat.NonReplicaUsed,
				}
			}
		case HistoricalReplicaUsed:
			{
				returnStats[i] = DateInt64Value{
					Date:  fmt.Sprintf("%04d-%02d-%02d", rawStat.Year, rawStat.Month, rawStat.Day),
					Value: rawStat.ReplicaUsed,
				}
			}
		}
	}

	return returnStats, nil

}
