package inc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/stitch"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type DeleteWorkerStatus struct {
	DataId       string
	Error        error
	ServerErrors []DeleteWorkerServerStatus
}

type DeleteWorkerServerStatus struct {
	server string
	err    error
}

// DeleteFragmentsByPath deletes the path of a shard file
func (i Inc) DeleteFragmentsByPath(pathOfShard string) error {
	return os.Remove(path.Join(i.ShardsLocation, pathOfShard))
}

// CronJobDeleteFragments scans the DB and nodes for any fragments that should be deleted
// to clear up space.
func (i Inc) CronJobDeleteFragments(manualRun bool) (string, error) {

	// Checks if the job is currently being done.

	// The job should be only handled by the server with the least free space LMAO.
	// Thus, we register again to get the latest info.

	err := i.RegisterServer(false)
	if err != nil {
		return "register server error", err
	}

	// Check if the job is currently running.

	server, timestamp, err := dbfs.IsCronDeleteRunning(i.Db)
	if err != nil {
		return "", err
	}
	if server != "" {

		errorMsg := "cron job is already running by server " + server
		errorMsg += "\nlast started at " + timestamp.String()

		// Give warning if it seems stuck
		if timestamp.Add(time.Hour).Before(time.Now()) {
			errorMsg += "\nWARNING: last started more than an hour ago"
			return errorMsg, JobCurrentlyRunningWarning
		}

		return errorMsg, JobCurrentlyRunning
	}
	// else it should not be running

	// Check if the server should handle it (least amount of data free)
	servers, err := dbfs.GetServers(i.Db)
	if err != nil {
		return "", err
	}

	if !manualRun {
		// Get the server with the least amount of data free
		var serverWithLeastFreeSpace string
		var leastFreeSpace uint64
		for _, server := range servers {
			freeSpace := server.FreeSpace
			if freeSpace < leastFreeSpace || leastFreeSpace == uint64(0) {
				leastFreeSpace = freeSpace
				serverWithLeastFreeSpace = server.Name
			}
		}

		if serverWithLeastFreeSpace != i.ServerName {
			return "", errors.New("Assigning server " + serverWithLeastFreeSpace + " to process it")
			// TODO
		}
	}

	// Start the job
	_, err = dbfs.MarkOldFileVersions(i.Db)
	if err != nil {
		return "", err
	}

	fragments, err := dbfs.GetToBeDeletedFragments(i.Db)
	if err != nil {
		return "", err
	}

	if len(fragments) == 0 {
		return "Finished. Deleted 0 fragments", nil
	}

	dataIdFragmentMap := make(map[string][]dbfs.Fragment)

	for _, fragment := range fragments {
		dataIdFragmentMap[fragment.FileVersionDataId] = append(dataIdFragmentMap[fragment.FileVersionDataId], fragment)
	}

	// "Deletion Code"

	const maxGoroutines = 10
	input := make(chan string, len(dataIdFragmentMap))
	output := make(chan DeleteWorkerStatus, len(dataIdFragmentMap))

	// Worker function

	for num := 0; num < maxGoroutines; num++ {
		go i.deleteWorker(dataIdFragmentMap, input, output, i.Db)
	}

	for dataId, _ := range dataIdFragmentMap {
		input <- dataId
	}
	close(input)

	AllDeleteWorkerErrors := make([]DeleteWorkerStatus, 0)

	for j := 0; j < len(dataIdFragmentMap); j++ {

		status := <-output
		if status.Error != nil {
			AllDeleteWorkerErrors = append(AllDeleteWorkerErrors, status)
		}

	}

	// If there are errors, return them
	if len(AllDeleteWorkerErrors) > 0 {
		return fmt.Sprintf("Finished (WARNING). Deleted %d, but %d errors occurred",
			len(dataIdFragmentMap)-len(AllDeleteWorkerErrors), len(AllDeleteWorkerErrors)), ErrJobFailed
	} else {
		return "Finished. Deleted " + fmt.Sprintf("%d", len(dataIdFragmentMap)) + " File Versions", nil
	}

}

