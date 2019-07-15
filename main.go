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

func drawSprite(x, y int, sprite *sprites.Sprite, rm *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x), float64(y))

	rm.DrawImage(sprite.Img(), opt)
}

func getLayerRow(row int, layer *room.Layer) []int {
	if row*tilesX > len(layer.Tiles()) {
		log.Fatal("out of bounds row on getTileRow")
	}

	return layer.Tiles()[row*tilesX : (row+1)*tilesX]
}

func drawRoom() *ebiten.Image {
	var spriteDrawn, shadowDrawn bool
	for _, layer := range scene.Layers() {
		pr := layer.Priority() // z-layer of this tile, basically
		mapTiles := layer.Tiles()
		for i := 0; i < len(mapTiles); i += tilesX {
			row := int(i / tilesX)
			for col := 0; col < tilesX; col++ {
				tile := tiles.GetSprite(mapTiles[i+col])
				drawSprite(col*tileDimX, row*tileDimY, tile, roomImage)
			}

			sx, sy, sz := charBlock.Pos()
			sd := charBlock.Depth()
			charPr := int(math.Ceil(float64(sz+1) / float64(tileDimY/2)))
			shadowPr := int(math.Floor(float64(shadowZ+sd) / float64(tileDimY/2)))
			charRow := int(math.Ceil(float64(sy) / float64(tileDimY)))
			doDraw := charPr == pr && charRow == (row+pr)
			shadowDraw := shadowPr == pr && charRow == (row+pr)

			if shadowDraw {
				drawSprite(sx-4, sy-shadowZ-8, shadowChar, roomImage)
				shadowDrawn = true
			}
			if doDraw {
				drawSprite(sx-4, sy-sz-8, girlChar, roomImage)
				spriteDrawn = true
			}
		}
	}

	if !spriteDrawn {
		sx, sy, sz := charBlock.Pos()
		if !shadowDrawn {
			drawSprite(sx-3, sy-shadowZ-8, shadowChar, roomImage)
		}
		drawSprite(sx-3, sy-sz-8, girlChar, roomImage)
	}
	return roomImage
}

// ---------

var logicFrame int
var slowdownSwitch bool

func checkInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && onGround {
		fallV = 4.0
		onGround = false
	}
	var dx, dy, dz int

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		slowdownSwitch = !slowdownSwitch
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

	logicFrame = (logicFrame + 1) % 4
	if (logicFrame == 0 && slowdownSwitch) || !slowdownSwitch {
		dz = int(math.Max(fallV, -6)) // per second?? frame?? :thinking_face:
		var ax, ay, az int = dx, dy, dz
		var hitGround, hitCeiling, xResolved, yResolved bool
		var unresolved = true
		for unresolved {
			ax, ay, az, hitGround, hitCeiling, xResolved, yResolved =
				collider.ResolveCollision(ax, ay, az, charBlock, scene.Colliders())
			unresolved = xResolved || yResolved
			if xResolved {
				// fmt.Printf("resolved x %d - ", ax)
				charBlock.Translate(ax, 0, 0)
				ax = 0
			} else if yResolved {
				// fmt.Printf("resolved y %d - ", ay)
				charBlock.Translate(0, ay, 0)
				ay = 0
			} else {
				// fmt.Printf("resolved all three %d %d %d\n", ax, ay, az)
				charBlock.Translate(ax, ay, az)
			}
		}
		// fmt.Printf("dz %d az %d\n fallV %f\n", dz, az, fallV)
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
			fallV -= 0.3
		}

		shadowZ = scene.Colliders().FindFloor(charBlock)
	}

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
	x, y, z := charBlock.Pos()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %d\ny: %d\n z: %d\n shadz: %d\nrow: %d\ncol: %d\n pr: %d",
		x, y, z, shadowZ, int(y/16), int(x/16), int(z/16)))
	return nil
}

func main() {
	if err := ebiten.Run(
		update, screenW*3, screenH*3, 1, "Jumpin' Game",
	); err != nil {
		log.Fatal(err)
	}
}
