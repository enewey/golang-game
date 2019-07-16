package actors

// import (
// 	"enewey.com/golang-game/types"
// )

// // ActorEv - event for actors
// type ActorEv interface {
// 	Command() int
// 	Target() int
// 	Origin() int
// 	Process(*Manager) bool
// }

// // Command Types
// const (
// 	MoveBy = iota
// 	Jump
// )

// // BaseActorEvent root anonymous struct
// type BaseActorEvent struct {
// 	cmd    int
// 	target int
// 	origin int
// }

// // Command w
// func (e *BaseActorEvent) Command() int { return e.cmd }

// // Target w
// func (e *BaseActorEvent) Target() int { return e.target }

// // Origin w
// func (e *BaseActorEvent) Origin() int { return e.origin }

// // MoveByActorEvent w
// type MoveByActorEvent struct {
// 	BaseActorEvent
// 	dx, dy, dz int
// 	duration   types.Frame
// }

// // NewMoveByActorEvent w
// func NewMoveByActorEvent(target, origin, dx, dy, dz int, duration types.Frame) ActorEv {
// 	return &MoveByActorEvent{
// 		BaseActorEvent{MoveBy, target, origin},
// 		dx, dy, dz, duration,
// 	}
// }

// // Process w
// func (ev *MoveByActorEvent) Process(mgr *Manager) bool {
// 	ta := mgr.Actors()[ev.target]
// 	mgr.Actions().Push(
// 		NewMoveByAction(ta, ev.dx, ev.dy, ev.dz, ev.duration),
// 	)

// 	return true
// }

// // JumpActorEvent w
// type JumpActorEvent struct {
// 	BaseActorEvent
// 	upV float64
// }

// // NewJumpActorEvent w
// func NewJumpActorEvent(target, origin int, upV float64) ActorEv {
// 	return &JumpActorEvent{
// 		BaseActorEvent{Jump, target, origin},
// 		upV,
// 	}
// }

// // Process w
// func (ev *JumpActorEvent) Process(mgr *Manager) bool {
// 	ta := mgr.Actors()[ev.target]
// 	mgr.Actions().Push(
// 		NewJumpAction(ta, ev.upV),
// 	)

// 	return true
// }
