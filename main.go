package main

import (
	// "fmt"
	_ "image/png"
	"log"
	// "math"

	"github.com/hajimehoshi/ebiten"
	// "github.com/hajimehoshi/ebiten/ebitenutil"
	// "github.com/hajimehoshi/ebiten/inpututil"

	"enewey.com/golang-game/cache"
	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/colliders"
	// "enewey.com/golang-game/room"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/scene"
)

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
var gameScene *scene.Scene
var roomImage *ebiten.Image

var charBlock *colliders.Collider

func init() {
	tiles = cache.Get().LoadSpritesheet("blue-walls.png", TILE_DIMX, TILE_DIMY)
	charas = cache.Get().LoadSpritesheet("hoodgirl.png", TILE_DIMX, TILE_DIMY)
	girlChar = charas.GetSprite(0)
	shadowChar = charas.GetSprite(1)
	charBlock = colliders.NewBlock(cX+3, cY+8, cZ, 10, 8, 12, "chara")
	girl := actors.NewActor("player", girlChar, shadowChar, charBlock)

	gameScene = scene.New(girl, cache.Get().LoadRoom("room2"), tiles)
	roomImage, _ = ebiten.NewImage(SCREEN_W*2, SCREEN_H*2, ebiten.FilterDefault)
}

// func drawSprite(x, y int, sprite *sprites.Sprite, rm *ebiten.Image) {
// 	opt := &ebiten.DrawImageOptions{}
// 	opt.GeoM.Translate(float64(x), float64(y))

// 	rm.DrawImage(sprite.Img(), opt)
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func drawRoom(img *ebiten.Image) *ebiten.Image {
// 	for _, layer := range scene.Layers() {
// 		pr := layer.Priority() // z-layer of this tile, basically
// 		for row := 0; row < TILESY; row++ {
// 			mapTiles := layer.TilesRow(row, TILESX)
// 			for col := 0; col < len(mapTiles); col++ {
// 				tile := tiles.GetSprite(mapTiles[col])
// 				drawSprite(col*TILE_DIMX, row*TILE_DIMY, tile, img)
// 			}

// 			sx, sy, sz := charBlock.Pos()
// 			sd := charBlock.Depth()
// 			charPr := int(math.Ceil(float64(sz+1) / float64(TILE_DIMY/2)))
// 			shadowPr := int(math.Floor(float64(shadowZ+sd) / float64(TILE_DIMY/2)))
// 			charRow := int(math.Ceil(float64(sy) / float64(TILE_DIMY)))

// 			if shadowPr == pr && max(charRow-pr, 0) == row {
// 				drawSprite(sx-4, sy-shadowZ-8, shadowChar, img)
// 			}
// 			if charPr == pr && max(charRow-pr, 0) == row {
// 				drawSprite(sx-4, sy-sz-8, girlChar, img)
// 			}
// 		}
// 	}

// 	return roomImage
// }

// ---------

// var logicFrame int
// var slowdownSwitch bool

// func checkInputs() {
// 	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && onGround {
// 		fallV = 4.0
// 		onGround = false
// 	}
// 	var dx, dy, dz int

// 	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
// 		ebiten.SetFullscreen(!ebiten.IsFullscreen())
// 	}
// 	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
// 		slowdownSwitch = !slowdownSwitch
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyRight) {
// 		dx++
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
// 		dx--
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyDown) {
// 		dy++
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyUp) {
// 		dy--
// 	}

// 	logicFrame = (logicFrame + 1) % 4
// 	if (logicFrame == 0 && slowdownSwitch) || !slowdownSwitch {
// 		dz = int(math.Max(fallV, -6)) // per second?? frame?? :thinking_face:
// 		var ax, ay, az int = dx, dy, dz
// 		var hitGround, hitCeiling, xResolved, yResolved bool
// 		var unresolved = true
// 		for unresolved {
// 			ax, ay, az, hitGround, hitCeiling, xResolved, yResolved =
// 				colliders.ResolveCollision(ax, ay, az, charBlock, scene.Colliders())
// 			unresolved = xResolved || yResolved
// 			if xResolved {
// 				// fmt.Printf("resolved x %d - ", ax)
// 				charBlock.Translate(ax, 0, 0)
// 				ax = 0
// 			} else if yResolved {
// 				// fmt.Printf("resolved y %d - ", ay)
// 				charBlock.Translate(0, ay, 0)
// 				ay = 0
// 			} else {
// 				// fmt.Printf("resolved all three %d %d %d\n", ax, ay, az)
// 				charBlock.Translate(ax, ay, az)
// 			}
// 		}
// 		// fmt.Printf("dz %d az %d\n fallV %f\n", dz, az, fallV)
// 		if hitGround {
// 			onGround = true
// 		} else if az != 0 {
// 			onGround = false
// 		}
// 		if onGround && fallV < -1 && az >= 0 {
// 			fallV = 0
// 		} else if hitCeiling && fallV > 0 && az <= 0 {
// 			fallV = 0
// 		} else {
// 			fallV -= 0.3
// 		}

// 		shadowZ = scene.Colliders().FindFloor(charBlock)
// 	}

// }

// -------

func update(screen *ebiten.Image) error {
	// checkInputs()
	gameScene.Update(1)

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	rm := gameScene.Render(roomImage)
	opt := &ebiten.DrawImageOptions{}
	// opt.GeoM.Translate(-float64(SCREEN_W)/2.0, -float64(SCREEN_H)/2.0)
	opt.GeoM.Scale(3, 3)
	screen.DrawImage(rm, opt)
	// x, y, z := charBlock.Pos()
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %d\ny: %d\n z: %d\n shadz: %d\nrow: %d\ncol: %d\n pr: %d",
	// 	x, y, z, shadowZ, int(math.Ceil(float64(y)/float64(TILE_DIMY))),
	// 	int(x/16), int(math.Ceil(float64(z+1)/float64(TILE_DIMY/2)))))
	return nil
}

func main() {
	if err := ebiten.Run(
		update, SCREEN_W*3, SCREEN_H*3, 1, "Jumpin' Game",
	); err != nil {
		log.Fatal(err)
	}
}
