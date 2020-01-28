package clock

import (
	"math/big"

	"enewey.com/golang-game/types"
)

// Clock - represents how many ticks of gamestate have passed
type Clock = *big.Int

var clock Clock

func init() {
	clock = big.NewInt(0)
}

// Inc - increment the game clock by a number of frames.
func Inc(f types.Frame) Clock {
	clock.Add(clock, big.NewInt(int64(f)))
	return clock
}

// Get - get the current game clock
func Get() Clock {
	return clock
}

// Diff - get the difference between the current game clock and the argument
func Diff(c Clock) Clock {
	return clock.Sub(c, clock)
}

// Cmp - returns -1 if c < i; 0 if c == i; 1 if c > i
func Cmp(c Clock, i int) int {
	return c.Cmp(big.NewInt(int64(i)))
}
