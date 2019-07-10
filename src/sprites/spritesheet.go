package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const dimX = 30
const dimY = 30

// Spritesheet woo
type Spritesheet struct {
	image   *ebiten.Image
	theight int
	twidth  int
}

// New woo
func New(image *ebiten.Image, dimX, dimY int) *Spritesheet {
	return &Spritesheet{
		image,
		dimX,
		dimY,
	}
}

// GetSpriteByNum woo
func (ts *Spritesheet) GetSpriteByNum(num int) *Sprite {
	x := num % dimX
	y := int(num / dimY)
	return ts.GetSprite(x, y)
}

// GetSprite woo
func (ts *Spritesheet) GetSprite(x, y int) *Sprite {
	tx := x * ts.twidth
	ty := y * ts.theight

	return &Sprite{ts.image.SubImage(
		image.Rect(tx, ty, tx+ts.twidth, ty+ts.theight),
	).(*ebiten.Image)}
}
