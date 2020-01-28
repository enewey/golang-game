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
	r.state["clock"] = c
}

// AfterConsecutiveReaction is a reaction that triggers its effect only after Tapped with the testFunc passing N frames in a row.
type AfterConsecutiveReaction struct {
	temporalReaction
	testFunc func(...interface{}) bool
	trigger  int
}

func (r *AfterConsecutiveReaction) getTicks() int {
	return r.state["ticks"].(types.Frame)
}
func (r *AfterConsecutiveReaction) setTicks(f types.Frame) {
	r.state["ticks"] = f
}

// Tap - trigger the reaction if Tap was called N frames in a row with the testFunc passing
func (r *AfterConsecutiveReaction) Tap(args ...interface{}) {
	prev := r.getClock()
	r.setClock(clock.Get())

	if !r.testFunc(args...) {
		r.setTicks(0)
	}

	diff := clock.Diff(prev)
	if clock.Cmp(diff, 1) > 0 {
		r.setTicks(1)
		return
	}

	r.setTicks(r.getTicks() + 1)
	if r.getTicks() > r.trigger {
		r.setTicks(0)
		r.reaction(args...)
	}
}