// deleteWorker is a worker function for the cron job.
// Takes in a channel of dataIds to delete and deletes them based on the map of dataId to fragments
// Returns a channel of DeleteWorkerStatus
// TODO: See if you can optimize by using an input channel of a compound struct with dataId and []fragment
// instead of a map of dataId to fragments
func (i Inc) deleteWorker(dataIdFragmentMap map[string][]dbfs.Fragment, input <-chan string,
	output chan<- DeleteWorkerStatus, db *gorm.DB) {

	// for each dataId
	for dataId := range input {

		// Create channels to receive the results of the goroutines
		serversBackChan := make(chan DeleteWorkerServerStatus, len(dataIdFragmentMap[dataId]))

		// Set to monitor timeout servers.
		serversPending := util.NewSet[string]()
		for _, fragment := range dataIdFragmentMap[dataId] {
			serversPending.Add(fragment.ServerName)
		}

		// spin up 'em routines
		for _, fragment := range dataIdFragmentMap[dataId] {
			go func(path, server, dataId string, dwss chan<- DeleteWorkerServerStatus) {

				if server == i.ServerName {
					// local
					err := i.DeleteFragmentsByPath(path)
					if err != nil {
						dwss <- DeleteWorkerServerStatus{server, err}
					} else {
						dwss <- DeleteWorkerServerStatus{server, nil}
					}
				} else {
					// call handling server
					req, err := http.NewRequest(http.MethodDelete,
						strings.Replace(FragmentPath,
							"{fragmentPath}", path, -1), nil)
					if err != nil {
						fmt.Println(err)
						return // TODO CHECK IF SHOULD BE RETURN
					}
					resp, err := i.HttpClient.Do(req)
					if err != nil {
						fmt.Println(err)
						return // TODO CHECK IF SHOULD BE RETURN
					}

					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						dwss <- DeleteWorkerServerStatus{server, errors.New("Server returned " + resp.Status)}
					} else {
						dwss <- DeleteWorkerServerStatus{server, nil}
					}

				}

			}(fragment.FileFragmentPath, fragment.ServerName, fragment.FileVersionDataId, serversBackChan)
		}

		failedServerErrors := make([]DeleteWorkerServerStatus, 0)

		// waiting for the output from the channel. timeout for each server is 60 sec.
		for i := 0; i < len(dataIdFragmentMap[dataId]); i++ {
			select {
			case serverBack := <-serversBackChan:
				if serverBack.err != nil {
					failedServerErrors = append(failedServerErrors, serverBack)
				} else {
					serversPending.Remove(serverBack.server)
				}
			case <-time.After(time.Second * 60):
				// Log Timeout
				continue
			}

		}

		// Checking if any servers timed out
		if serversPending.Size() > 0 {
			serversPending.Each(func(server string) {
				failedServerErrors = append(failedServerErrors, DeleteWorkerServerStatus{server, ErrServerTimeout})
			})
		}

		// Checking how many failed
		if len(failedServerErrors) > 0 {
			output <- DeleteWorkerStatus{dataId, ErrServerFailed, failedServerErrors}
		} else {
			// success
			err := dbfs.FinishDeleteDataId(db, dataId)
			if err != nil {
				output <- DeleteWorkerStatus{dataId, err, nil}
			}
			output <- DeleteWorkerStatus{dataId, nil, nil}
		}

	}
}

// LocalOrphanedShardsCheck checks if there are any orphaned shards files in the local server
// If there are, it returns the file paths of the orphaned shards
func (i Inc) LocalOrphanedShardsCheck() ([]string, error) {

	// Get all the shards/fragments belonging to this server
	fragments, err := dbfs.GetFragmentByServer(i.Db, i.ServerName)
	if err != nil {
		return nil, err
	}

	// Convert fragments to a set of fragments for fast lookup
	fragmentSet := util.NewSet[string]()
	for _, fragment := range fragments {
		fragmentSet.Add(fragment.FileFragmentPath)
	}

	// List of orphaned shards
	orphanedShards := make([]string, 0)

	dir, err := os.ReadDir(i.ShardsLocation)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		// Check if the file is in the list of fragments
		if !fragmentSet.Has(file.Name()) {
			orphanedShards = append(orphanedShards, file.Name())
		}
	}

	if len(orphanedShards) > 0 {
		return orphanedShards, ErrOrphanedShardsFound
	} else {
		return nil, nil
	}

}

