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

	select {
	case e := <-inst.Once(instance.ActionStart):
		resp.IsSuccess = true
		resp.InstanceID = e.Int(0)
		resp.PID = e.Int(1)
	case e := <-inst.Once(instance.ActionError):
		resp.IsSuccess = false
		resp.InstanceID = e.Int(0)
		resp.Error = e.String(2)
	}

	return nil
}

// StopInstance -
func (t *Tower) StopInstance(req *StopInstanceRequest, resp *StopInstanceResponse) error {
	master := t.master
	inst, err := master.StopInstance(req.ID)
	if err != nil {
		return err
	}
	select {
	case e := <-inst.Once(instance.ActionStop):
		resp.IsSuccess = true
		resp.InstanceID = e.Int(0)
		resp.ExitCode = e.Int(1)
	case e := <-inst.Once(instance.ActionError):
		resp.IsSuccess = false
		resp.InstanceID = e.Int(0)
		resp.Error = e.String(2)
	}
	return nil
}

// ListInstance -
func (t *Tower) ListInstance(req *ListInstanceRequest, resp *ListInstanceResponse) error {
	master := t.master
	var infos = []instance.Info{}
	insts := master.GetInstancesByFilter(req)
	for _, inst := range insts {
		infos = append(infos, inst.GetInfo())
	}

	resp.InstanceInfos = infos
	return nil
}
