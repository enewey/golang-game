package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"

	"enewey.com/golang-game/src/cache"
	"enewey.com/golang-game/src/room"
	"enewey.com/golang-game/src/sprites"
)

const tileDimX = 16
const tileDimY = 16
const tilesX = 10
const tilesY = 8

var screenW = tileDimX * tilesX
var screenH = tileDimY * tilesY

var charaX = 50
var charaY = 76
var charaH = 0
var charaZ = 1
var jumping = false
var jumpTime = 0

var tiles *sprites.Spritesheet
var charas *sprites.Spritesheet
var girlChar *sprites.Sprite
var scene *room.Room

func init() {

	var err error
	tiles = sprites.New(cache.Get().LoadImage("blue-walls.png"), tileDimX, tileDimY)
	if err != nil {
		log.Fatal(err)
	}

	charas = sprites.New(cache.Get().LoadImage("hoodgirl.png"), tileDimX, tileDimY)
	girlChar = charas.GetSpriteByNum(0)

	scene = cache.Get().LoadRoom("room2")
}

func drawTile(mapX, mapY, tileNum int, rm *ebiten.Image, tiles *sprites.Spritesheet) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(tileDimX*mapX), float64(tileDimY*mapY))
	// opt.GeoM.Scale(1.5, 1)

	rm.DrawImage(tiles.GetSpriteByNum(tileNum).Img(), opt)
}

func drawSprite(x, y int, sprite *ebiten.Image, rm *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x), float64(y))

	rm.DrawImage(sprite, opt)
}

func drawRoom() *ebiten.Image {
	rm, _ := ebiten.NewImage(screenW, screenH, ebiten.FilterDefault)

	for pr, layer := range scene.Layers() {
		mapTiles := layer.Tiles()
		for i := 0; i < len(mapTiles); i++ {
			row := int(i / tilesX)
			col := i % tilesX

			drawTile(col, row, mapTiles[i], rm, tiles)
			if col == 0 && row == int(charaY/16)+1 && pr <= charaZ {
				drawSprite(charaX, charaY-charaH, girlChar.Img(), rm)
			}
		}
	}

	return rm
}

func checkInputs() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !jumping {
		jumping = true
		jumpTime = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
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
	if jumping {
		jumpTime += 2
		rads := math.Pi * float64(jumpTime) / 100.0
		sine := math.Sin(rads)
		charaZ = 1 + int(math.Round(2.0*sine))
		charaH = int(1 + (sine * tileDimY))
		if jumpTime >= 100 {
			jumping = false
			charaH = 0
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	rm := drawRoom()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3, 3)
	screen.DrawImage(rm, opt)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("z: %d\nh: %d", charaZ, charaH))
	return nil
}

func main() {
	if err := ebiten.Run(update, screenW*3, screenH*3, 1, "Render tiles"); err != nil {
		log.Fatal(err)
	}
}
