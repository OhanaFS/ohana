package inc

import (
	"flag"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/uptime"
	"net"
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

type memUsage struct {
	MemoryUsed  uint64
	MemoryFree  uint64
	MemoryTotal uint64
}

type diskRW struct {
	Read  uint64
	Write uint64
}

func (i Inc) GetLocalServerStatusReport() (*LocalServerReport, error) {

	// Get the free space on the server.
	usedSpace := getUsedStorage(i.ShardsLocation)
	if usedSpace <= 0 {
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
	rx, tx, err := i.GetRXTX()
	if err != nil {
		return nil, err
	}

	// read writes

	//read, write, err := i.GetDriveReadWrite()
	//if err != nil {
	//	return nil, err
	//}

	// warnings errors fatal

	warnings, errors, fatales, err := i.GetWarningsErrorsFatal()

	// TODO check smart NOT IMPLEMENTED.

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

	// Need to add IPs to config file to know what to bind to.
	// From there we can pick it up.

	return "", ErrCannotFindNetworkAdaptor

}

// GetRXTX returns the RX and TX bytes/second for the network adapter inc is bound to.
func (i Inc) GetRXTX() (uint64, uint64, error) {

	// Doesn't know the IP in testing, so leave it alone lmao.
	if flag.Lookup("test.v") != nil {
		return 0, 0, nil
	}

	// First get the adaptor with the IP.
	adaptor, err := getAdaptorWithIP(i.BindIp)
	if err != nil {
		return 0, 0, err
	}

	// get the rx tx bytes.
	initalStats, err := network.Get()

	var initalRx, initalTx uint64

	for _, stat := range initalStats {
		if stat.Name == adaptor {
			// process crap
			initalRx = stat.RxBytes
			initalTx = stat.TxBytes
			break
		}
	}

	// sleep for a second.
	time.Sleep(time.Second)

	newStats, err := network.Get()

	for _, stat := range newStats {
		if stat.Name == adaptor {
			return stat.RxBytes - initalRx, stat.TxBytes - initalTx, nil
		}
	}

	return 0, 0, ErrCannotFindNetworkAdaptor
}

func (i Inc) GetFSDeviceName() ([]string, error) {

	cmd := exec.Command("df", "-h", i.ShardsLocation)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// get the device name from the output.
	rows := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(rows) < 2 {
		return nil, fmt.Errorf("No device name found in output %s ", rows)
	}

	arrayOfDeviceNames := make([]string, len(rows)-1)

	for j := 1; j < len(rows); j++ {
		cols := strings.Fields(rows[j])
		arrayOfDeviceNames[j-1] = cols[0]
	}

	return arrayOfDeviceNames, nil
}

func (i Inc) GetDriveReadWrite() (uint64, uint64, error) {
	deviceNames, err := i.GetFSDeviceName()
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

func (i Inc) GetWarningsErrorsFatal() (int64, int64, int64, error) {

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
