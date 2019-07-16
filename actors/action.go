package actors

import (
	"math"

	"enewey.com/golang-game/types"
)

// Action - a series of commands fed to a target Actor
type Action interface {
	Target() *Actor
	Elapsed() types.Frame
	Process(types.Frame) bool // return value denotes completion.
	// A completed Action is to be discarded.
}

// Actions woo
type Actions []Action

// BaseAction woo
type BaseAction struct {
	target   *Actor
	duration types.Frame // frames
	elapsed  types.Frame // frames
}

// Target woo
func (b *BaseAction) Target() *Actor       { return b.target }
// Elapsed woo
func (b *BaseAction) Elapsed() types.Frame { return b.elapsed }

// MoveToAction woo
type MoveToAction struct {
	BaseAction
	sx, sy, sz int     // starting x/y/z
	tx, ty, tz int     // target x/y/z
	speed      float64 // pixels per 0.0167 seconds
}

// Process woo
func (a *MoveToAction) Process(df int) bool {
	x, y, z := a.target.Pos()
	a.elapsed += df
	if (x == a.tx && y == a.ty && z == a.tz) || a.elapsed > a.duration {
		a.target.SetVel(0, 0, 0)
		return true
	}
	vx := calcMoveToVel(a.sx, a.tx, x, a.speed, a.elapsed)
	vy := calcMoveToVel(a.sy, a.ty, y, a.speed, a.elapsed)
	vz := calcMoveToVel(a.sz, a.tz, z, a.speed, a.elapsed)
	a.target.SetVel(vx, vy, vz)
	return false
}

// MoveByAction woo
type MoveByAction struct {
	BaseAction
	dx, dy, dz int // delta x/y/z
	vx, vy, vz float64
}

// NewMoveByAction woo
func NewMoveByAction(target *Actor, dx, dy, dz int, duration types.Frame) *MoveByAction {
	return &MoveByAction{
		BaseAction{
			target,
			duration,
			0,
		},
		dx, dy, dz,
		float64(dx) / float64(duration),
		float64(dy) / float64(duration),
		float64(dz) / float64(duration),
	}
}

// Process w
func (a *MoveByAction) Process(df types.Frame) bool {
	a.elapsed += df
	if a.elapsed > a.duration {
		a.target.SetVel(0, 0, 0)
		return true
	}
	a.target.SetVel(a.vx, a.vy, a.vz)
	return false
}

func calcMoveToVel(start, end, current int, speed float64, elapsed types.Frame) float64 {
	projectedDist := speed * float64(elapsed) // hmm.. rounding?
	actualDist := current - start
	destinationDist := end - start
	if actualDist < 0 {
		projectedDist *= -1
	}
	if math.Abs(projectedDist) > math.Abs(float64(destinationDist)) {
		return float64(destinationDist - actualDist)
	}
	return projectedDist - float64(actualDist)
}

// JumpAction w
type JumpAction struct {
	BaseAction
	v float64
}

// NewJumpAction w
func NewJumpAction(target *Actor, v float64) *JumpAction {
	return &JumpAction{BaseAction{target, 0, 0}, v}
}

// Process w
func (a *JumpAction) Process(df types.Frame) bool {
	if a.target.OnGround() {
		vx, vy, _ := a.target.Vel()
		a.target.SetVel(vx, vy, a.v)
		a.target.onGround = false
	}
	return true
}
