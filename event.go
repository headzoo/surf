package surf

const (
	OnError    = "error"
	OnLoad     = "load"
	OnUnload   = "unload"
	OnRequest  = "request"
	OnResponse = "response"
)

type EventArgValues map[string]interface{}

// EventArgs stores arguments to an event.
type EventArgs struct {
	Values         EventArgValues
	Error          error
	stopped        bool
	preventDefault bool
}

// NewEventArgs returns a new *EventArgs instance.
func NewEventArgs(values EventArgValues) *EventArgs {
	return &EventArgs{
		Values: values,
	}
}

// StopPropagation prevents further dispatching of the event.
func (a *EventArgs) StopPropagation() {
	a.stopped = true
}

// IsStopped returns true when StopPropagation() has been called.
func (a *EventArgs) IsStopped() bool {
	return a.stopped
}

// Cancels the event if it is cancelable, without stopping further dispatching of the event.
func (a *EventArgs) PreventDefault() {
	a.preventDefault = false
}

// IsDefaultPrevented returns true when PreventDefault() has been called.
func (a *EventArgs) IsDefaultPrevented() bool {
	return a.preventDefault
}

// GetString returns the value at key as a string.
func (a *EventArgs) GetString(key string) string {
	if v, ok := a.Values[key]; ok {
		return v.(string)
	}
	return ""
}

// GetBool returns the value at key as a bool.
func (a *EventArgs) GetBool(key string) bool {
	if v, ok := a.Values[key]; ok {
		return v.(bool)
	}
	return false
}

// GetInt returns the value at key as an int.
func (a *EventArgs) GetInt(key string) int {
	if v, ok := a.Values[key]; ok {
		return v.(int)
	}
	return 0
}

// GetInt64 returns the value at key as an int64.
func (a *EventArgs) GetInt64(key string) int64 {
	if v, ok := a.Values[key]; ok {
		return v.(int64)
	}
	return 0
}

// GetFloat64 returns the value at key as a float64.
func (a *EventArgs) GetFloat64(key string) float64 {
	if v, ok := a.Values[key]; ok {
		return v.(float64)
	}
	return 0.0
}

// Event stores the details of an event.
type Event struct {
	Name   string
	Target interface{}
	Args   *EventArgs
}

type EventListenerFunc func(e *Event)

// EventTarget represents an object which dispatches events.
type EventTarget struct {
	listeners map[string][]EventListenerFunc
}

// NewEventTarget returns a *EventTarget instance.
func NewEventTarget() *EventTarget {
	return &EventTarget{
		listeners: map[string][]EventListenerFunc{},
	}
}

// AddEventListener registers a listener on the target.
func (t *EventTarget) AddEventListener(event string, fn EventListenerFunc) {
	if _, ok := t.listeners[event]; !ok {
		t.listeners[event] = []EventListenerFunc{}
	}
	t.listeners[event] = append(t.listeners[event], fn)
}

// RemoveEventListener removes a registered listener on the target.
func (t *EventTarget) RemoveEventListener(event string, fn EventListenerFunc) {
	if _, ok := t.listeners[event]; ok {
		for i, f := range t.listeners[event] {
			if &f == &fn {
				t.listeners[event] = append(t.listeners[event][:i], t.listeners[event][i+1:]...)
				break
			}
		}
	}
}

// DispatchEvent dispatches the given event to any registered listeners.
func (t *EventTarget) DispatchEvent(event string, target interface{}, args *EventArgs) error {
	if _, ok := t.listeners[event]; ok {
		if args == nil {
			args = &EventArgs{}
		}
		e := &Event{
			Name:   event,
			Target: target,
			Args:   args,
		}
		debugMessage(`DispatchEvent "%s" to %d listeners`, event, len(t.listeners[event]))
		for _, fn := range t.listeners[event] {
			fn(e)
			if args.IsStopped() {
				break
			}
		}
	}
	return nil
}
