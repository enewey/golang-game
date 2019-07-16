package room

import "log"

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

// TilesRow - a single row of tile refs
// row - 0-indexed top-justified row
// dimX - length of a single row
func (lyr *Layer) TilesRow(row, dimX int) []int {
	if row*dimX > len(lyr.Tiles()) {
		log.Fatal("out of bounds row on getTileRow")
	}

	return lyr.Tiles()[row*dimX : (row+1)*dimX]
}

// Priority woo
func (lyr *Layer) Priority() int { return lyr.priority }

// ByPriority woo
type ByPriority []*Layer

func (y ByPriority) Len() int      { return len(y) }
func (y ByPriority) Swap(i, j int) { y[i], y[j] = y[j], y[i] }
func (y ByPriority) Less(i, j int) bool {
	return y[i].Priority() < y[j].Priority()
}
