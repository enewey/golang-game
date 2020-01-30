package events

import (
	"enewey.com/golang-game/clock"
	"enewey.com/golang-game/types"
)

type rFunc func(...interface{})

// Reaction are triggered functions in response to something happening, typically a collision.
type Reaction interface {
	Tap(...interface{})
}

// BasicReaction - simple passthrough-type reaction with no additional conditions or state
type BasicReaction struct {
	reaction rFunc
}

// Tap - passthrough the args to the reaction
func (r *BasicReaction) Tap(args ...interface{}) {
	r.reaction(args...)
}

// NewReaction - returns a new BasicReactions
func NewReaction(f rFunc) *BasicReaction {
	return &BasicReaction{f}
}

type statefulReaction struct {
	BasicReaction
	state map[string]interface{}
}

// TtmporalReaction is a stateful reaction that is aware of the time between consecutive reactions
type temporalReaction struct {
	statefulReaction
}

func (r *temporalReaction) getClock() clock.Clock {
	return r.state["clock"].(clock.Clock)
}
func (r *temporalReaction) setClock(c clock.Clock) {
	r.state["clock"] = r.getClock().Set(c)
}

func (r *AfterConsecutiveReaction) getLastTrigger() clock.Clock {
	return r.state["last_trigger"].(clock.Clock)
}

func (r *AfterConsecutiveReaction) setLastTrigger(c clock.Clock) {
	r.state["last_trigger"] = r.getLastTrigger().Set(c)
}

// AfterConsecutiveReaction is a reaction that triggers its effect only after Tapped with the testFunc passing N frames in a row.
type AfterConsecutiveReaction struct {
	temporalReaction
	testFunc                 func(...interface{}) bool
	triggerAfter, resetAfter int
}

// NewAfterConsecutiveReaction - creates a new Reaction where the reaction function is only triggered after the
// reaction is Tapped and the params pass the test function a number N frames in a row.
func NewAfterConsecutiveReaction(reaction rFunc, test func(...interface{}) bool, triggerAfter, resetAfter int) *AfterConsecutiveReaction {
	state := make(map[string]interface{})
	state["clock"] = clock.Copy()
	state["ticks"] = 0
	state["last_trigger"] = clock.Diff(clock.Get()) // a bit hacky to get a zero timestamp
	return &AfterConsecutiveReaction{
		temporalReaction{
			statefulReaction{
				BasicReaction{reaction},
				state,
			},
		},
		test,
		triggerAfter,
		resetAfter,
	}
}

func (r *AfterConsecutiveReaction) getTicks() int {
	return r.state["ticks"].(types.Frame)
}
func (r *AfterConsecutiveReaction) setTicks(f types.Frame) {
	r.state["ticks"] = f
}

// Tap - trigger the reaction if Tap was called N frames in a row with the testFunc passing
func (r *AfterConsecutiveReaction) Tap(args ...interface{}) {

	lastTrigger := r.getLastTrigger()

	if !r.testFunc(args...) && clock.Cmp(clock.Diff(lastTrigger), r.resetAfter) > 0 {
		r.setTicks(0)
		return
	}

	prevTap := r.getClock()
	diff := clock.Diff(prevTap)
	r.setClock(clock.Get())
	if clock.Cmp(diff, 1) > 0 {
		r.setTicks(1)
		return
	}

	r.setTicks(r.getTicks() + 1)
	if r.getTicks() > r.triggerAfter {
		r.setTicks(0)
		r.reaction(args...)
	}
}
