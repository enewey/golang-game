package sprites

import (
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
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

// DrawSprite draws the sprite with the given identifier as it relates to the
// Spritemap at the given coordinates on the given image.
func (sm *Spritemap) DrawSprite(id, x, y int, img *ebiten.Image) *ebiten.Image {
	return (*sm)[id].DrawSprite(x, y, img)
}
