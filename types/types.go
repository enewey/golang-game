package types

import (
	"math"

	"enewey.com/golang-game/utils"
)

// Frame represents a number of update cycles
type Frame = int

// Reaction is an adhoc function where the params passed in
// will have their state modified
type Reaction = func(...interface{})

// Point is a simple xy point
type Point struct {
	X, Y int
}

// NewPoint returns a pointer to a new Point struct
func NewPoint(x, y int) *Point { return &Point{x, y} }

// Directions for the actor 'direction' property
const (
	Up = iota
	UpRight
	Right
	DownRight
	Down
	DownLeft
	Left
	UpLeft
)

// AxisMap used to describe relevant axes, for example in a movement action
type AxisMap struct {
	X, Y, Z int
}

// VecToAxisMap - converts a vector to an AxisMap. The highest value(s)
// of the vector will be non-zero in the axis map.
func VecToAxisMap(x, y, z float64) *AxisMap {
	highest := math.Abs(x)
	if math.Abs(y) > highest {
		highest = math.Abs(y)
	}
	if math.Abs(z) > highest {
		highest = math.Abs(z)
	}

	if math.Abs(x) != highest {
		x = 0
	}
	if math.Abs(y) != highest {
		y = 0
	}
	if math.Abs(z) != highest {
		z = 0
	}
	return &AxisMap{
		X: int(utils.Normalize(x)),
		Y: int(utils.Normalize(y)),
		Z: int(utils.Normalize(z)),
	}
}

// FilterVec - takes in a vector and filters out values based on the AxisMap
// e.g. if AxisMap.X is -1, then only negative vx values will be returned
// e.g. if AxisMap.Y is 0, then the returned vy will be what was passed in
func (m *AxisMap) FilterVec(vx, vy, vz float64) (float64, float64, float64) {
	return utils.Pass(m.X, vx), utils.Pass(m.Y, vy), utils.Pass(m.Z, vz)
}
