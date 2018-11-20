package instance

import (
	"github.com/olebedev/emitter"
)

// EventConfig defines all switches of a command that triggers
// the corresponding event
type EventConfig struct {
	OnStart      bool
	OnStop       bool
	OnStdinData  bool
	OnStdoutData bool
	OnStderrData bool
}

// EventHandle manages all events
type EventHandle struct {
	config  *EventConfig
	emitter *emitter.Emitter
}

// Event -
type Event = emitter.Event

// Action -
type Action = string

const (
	// ActionStart -
	ActionStart Action = "action_start"
	// ActionStop -
	ActionStop Action = "action_stop"
	// ActionRestart -
	ActionRestart Action = "action_restart"
	// ActionError -
	ActionError Action = "action_error"
)

// methods
func newEventHandle() *EventHandle {
	defaultEventConfig := &EventConfig{
		OnStart:      true,
		OnStop:       true,
		OnStdinData:  false,
		OnStdoutData: false,
		OnStderrData: false,
	}

	eventHandle := &EventHandle{
		config:  defaultEventConfig,
		emitter: &emitter.Emitter{},
	}

	return eventHandle
}

// Close - close all event listeners
func (handle *EventHandle) Close() {
	handle.emitter.Off("*")
}

// SendEvent - send corresponding event to instance
func (handle *EventHandle) SendEvent(action Action, inst *Instance, err error, args ...interface{}) {
	emitter := handle.emitter
	// send error event
	if err != nil {
		// send error event with no other reasons
		// params: [id, action_name, error]
		emitter.Emit(ActionError, inst.ID, action, err.Error())
	} else {
		// send concrete events
		switch action {
		case ActionStart:
			// make sure the command EXISTS!
			pid := inst.Command.Process.Pid
			// params: [id, pid]
			emitter.Emit(ActionStart, inst.ID, pid)
		case ActionStop:
			exitCode := 0
			exitArg := args[0]
			if code, ok := exitArg.(int); ok {
				exitCode = code
			}
			emitter.Emit(ActionStop, inst.ID, exitCode)
		}
	}
}