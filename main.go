package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"

	"enewey.com/golang-game/src/room"
	"enewey.com/golang-game/src/tileset"
)

const tileWidth = 16
const tileHeight = 16

var tiles *tileset.Tileset
var scene *room.Room

func init() {
	var err error
	tiles, err = tileset.New("assets/img/blue-walls.png").Load()
	if err != nil {
		log.Fatal(err)
	}

	scene = room.NewRoomFromFile("assets/rooms/room1.room")
}

func drawTile(mapX, mapY, tileNum int, rm *ebiten.Image, tiles *tileset.Tileset) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(16*mapX), float64(16*mapY))
	// opt.GeoM.Scale(1.5, 1)

	rm.DrawImage(tiles.GetTileByNum(tileNum), opt)
}

func drawRoom() *ebiten.Image {
	rm, _ := ebiten.NewImage(320, 240, ebiten.FilterDefault)

	for _, layer := range scene.Layers() {
		mapTiles := layer.Tiles()
		for i := 0; i < len(mapTiles); i++ {
			row := int(i / 20)
			col := i % 20

			drawTile(col, row, mapTiles[i], rm, tiles)
		}
	}

	return rm
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	rm := drawRoom()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(2, 2)
	screen.DrawImage(rm, opt)
	return nil
}

func main() {
	if err := ebiten.Run(update, 640, 480, 1, "Render tiles"); err != nil {
		log.Fatal(err)
	}
}
