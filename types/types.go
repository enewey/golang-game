package types

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
