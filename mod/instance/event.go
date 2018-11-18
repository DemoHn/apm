package instance

// EventConfig defines all switches of a command that triggers
// the corresponding event
type EventConfig struct {
	OnStart      bool
	OnStop       bool
	OnStdinData  bool
	OnStdoutData bool
	OnStderrData bool
}

// Event is an interface
type Event interface {
	// InstanceID - get instancce ID
	InstanceID() int
	// Name - get event (Action) name
	Name() string
}

// EventHandle manages all events
type EventHandle struct {
	listeners []chan Event
	config    *EventConfig
}

// Action -
type Action = string

const (
	ActionStart          Action = "ACTION_START"
	ActionStop           Action = "ACTION_STOP"
	ActionRestart        Action = "ACTION_RESTART"
	ActionEncounterError Action = "ACTION_ENCOUNTER_ERROR"
)

// StartEvent -
type StartEvent struct {
	instanceID int
	Pid        int
}

// Name - get event name
func (evt StartEvent) Name() string {
	return ActionStart
}

// InstanceID -
func (evt StartEvent) InstanceID() int {
	return evt.instanceID
}

// StopEvent -
type StopEvent struct {
	ExitCode   int
	instanceID int
}

// Name - get event name
func (evt StopEvent) Name() string {
	return ActionStop
}

// InstanceID -
func (evt StopEvent) InstanceID() int {
	return evt.instanceID
}

// ErrorEvent -
type ErrorEvent struct {
	instanceID int
	Action     string
	Error      error
}

// Name - get event name
func (evt ErrorEvent) Name() string {
	return ActionEncounterError
}

// InstanceID -
func (evt ErrorEvent) InstanceID() int {
	return evt.instanceID
}

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
		config:    defaultEventConfig,
		listeners: make([]chan Event, 0),
	}

	return eventHandle
}

func (handle *EventHandle) newListener() <-chan Event {
	listener := make(chan Event)
	handle.listeners = append(handle.listeners, listener)

	return listener
}

func (handle *EventHandle) closeAll() {
	for _, cs := range handle.listeners {
		close(cs)
	}
}

// sendEvent - send corresponding event to instance
func (handle *EventHandle) sendEvent(Action Action, instance *Instance, err error) {
	var evt Event
	evt = nil

	// send error event
	if err != nil {
		evt = ErrorEvent{
			instanceID: instance.ID,
			Action:     Action,
			Error:      err,
		}
	} else {
		// send concrete events
		switch Action {
		case ActionStart:
			// make sure the command EXISTS!
			pid := instance.Command.Process.Pid
			evt = StartEvent{
				instanceID: instance.ID,
				Pid:        pid,
			}
		case ActionStop:
			// TODO: how to get the exit code?
			exitCode := 0
			evt = StopEvent{
				instanceID: instance.ID,
				ExitCode:   exitCode,
			}
		}
	}

	// fan-out to all event listeners
	for _, cs := range handle.listeners {
		if handle.filterEvent(evt) == true {
			cs <- evt
		}
	}
}

func (handle *EventHandle) filterEvent(evt Event) bool {
	switch evt.Name() {
	case ActionStart:
		return handle.config.OnStart
	case ActionStop:
		return handle.config.OnStop
	case ActionEncounterError:
		return true
	}

	return false
}
