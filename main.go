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
	"enewey.com/golang-game/config"
)

var cX = 10
var cY = 10
var cZ = 0
var shadowZ = 0
var girl *actors.Actor
var gameScene *scene.Scene
var roomImage *ebiten.Image
var cfg *config.Config

func init() {
	cfg = config.Get()

	tiles := cache.Get().LoadSpritesheet("blue-walls.png", cfg.TileDimX, cfg.TileDimY)
	charas := cache.Get().LoadSpritesheet("hoodgirl.png", cfg.TileDimX, cfg.TileDimY)
	girlChar := charas.GetSprite(0)
	shadowChar := charas.GetSprite(1)
	charBlock := colliders.NewBlock(cX+3, cY+8, cZ, 10, 8, 12, "chara")
	girl = actors.NewActor("player", girlChar, shadowChar, charBlock)

	gameScene = scene.New(girl, cache.Get().LoadRoom("room1"), tiles, 0, 0)
	roomImage, _ = ebiten.NewImage(cfg.ScreenWidth()*2, cfg.ScreenHeight()*2, ebiten.FilterDefault)
}

var debug bool
var slownum int
func update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
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
		update, cfg.ScreenWidth()*3, cfg.ScreenHeight()*3, 1, "Jumpin' Game",
	); err != nil {
		log.Fatal(err)
	}
}
