package master

// StartInstanceRequest defines the input parameters of `Tower.StartInstance`
type StartInstanceRequest struct {
	Name    string
	Command string
}

// StartInstanceResponse defines the reply of `Tower.StartInstance`
type StartInstanceResponse struct {
	IsSuccess  bool
	Error      error
	InstanceID int
	PID        int
}
