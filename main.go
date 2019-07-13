package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

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

var cX = 50
var cY = 76
var cZ = 0
var shadowZ = 0
var jumpZ, jumpTime int
var fallV float64
var onGround bool
var jumping = false

var tiles *sprites.Spritesheet
var charas *sprites.Spritesheet
var girlChar *sprites.Sprite
var shadowChar *sprites.Sprite
var scene *room.Room
var roomImage *ebiten.Image

var charBlock *collider.Collider

func init() {

	var err error
	tiles = cache.Get().LoadSpritesheet("blue-walls.png", tileDimX, tileDimY)
	if err != nil {
		log.Fatal(err)
	}

	charas = cache.Get().LoadSpritesheet("hoodgirl.png", tileDimX, tileDimY)
	girlChar = charas.GetSprite(0)
	shadowChar = charas.GetSprite(1)

	scene = cache.Get().LoadRoom("room2")
	charBlock = collider.NewBlock(cX+3, cY+8, cZ, 10, 8, 12, "chara")
	roomImage, _ = ebiten.NewImage(screenW*2, screenH*2, ebiten.FilterDefault)
}

func drawTile(mapX, mapY, tileNum int,
	rm *ebiten.Image, tiles *sprites.Spritesheet) {

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(tileDimX*mapX), float64(tileDimY*mapY))

	rm.DrawImage(tiles.GetSprite(tileNum).Img(), opt)
}

func drawSprite(x, y int, sprite *ebiten.Image, rm *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x), float64(y))

	rm.DrawImage(sprite, opt)
}

func drawRoom() *ebiten.Image {
	var spriteDrawn, shadowDrawn bool
	for pr, layer := range scene.Layers() {
		mapTiles := layer.Tiles()
		for i := 0; i < len(mapTiles); i++ {
			row := int(i / tilesX)
			col := i % tilesX

			drawTile(col, row, mapTiles[i], roomImage, tiles)
			sx, sy, sz := charBlock.GetPos()
			if col == 0 && row == int(sy/16)+1 && pr == int(shadowZ/16)+1 && !shadowDrawn {
				drawSprite(sx-3, sy-shadowZ-8, shadowChar.Img(), roomImage)
				shadowDrawn = true
			}
			if col == 0 && row == int(sy/16)+1 && pr == int(sz/16)+1 && !spriteDrawn {
				drawSprite(sx-3, sy-sz-8, girlChar.Img(), roomImage)
				spriteDrawn = true
			}
		}
	}

	if !spriteDrawn {
		sx, sy, sz := charBlock.GetPos()
		if !shadowDrawn {
			drawSprite(sx-3, sy-shadowZ-8, shadowChar.Img(), roomImage)
		}
		drawSprite(sx-3, sy-sz-8, girlChar.Img(), roomImage)
	}

	return roomImage
}

// ---------

func checkInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && onGround {
		fallV = 3.5
		onGround = false
	}
	var dx, dy, dz int

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

	dz = int(math.Max(fallV, -6)) // per second?? frame?? :thinking_face:
	ax, ay, az, hitGround, hitCeiling :=
		collider.ResolveCollision(dx, dy, dz, charBlock, scene.Colliders())
	// fmt.Printf("dz %d az %d\n fallV %f\n", dz, az, fallV)
	charBlock.Translate(ax, ay, az)
	if hitGround {
		onGround = true
	} else if az != 0 {
		onGround = false
	}
	if onGround && fallV < -1 && az >= 0 {
		fallV = 0
	} else if hitCeiling && fallV > 0 && az <= 0 {
		fallV = 0
	} else {
		fallV -= 0.25
	}

	shadowZ = scene.Colliders().FindFloor(charBlock)
}

// -------

func update(screen *ebiten.Image) error {
	checkInputs()

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	rm := drawRoom()
	opt := &ebiten.DrawImageOptions{}
	// opt.GeoM.Translate(-float64(screenW)/2.0, -float64(screenH)/2.0)
	opt.GeoM.Scale(3, 3)
	screen.DrawImage(rm, opt)
	x, y, z := charBlock.GetPos()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %d\ny: %d\n z: %d\n shadz: %d", x, y, z, shadowZ))
	return nil
}

func main() {
	if err := ebiten.Run(
		update, screenW*3, screenH*3, 1, "Jumpin' Game",
	); err != nil {
		log.Fatal(err)
	}
}
