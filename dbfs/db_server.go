package dbfs

import (
	"errors"
	"gorm.io/gorm"
)

const (
	ServerOnline       = int8(1)
	ServerOffline      = int8(2)
	ServerStarting     = int8(3)
	ServerStopping     = int8(4)
	ServerWarning      = int8(5)
	ServerError        = int8(6)
	ServerOfflineError = int8(7)
)

var (
	ErrServerNotFound = errors.New("server not found")
	ErrServerOffline  = errors.New("server is offline")
)

// Server needs a name. Don't need UUID.
// With that, comes with an IP or host address so it knows what to connect to someone with.
// also needs an struct to specify the status of it.
type Server struct {
	Name      string `gorm:"not null; primaryKey; unique" json:"name"`
	HostName  string `gorm:"not null; unique" json:"hostname"`
	Port      string `gorm:"not null" json:"port"`
	Status    int8   `gorm:"not null" json:"status"`
	FreeSpace uint64 `gorm:"not null" json:"free_space"`
	UsedSpace uint64 `gorm:"not null" json:"used_space"`
}

// GetServers returns all servers in the database.
func GetServers(tx *gorm.DB) ([]Server, error) {

	var servers []Server
	err := tx.Where("name <> ?", "").Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return servers, nil
}

// GetServerAddress returns the address of the server.
func GetServerAddress(tx *gorm.DB, serverName string) (string, error) {

	var server Server
	err := tx.Find(&server, "name = ?", serverName).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", ErrServerNotFound
		} else {
			return "", err
		}
	}

	if server.Status == ServerOffline {
		return "", ErrServerOffline
	}

	return server.HostName + ":" + server.Port, nil

}

// MarkServerOffline marks the server as offline. Used when turning off a server.
func MarkServerOffline(tx *gorm.DB, serverName string) error {

	var server Server
	err := tx.Find(&server, "name = ?", serverName).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrServerNotFound
		} else {
			return err
		}
	}

	if err := tx.Model(&server).Where("name = ?", serverName).Update("status", ServerOffline).Error;
		err != nil {
		return err
	}

	// Mark fragments as offline.

	err = tx.Model(&Fragment{}).Where("server_name = ?", serverName).
		Update("status", FragmentStatusOffline).Error

	if err != nil {
		return err
	}

	return nil
}
