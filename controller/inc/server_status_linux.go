//go:build linux
// +build linux

package inc

import (
	"fmt"
	"github.com/mackerelio/go-osstat/disk"
	memory2 "github.com/mackerelio/go-osstat/memory"
	"strings"
	"time"
)

func GetMemoryStats() (*memUsage, error) {

	mem, err := memory2.Get()
	if err != nil {
		return nil, err
	}

	memUsed := mem.Total - mem.Available
	memFree := mem.Available
	memTotal := mem.Total

	memory := memUsage{
		MemoryUsed:  memUsed,
		MemoryFree:  memFree,
		MemoryTotal: memTotal,
	}

	return &memory, nil
}

func GetDiskRW(dn string) (*diskRW, error) {

	initalStats, err := disk.Get()
	if err != nil {
		return nil, err
	}

	var initalRead, initalWrite uint64

	var disknames string

	for _, stat := range initalStats {
		disknames = disknames + stat.Name + " "
		if strings.Contains(stat.Name, dn) {
			initalRead = stat.ReadsCompleted
			initalWrite = stat.WritesCompleted
			break
		}
	}

	time.Sleep(time.Second)

	afterStats, err := disk.Get()
	if err != nil {
		return nil, err
	}

	for _, stat := range afterStats {
		if strings.Contains(stat.Name, dn) {
			diskRW := diskRW{
				Read:  stat.ReadsCompleted - initalRead,
				Write: stat.WritesCompleted - initalWrite,
			}
			return &diskRW, nil
		}
	}

	return nil, fmt.Errorf("disk %s not found, out of %s", dn, disknames)

	//return nil, ErrCannotFindDriveDeviceName
}
