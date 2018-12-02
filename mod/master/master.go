package master

import (
	"syscall"

	"github.com/DemoHn/apm/mod/instance"
	"github.com/DemoHn/apm/util"
	// loggers
	"github.com/DemoHn/apm/mod/logger"
)

// Master - the only one master that controls all instances
type Master struct {
	debugMode bool
	rpc       *rpcServer
	instances *instanceMap
}

var master *Master

// New -
func New(debugMode bool) *Master {
	master = &Master{
		debugMode: debugMode,
	}
	// new global logger
	logger.Init(debugMode)
	return master
}

// Init -
func (m *Master) Init(sockFile string) error {
	var err error
	// init RPC first
	err = m.initRPC(sockFile)
	if err != nil {
		return err
	}

	// init instance map to add/del instances
	err = m.initInstanceMap()
	if err != nil {
		return err
	}

	return nil
}

// StartInstance - create & start instance
func (m *Master) StartInstance(req *StartInstanceRequest) (*instance.Instance, error) {
	prog, args, err := util.SplitCommand(req.Command)
	if err != nil {
		return nil, err
	}
	// create instance
	inst := instance.New(prog, args)
	err2 := m.addInstance(req.Name, inst)
	if err2 != nil {
		return nil, err2
	}
	// start instnace - non-blocking
	inst.Run()
	return inst, nil
}

// StopInstance - stop instance
// Notice: still should wait for
func (m *Master) StopInstance(id int) (*instance.Instance, error) {
	inst, err := m.findInstance(id)
	if err != nil {
		return nil, err
	}

	err2 := inst.Stop(syscall.SIGTERM)
	if err2 != nil {
		return nil, err2
	}
	return inst, nil
}

// GetOneInstance - get instance
func (m *Master) GetOneInstance(id int) *instance.Instance {
	inst, err := m.findInstance(id)
	if err != nil {
		return nil
	}

	return inst
}

// GetInstancesByFilter -
func (m *Master) GetInstancesByFilter(req *ListInstanceRequest) []*instance.Instance {
	return m.findInstancesByFilter(req.ID, req.Name)
}

// Listen to the sockFile
func (m *Master) Listen() error {
	return m.listen()
}

// Shutdown -
func (m *Master) Shutdown() error {
	return nil
}
