package sprites

import (
	"enewey.com/golang-game/types"
)

// Spritemap allows mapping an int to sprites -- handy for use with the Direction iota
type Spritemap map[int]*Sprite

// NewCharaSpritemap returns a new 4 directional spritemap for an actor
func NewCharaSpritemap(d, r, u, l *Sprite) *Spritemap {
	return &Spritemap{
		types.Up:        u,
		types.Down:      d,
		types.Right:     r,
		types.Left:      l,
		types.UpRight:   u,
		types.UpLeft:    u,
		types.DownRight: d,
		types.DownLeft:  d,
	}
}

// Sprite returns the Sprite with the given ID
func (sm *Spritemap) Sprite(id int) *Sprite {
	return (*sm)[id]
}
