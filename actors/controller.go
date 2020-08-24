package actors

import (
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/utils"
)

type controller struct {
	target int
}

func (c *controller) SetTarget(t int) {
	c.target = t
}

// Controller w
type Controller interface {
	Tap(Actor, input.Input, types.Frame) bool
	SetTarget(t int)
}

// PlayerController w
type PlayerController struct {
	controller
}

// NewPlayerController t
func NewPlayerController() *PlayerController {
	return &PlayerController{controller{-1}}
}

// Tap w
func (c *PlayerController) Tap(target Actor, state input.Input, df types.Frame) bool {
	return controlPlayer(target, state)
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
		events.Enqueue(NewJumpEvent(target, 3.5))
	}

	if state[cfg.KeyDash()].JustPressed() && !player.(CanDash).Dashed() && player.OnGround() {
		vx, vy := utils.Normalize2(utils.Itof(DirToVec(player.Direction())))
		events.Enqueue(NewDashEvent(target, vx*2.5, vy*2.5, 0.0))
	}
	return true
}

type statefulController struct {
	controller
	state map[string]interface{}
}

// MoveSequenceController to move an actor in a repeated sequence
type MoveSequenceController struct {
	statefulController
	moves    []types.Direction // Directions
	dist     int               // distance of each move
	delay    types.Frame
	lastMove int
}

// NewMoveSequenceController w
func NewMoveSequenceController(moves []types.Direction, dist int, delay types.Frame) *MoveSequenceController {
	statemap := make(map[string]interface{})
	statemap["ticks"] = 0
	return &MoveSequenceController{
		statefulController{controller{-1}, statemap},
		moves, dist, delay, -1,
	}
}

// Tap w
func (c *MoveSequenceController) Tap(target Actor, state input.Input, df types.Frame) bool {
	c.state["ticks"] = c.state["ticks"].(int) + df
	ticks := c.state["ticks"].(int)
	move := int(ticks/c.delay) % len(c.moves)
	if c.lastMove != move {
		dx, dy := DirToVec(c.moves[move])
		events.Enqueue(NewMoveByEvent(target, float64(dx*c.dist), float64(dy*c.dist), 0, c.delay))
		c.lastMove = move
	}
	return true
}
