package dbfs

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

const (
	CronErrorTypeInternalError = 1
	CronErrorTypeNotAvailable  = 2
	CronErrorTypeCorrupted     = 3
)

var (
	ErrorCronJobStillRunning = errors.New("cron job is still running")
	ErrorCronJobDoesNotExist = errors.New("cron job does not exist")
)

// ResultsCffhc Result Current Files Fragment Health Check result
type ResultsCffhc struct {
	JobId     int
	FileName  string
	FileId    string
	DataId    string
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
	Msg        string
}

// ResultsAffhc Result All Fragments Health Check result
type ResultsAffhc struct {
	JobId     int
	FileName  string
	FileId    string
	DataId    string
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
	Msg        string
}

// ResultsMissingShard Result All Fragments Health Check result
type ResultsMissingShard struct {
	JobId     int
	FileName  string
	FileId    string
	DataId    string
	FragPath  string
	ServerId  string
	Error     string
	ErrorType int
}

// JobProgressMissingShard All Files Fragment Health Check job progress
type JobProgressMissingShard struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsOrphanedShard Result All Fragments Health Check result
type ResultsOrphanedShard struct {
	JobId    int
	ServerId string
	FileName string
}

// JobProgressOrphanedShard All Files Fragment Health Check job progress
type JobProgressOrphanedShard struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

func GetResultsCffhc(tx *gorm.DB, jobId int) ([]ResultsCffhc, error) {

	// Check if the job is still running or exists
	var job JobprogressCffhc
	err := tx.Where("job_id = ?", jobId).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrorCronJobDoesNotExist
		}
	} else {
		if job.InProgress {
			return nil, ErrorCronJobStillRunning
		}
	}

	var results []ResultsCffhc
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}

func GetResultsAffhc(tx *gorm.DB, jobId int) ([]ResultsAffhc, error) {

	// Check if the job is still running or exists
	var job JobprogressAffhc
	err := tx.Where("job_id = ?", jobId).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrorCronJobDoesNotExist
		}
	} else {
		if job.InProgress {
			return nil, ErrorCronJobStillRunning
		}
	}

	var results []ResultsAffhc
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}

func GetResultsMissingShard(tx *gorm.DB, jobId int) ([]ResultsMissingShard, error) {

	var job JobProgressMissingShard
	err := tx.Where("job_id = ?", jobId).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrorCronJobDoesNotExist
		}
	} else {
		if job.InProgress {
			return nil, ErrorCronJobStillRunning
		}
	}

	var results []ResultsMissingShard
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}

func GetResultsOrphanedShard(tx *gorm.DB, jobId int) ([]ResultsOrphanedShard, error) {

	// Check if the job is still running or exists
	var job JobProgressOrphanedShard
	err := tx.Where("job_id = ?", jobId).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrorCronJobDoesNotExist
		}
	} else {
		if job.InProgress {
			return nil, ErrorCronJobStillRunning
		}
	}

	var results []ResultsOrphanedShard
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}
