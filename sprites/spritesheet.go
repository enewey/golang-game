package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Spritesheet woo
type Spritesheet struct {
	cache   map[int]*Sprite
	image   *ebiten.Image
	theight int
	twidth  int
	rows    int
	cols    int
}

// New woo
func New(image *ebiten.Image, dimX, dimY, rows, cols int) *Spritesheet {
	return &Spritesheet{
		make(map[int]*Sprite),
		image,
		dimX,
		dimY,
		rows,
		cols,
	}
}

// GetSprite woo
func (ts *Spritesheet) GetSprite(num int) *Sprite {
	if ts.cache[num] == nil {
		x := num % ts.cols
		y := int(num / ts.rows)
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
