package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"

	"enewey.com/golang-game/src/cache"
	"enewey.com/golang-game/src/room"
	"enewey.com/golang-game/src/sprites"
)

const tileWidth = 16
const tileHeight = 16

var charaX = 50
var charaY = 76
var charaPr = 1

var tiles *sprites.Spritesheet
var charas *sprites.Spritesheet
var girlChar *sprites.Sprite
var scene *room.Room

func init() {

	var err error
	tiles = sprites.New(cache.Get().LoadImage("blue-walls.png"), 16, 16)
	if err != nil {
		log.Fatal(err)
	}

	charas = sprites.New(cache.Get().LoadImage("hoodgirl.png"), 16, 16)
	girlChar = charas.GetSpriteByNum(0)

	scene = cache.Get().LoadRoom("room1")
}

func drawTile(mapX, mapY, tileNum int, rm *ebiten.Image, tiles *sprites.Spritesheet) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(16*mapX), float64(16*mapY))
	// opt.GeoM.Scale(1.5, 1)

	rm.DrawImage(tiles.GetSpriteByNum(tileNum).Img(), opt)
}

func drawSprite(x, y int, sprite *ebiten.Image, rm *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x), float64(y))

	rm.DrawImage(sprite, opt)
}

func drawRoom() *ebiten.Image {
	rm, _ := ebiten.NewImage(320, 240, ebiten.FilterDefault)

	for pr, layer := range scene.Layers() {
		mapTiles := layer.Tiles()
		for i := 0; i < len(mapTiles); i++ {
			row := int(i / 20)
			col := i % 20

			drawTile(col, row, mapTiles[i], rm, tiles)
			if col == 0 && row == int(charaY/16)+1 && pr == charaPr {
				drawSprite(charaX, charaY, girlChar.Img(), rm)
			}
		}
	}

	return rm
}

func checkInputs() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		charaX++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		charaX--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		charaY++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		charaY--
	}
}

func update(screen *ebiten.Image) error {
	checkInputs()

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	rm := drawRoom()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3, 3)
	screen.DrawImage(rm, opt)
	return nil
}

func main() {
	if err := ebiten.Run(update, 960, 720, 1, "Render tiles"); err != nil {
		log.Fatal(err)
	}
}
