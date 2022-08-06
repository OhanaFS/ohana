//go:build darwin
// +build darwin

package inc

import (
	"fmt"
	memory2 "github.com/mackerelio/go-osstat/memory"
	"os/exec"
	"strconv"
	"strings"
)

func GetMemoryStats() (*memUsage, error) {

	mem, err := memory2.Get()
	if err != nil {
		return nil, err
	}

	memUsed := mem.Used
	memFree := mem.Total - memUsed
	memTotal := mem.Total

	memory := memUsage{
		MemoryUsed:  memUsed,
		MemoryFree:  memFree,
		MemoryTotal: memTotal,
	}

	return &memory, nil
}

func GetDiskRW(dn string) (*diskRW, error) {
	out, err := exec.Command("iostat", "-d", "disk0").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	// I'm only caring for the last line, 2nd column :)
	if len(lines) != 3 {
		return nil, fmt.Errorf("iostat output is too short")
	}
	lastLine := lines[len(lines)-1]
	fields := strings.Fields(lastLine)
	if len(fields) != 3 {
		return nil, fmt.Errorf("iostat output is too short")
	}
	rwValueString := fields[1]
	rwValue, err := strconv.ParseUint(rwValueString, 10, 64)
	if err != nil {
		return nil, err
	}

	// need to convert MB into bytes
	diskRW := diskRW{
		Read:  rwValue,
		Write: rwValue,
	}

	return &diskRW, nil
}
