package master

import (
	"github.com/DemoHn/apm/mod/instance"
)

// StartInstanceRequest defines the input parameters of `Tower.StartInstance`
type StartInstanceRequest struct {
	Name    string
	Command string
	ID      *int
}

// StartInstanceResponse defines the reply of `Tower.StartInstance`
type StartInstanceResponse struct {
	IsSuccess  bool
	Error      string
	InstanceID int
	PID        int
}

// StopInstanceRequest defines the input parameters of `Tower.StopInstance`
type StopInstanceRequest struct {
	ID int
}

// StopInstanceResponse defines the reply of `Tower.StopInstance`
type StopInstanceResponse struct {
	IsSuccess  bool
	ExitCode   int
	Error      string
	InstanceID int
}

// ListInstanceRequest defines the payload
type ListInstanceRequest struct {
	// ID - [optional]
	ID   *int
	Name *string
}

// ListInstanceResponse - response
type ListInstanceResponse struct {
	InstanceInfos []instance.Info
}
