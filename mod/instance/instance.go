package instance

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/olebedev/emitter"
)

// Instance type defines a general instance type
type Instance struct {
	// ID - instance ID
	ID int
	// Name - instance name
	Name string
	// Path - instance command path
	Path string
	// Args - command arguments
	Args []string
	// Command - command instance
	Command *exec.Cmd
	// Status - instance status READY | RUNNING | STOPPED
	status string
	// event
	eventHandle *EventHandle
	// mutex
	sync.RWMutex
}

// New apm instance (the basic unit of apm management, may contains multiple processes)
func New(path string, args []string) *Instance {
	return &Instance{
		Path:        path,
		Args:        args,
		eventHandle: newEventHandle(),
		// initial status
		status: statusReady,
	}
}

// setters

// SetID -
func (inst *Instance) SetID(id int) {
	inst.ID = id
}

// SetName -
func (inst *Instance) SetName(name string) {
	inst.Name = name
}

// Run a instance with events registered
func (inst *Instance) Run() {
	var err error

	eventHandle := inst.eventHandle
	// init cmd
	cmd := exec.Command(inst.Path, inst.Args...)
	// TODO for debugging
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	inst.Command = cmd
	// auto-start
	err = cmd.Start()
	eventHandle.SendEvent(ActionStart, inst, err)

	if err == nil {
		setStatus(inst, statusRunning)
		err = cmd.Wait()
		// if err = *exec.ExitError, that means the process returned
		// with non-zero value
		if err == nil {
			eventHandle.SendEvent(ActionStop, inst, nil, 0)
		} else {
			if exitError, ok := err.(*exec.ExitError); ok {
				ws := exitError.Sys().(syscall.WaitStatus)
				exitCode := ws.ExitStatus()
				eventHandle.SendEvent(ActionStop, inst, nil, exitCode)
			} else {
				eventHandle.SendEvent(ActionStop, inst, err)
			}
		}
	}

	setStatus(inst, statusStopped)
	// finish and cancel sendEvent() go-rountines
	eventHandle.Close()
}

// Stop - stop instance.
// Notice: It will just send a SIGTERM signal to the running process
// and will not stop it immediately.
func (inst *Instance) Stop(signal os.Signal) error {
	if inst.status != statusRunning {
		return fmt.Errorf("[apm] instance is not running, thus stop failed")
	}
	// send stop signal
	err := inst.Command.Process.Signal(signal)
	return err
}

// ForceStop - stop the instnace by force
func (inst *Instance) ForceStop() error {
	if inst.status != statusRunning {
		return fmt.Errorf("[apm] instance is not running, thus forceStop failed")
	}
	err := inst.Command.Process.Kill()
	return err
}

// Once - add listener to receive events
func (inst *Instance) Once(topic string) <-chan Event {
	return inst.eventHandle.emitter.Once(topic, emitter.Sync)
}

// GetStatus - get status
func (inst *Instance) GetStatus() Status {
	return getStatus(inst)
}
