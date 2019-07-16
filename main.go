package main

import (
	_ "image/png"
	"log"
	"fmt"
	
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"enewey.com/golang-game/cache"
	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/scene"
)

var cX = 50
var cY = 76
var cZ = 0
var shadowZ = 0
var girl *actors.Actor
var gameScene *scene.Scene
var roomImage *ebiten.Image

func init() {
	tiles := cache.Get().LoadSpritesheet("blue-walls.png", TILE_DIMX, TILE_DIMY)
	charas := cache.Get().LoadSpritesheet("hoodgirl.png", TILE_DIMX, TILE_DIMY)
	girlChar := charas.GetSprite(0)
	shadowChar := charas.GetSprite(1)
	charBlock := colliders.NewBlock(cX+3, cY+8, cZ, 10, 8, 12, "chara")
	girl = actors.NewActor("player", girlChar, shadowChar, charBlock)

	gameScene = scene.New(girl, cache.Get().LoadRoom("room2"), tiles)
	roomImage, _ = ebiten.NewImage(SCREEN_W*2, SCREEN_H*2, ebiten.FilterDefault)
}

var debug bool
var slownum int
func update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		debug = !debug
		slownum = 0
	}
	if debug { slownum++ }
	if slownum % 4 == 0 {
		slownum = 0
		gameScene.Update(1)
	}
	
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	
	rm := gameScene.Render(roomImage)
	x,y,z := girl.Pos()

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3, 3)
	screen.DrawImage(rm, opt)
	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %d\ny: %d\nz: %d", x, y, z))
	}
		
	return nil
}

func main() {
	if err := ebiten.Run(
		update, SCREEN_W*3, SCREEN_H*3, 1, "Jumpin' Game",
	); err != nil {
		log.Fatal(err)
	}
}
