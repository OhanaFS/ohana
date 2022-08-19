package inc

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"golang.org/x/sys/unix"
	"gorm.io/gorm"
)

// RegisterServer registers a server as online in the database.
// If server already exists, it will update the hostname and status.
// Will attempt to connect to every node in the cluster.
func (i Inc) RegisterServer(initialRun bool) error {

	spaceFree := getFreeSpace(i.ShardsLocation)
	spaceUsed := uint64(getUsedStorage(i.ShardsLocation))

	// Register as Starting

	isNewServer := false

	// Check if server exists.
	var server dbfs.Server
	err := i.Db.First(&server, "name = ?", i.ServerName).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			isNewServer = true
			fmt.Println("Registering as new server...")
		} else {
			return err
		}
	}

	server.Name = i.ServerName
	server.HostName = i.HostName
	server.Port = i.Port
	server.Status = dbfs.ServerStarting
	server.FreeSpace = spaceFree
	server.UsedSpace = spaceUsed

	if err := i.Db.Save(&server).Error; err != nil {
		return err
	}

	// Check if server can ping every node in the cluster.
	// If it can't, it will be marked as offline.

	if initialRun {
		fmt.Println("Checking if server can ping every node in the cluster...")
	}

	var servers []dbfs.Server
	i.Db.Where("status = ? ", dbfs.ServerOnline).Find(&servers)

	for _, server := range servers {
		if !i.Ping(server.HostName, server.Port) {
			fmt.Println("Server", server.HostName, "is unreachable.")
			err := MarkServerOffline(i.Db, i.HostName, server.HostName)
			if err != nil {
				if isNewServer {
					if i.Db.Delete(&server).Error != nil {
						return errors.New("ERROR. CANNOT DELETE SERVER FROM DATABASE. " + err.Error())
					}
				} else {
					// mark server as offline instead
					server.Status = dbfs.ServerOfflineError
					if i.Db.Save(&server).Error != nil {
						return errors.New("ERROR. CANNOT UPDATE SERVER FROM DATABASE. " + err.Error())
					}
				}
				return err
			}
		}
	}

	// Register as Online

	if initialRun {
		fmt.Println("Registering server as online...")
	}

	err = i.Db.Model(&server).Where("name = ?", i.ServerName).Update("status", dbfs.ServerOnline).Error
	if err != nil {
		return err
	}

	if initialRun {
		fmt.Println("Updated. Server registered.")
	}

	return err
}

// getFreeSpace returns the free space of the folder
// TODO: Fix on windows
func getFreeSpace(path string) uint64 {
	var stat unix.Statfs_t

	w, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic("ERROR. CANNOT CREATE SHARDS FOLDER.")
		}
	} else if !w.IsDir() {
		panic("ERROR. SHARDS FOLDER IS NOT A DIRECTORY.")
	}

	err = unix.Statfs(path, &stat)
	if err != nil {
		panic("ERROR. CANNOT GET SHARDS FOLDER STATUS.")
	}

	return stat.Bavail * uint64(stat.Bsize)
}

// MarkServerOffline marks a server as offline in the database.
func MarkServerOffline(tx *gorm.DB, requestServer, destServer string) error {

	err := dbfs.MarkServerOffline(tx, destServer)
	if err != nil {
		return err
	}

	// TODO: Log this event.

	return err
}

// Ping returns true if the server is online.
func (i Inc) Ping(hostname, port string) bool {

	res, err := i.HttpClient.Get("http://" + hostname + ":" + port + "/inc/ping")
	if err != nil {
		return false
	}

	return res.StatusCode == http.StatusOK
}

// Pong returns true if the server is online.
// This function needs to be attached to a http server.
func Pong(w http.ResponseWriter, r *http.Request) {
	util.HttpJson(w, http.StatusOK, map[string]string{"status": "online"})
}

// ShutdownServer will try to gracefully shutdown the server.
func (i Inc) ShutdownServer(w http.ResponseWriter, r *http.Request) {

	i.Shutdown <- true

	util.HttpJson(w, http.StatusOK, map[string]string{"status": "shutdown"})
}

// RegisterIncServices registers the inc services.
func RegisterIncServices(i *Inc) {

	// Register server service.
	go func() {
		time.Sleep(time.Second * 2)
		registerTicker := time.NewTicker(5 * time.Minute)
		err := i.RegisterServer(true)
		if err != nil {
			i.Shutdown <- true
		}

		// Register service / shutdown handler
		go func() {
			for {
				select {
				case <-i.Shutdown:
					registerTicker.Stop()
					fmt.Println("Shutdown signal received. Exiting in 5 seconds...")
					time.Sleep(time.Second * 5)
					os.Exit(0)
					// deregister services
					return
				case <-registerTicker.C:
					_ = i.RegisterServer(false)
				}
			}
		}()

		// Daily Update handler
		dailyUpdateTicker := time.NewTicker(time.Hour)
		go func() {
			for {
				select {
				case <-dailyUpdateTicker.C:
					i.DailyUpdate()
				}
			}
		}()

	}()

}

func (i Inc) DailyUpdate() {

	// Daily Update of the stats of the server
	// Check if the server is the one to do the update

	logger := dbfs.NewLogger(i.Db, i.ServerName)
	logger.LogInfo("Daily update started.")

	servers, err := dbfs.GetServers(i.Db)
	if err != nil {
		logger.LogError("ERROR. CANNOT GET SERVERS FROM DATABASE." +
			err.Error())
	}

	var serverWithLeastFreeSpace string
	var leastFreeSpace uint64
	for _, server := range servers {
		freeSpace := server.FreeSpace
		if freeSpace < leastFreeSpace || leastFreeSpace == uint64(0) {
			leastFreeSpace = freeSpace
			serverWithLeastFreeSpace = server.Name
		}
	}

	if serverWithLeastFreeSpace == i.ServerName {
		// Log the update
		logger.LogInfo("Updated server stats...")
		if dbfs.DumpDailyStats(i.Db) != nil {
			logger.LogError("ERROR. CANNOT UPDATE DAILY STATS. " + err.Error())
		}
	}

}
