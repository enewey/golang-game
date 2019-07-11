package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"

	"enewey.com/golang-game/cache"
	"enewey.com/golang-game/collider"
	"enewey.com/golang-game/room"
	"enewey.com/golang-game/sprites"
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
var roomImage *ebiten.Image
var collGroupCache map[int]*resolv.Space

var block *resolv.Rectangle
var charBlock *resolv.Rectangle

func init() {

	var err error
	tiles = cache.Get().LoadSpritesheet("blue-walls.png", tileDimX, tileDimY)
	if err != nil {
		log.Fatal(err)
	}

	charas = cache.Get().LoadSpritesheet("hoodgirl.png", tileDimX, tileDimY)
	girlChar = charas.GetSprite(0)

	scene = cache.Get().LoadRoom("room2")

	block = resolv.NewRectangle(48, 32, 16, 16)
	charBlock = resolv.NewRectangle(int32(charaX)+2, int32(charaY)+8, 12, 8)
	roomImage, _ = ebiten.NewImage(screenW, screenH, ebiten.FilterDefault)
	collGroupCache = make(map[int]*resolv.Space)
}

func drawTile(mapX, mapY, tileNum int,
	rm *ebiten.Image, tiles *sprites.Spritesheet) {

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(tileDimX*mapX), float64(tileDimY*mapY))
	// opt.GeoM.Scale(1.5, 1)

	rm.DrawImage(tiles.GetSprite(tileNum).Img(), opt)
}

func drawSprite(x, y int, sprite *ebiten.Image, rm *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x), float64(y))

	rm.DrawImage(sprite, opt)
}

func drawRoom() *ebiten.Image {
	for pr, layer := range scene.Layers() {
		mapTiles := layer.Tiles()
		for i := 0; i < len(mapTiles); i++ {
			row := int(i / tilesX)
			col := i % tilesX

			drawTile(col, row, mapTiles[i], roomImage, tiles)
			if col == 0 && row == int(charaY/16)+1 && pr <= charaZ {
				drawSprite(charaX, charaY-charaH, girlChar.Img(), roomImage)
			}
		}
	}

	return roomImage
}

func checkInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !jumping {
		jumping = true
		jumpTime = 0
	}
	var dx, dy int

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy--
	}

	if !(dx == 0 && dy == 0) {
		ax, ay := check2DCollision(dx, dy, charaZ, scene.Colliders())
		charaX += ax
		charaY += ay
	}
}

func check2DCollision(dx, dy, currZ int, blocks collider.Colliders) (int, int) {
	var rx, ry int = dx, dy
	if collGroupCache[currZ] == nil {
		collGroupCache[currZ] = blocks.GetGroup(currZ, "walls")
	}
	walls := collGroupCache[currZ]
	// check dx first
	resX := walls.Resolve(charBlock, int32(dx), 0)
	if resX.Colliding() {
		charBlock.X += resX.ResolveX
		rx = int(resX.ResolveX)
	} else {
		charBlock.X += int32(dx)
	}

	// then check dy
	resY := walls.Resolve(charBlock, 0, int32(dy))
	if resY.Colliding() {
		charBlock.Y += resY.ResolveY
		ry = int(resY.ResolveY)
	} else {
		charBlock.Y += int32(dy)
	}

	return rx, ry
}

// func checkZCollision(dz, currZ int) int {

// }

func update(screen *ebiten.Image) error {
	checkInputs()
	if jumping {
		jumpTime += 3
		rads := math.Pi * float64(jumpTime) / 100.0
		sine := math.Sin(rads)
		charaZ = 1 + int(3.0*sine)
		charaH = int(1 + (sine * tileDimY * 1.25))
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
	if err := ebiten.Run(
		update, screenW*3, screenH*3, 1, "Jumpin' Game",
	); err != nil {
		log.Fatal(err)
	}
}
