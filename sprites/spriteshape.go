package sprites

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// SpriteRect is a rectangle shape
type SpriteRect struct {
	Sprite *Sprite
	w, h   int
}

// NewSpriteRect returns a filled color rectangle
func NewSpriteRect(w, h int, c color.Color) *SpriteRect {
	shape, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	shape.Fill(c)
	return &SpriteRect{&Sprite{shape}, w, h}
}
