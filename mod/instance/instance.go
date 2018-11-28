package instance

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/AlekSi/pointer"

	"github.com/DemoHn/apm/mod/process"
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

	// command - process object
	command *process.Process
	// state - instance status struct
	status *Status
	// event
	eventHandle *EventHandle
}

// Info shows all informations
type Info struct {
	ID           int
	Name         string
	Status       StatusFlag
	RestartTimes int
	PID          *int
	// CPU time - in ratio * core
	CPU *float64
	// Memory occupied - in bytes
	Memory *int64
	// LaunchTime - in seconds
	LaunchTime *float64
}

// New apm instance (the basic unit of apm management, may contains multiple processes)
func New(path string, args []string) *Instance {
	return &Instance{
		Path:        path,
		Args:        args,
		eventHandle: newEventHandle(),
		// initial status
		status: initStatus(),
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

	status := inst.status
	eventHandle := inst.eventHandle
	// status check
	if status.getStatus() == StatusRunning {
		err = fmt.Errorf("instance has already been started")
		eventHandle.sendEvent(ActionStart, inst, err)
		return
	}
	// init cmd
	cmd := process.New(inst.Path, inst.Args...)
	// TODO for debugging
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	inst.command = cmd
	// auto-start
	fmt.Println("cmd", cmd)
	err = cmd.Start()
	if err != nil {
		eventHandle.sendEvent(ActionStart, inst, err)
		return
	}

	// send start event
	status.setStatus(StatusRunning)
	status.addRestartCounter()
	eventHandle.sendEvent(ActionStart, inst, err)

	err = cmd.Wait()
	// if err = *exec.ExitError, that means the process returned
	// with non-zero value
	status.setStatus(StatusStopped)
	if err == nil {
		eventHandle.sendEvent(ActionStop, inst, nil, 0)
		return
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		exitCode := ws.ExitStatus()
		eventHandle.sendEvent(ActionStop, inst, nil, exitCode)
	} else {
		eventHandle.sendEvent(ActionStop, inst, err)
	}

	// finish and cancel sendEvent() go-rountines
	eventHandle.close()
}

// Stop - stop instance.
// Notice: It will just send a SIGTERM signal to the running process
// and will not stop it immediately.
func (inst *Instance) Stop(signal os.Signal) error {
	status := inst.status
	if status.getStatus() != StatusRunning {
		return fmt.Errorf("[apm] instance is not running, thus stop failed")
	}
	// send stop signal
	err := inst.command.Stop(signal)
	return err
}

// ForceStop - stop the instnace by force
func (inst *Instance) ForceStop() error {
	status := inst.status
	if status.getStatus() != StatusRunning {
		return fmt.Errorf("[apm] instance is not running, thus forceStop failed")
	}
	err := inst.command.Kill()
	return err
}

// GetInfo - get current instance running information
func (inst *Instance) GetInfo() Info {
	status := inst.status
	command := inst.command
	info := Info{
		ID:           inst.ID,
		Name:         inst.Name,
		Status:       status.getStatus(),
		RestartTimes: status.getRestartCounter(),
		PID:          nil,
		CPU:          nil,
		Memory:       nil,
		LaunchTime:   nil,
	}

	if status.getStatus() == StatusRunning {
		// pid
		var pid int
		pid = command.GetPID()
		info.PID = &pid

		pidusage := status.getPidUsage(pid)
		if pidusage != nil {
			info.CPU = pointer.ToFloat64(pidusage.CPU)
			info.Memory = pointer.ToInt64(pidusage.Memory)
			info.LaunchTime = pointer.ToFloat64(pidusage.Elapsed)
		}
	}
	return info
}

// Once - add listener to receive events
func (inst *Instance) Once(topic string) <-chan Event {
	eventHandle := inst.eventHandle
	return eventHandle.Once(topic, emitter.Sync)
}
