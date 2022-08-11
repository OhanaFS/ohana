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
						return
					}
					resp, err := i.HttpClient.Do(req)
					if err != nil {
						fmt.Println(err)
						return
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
func (i Inc) LocalOrphanedShardsCheck(jobId int, storeResults bool) ([]string, error) {

	if storeResults {
		// create JobProgressOrphanedShard
		jobProgressOrphanedShard := dbfs.JobProgressOrphanedShard{
			JobId:      jobId,
			StartTime:  time.Now(),
			ServerId:   i.ServerName,
			InProgress: true,
		}
		if i.Db.Create(&jobProgressOrphanedShard).Error != nil {
			return nil, errors.New("failed to create JobProgressOrphanedShard")
		}
	}

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

	if storeResults {
		// dump it into ResultsOrphanedShard

		var resultsOrphanedShards = make([]dbfs.ResultsOrphanedShard, len(orphanedShards))

		for j, orphanedShard := range orphanedShards {
			resultsOrphanedShards[j] = dbfs.ResultsOrphanedShard{
				JobId:    jobId,
				ServerId: i.ServerName,
				FileName: orphanedShard,
			}
		}

		if len(resultsOrphanedShards) > 0 {
			if i.Db.Create(&resultsOrphanedShards).Error != nil {
				return nil, errors.New("Failed to create ResultsOrphanedShard")
			}
		}

		// update JobProgressOrphanedShard
		if i.Db.Model(&dbfs.JobProgressOrphanedShard{}).
			Where("job_id = ? and server_id = ?", jobId, i.ServerName).
			Update("in_progress", false).Error != nil {
			return nil, errors.New("Failed to update JobProgressOrphanedShard")
		}

	}

	if len(orphanedShards) > 0 {
		return orphanedShards, nil
	} else {
		return nil, nil
	}

}

// LocalMissingShardsCheck checks if there are any shards or fragment files missing in the local server
// doesn't store results into database.
func (i Inc) LocalMissingShardsCheck(jobId int, storeResults bool) ([]dbfs.ResultsMissingShard, error) {

	if storeResults {
		// create JobProgressMissingShard
		jobProgressMissingShard := dbfs.JobProgressMissingShard{
			JobId:      jobId,
			StartTime:  time.Now(),
			ServerId:   i.ServerName,
			InProgress: true,
		}
		if i.Db.Create(&jobProgressMissingShard).Error != nil {
			return nil, errors.New("failed to create JobProgressMissingShard")
		}
	}

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
	var missingShards []dbfs.ResultsMissingShard

	// Check if the local files are in the list of fragments
	for _, fragment := range fragments {
		if !localFiles.Has(fragment.FileFragmentPath) {

			// Get FileName if shard is bad

			var fileVersions []dbfs.FileVersion
			err := i.Db.Where("data_id = ?", fragment.FileVersionDataId).Find(&fileVersions).Error
			if err != nil {
				return nil, err
			}

			filenames := make([]string, len(fileVersions))
			fileids := make([]string, len(fileVersions))
			for i, fileVersion := range fileVersions {
				filenames[i] = fileVersion.FileName
				fileids[i] = fileVersion.FileId
			}

			jsonFilenameBytes, _ := json.Marshal(filenames)
			jsonFileIdBytes, _ := json.Marshal(fileids)

			missingShards = append(missingShards, dbfs.ResultsMissingShard{
				JobId:     jobId,
				FileName:  string(jsonFilenameBytes),
				FileId:    string(jsonFileIdBytes),
				DataId:    fragment.FileVersionDataId,
				FragPath:  fragment.FileFragmentPath,
				ServerId:  i.ServerName,
				Error:     "Missing shard",
				ErrorType: 1,
			})
		}
	}

	if storeResults {
		// dump it into ResultsMissingShard
		if len(missingShards) > 0 {
			if i.Db.Create(&missingShards).Error != nil {
				return nil, errors.New("Failed to create ResultsMissingShard")
			}
		}
		// update JobProgressMissingShard
		if i.Db.Model(&dbfs.JobProgressMissingShard{}).
			Where("job_id = ? and server_id = ?", jobId, i.ServerName).
			Update("in_progress", false).Error != nil {
			return nil, errors.New("Failed to update JobProgressMissingShard")
		}
	}

	if len(missingShards) > 0 {
		return missingShards, nil
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
		if err != nil {
			return nil, err
		}

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

func (i Inc) LocalCurrentFilesFragmentsHealthCheck(jobId int) error {
	// db join query to get all the fragments belonging to this server current versions

	type result struct {
		FileId           string
		FileName         string
		DataId           string
		FileFragmentPath string
		ServerName       string
	}

	// store start in the JobProgress_CFFHC.
	err := i.Db.Create(&dbfs.JobProgressCFSHC{
		JobId:      jobId,
		StartTime:  time.Now(),
		ServerId:   i.ServerName,
		InProgress: true,
	}).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var fragmentsToCheck []result

	// Get all the fragments belonging to this server
	// We are going to use a join

	err = i.Db.Model(&dbfs.File{}).Select(
		"files.file_id, files.file_name, files.data_id, fragments.file_fragment_path, fragments.server_name").
		Joins("JOIN fragments ON files.data_id = fragments.file_version_data_id").
		Where("fragments.server_name = ? AND files.entry_type = ?", i.ServerName, dbfs.IsFile).
		Find(&fragmentsToCheck).Error
	if err != nil {
		log.Println(err)
		return err
	}

	type keyStruct struct {
		FileFragmentPath string
		FragmentServer   string
	}

	resultsMap := make(map[keyStruct]*dbfs.ResultsCFSHC)

	for _, frag := range fragmentsToCheck {

		key := keyStruct{
			FileFragmentPath: frag.FileFragmentPath,
			FragmentServer:   frag.ServerName,
		}

		existsFlag := false
		if _, ok := resultsMap[key]; ok {
			existsFlag = true
		}

		fileIdsJSONBytes, _ := json.Marshal([]string{frag.FileId})
		fileNamesJSONBytes, _ := json.Marshal([]string{frag.FileName})

		fileIdsJSON := string(fileIdsJSONBytes)
		fileNamesJSON := string(fileNamesJSONBytes)

		if !existsFlag {
			verificationResult, err := i.LocalIndividualFragHealthCheck(frag.FileFragmentPath)

			if err != nil {
				// mark the fragment as bad
				resultsMap[key] = &dbfs.ResultsCFSHC{
					JobId:     jobId,
					FileName:  fileNamesJSON,
					FileId:    fileIdsJSON,
					DataId:    frag.DataId,
					FragPath:  frag.FileFragmentPath,
					ServerId:  frag.ServerName,
					Error:     err.Error(),
					ErrorType: 1,
				}
			} else if !verificationResult.IsAvailable {
				resultsMap[key] = &dbfs.ResultsCFSHC{
					JobId:     jobId,
					FileName:  fileNamesJSON,
					FileId:    fileIdsJSON,
					DataId:    frag.DataId,
					FragPath:  frag.FileFragmentPath,
					ServerId:  frag.ServerName,
					Error:     "Fragment is not available",
					ErrorType: 2,
				}
			} else if len(verificationResult.BrokenBlocks) > 0 {
				resultsMap[key] = &dbfs.ResultsCFSHC{
					JobId:     jobId,
					FileName:  fileNamesJSON,
					FileId:    fileIdsJSON,
					DataId:    frag.DataId,
					FragPath:  frag.FileFragmentPath,
					ServerId:  frag.ServerName,
					Error:     "Fragment is broken",
					ErrorType: 3,
				}
			}

		} else {
			existingResult := resultsMap[key]

			tempStrings := make([]string, 0)

			// FileName
			if json.Unmarshal([]byte(existingResult.FileName), &tempStrings) != nil {
				continue
			}

			tempStrings = append(tempStrings, frag.FileName)
			tempStringsByte, err := json.Marshal(tempStrings)
			if err != nil {
				continue
			}

			existingResult.FileName = string(tempStringsByte)

			// FileId

			if json.Unmarshal([]byte(existingResult.FileId), &tempStrings) != nil {
				continue
			}

			tempStrings = append(tempStrings, frag.FileId)
			tempStringsByte, err = json.Marshal(tempStrings)
			if err != nil {
				continue
			}

			existingResult.FileId = string(tempStringsByte)

		}
	}

	// insert the results into the database
	for _, result := range resultsMap {
		if i.Db.Create(result).Error != nil {
			log.Println(err)
			return err
		}
	}

	// Close the job progress
	i.Db.Model(&dbfs.JobProgressCFSHC{}).
		Where("job_id = ? AND server_id = ?", jobId, i.ServerName).
		Update("in_progress", false)

	return nil
}

func (i Inc) LocalAllFilesFragmentsHealthCheck(jobId int) error {
	// db join query to get all the fragments belonging to this server current versions

	type result struct {
		FileId            string
		FileName          string
		FileVersionDataId string
		FileFragmentPath  string
		ServerName        string
	}

	// store start in the JobProgressAFSHC.
	err := i.Db.Create(&dbfs.JobProgressAFSHC{
		JobId:      jobId,
		StartTime:  time.Now(),
		ServerId:   i.ServerName,
		InProgress: true,
	}).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var fragmentsToCheck []result

	// Get all the fragments belonging to this server
	// We are going to use a join

	// We will first reference the files table so that we get the latest file names
	// If there are stuff no longer there, it will be "" and we will find it later.

	err = i.Db.Model(&dbfs.Fragment{}).Select(
		"files.file_id, files.file_name, fragments.file_version_data_id, fragments.file_fragment_path, fragments.server_name").
		Joins("JOIN files ON fragments.file_version_file_id = files.file_id").
		Where("fragments.server_name = ? AND files.entry_type = ?", i.ServerName, dbfs.IsFile).Find(&fragmentsToCheck).Error
	if err != nil {
		log.Println(err)
		return err
	}

	type keyStruct struct {
		FileFragmentPath string
		FragmentServer   string
	}

	resultsMap := make(map[keyStruct]*dbfs.ResultsAFSHC)

	for _, frag := range fragmentsToCheck {

		// Check if any of the fragments are missing filenames or file_ids, in which, we grab them from
		// the file version table. This is done to ensure that we always get the latest version of the file
		// in the case that the file has been updated. And to ensure that it works properly with postgres

		key := keyStruct{
			FileFragmentPath: frag.FileFragmentPath,
			FragmentServer:   frag.ServerName,
		}

		multipleFileIds := make([]string, 0)
		multipleFileNames := make([]string, 0)

		var fileIdsJSON, fileNamesJSON string

		if frag.FileId == "" || frag.FileName == "" {

			var fileVersions []dbfs.FileVersion
			err = i.Db.Where("data_id = ? and status NOT IN ?",
				frag.FileVersionDataId, []int8{dbfs.FileStatusToBeDeleted, dbfs.FileStatusDeleted}).
				Order("created_at desc").Find(&fileVersions).Error
			if err != nil {
				log.Println(err)
				return err
			}
			if len(fileVersions) == 0 {
				continue // likely deleted file so we don't care about it
			} else if len(fileVersions) == 1 {
				frag.FileId = fileVersions[0].FileId
				frag.FileName = fileVersions[0].FileName

				fileNamesJSONBytes, _ := json.Marshal([]string{frag.FileName})
				fileIdsJSONBytes, _ := json.Marshal([]string{frag.FileId})
				fileNamesJSON = string(fileNamesJSONBytes)
				fileIdsJSON = string(fileIdsJSONBytes)

			} else if len(fileVersions) > 1 {
				for _, fileVersion := range fileVersions {
					multipleFileIds = append(multipleFileIds, fileVersion.FileId)
					multipleFileNames = append(multipleFileNames, fileVersion.FileName)
				}

				// JSON-ify FileName and FileId
				fileNamesJSONBytes, _ := json.Marshal(multipleFileNames)
				fileIdsJSONBytes, _ := json.Marshal(multipleFileIds)

				fileNamesJSON = string(fileNamesJSONBytes)
				fileIdsJSON = string(fileIdsJSONBytes)

			}
		} else { // still need to json-ify the file_id and file_name
			fileIdsJSONBytes, _ := json.Marshal([]string{frag.FileId})
			fileNamesJSONBytes, _ := json.Marshal([]string{frag.FileName})

			fileIdsJSON = string(fileIdsJSONBytes)
			fileNamesJSON = string(fileNamesJSONBytes)

		}

		// If we have checked before, skip and just append file name and file id to the existing results
		existsFlag := false
		if _, ok := resultsMap[key]; ok {
			existsFlag = true
		}

		if !existsFlag {

			// check the health of the fragment
			verificationResult, err := i.LocalIndividualFragHealthCheck(frag.FileFragmentPath)
			if err != nil {
				// mark the fragment as bad
				resultsMap[key] = &dbfs.ResultsAFSHC{
					JobId:     jobId,
					FileName:  fileNamesJSON,
					FileId:    fileIdsJSON,
					DataId:    frag.FileVersionDataId,
					FragPath:  frag.FileFragmentPath,
					ServerId:  frag.ServerName,
					Error:     err.Error(),
					ErrorType: 1,
				}

			} else if !verificationResult.IsAvailable {
				resultsMap[key] = &dbfs.ResultsAFSHC{
					JobId:     jobId,
					FileName:  fileNamesJSON,
					FileId:    fileIdsJSON,
					DataId:    frag.FileVersionDataId,
					FragPath:  frag.FileFragmentPath,
					ServerId:  frag.ServerName,
					Error:     "Fragment is not available",
					ErrorType: 2,
				}

			} else if len(verificationResult.BrokenBlocks) > 0 {

				resultsMap[key] = &dbfs.ResultsAFSHC{
					JobId:     jobId,
					FileName:  fileNamesJSON,
					FileId:    fileIdsJSON,
					DataId:    frag.FileVersionDataId,
					FragPath:  frag.FileFragmentPath,
					ServerId:  frag.ServerName,
					Error:     "Fragment is broken",
					ErrorType: 3,
				}
			}

		} else {

			// get the previous result
			existingResult := resultsMap[key]

			tempStrings := make([]string, 0)

			// FileName
			err := json.Unmarshal([]byte(existingResult.FileName), &tempStrings)
			if err != nil {
				continue
			}

			tempStrings = append(tempStrings, frag.FileName)
			tempStringsByte, err := json.Marshal(tempStrings)
			if err != nil {
				continue
			}

			existingResult.FileName = string(tempStringsByte)

			// FileId

			err = json.Unmarshal([]byte(existingResult.FileId), &tempStrings)
			if err != nil {
				continue
			}

			tempStrings = append(tempStrings, frag.FileId)
			tempStringsByte, err = json.Marshal(tempStrings)
			if err != nil {
				continue
			}

			existingResult.FileId = string(tempStringsByte)

		}

	}

	// Insert the results into the database
	for _, result := range resultsMap {
		err = i.Db.Create(result).Error
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Close the job progress
	i.Db.Model(&dbfs.JobProgressAFSHC{}).
		Where("job_id = ? AND server_id = ?", jobId, i.ServerName).
		Update("in_progress", false)

	return nil
}
