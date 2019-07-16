package scene

import (
	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/room"
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
)

// Scene -	coordinates window, actor, and room entities.
// 			processes inputs, delegates queued events, triggers actions, and
//			resolves collisions.
type Scene struct {
	actorM *actors.Manager
	room   *room.Room
}

func New(actorM *actors.Manager, room *room.Room) *Scene {
	return &Scene{actorM, room}
}

func (s *Scene) Update(df types.Frame) {
	// first process inputs
	input.State().Tick(df)
	// then process/delegate events

	// then call the manager act() functions

	// resolve collisions of actor against room based on staged actor velocities

	// allow actor manager to resolve collisions between actors
	// (which may generate more events)

	// finally render

}

func (s *Scene) processEvents() {

}

func (s *Scene) act() {

}

func (s *Scene) resolveCollisions() {

}

func queueEventFromInput(state input.Input) {
	var dx, dy int
	if state[ebiten.KeyUp].Pressed() {
		dy++
	}
	if state[ebiten.KeyDown].Pressed() {
		dy--
	}
	if state[ebiten.KeyLeft].Pressed() {
		dx--
	}
	if state[ebiten.KeyRight].Pressed() {
		dx++
	}
	if dx != 0 || dy != 0 {
		ev := events.NewMoveByActorEvent(0, -1, dx, dy, 0, 1)
		events.Hub().ActorEvents().Enqueue(ev)
	}

	if state[ebiten.KeySpace].JustPressed() {
		ev := events.NewJumpActorEvent(0, -1, 4.0)
		events.Hub().ActorEvents().Enqueue(ev)
	}
}
