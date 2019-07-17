package actors

import (
	"fmt"
	"math"

	"enewey.com/golang-game/types"
)

// Action - something that happens to a target Actor over a number of frames.
type Action interface {
	Target() *Actor
	Elapsed() types.Frame
	Process(types.Frame) bool // return value denotes completion.
	// A completed Action is to be discarded.
}

// Actions woo
type Actions []Action

// Add - Add a new action. Keeps the slice slim.
func (acts Actions) Add(a Action) {
	for i, v := range acts {
		if v == nil {
			acts[i] = a
			return
		}
	}
	acts = append(acts, a)
}

// BaseAction woo
type BaseAction struct {
	target   *Actor
	duration types.Frame // frames
	elapsed  types.Frame // frames
}

// Target woo
func (b *BaseAction) Target() *Actor { return b.target }

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
	fmt.Printf("Process dx %d dy %d dz %d\nvx %f vy %f vz %f\n duration %d elapsed %d\n", a.dx, a.dy, a.dz, a.vx, a.vy, a.vz, a.duration, a.elapsed)
	a.elapsed += df
	a.target.vx, a.target.vy = a.vx, a.vy
	if a.elapsed >= a.duration {
		a.target.vx -= a.vx
		a.target.vy -= a.vy
		return true
	}
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

// DashAction w
type DashAction struct {
	BaseAction
	vx, vy, vz float64
}

// NewDashAction w
func NewDashAction(target *Actor, vx, vy float64) *DashAction {
	return &DashAction{BaseAction{target, 15, 0}, vx, vy, 0.3}
}

// Process w
func (a *DashAction) Process(df types.Frame) bool {
	a.elapsed += df
	a.target.dashed = true

	dz := a.vz * float64((a.duration-a.elapsed)/a.duration)
	if dz <= 0.1 {
		dz = a.target.vz
	}
	a.target.SetVel(a.vx, a.vy, dz)

	a.target.controlled = false
	return a.elapsed >= a.duration && a.target.OnGround()
}
