package instance

import (
	"os"
	"os/exec"
	"sync"
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
	eventHandle.sendEvent(ActionStart, inst, err)

	if err == nil {
		setStatus(inst, statusRunning)
		err = cmd.Wait()
		eventHandle.sendEvent(ActionStop, inst, err)
	}

	setStatus(inst, statusStopped)
	// finish and cancel sendEvent() go-rountines
	eventHandle.closeAll()
}

// NewListener - add listener to receive events
func (inst *Instance) NewListener() <-chan Event {
	eventHandle := inst.eventHandle
	return eventHandle.newListener()
}

// GetStatus - get status
func (inst *Instance) GetStatus() Status {
	return getStatus(inst)
}
