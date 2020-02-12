package actors

import (
	"fmt"
	"math"

	"enewey.com/golang-game/config"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/utils"
)

// ActionType, expressed from an Event.command
const (
	MoveToActionType = iota
	MoveByActionType
	JumpActionType
	DashActionType
	ChangePosActionType
)

// InterpretEvent - translate an event into an action
func InterpretEvent(ev *events.Event) Action {
	fmt.Printf("interpreting event %d :: ", ev.Code())
	p := ev.Payload()
	switch ev.Code() {
	case MoveToActionType:
	case MoveByActionType:
		fmt.Printf("moveby action interpreted %v\n", ev.Payload())
		return NewMoveByAction(p[0].(Actor), p[1].(int), p[2].(int), p[3].(int), p[4].(int))
	case JumpActionType:
		fmt.Printf("jump action interpreted %v\n", ev.Payload())
		return NewJumpAction(p[0].(Actor), p[1].(float64))
	case DashActionType:
		fmt.Printf("dash action interpreted %v\n", ev.Payload())
		return NewDashAction(p[0].(Actor), p[1].(float64), p[2].(float64), p[3].(float64))
	case ChangePosActionType:

	default:
		fmt.Printf("unknown actor event code %d\n", ev.Code())
	}

	return nil
}

// Action - something that happens to a target Actor over a number of frames.
type Action interface {
	Target() Actor
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
	target   Actor
	duration types.Frame // frames
	elapsed  types.Frame // frames
}

// Target woo
func (b *BaseAction) Target() Actor { return b.target }

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
	target := a.target.(CanMove)
	a.elapsed += df
	if (x == a.tx && y == a.ty && z == a.tz) || a.elapsed > a.duration {
		target.SetVel(0, 0, 0)
		return true
	}
	vx := calcMoveToVel(a.sx, a.tx, x, a.speed, a.elapsed)
	vy := calcMoveToVel(a.sy, a.ty, y, a.speed, a.elapsed)
	vz := calcMoveToVel(a.sz, a.tz, z, a.speed, a.elapsed)
	target.SetVel(vx, vy, vz)
	return false
}

// MoveByAction woo
type MoveByAction struct {
	BaseAction
	dx, dy, dz int // delta x/y/z
	vx, vy, vz float64
	cx, cy, cz float64 // carry-over values between frames
}

// NewMoveByAction woo
func NewMoveByAction(target Actor, dx, dy, dz int, duration types.Frame) *MoveByAction {
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
		0.0, 0.0, 0.0,
	}
}

// Process w
func (a *MoveByAction) Process(df types.Frame) bool {
	// _, _, vz := target.Vel()
	a.elapsed += df
	target := a.target.(CanMove)
	if a.elapsed > a.duration {
		target.SetVel(0, 0, 0)
		return true
	}
	target.SetVel(a.vx+a.cx, a.vy+a.cy, a.vz+a.cz)
	a.cx, a.cy, a.cz = utils.Carry(a.vx+a.cx, a.vy+a.cy, a.vz+a.cz)
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
func NewJumpAction(target Actor, v float64) *JumpAction {
	return &JumpAction{BaseAction{target, 0, 0}, v}
}

// Process w
func (a *JumpAction) Process(df types.Frame) bool {
	target := a.target.(CanMove)
	if target.OnGround() {
		vx, vy, _ := target.Vel()
		target.SetVel(vx, vy, a.v)
		target.SetOnGround(false)
	}
	return true
}

// DashAction w
type DashAction struct {
	BaseAction
	vx, vy, vz float64
	axes       *types.AxisMap
}

// NewDashAction w
func NewDashAction(target Actor, vx, vy, vz float64) *DashAction {
	axes := types.VecToAxisMap(utils.Normalize3(vx, vy, vz))
	return &DashAction{BaseAction{target, 15, 0}, vx, vy, vz, axes}
}

// Process w
func (a *DashAction) Process(df types.Frame) bool {
	target := a.target.(CanMove)
	if dasher, ok := target.(CanDash); ok {
		dasher.SetDashed(true)
	} else {
		return true
	}

	// target.SetControlled(false)

	vx, vy, vz := a.axes.RejectVec(target.Vel())
	if a.axes.IsZ() {
		a.vz += config.Get().Gravity()
	}

	target.SetVel(
		a.vx+vx,
		a.vy+vy,
		a.vz+vz,
	)

	a.elapsed += df
	return a.elapsed >= a.duration && (target.OnGround() || a.axes.IsZ())
}

// ChangePosAction will change the position of the target actor immediately.
type ChangePosAction struct {
	BaseAction
	x, y, z int
}

// NewChangePosAction creates a new change pos action.
// 'Persists' says that this action will continue to change
func NewChangePosAction(target Actor, x, y, z int) *ChangePosAction {
	return &ChangePosAction{BaseAction{target, 0, 0}, x, y, z}
}

// Process function for the ChangePosAction
func (a *ChangePosAction) Process(df types.Frame) bool {
	target := a.target.(CanMove)
	target.Collider().SetPos(a.x, a.y, a.z)
	return true
}
