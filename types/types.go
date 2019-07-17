package types

// Frame represents a number of update cycles
type Frame = int

// Point is a simple xy point
type Point struct {
	X, Y int
}

// NewPoint returns a pointer to a new Point struct
func NewPoint(x, y int) *Point { return &Point{x, y} }
