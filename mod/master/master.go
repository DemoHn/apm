package master

import (
	"sync"

	"github.com/DemoHn/apm/mod/instance"
	"github.com/DemoHn/apm/util"
)

// Master - the only one master that controls all instances
type Master struct {
	rpc       *rpcServer
	instances *instanceMap
}

var master *Master

// New -
func New() *Master {
	var once sync.Once
	once.Do(func() {
		if master == nil {
			master = &Master{}
		}
	})

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

	inst := instance.New(prog, args)
	err2 := m.addInstance(req.Name, inst)
	if err2 != nil {
		return nil, err2
	}
	// start instnace
	go func() {
		inst.Run()
	}()
	return inst, nil
}

// Listen to the sockFile
func (m *Master) Listen() error {
	return m.listen()
}

// Shutdown -
func (m *Master) Shutdown() error {
	return nil
}