// LocalMissingShardsCheck checks if there are any shards or fragment files missing in the local server
// returns a list of missing files if any
func (i Inc) LocalMissingShardsCheck() ([]string, error) {

	// Get all the shards/fragments belonging to this server
	fragments, err := dbfs.GetFragmentByServer(i.Db, i.ServerName)
	if err != nil {
		return nil, err
	}

	// Convert local files to a set of filenames for fast lookup
	dir, err := os.ReadDir(i.ShardsLocation)
	if err != nil {
		return nil, err
	}
	localFiles := util.NewSet[string]()
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		localFiles.Add(file.Name())
	}

	// Array of missing Fragment data id
	missingFragments := make([]string, 0)

	// Check if the local files are in the list of fragments
	for _, fragment := range fragments {
		if !localFiles.Has(fragment.FileFragmentPath) {
			missingFragments = append(missingFragments, fragment.FileVersionDataId)
		}
	}

	if len(missingFragments) > 0 {
		return missingFragments, ErrMissingShardsFound
	} else {
		return nil, nil
	}

}

// LocalIndividualFragHealthCheck checks if the local fragment is in good condition
func (i Inc) LocalIndividualFragHealthCheck(fragPath string) (*stitch.ShardVerificationResult, error) {

	shardFile, err := os.Open(path.Join(i.ShardsLocation, fragPath))
	// err handling
	if err != nil {
		return nil, err
	}

	integrity, err := stitch.VerifyShardIntegrity(shardFile)

	return integrity, err
}

func (i Inc) IndividualFragHealthCheck(fragment dbfs.Fragment) (*stitch.ShardVerificationResult, error) {

	var result *stitch.ShardVerificationResult
	var err error

	// Check who the /Users/adrieltan/github/stitch/cmd/stitch/cmd/pipeline_cmd.gofragment belongs to
	if fragment.ServerName == i.ServerName {
		// local
		result, err = i.LocalIndividualFragHealthCheck(fragment.FileFragmentPath)
		if err != nil {
			return nil, err
		}
	} else {
		// get server to check
		server, err := dbfs.GetServerAddress(i.Db, fragment.ServerName)
		if err != nil {
			return nil, err
		}

		// call handling server

		resp, err := i.HttpClient.Get(server + strings.Replace(FragmentHealthCheckPath,
			"{fragmentPath}", fragment.FileFragmentPath, -1))

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, errors.New("Error: " + resp.Status)
		}

		// decode the response
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}

	}

	// mark shard health as bad
	if len(result.BrokenBlocks) > 0 {
		err := fragment.UpdateStatus(i.Db, dbfs.FragmentStatusBad)
		if err != nil {
			return nil, err
		}
	}
	return result, nil

}

