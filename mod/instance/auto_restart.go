package instance

import (
	"time"
)

const (
	defaultInterval = 3 * time.Second // 3s
)

// AutoRestartHandle - handle all issues about auto restart of the instance
type AutoRestartHandle struct {
	Enable bool
	// AutoRestartInterval - the delay time to restart after which the process exited
	// This is to prevent restart too frequent
	interval time.Duration
	// the original instance
	instance *Instance
	maskLock bool
}

func newAutoRestartHandle(enable bool) *AutoRestartHandle {
	return &AutoRestartHandle{
		Enable:   enable,
		interval: defaultInterval,
		instance: nil,
		maskLock: false,
	}
}

// setters
func (ar *AutoRestartHandle) setInstance(inst *Instance) {
	ar.instance = inst
}

func (ar *AutoRestartHandle) setInterval(interval time.Duration) {
	ar.interval = interval
}

// Tick - trigger restart operation
func (ar *AutoRestartHandle) tick() {
	if ar.Enable && ar.maskLock == false {
		go func() {
			<-time.After(ar.interval)
			ar.instance.Run()
		}()
	}
}

// Mask - hide auto-restart operation temperaily
// It will work only once.
// This is usually used in restart operation
func (ar *AutoRestartHandle) mask() {
	ar.maskLock = true
}

// Unmask - enable auto-restart Tick again
func (ar *AutoRestartHandle) unmask() {
	ar.maskLock = false
}
