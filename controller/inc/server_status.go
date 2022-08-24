package inc

import (
	"encoding/json"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/uptime"
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type LocalServerReport struct {
	ServerName     string `json:"name"`
	Hostname       string `json:"hostname"`
	Port           string `json:"port"`
	Status         int    `json:"status"`
	FreeSpace      int64  `json:"free_space"`
	UsedSpace      int64  `json:"used_space"`
	LoadAvg        string `json:"load_avg"`
	Uptime         int64  `json:"uptime"`
	Cpu            int    `json:"cpu"`
	MemoryUsed     int64  `json:"memory_used"`
	MemoryFree     int64  `json:"memory_free"`
	NetworkRxBytes int64  `json:"network_rx_bytes"`
	NetworkTxBytes int64  `json:"network_tx_bytes"`
	Reads          int64  `json:"reads"`
	Writes         int64  `json:"writes"`
	Fatales        int64  `json:"fatal"`
	Warnings       int64  `json:"warnings"`
	Errors         int64  `json:"errors"`
	SmartGood      bool   `json:"smart_good"`
	SmartStatus    string `json:"smart_status"`
}

// memUsage is a struct that holds the memory usage for the server
// required to support multi-platforms (darwin, linux)
type memUsage struct {
	MemoryUsed  uint64
	MemoryFree  uint64
	MemoryTotal uint64
}

// diskRW is a struct that holds the read and write stats for the server
// required to support multi-platforms (darwin, linux)... tho broken on docker lol
type diskRW struct {
	Read  uint64
	Write uint64
}

// getLocalServerStatusReport returns a report of the local server.
func (i Inc) getLocalServerStatusReport() (*LocalServerReport, error) {

	// Get the free space on the server.
	usedSpace := getUsedStorage(i.ShardsLocation)
	if usedSpace < 0 {
		return nil, ErrServerFailed
	}

	// Get things from os system stat library

	loadAvg, err := loadavg.Get()
	if err != nil {
		return nil, err
	}
	loadAvgString := fmt.Sprintf("Load Avg 1 min: %f , Load Avg 5 min: %f, Load Avg 15 min: %f",
		loadAvg.Loadavg1, loadAvg.Loadavg5, loadAvg.Loadavg15)

	uptimeResult, err := uptime.Get()
	if err != nil {
		return nil, err
	}

	cpuStats, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	cpuUsage := float64(cpuStats.User+cpuStats.Nice+cpuStats.System) / float64(cpuStats.Total) * 100

	// Get the memory usage of the server.

	memUsage, err := GetMemoryStats()

	// TODO: Convert this into a go routine and a channel to get the stats.
	// rx tx
	rx, tx, err := i.getRXTX()
	if err != nil {
		return nil, err
	}

	// read writes... broken on docker so we can't use it.	:/

	//read, write, err := i.getDriveReadWrite()
	//if err != nil {
	//	return nil, err
	//}

	// warnings errors fatal

	warnings, errors, fatales, err := i.getNumWarningsErrorsFatal()

	// TODO SMART CHECK NOT IMPLEMENTED.
	/*
		The main issue with SMART is that
		a. A lot of the smart commands are not implemented on all platforms, and they aren't clear like return bad status
		b. Most SMART commands require root access, and we don't necessarily want dbfs to run at that level
		c. Good SMART utils like smartctl need to be installed.
	*/

	lsr := LocalServerReport{
		ServerName:     i.ServerName,
		Hostname:       i.HostName,
		Port:           i.Port,
		Status:         int(dbfs.ServerOnline), // TODO: Might need to check this
		FreeSpace:      int64(getFreeSpace(i.ShardsLocation)),
		UsedSpace:      usedSpace,
		LoadAvg:        loadAvgString,
		Uptime:         int64(uptimeResult.Seconds()),
		Cpu:            int(cpuUsage),
		MemoryUsed:     int64(memUsage.MemoryUsed),
		MemoryFree:     int64(memUsage.MemoryFree),
		NetworkRxBytes: int64(rx),
		NetworkTxBytes: int64(tx),
		Warnings:       warnings,
		Errors:         errors,
		Fatales:        fatales,
		SmartGood:      true,
		SmartStatus:    "",
	}

	return &lsr, nil

}

// ReturnServerDetails returns the details for the server. Meant to be called if handling request server is not local
func (i *Inc) ReturnServerDetails(w http.ResponseWriter, r *http.Request) {

	serverDeets, err := i.getLocalServerStatusReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json encode file
	util.HttpJson(w, http.StatusOK, serverDeets)
}

// GetServerStatusReport returns the status report for the server.
func (i *Inc) GetServerStatusReport(serverName string) (*LocalServerReport, error) {

	var serverDeets *LocalServerReport
	var err error

	if serverName == i.ServerName {
		serverDeets, err = i.getLocalServerStatusReport()
		if err != nil {
			return nil, err
		}
	} else {
		// call other server for status
		// Check if the server exists
		var server dbfs.Server
		err = i.Db.Where("name = ?", serverName).First(&server).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, ErrServerNotFound
			}
			return nil, err
		}

		// get the server details from the other server
		request, err := http.NewRequest("GET", fmt.Sprintf("https://%s:%s/api/v1/node/details",
			server.HostName, server.Port), nil)
		if err != nil {
			return nil, err
		}
		response, err := i.HttpClient.Do(request)
		if err != nil {
			fmt.Println(err)

			offlineServerRequest := LocalServerReport{
				ServerName: serverName,
				Hostname:   server.HostName,
				Port:       server.Port,
				Status:     int(dbfs.ServerOffline),
			}

			return &offlineServerRequest, nil
		}

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("server %s returned status %d", serverName, response.StatusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &serverDeets)
		if err != nil {
			return nil, err
		}

		return serverDeets, nil

	}

	return serverDeets, nil
}

