package process

import (
	"fmt"
	"os"
	"os/exec"
)

// IProcess - general Process Interface
type IProcess interface {
	GetPID() int
	Start() error
	Stop(os.Signal) error
	Kill() error
	IsExited() bool
}

// Process wraps and standardize the actual *exec.Cmd object
type Process struct {
	// Cmd - currently we are using the standard lib
	// maybe replace to our version later [TODO]
	*exec.Cmd
}

// New - init a new process object
func New(name string, args ...string) *Process {
	cmd := exec.Command(name, args...)

	return &Process{
		Cmd: cmd,
	}
}

// GetPID - get PID of **running** process command
func (proc *Process) GetPID() int {
	osProc := proc.Cmd.Process
	if osProc != nil {
		return osProc.Pid
	}
	// if command is not found
	return 0
}

// Start - start the command
func (proc *Process) Start() error {
	return proc.Cmd.Start()
}

// Stop - send linux signal to stop the process
func (proc *Process) Stop(signal os.Signal) error {
	if proc.Cmd.Process != nil {
		return proc.Cmd.Process.Signal(signal)
	}
	return fmt.Errorf("Stop process error - no `Cmd.Process`")
}

// Kill - kill the process
func (proc *Process) Kill() error {
	if proc.Cmd.Process != nil {
		return proc.Cmd.Process.Kill()
	}
	return fmt.Errorf("Kill process error - no `Cmd.Process`")
}

// IsExited - to judge if a process is really exited
func (proc *Process) IsExited() bool {
	procState := proc.Cmd.ProcessState
	if procState != nil {
		return procState.Exited()
	}
	return true
}
