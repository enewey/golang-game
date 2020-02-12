package actors

import (
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/utils"
)

// Controller types
const (
	PlayerController = iota
)

// Control is the main function to be called.
// Args are the controller type, the actor being controlled, and whatever argument the controller requires.
// e.g. the PlayerController requires the input state as its arg.
func Control(t int, target Actor, arg interface{}) bool {
	switch t {
	case PlayerController:
		if in, ok := arg.(input.Input); ok {
			return controlPlayer(target, in)
		}
	}
	return false
}

func controlPlayer(target Actor, state input.Input) bool {
	cfg := config.Get()
	player := target.(CanMove)

	if player.OnGround() {
		_, _, vz := player.Vel()
		player.SetVel(0, 0, vz)
	}
	var dx, dy float64
	if state[cfg.KeyUp()].Pressed() {
		dy--
	}
	if state[cfg.KeyDown()].Pressed() {
		dy++
	}
	if state[cfg.KeyLeft()].Pressed() {
		dx--
	}
	if state[cfg.KeyRight()].Pressed() {
		dx++
	}
	player.SetVelX(dx)
	player.SetVelY(dy)
	player.CalcDirection()

	if state[cfg.KeyConfirm()].JustPressed() {
		events.Enqueue(NewInteractEvent(target))
	}

	if state[cfg.KeyJump()].JustPressed() && player.OnGround() {
		events.Enqueue(NewJumpEvent(target, 4.0))
	}

	if state[cfg.KeyDash()].JustPressed() && !player.(CanDash).Dashed() && player.OnGround() {
		vx, vy := utils.Normalize2(utils.Itof(DirToVec(player.Direction())))
		events.Enqueue(NewDashEvent(target, vx*2.5, vy*2.5, 0.0))
	}
	return true
}
