package instance

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/AlekSi/pointer"
	"github.com/DemoHn/apm/mod/logger"
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
	// AutoRestart - enable "AutoRestart" feature or not
	AutoRestart bool
	// command - process object
	command           process.IProcess
	status            *Status
	eventHandle       *EventHandle
	autoRestartHandle *AutoRestartHandle
}

var log *logger.Logger

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
	log = logger.Get()
	inst := &Instance{
		Path:        path,
		Args:        args,
		eventHandle: newEventHandle(),
		// TODO: more config
		AutoRestart: true,
		// initial status
		status:            initStatus(),
		autoRestartHandle: newAutoRestartHandle(),
	}

	return inst
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
func (inst *Instance) Run() error {
	/*autoRestartHandle := inst.autoRestartHandle
	status := inst.status

		autoRestartOn := inst.AutoRestart || autoRestartHandle.restartLock
		if autoRestartOn && status.getStatus() != StatusReady ||
			!autoRestartOn && status.getStatus() == StatusRunning {
			return fmt.Errorf("duplicate Run()")
		}*/
	go func() {
		inst.spawnProcess()
		inst.eventHandle.close()
	}()

	return nil
}
func (inst *Instance) spawnProcess() {
	var err error

	status := inst.status
	eventHandle := inst.eventHandle
	autoRestartHandle := inst.autoRestartHandle

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
	// unmark to enable auto-restart again
	autoRestartHandle.unmask()

	log.Debugf("[apm] ID(%d) going to start", inst.ID)
	err = cmd.Start()
	if err != nil {
		log.Debugf("[apm] ID(%d) start failed err=%s", inst.ID, err)
		eventHandle.sendEvent(ActionStart, inst, err)
		autoRestartHandle.tick(inst)
		return
	}

	// send start event
	status.setStatus(StatusRunning)
	status.addRestartCounter()
	eventHandle.sendEvent(ActionStart, inst, err)
	log.Debugf("[apm] ID(%d) instance is running", inst.ID)
	err = cmd.Wait()

	log.Debugf("[apm] ID(%d) going to stop", inst.ID)
	status.setStatus(StatusStopped)

	autoRestartHandle.tick(inst)
	if err == nil {
		log.Debugf("[apm] ID(%d) stop succeed", inst.ID)
		eventHandle.sendEvent(ActionStop, inst, nil, 0)
		return
	}

	log.Debugf("[apm] ID(%d) stop with err=%s", inst.ID, err)
	// if err = *exec.ExitError, that means the process returned
	// with non-zero value
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		exitCode := ws.ExitStatus()
		eventHandle.sendEvent(ActionStop, inst, nil, exitCode)
	} else {
		eventHandle.sendEvent(ActionStop, inst, err)
	}
}

// Stop - stop instance.
// Notice: It will just send a SIGTERM signal to the running process
// and will not stop it immediately.
func (inst *Instance) Stop(signal os.Signal) error {
	autoRestartHandle := inst.autoRestartHandle

	// send stop signal
	autoRestartHandle.mask()
	err := inst.command.Stop(signal)
	return err
}

// ForceStop - stop the instance by force
func (inst *Instance) ForceStop() error {
	autoRestartHandle := inst.autoRestartHandle

	autoRestartHandle.mask()
	err := inst.command.Kill()
	return err
}

// Restart - restart the instance
func (inst *Instance) Restart(signal os.Signal) error {
	autoRestartHandle := inst.autoRestartHandle
	// acquire restart lock to make auto-restart work by force
	// it will be automatically released after tick()
	autoRestartHandle.forceRestart()

	return inst.command.Stop(signal)
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

	if info.Status == StatusRunning {
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
