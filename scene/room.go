package scene

import (
	"fmt"

	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/cache"
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/sprites"
)

type room struct {
	Width, Height int
	actors        []actors.Actor
}

func createRoom(dat *Data) *room {
	var guys = make([]actors.Actor, len(dat.Actors))
	for i, adat := range dat.Actors {
		sprite := loadSpriteData(adat.Sprite)
		collider := loadColliderData(adat.Collider)

		var a actors.Actor
		switch adat.Kind {
		case "static":
			a = actors.NewStaticActor(
				adat.Name,
				sprite,
				collider,
				adat.OffsetX,
				adat.OffsetY,
			)
		}
		guys[i] = a
	}
	return &room{dat.Width, dat.Height, guys}
}

func loadSpriteData(dat *spriteData) sprites.Spritemap {
	cfg := config.Get()
	var s sprites.Spritemap

	switch dat.Kind {
	case "compound":
		tiles := make([]*sprites.Sprite, len(dat.Tiles))
		for i, tile := range dat.Tiles {
			tiles[i] = cache.Get().LoadSpritesheet(dat.Sheet, cfg.TileDimX, cfg.TileDimY).GetSprite(tile)
		}
		s = sprites.NewStaticSpritemap(sprites.NewCompoundSprite(tiles, dat.Rows, dat.Cols, cfg.TileDimX, cfg.TileDimY))
	}

	return s
}

func loadColliderData(dat *colliderData) colliders.Collider {
	fmt.Printf("loading collider data: %+v : ", dat)
	if dat.BlockColliderData != nil {

		data := dat.BlockColliderData
		fmt.Printf("block collider %+v\n", data)
		return colliders.NewBlock(dat.X, dat.Y, dat.Z, data.W, data.H, dat.D, dat.Blocking, dat.Name)
	}

	if dat.TriangleColliderData != nil {
		data := dat.TriangleColliderData
		return colliders.NewTriangle(
			dat.X, dat.Y, dat.Z,
			data.Rx2, data.Ry2, data.Rx3, data.Ry3,
			dat.D, data.Axis, dat.Blocking, dat.Name)
	}

	return nil
}
