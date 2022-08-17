package dbfs

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"time"
)

const (
	CronErrorTypeInternalError = 1
	CronErrorTypeNotAvailable  = 2
	CronErrorTypeCorrupted     = 3
	CronErrorTypeSolved        = 4 // The error was solved
	JobStatusRunning           = 1
	JobNoErrors                = 2
	JobHasErrors               = 3
	CronErrorTypeMissingFile   = 1
)

var (
	ErrorCronJobStillRunning = errors.New("cron job is still running")
	ErrorCronJobDoesNotExist = errors.New("cron job does not exist")
)

type Job struct {
	JobId              uint `gorm:"primaryKey; not null"`
	StartTime          time.Time
	EndTime            time.Time
	TotalTimeTaken     time.Duration
	TotalShardsScanned int
	TotalFilesScanned  int
	// MissingShardsCheck has a weightage of 10 in the progress calculation
	MissingShardsCheck    bool
	MissingShardsProgress []JobProgressMissingShard `gorm:"foreignkey:JobId"`
	// OrphanedShardsCheck Check has a weightage of 10 in the progress calculation
	OrphanedShardsCheck    bool
	OrphanedShardsProgress []JobProgressOrphanedShard `gorm:"foreignkey:JobId"`
	// QuickShardsCheck Check has a weightage of 50 in the progress calculation
	QuickShardsHealthCheck    bool
	QuickShardsHealthProgress []JobProgressCFSHC `gorm:"foreignkey:JobId"`
	// AllFilesShardsHealthCheck Check has a weightage of 100 in the progress calculation
	AllFilesShardsHealthCheck    bool
	AllFilesShardsHealthProgress []JobProgressAFSHC `gorm:"foreignkey:JobId"`
	// PermissionCheck Check has a weightage of 20 in the progress calculation
	PermissionCheck   bool
	PermissionResults *JobProgressPermissionCheck `gorm:"foreignkey:JobId"`
	// DeleteFragments Check has a weightage of 10 in the progress calculation
	DeleteFragments        bool
	DeleteFragmentsResults []JobProgressDeleteFragments `gorm:"foreignkey:JobId"`
	// OrphanedFilesCheck has a weightage of 20 in the progress calculation
	OrphanedFilesCheck   bool
	OrphanedFilesResults []JobProgressOrphanedFile `gorm:"foreignkey:JobId"`
	// Progress is the percentage of the job that is complete.
	Progress  int
	StatusMsg string
	Status    int
	// Status is computed on the fly based on the progress of the job.
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

type ShardActions struct {
	DataId   string
	Fix      bool
	Delete   bool
	Password string
}

type OrphanedShardActions struct {
	ServerId string
	FileName string
	Delete   bool
}

type OrphanedFilesActions struct {
	ParentFolderId string
	Delete         bool
	Move           bool
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
	JobId     uint
	ServerId  string
	FileName  string
	Error     string
	ErrorType int
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
	OrphanedFilesCheck  bool
}

// ResultsOrphanedFile Orphaned File result
type ResultsOrphanedFile struct {
	JobId          uint   `gorm:"primary_key"`
	ParentFolderId string `gorm:"primary_key"`
	Contents       string
	Error          string
	// Error will store the path route it took to get to the error
	ErrorType int
	// ErrorType will store what happened with the file (orphaned, moved, deleted)
}

// JobProgressOrphanedFile Orphaned File job progress
type JobProgressOrphanedFile struct {
	JobId      uint `gorm:"primary_key"`
	StartTime  time.Time
	EndTime    time.Time
	Processed  int64
	Total      int64
	InProgress bool
	Msg        string
}

// GetAllJobs Returns all jobs in the database based on the paramters passed in
func GetAllJobs(tx *gorm.DB, startNum int, startDate, endDate time.Time, filter int) ([]Job, error) {

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

	// get number of servers working on it
	servers, err := GetServers(tx)
	if err != nil {
		return nil, fmt.Errorf("Error getting servers: %v", err)
	}

	for _, job := range jobs {
		job.Progress = calculateProgress(job, len(servers))
		if job.Progress == 100 {
			// update job via gorm to set status to complete
			tx.Model(&job).Where("job_id = ?", job.JobId).Update("progress", 100)
		}
	}

	return jobs, nil
}

func calculateProgress(job Job, numOfServers int) int {

	// map of string to int
	progressMap := map[string]int{
		"MissingShardsCheck":        10,
		"OrphanedShardsCheck":       10,
		"QuickShardsHealthCheck":    50,
		"AllFilesShardsHealthCheck": 100,
		"PermissionCheck":           20,
		"DeleteFragments":           10,
		"OrphanedFilesCheck":        20,
	}
	var totalProgress int

	if job.MissingShardsCheck {
		totalProgress += progressMap["MissingShardsCheck"]
	}
	if job.OrphanedShardsCheck {
		totalProgress += progressMap["OrphanedShardsCheck"]
	}
	if job.QuickShardsHealthCheck {
		totalProgress += progressMap["QuickShardsHealthCheck"]
	}
	if job.AllFilesShardsHealthCheck {
		totalProgress += progressMap["AllFilesShardsHealthCheck"]
	}
	if job.PermissionCheck {
		totalProgress += progressMap["PermissionCheck"]
	}
	if job.DeleteFragments {
		totalProgress += progressMap["DeleteFragments"]
	}
	if job.OrphanedFilesCheck {
		totalProgress += progressMap["OrphanedFilesCheck"]
	}

	currentProgress := float64(0)
	if job.MissingShardsCheck {
		for _, progress := range job.MissingShardsProgress {
			if !progress.InProgress {
				currentProgress += float64(progressMap["MissingShardsCheck"]) / float64(numOfServers)
			}
		}
	}
	if job.OrphanedShardsCheck {
		for _, progress := range job.OrphanedShardsProgress {
			if !progress.InProgress {
				currentProgress += float64(progressMap["OrphanedShardsCheck"]) / float64(numOfServers)
			}
		}
	}
	if job.QuickShardsHealthCheck {
		for _, progress := range job.QuickShardsHealthProgress {
			if !progress.InProgress {
				currentProgress += float64(progressMap["QuickShardsHealthCheck"]) / float64(numOfServers)
			}
		}
	}
	if job.AllFilesShardsHealthCheck {
		for _, progress := range job.AllFilesShardsHealthProgress {
			if !progress.InProgress {
				currentProgress += float64(progressMap["AllFilesShardsHealthCheck"]) / float64(numOfServers)
			}
		}
	}
	if job.PermissionCheck {
		if job.PermissionResults != nil {
			if !job.PermissionResults.InProgress {
				currentProgress += float64(progressMap["PermissionCheck"])
			}
		}
	}
	if job.DeleteFragments {
		for _, progress := range job.DeleteFragmentsResults {
			if !progress.InProgress {
				currentProgress += float64(progressMap["DeleteFragments"]) / float64(numOfServers)
			}
		}
	}
	if job.OrphanedFilesCheck {
		for _, progress := range job.OrphanedFilesResults {
			if !progress.InProgress {
				currentProgress += float64(progressMap["OrphanedFilesCheck"]) / float64(numOfServers)
			}
		}
	}

	return int(math.Ceil(currentProgress / float64(totalProgress) * float64(100)))
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

	// Progress

	// get number of servers working on it
	servers, err := GetServers(tx)
	if err != nil {
		return nil, fmt.Errorf("Error getting servers: %v", err)
	}
	job.Progress = calculateProgress(job, len(servers))
	if job.Progress == 100 {
		// update job via gorm to set status to complete
		tx.Model(&job).Where("job_id = ?", job.JobId).Update("progress", 100)
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
		if job.PermissionCheck {
			err := tx.Delete(job.PermissionResults).Error
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
		OrphanedFilesCheck:        parameters.OrphanedFilesCheck,
	}
	err := tx.Create(&job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}

func StartJob(tx *gorm.DB, job *Job) error {

	// TODO: NOTE. If you are testing with this, sleep for at least 10 seconds.
	// Otherwise, sqlite3 and gorm will get locked.
	if job.OrphanedFilesCheck {
		go CheckOrphanedFiles(tx, int(job.JobId), true)

	}

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

// GetResultsOrphanedFile Returns the results of an orphaned file job based on jobId
// Will return error if the job doesn't exist or if the job is still running
func GetResultsOrphanedFile(tx *gorm.DB, jobId int) ([]ResultsOrphanedFile, error) {

	// Check if the job is still running or exists
	var job JobProgressOrphanedFile
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

	var results []ResultsOrphanedFile
	err = tx.Where("job_id = ?", jobId).Find(&results).Error
	return results, err
}
