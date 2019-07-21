package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const dimX = 30
const dimY = 30

// Spritesheet woo
type Spritesheet struct {
	cache   map[int]*Sprite
	image   *ebiten.Image
	theight int
	twidth  int
}

// New woo
func New(image *ebiten.Image, dimX, dimY int) *Spritesheet {
	return &Spritesheet{
		make(map[int]*Sprite),
		image,
		dimX,
		dimY,
	}
}

// GetSprite woo
func (ts *Spritesheet) GetSprite(num int) *Sprite {
	if ts.cache[num] == nil {
		x := num % dimX
		y := int(num / dimY)
		ts.cache[num] = ts.getSpriteByCoord(x, y)
	}
	return ts.cache[num]
}

func (ts *Spritesheet) getSpriteByCoord(x, y int) *Sprite {
	tx := x * ts.twidth
	ty := y * ts.theight

	return &Sprite{ts.image.SubImage(
		image.Rect(tx, ty, tx+ts.twidth, ty+ts.theight),
	).(*ebiten.Image)}
}
