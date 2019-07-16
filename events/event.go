package events

// import (
// 	"enewey.com/golang-game/actors"
// )

// // Event types
// const (
// 	Global = iota
// 	Actor
// 	Window
// )

// var hub *EventHub

// func init() {
// 	hub = &EventHub{
// 		[]ActorEv{},
// 	}
// }

// // Hub w
// func Hub() *EventHub { return hub }

// // EventHub wo
// type EventHub struct {
// 	actorQ ActorEventQueue
// 	// globalQ []*GlobalEv
// 	// windowQ []*WindowEv
// }

// // ActorEvents a
// func (h *EventHub) ActorEvents() ActorEventQueue { return h.actorQ }

// // Flush - clear out all events from the event hub
// func (h *EventHub) Flush() {
// 	h.actorQ = []ActorEv{}
// }

// // ActorEventQueue w
// type ActorEventQueue []actors.ActorEv

// // Enqueue - queue up an event to be processed on the next tick
// func (aq ActorEventQueue) Enqueue(a ActorEv) {
// 	aq = append(aq, a)
// }

// // Read - reads the next event in the queue
// func (aq ActorEventQueue) Read() ActorEv {
// 	if len(aq) == 0 {
// 		return nil
// 	}
// 	pop := aq[0]
// 	aq = aq[1:]
// 	return pop
// }

// // HasNext - bool
// func (aq ActorEventQueue) HasNext() bool { return len(aq) > 0 }
