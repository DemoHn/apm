package instance

import (
	"sync"
)

// Status - determine the actual status
type Status = string

// StatusInfo - get current status (including cpu, memory, restart time)
type StatusInfo struct {
	status         Status
	firstStart     bool
	restartCounter int
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
func (inst *Instance) getCPUUsage() float64 {
	// only running instance could get CPU Info
	if inst.getStatus() == StatusRunning {
		// use specific code to read CPU usage
	}

	return 0.0
}

// memory
func (inst *Instance) getMemory() int64 {
	return 0
}
