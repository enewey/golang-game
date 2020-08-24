package actors

import (
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/types"
)

// This file contains a bunch of prefabricated events to be queued.
// They are not necessarily all actor-scoped events.

// ==== Actor Events

// NewMoveByEvent w
func NewMoveByEvent(target Actor, dx, dy, dz float64, duration types.Frame) *events.Event {
	return events.New(events.Actor, MoveByActionType, []interface{}{target, dx, dy, dz, duration})
}

// NewDashEvent creates an event that interprets as an actor Dash event
func NewDashEvent(target Actor, x, y, z float64) *events.Event {
	return events.New(events.Actor, DashActionType, []interface{}{target, x, y, z})
}

// NewJumpEvent creates an event that interprets as an actor Jump event
func NewJumpEvent(target Actor, jump float64) *events.Event {
	return events.New(events.Actor, JumpActionType, []interface{}{target, jump})
}

// ==== Global Events

// NewInteractEvent creates an event to be interpreted at a global level
func NewInteractEvent(subject Actor) *events.Event {
	return events.New(0, 0, []interface{}{subject})
}
