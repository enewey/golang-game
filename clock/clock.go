package clock

import (
	"math/big"

	"enewey.com/golang-game/types"
)

// Clock - represents how many ticks of gamestate have passed
type Clock = *big.Int

var singleton Clock

func init() {
	singleton = big.NewInt(0)
}

// Inc - increment the game clock by a number of frames.
func Inc(f types.Frame) Clock {
	singleton = singleton.Add(singleton, big.NewInt(int64(f)))
	return singleton
}

// Get - get the current game clock
func Get() Clock {
	return singleton
}

// Copy - gets a copy of the current clock
func Copy() Clock {
	return big.NewInt(0).Set(singleton)
}

// Diff - get the difference between the current game clock and the argument
func Diff(c Clock) Clock {
	return big.NewInt(0).Sub(singleton, c)
}

// Cmp - returns -1 if c < i; 0 if c == i; 1 if c > i
func Cmp(c Clock, i int) int {
	return c.Cmp(big.NewInt(int64(i)))
}
