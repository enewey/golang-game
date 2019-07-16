package events

import (
	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/types"
)

const (
	MoveByActor = iota
	JumpActor
)

type BaseActorEvent struct {
	cmd    int
	target int
	origin int
}

func (e *BaseActorEvent) Command() int { return e.cmd }
func (e *BaseActorEvent) Target() int  { return e.target }
func (e *BaseActorEvent) Origin() int  { return e.origin }

type MoveByActorEvent struct {
	BaseActorEvent
	dx, dy, dz int
	duration   types.Frame
}

func NewMoveByActorEvent(target, origin, dx, dy, dz int, duration types.Frame) ActorEv {
	return &MoveByActorEvent{
		BaseActorEvent{MoveByActor, target, origin},
		dx, dy, dz, duration,
	}
}

func (ev *MoveByActorEvent) Process(mgr *actors.Manager) bool {
	ta := mgr.Actors()[ev.target]
	mgr.Actions().Push(
		actors.NewMoveByAction(ta, ev.dx, ev.dy, ev.dz, ev.duration),
	)

	return true
}

type JumpActorEvent struct {
	BaseActorEvent
	upV float64
}

func NewJumpActorEvent(target, origin int, upV float64) ActorEv {
	return &JumpActorEvent{
		BaseActorEvent{JumpActor, target, origin},
		upV,
	}
}

func (ev *JumpActorEvent) Process(mgr *actors.Manager) bool {
	ta := mgr.Actors()[ev.target]
	mgr.Actions().Push(
		actors.NewJumpAction(ta, ev.upV),
	)

	return true
}
