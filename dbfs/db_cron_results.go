package dbfs

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	CronErrorTypeInternalError = 1
	CronErrorTypeNotAvailable  = 2
	CronErrorTypeCorrupted     = 3
	JobStatusRunning           = 1
	JobStatusCompleteErrors    = 2
	JobStatusCompleteNoErrors  = 3
)

var (
	ErrorCronJobStillRunning = errors.New("cron job is still running")
	ErrorCronJobDoesNotExist = errors.New("cron job does not exist")
)

type Job struct {
	JobId                        uint `gorm:"primaryKey; not null"`
	StartTime                    time.Time
	EndTime                      time.Time
	TotalTimeTaken               time.Duration
	TotalShardsScanned           int
	TotalFilesScanned            int
	MissingShardsCheck           bool
	MissingShardsProgress        []JobProgressMissingShard `gorm:"foreignkey:JobId"`
	OrphanedShardsCheck          bool
	OrphanedShardsProgress       []JobProgressOrphanedShard `gorm:"foreignkey:JobId"`
	QuickShardsHealthCheck       bool
	QuickShardsHealthProgress    []JobProgressCFSHC `gorm:"foreignkey:JobId"`
	AllFilesShardsHealthCheck    bool
	AllFilesShardsHealthProgress []JobProgressAFSHC `gorm:"foreignkey:JobId"`
	PermissionCheck              bool
	PermissionResults            []JobProgressPermissionCheck `gorm:"foreignkey:JobId"`
	DeleteFragments              bool
	DeleteFragmentsResults       []JobProgressDeleteFragments `gorm:"foreignkey:JobId"`
	Progress                     int
	StatusMsg                    string
	Status                       int
}

// ResultsCFSHC Current files shards health check result
type ResultsCFSHC struct {
	JobId     uint
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
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsAFSHC All file shards health check result
type ResultsAFSHC struct {
	JobId     uint
	FileName  string
	FileId    string
	DataId    string
	FragPath  string
	ServerId  string
	Error     string
	ErrorType int
}

type FixAFSHC struct {
	DataId   string
	Fix      bool
	Delete   bool
	Password string
}

// JobProgressAFSHC All files fragment health check job progress
type JobProgressAFSHC struct {
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsMissingShard All Fragments health check result
type ResultsMissingShard struct {
	JobId     uint
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
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// ResultsOrphanedShard Orphaned shards result
type ResultsOrphanedShard struct {
	JobId    uint
	ServerId string
	FileName string
}

// JobProgressOrphanedShard Orphaned shards job progress
type JobProgressOrphanedShard struct {
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// JobProgressPermissionCheck reports the progress of the permission check
type JobProgressPermissionCheck struct {
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

// JobProgressDeleteFragments reports the progress of deleting fragments
type JobProgressDeleteFragments struct {
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	ServerId   string
	InProgress bool
	Msg        string
}

type JobParameters struct {
	MissingShardsCheck  bool
	OrphanedShardsCheck bool
	QuickShardsCheck    bool
	AllFilesShardsCheck bool
	PermissionCheck     bool
	DeleteFragments     bool
}

// GetAllJobs Returns all jobs in the database based on the paramters passed in
func GetAllJobs(tx *gorm.DB, startNum int, startDate, endDate time.Time, filter int) ([]Job, error) {

	// TODO Calculate the Progress
	var jobs []Job
	var err error
	if filter == 0 {
		err = tx.Where("start_time >= ? AND start_time <= ? ",
			startDate, endDate).
			Order("start_time desc").
			Offset(startNum).
			Limit(10).
			Preload(clause.Associations).
			Find(&jobs).Error
	} else {
		err = tx.Where("start_time >= ? AND start_time <= ? AND status = ?",
			startDate, endDate, filter).
			Order("start_time desc").
			Offset(startNum).
			Limit(10).
			Preload(clause.Associations).
			Find(&jobs).Error
	}

	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJob Returns a job by id
func GetJob(tx *gorm.DB, jobId int) (*Job, error) {
	var job Job
	err := tx.Where("job_id = ?", jobId).Preload(clause.Associations).
		First(&job).Preload(clause.Associations).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrorCronJobDoesNotExist
		}
		return nil, err
	}
	return &job, nil
}

// DeleteJob Deletes a job by id and all the associated results
func DeleteJob(tx *gorm.DB, jobId int) error {
	var job Job
	err := tx.Where("job_id = ?", jobId).Preload(clause.Associations).
		First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrorCronJobDoesNotExist
		}
		return err
	}

	// Getting associated results and deleting it.
	err = tx.Transaction(func(tx *gorm.DB) error {

		for _, result := range job.MissingShardsProgress {
			err := tx.Delete(&result).Error
			if err != nil {
				return err
			}
		}
		for _, result := range job.OrphanedShardsProgress {
			err := tx.Delete(&result).Error
			if err != nil {
				return err
			}
		}
		for _, result := range job.QuickShardsHealthProgress {
			err := tx.Delete(&result).Error
			if err != nil {
				return err
			}
		}
		for _, result := range job.AllFilesShardsHealthProgress {
			err := tx.Delete(&result).Error
			if err != nil {
				return err
			}
		}
		for _, result := range job.PermissionResults {
			err := tx.Delete(&result).Error
			if err != nil {
				return err
			}
		}
		for _, result := range job.DeleteFragmentsResults {
			err := tx.Delete(&result).Error
			if err != nil {
				return err
			}
		}

		if tx.Delete(&job).Error != nil {
			return err
		}

		return nil

	})

	return err
}

// InitializeJob creates a job based on the parameters given.
// Will return a job object which contains an ID that can be used to get the job progress.
// It does not communicate with Inc to start the job.
func InitializeJob(tx *gorm.DB, parameters JobParameters) (*Job, error) {

	if parameters.AllFilesShardsCheck {
		parameters.QuickShardsCheck = false
		parameters.MissingShardsCheck = false
	}

	// get an id
	job := &Job{
		StartTime:                 time.Now(),
		MissingShardsCheck:        parameters.MissingShardsCheck,
		OrphanedShardsCheck:       parameters.OrphanedShardsCheck,
		QuickShardsHealthCheck:    parameters.QuickShardsCheck,
		AllFilesShardsHealthCheck: parameters.AllFilesShardsCheck,
		PermissionCheck:           parameters.PermissionCheck,
		DeleteFragments:           parameters.DeleteFragments,
	}
	err := tx.Create(&job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}

func StartJob(tx *gorm.DB, job *Job) error {

	return nil
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
