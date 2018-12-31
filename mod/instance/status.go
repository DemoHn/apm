package instance

import (
	"sync"

	"github.com/DemoHn/apm/infra/logger"
	"github.com/DemoHn/apm/util"
	deadlock "github.com/sasha-s/go-deadlock"
)

// StatusFlag - determine the actual status
type StatusFlag = string

type rwLocker interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

// Status - get current status (including cpu, memory, restart time)
type Status struct {
	flag           StatusFlag
	firstStart     bool
	restartCounter int
	pidusage       util.IPidUsage
	// read-write lock
	mu rwLocker
}

const (
	// StatusReady - the instance has not started yet
	StatusReady StatusFlag = "status_ready"
	// StatusRunning - the instance is running
	StatusRunning StatusFlag = "status_running"
	// StatusStopped - the instance has stopped (by signal or program is down)
	StatusStopped StatusFlag = "status_stopped"
)

func initStatus() *Status {
	s := &Status{
		flag:           StatusReady,
		firstStart:     false,
		restartCounter: 0,
	}

	if logger.DebugMode() {
		s.mu = new(deadlock.RWMutex)
	} else {
		s.mu = new(sync.RWMutex)
	}
	return s
}

// IsRunning - whether current status is `StatusRunning`
func (s *Status) IsRunning() bool {
	return s.getStatus() == StatusRunning
}

// IsReady - whether current status is `StatusReady`
func (s *Status) IsReady() bool {
	return s.getStatus() == StatusReady
}

// IsStopped - whether current status is `StatusStopped`
func (s *Status) IsStopped() bool {
	return s.getStatus() == StatusStopped
}

// internal functions
// getStatus
func (s *Status) setStatus(flag StatusFlag) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flag = flag
}

func (s *Status) getStatus() StatusFlag {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.flag
}

// restartCounter
func (s *Status) addRestartCounter() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.firstStart == false {
		s.firstStart = true
	} else {
		s.restartCounter = s.restartCounter + 1
	}
}

func (s *Status) getRestartCounter() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.restartCounter
}

// getPidUsage - stat current process' CPU time
func (s *Status) getPidUsage(pid int) *util.PidStat {
	// only running instance could get Pid Usage
	if s.getStatus() == StatusRunning {
		if s.pidusage != nil && s.pidusage.GetPid() == pid {
			// get recent stat
			return s.pidusage.GetStat()
		}
		// new pidUsage object
		s.pidusage = util.NewPidUsage(pid)
		return s.pidusage.GetStat()
	}

	return nil
}
