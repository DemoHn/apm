package instance

import (
	"sync"

	"github.com/DemoHn/apm/util"
)

// Status - determine the actual status
type Status = string

// StatusInfo - get current status (including cpu, memory, restart time)
type StatusInfo struct {
	status         Status
	firstStart     bool
	restartCounter int
	pidusage       *util.PidUsage
	// read-write lock
	mu sync.RWMutex
}

const (
	// StatusReady - the instance has not started yet
	StatusReady Status = "status_ready"
	// StatusRunning - the instance is running
	StatusRunning Status = "status_running"
	// StatusStopped - the instance has stopped (by signal or program is down)
	StatusStopped Status = "status_stopped"
)

func initStatusInfo() *StatusInfo {
	return &StatusInfo{
		status:         StatusReady,
		firstStart:     false,
		restartCounter: 0,
	}
}

// getStatus
func (inst *Instance) setStatus(status Status) {
	s := inst.status
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = status
}

func (inst *Instance) getStatus() Status {
	s := inst.status
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

// restartCounter
func (inst *Instance) addRestartCounter() {
	s := inst.status
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.firstStart == false {
		s.firstStart = true
	} else {
		s.restartCounter = s.restartCounter + 1
	}
}

func (inst *Instance) getRestartCounter() int {
	s := inst.status
	s.mu.RLock()
	defer s.mu.RLock()

	return inst.status.restartCounter
}

// getCPUUsage - stat current process' CPU time
func (inst *Instance) getPidUsage(pid int) *util.PidStat {
	s := inst.status
	// only running instance could get Pid Usage
	if inst.getStatus() == StatusRunning {
		if s.pidusage != nil && s.pidusage.Pid == pid {
			// get recent stat
			return s.pidusage.GetStat()
		}
		// new pidUsage object
		s.pidusage = util.NewPidUsage(pid)
		return s.pidusage.GetStat()
	}

	return nil
}
