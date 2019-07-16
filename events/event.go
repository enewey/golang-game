package events

import "enewey.com/golang-game/actors"

const (
	Global = iota
	Actor
	Window
)

var hub *EventHub

func init() {
	hub = &EventHub{
		[]ActorEv{},
	}
}
func Hub() *EventHub { return hub }

type EventHub struct {
	actorQ ActorEventQueue
	// globalQ []*GlobalEv
	// windowQ []*WindowEv
}

func (h *EventHub) ActorEvents() ActorEventQueue { return h.actorQ }

// Flush - clear out all events from the event hub
func (h *EventHub) Flush() {
	h.actorQ = []ActorEv{}
}

type ActorEventQueue []ActorEv

// Enqueue - queue up an event to be processed on the next tick
func (aq ActorEventQueue) Enqueue(a ActorEv) {
	aq = append(aq, a)
}

// Read - reads the next event in the queue
func (aq ActorEventQueue) Read() ActorEv {
	if len(aq) == 0 {
		return nil
	}
	pop := aq[0]
	aq = aq[1:]
	return pop
}

func (aq ActorEventQueue) HasNext() bool { return len(aq) > 0 }

type ActorEv interface {
	Command() int
	Target() int
	Origin() int
	Process(*actors.Manager) bool
}
