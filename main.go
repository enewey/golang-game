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
	"enewey.com/golang-game/clock"
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/scene"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/types"
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
	// game initialization

	cfg = config.Get()
	roomImage, _ = ebiten.NewImage(cfg.ScreenWidth()*2, cfg.ScreenHeight()*2, ebiten.FilterDefault)

	// begin scene initialization

	tiles := cache.Get().LoadSpritesheet("blue-walls.png", cfg.TileDimX, cfg.TileDimY)
	charas := cache.Get().LoadSpritesheet("hoodgirl.png", cfg.TileDimX, cfg.TileDimY)
	girlChar := sprites.NewCharaSpritemap(
		charas.GetSprite(0),
		charas.GetSprite(30),
		charas.GetSprite(60),
		charas.GetSprite(90),
	)
	charBlock := colliders.NewBlock(cX, cY, cZ, 10, 10, 14, true, "chara")
	girl = actors.NewCharActor("player", girlChar, charBlock, -4, -8, 1)
	gameScene = scene.New(girl, "assets/rooms/v2.room.json")

	shadowChar := charas.GetSprite(1)
	shadow, hook := scene.CreateShadow(girl, shadowChar)
	gameScene.ActorM.AddActor(shadow)
	gameScene.ActorM.AddHook(hook)

	// end scene initialization

	// adding extraneous actors

	rock := scene.NewTrampoline(81, 150, 0, sprites.NewStaticSpritemap(tiles.GetSprite(441)))
	gameScene.ActorM.AddActor(rock)

	pushy := scene.NewPushBlock(128, 160, 32, "push-block-1", sprites.Create2by1Block(tiles.GetSprite(366), tiles.GetSprite(133)))
	gameScene.ActorM.AddActor(pushy)

	pushy2 := scene.NewPushBlock(200, 80, 0, "push-block-2", sprites.Create2by1Block(tiles.GetSprite(366), tiles.GetSprite(133)))
	rockCtrl := actors.NewMoveSequenceController([]types.Direction{types.Up, types.Left, types.Down, types.Right}, 16, 32)
	gameScene.ActorM.AddActorWithController(pushy2, rockCtrl)
}

var debug bool
var slownum int

type game struct {
	width, height int
}

func (g *game) Layout(ow, oh int) (int, int) {
	return cfg.ScreenWidth() * 3, cfg.ScreenHeight() * 3
}

func (g *game) Update(screen *ebiten.Image) error {
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
		clock.Inc(1)
		gameScene.Update(1)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	rm := gameScene.Render(roomImage)

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3, 3)
	screen.DrawImage(rm, opt)
	if debug {
		x, y, z := girl.Pos()
		vx, vy, vz := girl.(actors.CanMove).Vel()
		ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %d - %f\ny: %d - %f\nz: %d - %f", x, vx, y, vy, z, vz))
	}

	return nil
}

func main() {
	g := &game{cfg.ScreenWidth() * 3, cfg.ScreenHeight() * 3}
	w := g.width
	h := g.height

	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Jumpin' Game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
