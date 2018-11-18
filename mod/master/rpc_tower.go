package master

import (
	"github.com/DemoHn/apm/mod/instance"
)

// Tower receives and handles all incoming commands.Tower
// P.S.: The term "tower" derives from "Control Tower" in aviation.
type Tower struct {
	master *Master
}

// StartInstance - create & run an instance
func (t *Tower) StartInstance(req *StartInstanceRequest, resp *StartInstanceResponse) error {
	master := t.master
	inst, err := master.StartInstance(req)
	if err != nil {
		return err
	}

	// listen to events
	l := inst.NewListener()
	for evt := range l {
		switch v := evt.(type) {
		case instance.StartEvent:
			resp.IsSuccess = true
			resp.InstanceID = evt.InstanceID()
			resp.PID = v.Pid

			return nil
		case instance.ErrorEvent:
			resp.IsSuccess = false
			resp.Error = v.Error
			resp.InstanceID = evt.InstanceID()

			return nil
		}
	}
	return nil
}

// Echo returns the same message - just for testing RPC
func (t *Tower) Echo(input string, output *string) error {
	*output = input
	return nil
}
