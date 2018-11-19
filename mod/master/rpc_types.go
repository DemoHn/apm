package master

// StartInstanceRequest defines the input parameters of `Tower.StartInstance`
type StartInstanceRequest struct {
	Name    string
	Command string
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
	Error      string
	IsSuccess  bool
	InstanceID int
}
