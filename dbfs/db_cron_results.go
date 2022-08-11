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

// ResultsCFSHC Current files shards health check result
type ResultsCFSHC struct {
	JobId     int
	FileName  string
	FileId    string
	DataId    string
	FragPath  string
	ServerId  string
	Error     string
	ErrorType int
}

// JobProgressCFSHC Current files shards health check job progress
type JobProgressCFSHC struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsAFSHC All file shards health check result
type ResultsAFSHC struct {
	JobId     int
	FileName  string
	FileId    string
	DataId    string
	FragPath  string
	ServerId  string
	Error     string
	ErrorType int
}

// JobProgressAFSHC All files fragment health check job progress
type JobProgressAFSHC struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsMissingShard All Fragments health check result
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

// JobProgressMissingShard Missing shards job progress
type JobProgressMissingShard struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsOrphanedShard Orphaned shards result
type ResultsOrphanedShard struct {
	JobId    int
	ServerId string
	FileName string
}

// JobProgressOrphanedShard Orphaned shards job progress
type JobProgressOrphanedShard struct {
	JobId      int
	StartTime  time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// GetResultsCFSHC Returns the results of a Current Files Shards Health Check job based on jobId
// Will return error if the job doesn't exist or if the job is still running
func GetResultsCFSHC(tx *gorm.DB, jobId int) ([]ResultsCFSHC, error) {

	// Check if the job is still running or exists
	var job JobProgressCFSHC
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

	var results []ResultsCFSHC
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}

// GetResultsAFSHC Returns the results of an All Files Shards Health Check job based on jobId
// Will return error if the job doesn't exist or if the job is still running
func GetResultsAFSHC(tx *gorm.DB, jobId int) ([]ResultsAFSHC, error) {

	// Check if the job is still running or exists
	var job JobProgressAFSHC
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

	var results []ResultsAFSHC
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}

// GetResultsMissingShard Returns the results of a missing shard job based on jobId
// Will return error if the job doesn't exist or if the job is still running
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

// GetResultsOrphanedShard Returns the results of an orphaned shard job based on jobId
// Will return error if the job doesn't exist or if the job is still running
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
