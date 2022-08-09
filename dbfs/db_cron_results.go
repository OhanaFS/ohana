package dbfs

import (
	"gorm.io/gorm"
	"time"
)

const (
	CronErrorTypeInternalError = 1
	CronErrorTypeNotAvailable  = 2
	CronErrorTypeCorrupted     = 3
)

// ResultsCffhc Result Current Files Fragment Health Check result
type ResultsCffhc struct {
	JobId     int
	FileName  string
	FileId    string
	FragPath  string
	ServerId  string
	Error     string
	ErrorType int
}

// JobprogressCffhc Current Files Fragment Health Check job progress
type JobprogressCffhc struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
}

// ResultsAffhc Result All Fragments Health Check result
type ResultsAffhc struct {
	JobId     int
	FileName  string
	FileId    string
	FragPath  string
	ServerId  string
	Error     string
	ErrorType int
}

// JobprogressAffhc All Files Fragment Health Check job progress
type JobprogressAffhc struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
}

func GetResultsCffhc(tx *gorm.DB, jobId int) ([]ResultsCffhc, error) {
	var results []ResultsCffhc
	err := tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}

func GetResultsAffhc(tx *gorm.DB, jobId int) ([]ResultsAffhc, error) {
	var results []ResultsAffhc
	err := tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}
