package sprites

import "github.com/hajimehoshi/ebiten"

type SpriteRect struct {
	Sprite *Sprite
	w, h   int
}

func NewSpriteRect(w, h int) *SpriteRect {
	shape, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	return &SpriteRect{&Sprite{shape}, w, h}
}