func (i Inc) LocalCurrentFilesFragmentsHealthCheck(jobId int) bool {
	// db join query to get all the fragments belonging to this server current versions

	type result struct {
		FileId           string
		FileName         string
		DataId           string
		FileFragmentPath string
		ServerName       string
	}

	// store start in the JobProgress_CFFHC
	err := i.Db.Create(&dbfs.JobprogressCffhc{
		JobId:      jobId,
		StartTime:  time.Now(),
		ServerId:   i.ServerName,
		InProgress: true,
	}).Error
	if err != nil {
		log.Println(err)
		return false
	}

	var results []result

	var resultsCffhc []dbfs.ResultsCffhc

	// Get all the fragments belonging to this server
	// We are going to use a join

	err = i.Db.Model(&dbfs.File{}).Select(
		"files.file_id, files.file_name, files.data_id, fragments.file_fragment_path, fragments.server_name").
		Joins("JOIN fragments ON files.data_id = fragments.file_version_data_id").
		Where("fragments.server_name = ?", i.ServerName).Find(&results).Error
	if err != nil {
		log.Println(err)
		return false
	}

	for _, result := range results {
		// check the health of the fragment
		verificationResult, err := i.LocalIndividualFragHealthCheck(result.FileFragmentPath)
		if err != nil {
			// mark the fragment as bad
			resultsCffhc = append(resultsCffhc, dbfs.ResultsCffhc{
				JobId:     jobId,
				FileName:  result.FileName,
				FileId:    result.FileId,
				FragPath:  result.FileFragmentPath,
				ServerId:  result.ServerName,
				Error:     err.Error(),
				ErrorType: 1,
			})
		} else if !verificationResult.IsAvailable {
			resultsCffhc = append(resultsCffhc, dbfs.ResultsCffhc{
				JobId:     jobId,
				FileName:  result.FileName,
				FileId:    result.FileId,
				FragPath:  result.FileFragmentPath,
				ServerId:  result.ServerName,
				Error:     "Fragment is not available",
				ErrorType: 2,
			})
		} else if len(verificationResult.BrokenBlocks) > 0 {
			resultsCffhc = append(resultsCffhc, dbfs.ResultsCffhc{
				JobId:     jobId,
				FileName:  result.FileName,
				FileId:    result.FileId,
				FragPath:  result.FileFragmentPath,
				ServerId:  result.ServerName,
				Error:     "Fragment is broken",
				ErrorType: 3,
			})
		}
	}

	// insert the results into the database
	err = i.Db.Create(&resultsCffhc).Error
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (i Inc) LocalAllFilesFragmentsHealthCheck(jobId int) bool {
	// db join query to get all the fragments belonging to this server current versions

	type result struct {
		FileId           string
		FileName         string
		DataId           string
		FileFragmentPath string
		ServerName       string
	}

	// store start in the JobprogressAffhc
	err := i.Db.Create(&dbfs.JobprogressAffhc{
		JobId:      jobId,
		StartTime:  time.Now(),
		ServerId:   i.ServerName,
		InProgress: true,
	}).Error
	if err != nil {
		log.Println(err)
		return false
	}

	var results []result

	var resultsAffhc []dbfs.ResultsAffhc

	// Get all the fragments belonging to this server
	// We are going to use a join

	err = i.Db.Model(&dbfs.FileVersion{}).Select(
		"file_versions.file_id, file_versions.file_name, file_versions.data_id, fragments.file_fragment_path, fragments.server_name").
		Joins("JOIN fragments ON file_versions.data_id = fragments.file_version_data_id").
		Where("fragments.server_name = ?", i.ServerName).Find(&results).Error
	if err != nil {
		log.Println(err)
		return false
	}

	for _, result := range results {
		// check the health of the fragment
		verificationResult, err := i.LocalIndividualFragHealthCheck(result.FileFragmentPath)
		if err != nil {
			// mark the fragment as bad
			resultsAffhc = append(resultsAffhc, dbfs.ResultsAffhc{
				JobId:     jobId,
				FileName:  result.FileName,
				FileId:    result.FileId,
				FragPath:  result.FileFragmentPath,
				ServerId:  result.ServerName,
				Error:     err.Error(),
				ErrorType: 1,
			})
		} else if !verificationResult.IsAvailable {
			resultsAffhc = append(resultsAffhc, dbfs.ResultsAffhc{
				JobId:     jobId,
				FileName:  result.FileName,
				FileId:    result.FileId,
				FragPath:  result.FileFragmentPath,
				ServerId:  result.ServerName,
				Error:     "Fragment is not available",
				ErrorType: 2,
			})
		} else if len(verificationResult.BrokenBlocks) > 0 {
			resultsAffhc = append(resultsAffhc, dbfs.ResultsAffhc{
				JobId:     jobId,
				FileName:  result.FileName,
				FileId:    result.FileId,
				FragPath:  result.FileFragmentPath,
				ServerId:  result.ServerName,
				Error:     "Fragment is broken",
				ErrorType: 3,
			})
		}
	}

	// insert the results into the database
	err = i.Db.Create(&resultsAffhc).Error
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
