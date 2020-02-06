package events

// Event Scope categories
const (
	Global = iota
	Actor
	Window
)

// Event - generic event that can originate from anywhere
type Event struct {
	// scope is the Global/Actor/Window iota
	scope int
	// command is specific to the scope
	code int
	// payload contains the parameters required of the command
	payload []interface{}
}

// New - create a new event struct
func New(scope, code int, payload []interface{}) *Event {
	return &Event{scope, code, payload}
}

// Scope - Global (0), Actor (1) or Window (2)
func (e *Event) Scope() int { return e.scope }

// Code - Coded event identifier, describing what happened based on the scope
func (e *Event) Code() int { return e.code }

// Payload - parameters/arguments for the event
func (e *Event) Payload() []interface{} { return e.payload }

var bus []*Event

func init() {
	bus = []*Event{}
}

// Bus w
func Bus() EventBus { return bus }

// EventBus wo
type EventBus []*Event

// Flush - clear out all events from the event hub
func Flush() {
	bus = []*Event{}
}

// Enqueue - queue up an event to be processed on the next tick
func Enqueue(ev *Event) {
	bus = append(bus, ev)
}

// EnqueueAll - queue up multiple events, with the assumption that they'll be executed in the order of events
// This requires reversing the input array, so that the array is appended to provide a FIFO ordering
func EnqueueAll(evs []*Event) {
	a := append(evs[:0:0], evs...)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	bus = append(bus, a...)
}

// Read - reads the next event in the queue
func Read() *Event {
	if len(bus) == 0 {
		return nil
	}
	pop := bus[0]
	bus = bus[1:]
	return pop
}

// HasNext - bool
func HasNext() bool { return len(bus) > 0 }
