package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ref: https://github.com/soyuka/pidusage
// Current State: only support *ix!

// PidStat - pid stat result
type PidStat struct {
	// PPid = parent pid
	PPid int
	// Process actual pid
	Pid int
	// CPU time in ratio * cores
	CPU float64
	// Memory (maxrss) in bytes
	Memory float64
	// Elapsed time since process start in second
	Elapsed float64
}

// internal cpuTime stat for recording historial data
type cpuTimeStat struct {
	stime  float64
	utime  float64
	uptime float64
}

// PidUsage - main
type PidUsage struct {
	Pid        int
	oldCPUStat *cpuTimeStat
}

// NewPidUsage - new PidUsage object
func NewPidUsage(pid int) *PidUsage {
	return &PidUsage{
		Pid: pid,
	}
}

// GetStat - get pid stat
func (usage *PidUsage) GetStat() *PidStat {
	var availableOS = []string{
		"linux",
		"freebsd",
	}
	for _, os := range availableOS {
		if runtime.GOOS == os {
			stat, err := usage.getStatOnNix()
			if err != nil {
				return nil
			}
			return stat
		}
	}
	log.Printf("OS:%s is not supported to get stat now!", runtime.GOOS)
	return nil
}

// internal functions
// get stat On *ix like systems (linux, unix)
func (usage *PidUsage) getStatOnNix() (*PidStat, error) {
	// TODO - new method to make memory RSS more precise!
	cpu := 0.0

	nInfo, err := parseProcfile(usage.Pid)
	if err != nil {
		return nil, err
	}
	cInfo, err := parseCPUInfo()
	if err != nil {
		return nil, err
	}

	// calculate cpu
	var total = 0.0
	var seconds = 0.0
	oldStat := usage.oldCPUStat
	childrens := nInfo.cutime + nInfo.cstime

	if oldStat != nil {
		total = (nInfo.stime - oldStat.stime + nInfo.utime - oldStat.utime + childrens) / cInfo.clockTick
		seconds = math.Abs(float64(cInfo.uptime - oldStat.uptime))
	} else {
		total = (nInfo.stime + nInfo.utime + childrens) / cInfo.clockTick
		seconds = math.Abs(nInfo.start/cInfo.clockTick - cInfo.uptime)
	}

	if seconds > 0 {
		cpu = float64(total) / seconds
	}

	// update oldstat
	usage.oldCPUStat = &cpuTimeStat{
		utime:  nInfo.utime,
		stime:  nInfo.stime,
		uptime: cInfo.uptime,
	}
	return &PidStat{
		PPid:    int(nInfo.ppid),
		Pid:     usage.Pid,
		CPU:     cpu,
		Memory:  nInfo.rss * cInfo.pageSize, // TODO: more precise calculation!
		Elapsed: cInfo.uptime - nInfo.start/cInfo.clockTick,
	}, nil
}

// CPU
// read procfile info

type procfileInfo struct {
	ppid   float64
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	start  float64
	rss    float64
}

func parseProcfile(pid int) (*procfileInfo, error) {
	procFile := fmt.Sprintf("/proc/%d/stat", pid)

	data, err := ioutil.ReadFile(procFile)
	// e.g
	if err != nil {
		return nil, err
	}

	pInfo := &procfileInfo{}
	// parse proc file
	procstr := string(data[:])
	lastP := strings.LastIndex(procstr, ")")
	infos := strings.Split(procstr[lastP+2:], " ")

	for index, info := range infos {
		n, _ := strconv.ParseFloat(info, 64)
		switch index {
		case 1:
			pInfo.ppid = n
		case 11:
			pInfo.utime = n
		case 12:
			pInfo.stime = n
		case 13:
			pInfo.cutime = n
		case 14:
			pInfo.cstime = n
		case 19:
			pInfo.start = n
		case 21:
			pInfo.rss = n
		}
	}

	return pInfo, nil
}

type procCPUInfo struct {
	uptime    float64
	pageSize  float64
	clockTick float64
}

func parseCPUInfo() (*procCPUInfo, error) {
	var clkTck float64 = 100
	var pageSize float64 = 4096

	uptimeFileBytes, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return nil, err
	}
	uptime, _ := strconv.ParseFloat(strings.Split(string(uptimeFileBytes), " ")[0], 64)
	clkTckStdout, err := exec.Command("getconf", "CLK_TCK").Output()
	if err == nil {
		clkTck, _ = strconv.ParseFloat(strings.Split(string(clkTckStdout), "\n")[0], 64)
	}

	pageSizeStdout, err := exec.Command("getconf", "PAGESIZE").Output()
	if err == nil {
		pageSize, _ = strconv.ParseFloat(strings.Split(string(pageSizeStdout), "\n")[0], 64)
	}

	return &procCPUInfo{
		uptime:    uptime,
		clockTick: clkTck,
		pageSize:  pageSize,
	}, nil
}
