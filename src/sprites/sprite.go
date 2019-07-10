package sprites

import "github.com/hajimehoshi/ebiten"

// Sprite woo
type Sprite struct {
	img *ebiten.Image
}

// Img woo
func (s *Sprite) Img() *ebiten.Image {
	return s.img
}
