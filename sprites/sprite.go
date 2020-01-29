package sprites

import (
	"github.com/hajimehoshi/ebiten"
)

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

// NewCompoundSprite - create a single sprite composed of tiles, specifying an array
// of tiles to draw, the number of rows and columns, and the x/y dimensions of each tile in pixels.
func NewCompoundSprite(sprites []*Sprite, rows, cols, tilex, tiley int) *Sprite {
	if rows*cols != len(sprites) {
		panic("tried to create compound sprite with improper tile dimensions")
	}
	img, err := ebiten.NewImage(cols*tilex, rows*tiley, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	for i, v := range sprites {
		row := i / cols
		col := i % cols

		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(col*tilex), float64(row*tiley))

		img.DrawImage(v.Img(), opt)
	}

	return &Sprite{img}
}
