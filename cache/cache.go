package cache

import (
	"io/ioutil"
	"log"

	"golang.org/x/image/font"

	"enewey.com/golang-game/sprites"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Cache woo
type Cache struct {
	images map[string]*ebiten.Image
	sheets map[string]*sprites.Spritesheet
	fonts  map[string]font.Face
}

const imgLoc = "assets/img/"
const roomLoc = "assets/rooms/"
const fontLoc = "assets/fonts/"

var singer *Cache

// Get - get the cache singleton
func Get() *Cache {
	if singer == nil {
		singer = &Cache{
			make(map[string]*ebiten.Image),
			make(map[string]*sprites.Spritesheet),
			make(map[string]font.Face),
		}
	}

	return singer
}

// LoadSpritesheet woo
func (c *Cache) LoadSpritesheet(src string, th, tw int) *sprites.Spritesheet {
	if c.sheets[src] == nil {
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

// LoadFont loads a font face from the cache.
// TODO: make font size adjustable and dpi and shit
func (c *Cache) LoadFont(src string) font.Face {
	if c.fonts[src] == nil {
		bytes, err := ioutil.ReadFile(fontLoc + src)
		if err != nil {
			panic(err)
		}
		tt, err := truetype.Parse(bytes)
		if err != nil {
			panic(err)
		}
		face := truetype.NewFace(tt, &truetype.Options{
			Size:       10,
			DPI:        72,
			Hinting:    font.HintingFull,
			SubPixelsX: 0,
			SubPixelsY: 0,
		})
		c.fonts[src] = face
	}
	return c.fonts[src]
}
