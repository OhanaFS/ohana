package inc

import (
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"golang.org/x/sys/unix"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

// RegisterServer registers a server as online in the database.
// If server already exists, it will update the hostname and status.
// Will attempt to connect to every node in the cluster.
func RegisterServer(tx *gorm.DB, config *config.Config, initalRun bool) error {

	// Race conditions to ensure no other server is registering atm.

	serverNotReady := true
	attempts := 0

	for serverNotReady {

		var server dbfs.Server
		err := tx.Model(&dbfs.Server{}).Where("status = ?", dbfs.ServerStarting).First(&server).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				serverNotReady = false
				continue
			}
			return err
		}
		attempts += 1

		fmt.Println("Waiting for other server to finish registering... attempt", attempts)
		time.Sleep(time.Second * 4)

	}

	spaceFree := getFreeSpace(config.Stitch.ShardsLocation)

	// Register as Starting

	// Check if server exists.
	var server dbfs.Server
	err := tx.First(&server, "name = ?", config.Database.ServerName).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Registering as new server...")
		} else {
			return err
		}
	}

	server = dbfs.Server{
		Name:      config.Database.ServerName,
		HostName:  config.Database.HostName,
		Port:      config.Database.Port,
		Status:    dbfs.ServerStarting,
		FreeSpace: spaceFree,
	}

	if err := tx.Save(&server).Error; err != nil {
		return err
	}

	// TODO: Check if server can ping every node in the cluster.
	// If it can't, it will be marked as offline.

	if initalRun {
		fmt.Println("Checking if server can ping every node in the cluster...")
	}

	var servers []dbfs.Server
	tx.Find(&servers).Where("status = ? ", dbfs.ServerOnline)

	for _, server := range servers {
		if !Ping(server.HostName, server.Port) {
			fmt.Println("Server", server.HostName, "is unreachable.")
			err := MarkServerOffline(tx, config, server.HostName)
			if err != nil {
				tx.Delete(&server)
				return err
			}
		}
	}

	// Register as Online

	server.Status = dbfs.ServerOnline

	if initalRun {
		fmt.Println("Registering server as online...")
	}

	err = tx.Save(&server).Error
	if err != nil {
		return err
	}

	if initalRun {
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

func MarkServerOffline(tx *gorm.DB, config *config.Config, server string) error {

	err := dbfs.MarkServerOffline(tx, server)
	if err != nil {
		return err
	}

	// TODO: Log this event.

	return err
}

func Ping(hostname, port string) bool {

	client := &http.Client{} // TODO: Put TLS config here. Need to configure holding object.

	res, err := client.Get("http://" + hostname + ":" + port + "/inc/ping")
	if err != nil {
		return false
	}

	return res.StatusCode == http.StatusOK
}

// Pong returns true if the server is online.
func Pong(w http.ResponseWriter, r *http.Request) {
	util.HttpJson(w, http.StatusOK, map[string]string{"status": "online"})
}
