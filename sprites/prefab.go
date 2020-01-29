package sprites

import "enewey.com/golang-game/config"

// Create2by1Block - creates a compound sprite shaped like a block,
// which is just two tiles imposed on top of each other.
func Create2by1Block(top, bottom *Sprite) Spritemap {
	xdim := config.Get().TileDimX
	ydim := config.Get().TileDimY

	return NewStaticSpritemap(NewCompoundSprite([]*Sprite{top, bottom}, 2, 1, xdim, ydim))
}
