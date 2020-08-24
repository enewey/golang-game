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
// Each action has a "Process" function and a "New___Action" function
// Creating an action should rarely happen directly, and instead should be the result of processing an event.
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
		// TODO: Do we need this? I mean probably, but...
	case MoveByActionType:
		fmt.Printf("moveby action interpreted %v\n", ev.Payload())
		return NewMoveByAction(p[0].(Actor), p[1].(float64), p[2].(float64), p[3].(float64), p[4].(int))
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

// Actions - an array of actions
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

// BaseAction - common data type for all actions.
type BaseAction struct {
	target   Actor
	duration types.Frame // frames
	elapsed  types.Frame // frames
}

// Target - the subject Actor of this Action
func (b *BaseAction) Target() Actor { return b.target }

// Elapsed - the number of frames for which this Action has been processed
func (b *BaseAction) Elapsed() types.Frame { return b.elapsed }

// --- Action Definitions

// MoveToAction - moves the Actor to a specific position at a given speed.
type MoveToAction struct {
	BaseAction
	sx, sy, sz float64 // starting x/y/z
	tx, ty, tz float64 // target x/y/z
	speed      float64 // pixels per 0.0167 seconds
}

// Process - processes a MoveToAction
func (a *MoveToAction) Process(df int) bool {
	x, y, z := a.target.Pos()
	target := a.target.(CanMove)
	a.elapsed += df
	if (float64(x) == a.tx && float64(y) == a.ty && float64(z) == a.tz) || a.elapsed > a.duration {
		target.SetVel(0, 0, 0)
		return true
	}
	vx := calcMoveToVel(a.sx, a.tx, float64(x), a.speed, a.elapsed)
	vy := calcMoveToVel(a.sy, a.ty, float64(y), a.speed, a.elapsed)
	vz := calcMoveToVel(a.sz, a.tz, float64(z), a.speed, a.elapsed)
	target.SetVel(vx, vy, vz)
	return false
}

// MoveByAction woo
type MoveByAction struct {
	BaseAction
	dx, dy, dz float64 // delta x/y/z
	vx, vy, vz float64
}

// NewMoveByAction woo
func NewMoveByAction(target Actor, dx, dy, dz float64, duration types.Frame) *MoveByAction {
	return &MoveByAction{
		BaseAction{
			target,
			duration,
			0,
		},
		dx, dy, dz,
		dx / float64(duration),
		dy / float64(duration),
		dz / float64(duration),
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
	target.SetVel(a.vx, a.vy, a.vz)
	return false
}

func calcMoveToVel(start, end, current float64, speed float64, elapsed types.Frame) float64 {
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
