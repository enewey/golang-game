package cache

import (
	"log"

	"enewey.com/golang-game/room"
	"enewey.com/golang-game/sprites"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Cache woo
type Cache struct {
	images map[string]*ebiten.Image
	rooms  map[string]*room.Room
	sheets map[string]*sprites.Spritesheet
}

const imgLoc = "assets/img/"
const roomLoc = "assets/rooms/"

var singer *Cache

// Get - get the cache singleton
func Get() *Cache {
	if singer == nil {
		singer = &Cache{
			make(map[string]*ebiten.Image),
			make(map[string]*room.Room),
			make(map[string]*sprites.Spritesheet),
		}
	}

	return singer
}

// LoadSpritesheet woo
func (c *Cache) LoadSpritesheet(src string, th, tw int) *sprites.Spritesheet {
	if (c.sheets[src] == nil) {
		c.sheets[src] = sprites.New(c.LoadImage(src), th, tw, 30, 30)
	}
	return c.sheets[src]
}

// LoadImage woo
func (c *Cache) LoadImage(src string) *ebiten.Image {
	if c.images[src] == nil {
		img, _, err := ebitenutil.NewImageFromFile(imgLoc+src, ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		c.images[src] = img
	}

	return c.images[src]
}

// LoadRoom woo
func (c *Cache) LoadRoom(name string) *room.Room {
	if c.rooms[name] == nil {
		c.rooms[name] = room.NewRoomFromFile(roomLoc + name + ".room")
	}

	return c.rooms[name]
}
