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

// Dims - returns the width and height of the Sprite image
func (s *Sprite) Dims() (int, int) {
	return s.img.Size()
}

// Draw - draw this sprite at x/y on the given image
func (s *Sprite) Draw(x, y int, img *ebiten.Image) *ebiten.Image {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x), float64(y))

	img.DrawImage(s.img, opt)
	return img
}