// getUsedStorage returns the used space on the server based on the shards' location.
func getUsedStorage(path string) int64 {
	usedSpace := int64(0)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			usedSpace += info.Size()
		}
		return err
	})
	if err != nil {
		return -1
	}
	return usedSpace
}

// getAdaptorWithIP returns the adaptor (eth0 for example) with the given IP.
func getAdaptorWithIP(ip string) (string, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			if strings.Contains(a.String(), ip) {
				return i.Name, nil
			}
		}
	}

	return "", ErrCannotFindNetworkAdaptor

}

// getRXTX returns the RX and TX bytes/second for the network adapter inc is bound to.
func (i Inc) getRXTX() (uint64, uint64, error) {

	// Doesn't know the IP in testing, so leave it alone lmao.
	//if flag.Lookup("test.v") != nil {
	//	return 0, 0, nil
	//}

	// First get the adaptor with the IP.
	//adaptor, err := getAdaptorWithIP(i.BindIp)
	//if err != nil {
	//	return 0, 0, err
	//}

	// get the rx tx bytes.
	initalStats, err := network.Get()
	if err != nil {
		return 0, 0, err
	}

	var initalRx, initalTx uint64

	initalRx = 0
	initalTx = 0

	for _, stat := range initalStats {
		initalRx += stat.RxBytes
		initalTx += stat.TxBytes
	}

	// sleep for a second.
	time.Sleep(time.Second)

	newStats, err := network.Get()
	if err != nil {
		return 0, 0, err
	}

	for _, stat := range newStats {

		initalRx -= stat.RxBytes
		initalTx -= stat.TxBytes
	}

	return initalRx, initalTx, nil
}

// getFSDeviceName returns the name of the device that is the root of the filesystem.
// BROKEN ON DOCKER.
func (i Inc) getFSDeviceName() ([]string, error) {

	cmd := exec.Command("df", "-h", i.ShardsLocation)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// get the device name from the output.
	rows := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(rows) < 2 {
		return nil, ErrCannotFindDriveDeviceName
	}

	arrayOfDeviceNames := make([]string, len(rows)-1)

	for j := 1; j < len(rows); j++ {
		cols := strings.Fields(rows[j])
		arrayOfDeviceNames[j-1] = cols[0]
	}

	return arrayOfDeviceNames, nil
}

// getDriveReadWrite returns the read and write actions/second for the drive.
func (i Inc) getDriveReadWrite() (uint64, uint64, error) {
	deviceNames, err := i.getFSDeviceName()
	if err != nil {
		return 0, 0, err
	}
	var rx, tx uint64

	for _, deviceName := range deviceNames {
		stats, err := GetDiskRW(deviceName)
		if err != nil {
			return 0, 0, err
		}
		rx += stats.Read
		tx += stats.Read
	}
	return rx, tx, nil
}

// getNumWarningsErrorsFatal returns the number of warnings, errors and fatal errors from alerts.
func (i Inc) getNumWarningsErrorsFatal() (int64, int64, int64, error) {

	var warnings, errors, fatal int64

	err := i.Db.Model(&dbfs.Alert{}).Where("log_type = ? AND server_name = ?", dbfs.LogServerWarning, i.ServerName).
		Count(&warnings).Error
	if err != nil {
		return 0, 0, 0, err
	}
	err = i.Db.Model(&dbfs.Alert{}).Where("log_type = ? AND server_name = ?", dbfs.LogServerError, i.ServerName).
		Count(&errors).Error
	if err != nil {
		return 0, 0, 0, err
	}
	err = i.Db.Model(&dbfs.Alert{}).Where("log_type = ? AND server_name = ?", dbfs.LogServerFatal, i.ServerName).
		Count(&fatal).Error
	if err != nil {
		return 0, 0, 0, err
	}
	return warnings, errors, fatal, nil

}

// getAlertsForServer returns all the alerts for that server
func (i Inc) getAlertsForServer() ([]dbfs.Alert, error) {

	var alerts []dbfs.Alert
	err := i.Db.Model(&dbfs.Alert{}).Where("server_name = ?", i.ServerName).Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}
