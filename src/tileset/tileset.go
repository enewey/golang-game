package tileset

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const dimX = 30
const dimY = 30

// Tileset woo
type Tileset struct {
	source  string
	image   *ebiten.Image
	theight int
	twidth  int
}

// New woo
func New(source string) *Tileset {
	return &Tileset{
		source,
		nil,
		16.0,
		16.0,
	}
}

// Load woo
func (ts *Tileset) Load() (*Tileset, error) {
	if ts.image != nil {
		return ts, nil
	}
	img, _, err := ebitenutil.NewImageFromFile("assets/img/blue-walls.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	ts.image = img

	return ts, nil
}

// GetTileByNum woo
func (ts *Tileset) GetTileByNum(num int) *ebiten.Image {
	x := num % dimX
	y := int(num / dimY)
	return ts.GetTile(x, y)
}

// GetTile woo
func (ts *Tileset) GetTile(x, y int) *ebiten.Image {
	if ts.image == nil {
		_, err := ts.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	tx := x * ts.twidth
	ty := y * ts.theight

	return ts.image.SubImage(
		image.Rect(tx, ty, tx+ts.twidth, ty+ts.theight),
	).(*ebiten.Image)
}
