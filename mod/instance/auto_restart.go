package instance

import (
	"sync"
	"time"
)

const (
	defaultInterval = 3 * time.Second // 3s
)

// AutoRestartHandle - handle all issues about auto restart of the instance
type AutoRestartHandle struct {
	// AutoRestartInterval - the delay time to restart after which the process exited
	// This is to prevent restart too frequent
	interval time.Duration
	// the original instance
	instance    *Instance
	maskLock    bool
	restartLock bool
	mu          sync.RWMutex
}

func newAutoRestartHandle() *AutoRestartHandle {
	return &AutoRestartHandle{
		interval:    defaultInterval,
		maskLock:    false,
		restartLock: false,
	}
}

// setters
func (ar *AutoRestartHandle) setInterval(interval time.Duration) {
	ar.interval = interval
}

// Tick - trigger restart operation
func (ar *AutoRestartHandle) tick(inst *Instance) {
	autoRestart := inst.AutoRestart
	if autoRestart || ar.restartLock {
		// release restart lock
		ar.unforceRestart()
		go func() {
			<-time.After(ar.interval)
			if ar.maskLock == false {
				inst.Run()
			} else {
				ar.unmask()
			}
		}()
	}
}

// Mask - hide auto-restart operation temperaily
// It will work only once.
// This is usually used in restart operation
func (ar *AutoRestartHandle) mask() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.maskLock = true
}

// Unmask - enable auto-restart Tick again
func (ar *AutoRestartHandle) unmask() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.maskLock = false
}

func (ar *AutoRestartHandle) forceRestart() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.restartLock = true
}

func (ar *AutoRestartHandle) unforceRestart() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.restartLock = false
}
