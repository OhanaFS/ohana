package inc

import (
	"errors"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"gorm.io/gorm"
	"os"
	"path"
	"time"
)

var (
	JobCurrentlyRunning        = errors.New("job is currently running")
	JobCurrentlyRunningWarning = errors.New("job is currently running. warning")
	ErrServerFailed            = errors.New("server failed")
	ErrServerTimeout           = errors.New("server timeout")
	ErrJobFailed               = errors.New("job failed")
)

func (i Inc) DeleteFragmentsByPath(pathOfShard string) error {
	return os.Remove(path.Join(i.ShardsLocation, pathOfShard))
}

func (i Inc) CronJobDeleteFragments(manualRun bool) (error, string) {

	// Checks if the job is currently being done.

	// The job should be only handled by the server with the least free space LMAO.
	// Thus, we register again to get the latest info.

	err := i.RegisterServer(false)
	if err != nil {
		return err, "register server error"
	}

	// Check if the job is currently running.

	server, timestamp, err := dbfs.IsCronDeleteRunning(i.Db)
	if err != nil {
		return err, ""
	}
	if server != "" {

		errorMsg := "cron job is already running by server " + server
		errorMsg += "\nlast started at " + timestamp.String()

		// Give warning if it seems stuck
		if timestamp.Add(time.Hour).Before(time.Now()) {
			errorMsg += "\nWARNING: last started more than an hour ago"
			return JobCurrentlyRunningWarning, errorMsg
		}

		return JobCurrentlyRunning, errorMsg
	}
	// else it should not be running

	// Check if the server should handle it (least amount of data free)
	servers, err := dbfs.GetServers(i.Db)
	if err != nil {
		return err, ""
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
			return errors.New("Assigning server " + serverWithLeastFreeSpace + " to process it"), ""
			// TODO
		}
	}

	// Start the job
	_, err = dbfs.MarkOldFileVersions(i.Db)
	if err != nil {
		return err, ""
	}

	fragments, err := dbfs.GetToBeDeletedFragments(i.Db)
	if err != nil {
		return err, ""
	}

	if len(fragments) == 0 {
		return nil, "Finished. Deleted 0 fragments"
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
		return ErrJobFailed, fmt.Sprintf("Finished (WARNING). Deleted %d, but %d errors occurred",
			len(dataIdFragmentMap)-len(AllDeleteWorkerErrors), len(AllDeleteWorkerErrors))
	} else {
		return nil, "Finished. Deleted " + fmt.Sprintf("%d", len(dataIdFragmentMap)) + " fragments"
	}

}

type DeleteWorkerStatus struct {
	DataId       string
	Error        error
	ServerErrors []DeleteWorkerServerStatus
}

type DeleteWorkerServerStatus struct {
	server string
	err    error
}

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
					fmt.Println("Deleting fragment:", path, server, dataId)
					dwss <- DeleteWorkerServerStatus{server, nil}
				}

			}(fragment.FileFragmentPath, fragment.ServerName, fragment.FileVersionDataId, serversBackChan)
		}

		failedServerErrors := make([]DeleteWorkerServerStatus, 0)

		// waiting for the output from the channel. timeout for each server is 60 sec.
		for i := 0; i < len(dataIdFragmentMap[dataId]); i++ {
			serverBack := <-serversBackChan
			if serverBack.err != nil {
				failedServerErrors = append(failedServerErrors, serverBack)
			} else {
				serversPending.Remove(serverBack.server)
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
