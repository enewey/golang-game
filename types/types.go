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

// Direction is for the actor 'direction' property
type Direction int

// Direction types
const (
	Up Direction = iota
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

// IsX - convenience method tells if the axis map is flagged for the X axis
func (a *AxisMap) IsX() bool { return a.X != 0 }

// IsY - convenience method tells if the axis map is flagged for the Y axis
func (a *AxisMap) IsY() bool { return a.Y != 0 }

// IsZ - convenience method tells if the axis map is flagged for the Z axis
func (a *AxisMap) IsZ() bool { return a.Z != 0 }

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

// RejectVec - takes in a vector and returns zeroes for every axis of the AxisMap that isn't zero.
// e.g. if AxisMap.X is -1, then vx will always be 0.
// e.g. if AxisMap.Y is 0, then vy will pass through as vy.
func (m *AxisMap) RejectVec(vx, vy, vz float64) (float64, float64, float64) {
	return utils.Unpass(m.X, vx), utils.Unpass(m.Y, vy), utils.Unpass(m.Z, vz)
}
