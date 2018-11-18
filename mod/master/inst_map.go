package master

import (
	"fmt"

	"github.com/DemoHn/apm/mod/instance"
)

type instanceMap struct {
	counter int
	// instance name -> id map
	nameMap map[string][]int
	// instance id -> obj map
	instanceMap map[int]*instance.Instance
}

func (m *Master) initInstanceMap() error {
	m.instances = &instanceMap{
		counter:     0,
		nameMap:     make(map[string][]int),
		instanceMap: make(map[int]*instance.Instance),
	}
	return nil
}

func (m *Master) addInstance(name string, instance *instance.Instance) error {
	if m.instances == nil {
		return fmt.Errorf("master.instances not found! initInstanceMap() did executed?")
	}

	// new instance
	insts := m.instances
	insts.counter = insts.counter + 1
	id := insts.counter
	// set instance
	instance.SetID(id)
	instance.SetName(name)

	// insert instance to map
	insts.instanceMap[id] = instance

	// add it to nameMap
	if _, ok := insts.nameMap[name]; ok == false {
		insts.nameMap[name] = []int{}
	}
	insts.nameMap[name] = append(insts.nameMap[name], id)
	return nil
}

func (m *Master) removeInstance(id int) error {
	if m.instances == nil {
		return fmt.Errorf("master.instances not found! initInstanceMap() did executed?")
	}

	// find obj first
	insts := m.instances
	if inst, ok := insts.instanceMap[id]; ok == true {
		name := inst.Name

		// remove obj from instanceMap
		delete(insts.instanceMap, id)
		// remove id from nameMap
		nArr := []int{}
		for _, item := range insts.nameMap[name] {
			if item != id {
				nArr = append(nArr, item)
			}
		}
		insts.nameMap[name] = nArr
		if len(insts.nameMap[name]) == 0 {
			delete(insts.nameMap, name)
		}
	}
	return nil
}
