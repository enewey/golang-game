package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"

	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/cache"
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/scene"
	"enewey.com/golang-game/sprites"
)

var cX = 120
var cY = 100
var cZ = 0
var shadowZ = 0
var girl actors.Actor
var gameScene *scene.Scene
var roomImage *ebiten.Image
var cfg *config.Config

func init() {
	cfg = config.Get()

	tiles := cache.Get().LoadSpritesheet("blue-walls.png", cfg.TileDimX, cfg.TileDimY)
	charas := cache.Get().LoadSpritesheet("hoodgirl.png", cfg.TileDimX, cfg.TileDimY)
	girlChar := sprites.NewCharaSpritemap(
		charas.GetSprite(0),
		charas.GetSprite(30),
		charas.GetSprite(60),
		charas.GetSprite(90),
	)
	shadowChar := charas.GetSprite(1)
	charBlock := colliders.NewBlock(cX, cY, cZ, 10, 10, 14, 1, "chara")
	girl = actors.NewCharActor("player", girlChar, shadowChar, charBlock, -4, -8)

	gameScene = scene.New(girl, cache.Get().LoadRoom("longboy"), tiles)

	gameScene.AddActor(actors.NewStaticActor(
		"wall",
		sprites.NewStaticSpritemap(tiles.GetSprite(441)),
		colliders.NewBlock(81, 150, 0, 12, 8, 8, 1, fmt.Sprintf("manual-rock")),
		-2, -8,
	))
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
	if debug {
		slownum++
	}
	if slownum%4 == 0 {
		slownum = 0
		gameScene.Update(1)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	rm := gameScene.Render(roomImage)
	x, y, z := girl.Pos()

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
