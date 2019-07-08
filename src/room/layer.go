package room

// Layer woo
type Layer struct {
	tiles    []int
	priority int
}

// NewLayer woo
func NewLayer(tiles []int, priority int) *Layer {
	return &Layer{tiles, priority}
}

// Tiles woo
func (lyr *Layer) Tiles() []int { return lyr.tiles }
