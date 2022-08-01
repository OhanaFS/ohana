package inc

import (
	"errors"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"gorm.io/gorm"
	"os"
	"sync"
	"time"
)

var (
	JobCurrentlyRunning        = errors.New("job is currently running")
	JobCurrentlyRunningWarning = errors.New("job is currently running. warning")
)

func (i Inc) DeleteFragmentsByPath(path string) error {
	return os.Remove(i.ShardsLocation + path)
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

	server, timestamp, err := dbfs.IsCronDeleteRunning(i.db)
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
	servers, err := dbfs.GetServers(i.db)
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

	// TODO:
	dbfs.MarkOldFileVersions(i.db)

	fragments, err := dbfs.GetToBeDeletedFragments(i.db)
	if err != nil {
		return err, ""
	}

	dataIdFragmentMap := make(map[string][]dbfs.Fragment)

	for _, fragment := range fragments {
		dataIdFragmentMap[fragment.FileVersionDataId] = append(dataIdFragmentMap[fragment.FileVersionDataId], fragment)
	}

	// "Deletion Code"

	const maxGoroutines = 10
	input := make(chan string, len(dataIdFragmentMap))
	output := make(chan string, len(dataIdFragmentMap))

	// Worker function

	for num := 0; num < maxGoroutines; num++ {
		go i.deleteWorker(dataIdFragmentMap, input, output)
	}

	for dataId, _ := range dataIdFragmentMap {
		input <- dataId
	}
	close(input)

	// TODO: ?? Will this freeze the db?
	err = i.db.Transaction(func(tx *gorm.DB) error {

		for i := 0; i < len(dataIdFragmentMap); i++ {
			dataIdProcessed := <-output

			// Create transaction
			err2 := dbfs.FinishDeleteDataId(tx, dataIdProcessed)
			if err2 != nil {
				return err2
			}
		}
		return nil
	})

	return nil, "Finished. Deleted " + fmt.Sprintf("%d", len(dataIdFragmentMap)) + " fragments"
}

func (i Inc) deleteWorker(dataIdFragmentMap map[string][]dbfs.Fragment, input <-chan string, output chan<- string) {

	for j := range input {
		var fragWg sync.WaitGroup

		for _, fragment := range dataIdFragmentMap[j] {
			fragWg.Add(1)
			go func(path, server, dataId string) {

				if server == i.ServerName {
					err := i.DeleteFragmentsByPath(path)
					if err != nil {
						panic(err) // TODO: Figure out how to handle this to make sure no panics.
					}
				} else {
					fmt.Println("Deleting fragment:", path, server, dataId)
				}

				// Assume this is a delete fragment call.

				defer fragWg.Done()
			}(fragment.FileFragmentPath, fragment.ServerName, fragment.FileVersionDataId)
		}

		fragWg.Wait()
		output <- j
	}

}
