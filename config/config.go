package config

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// Config provides global configurations for the game, primarily related to graphics.
type Config struct {
	TileDimX, TileDimY, TilesX, TilesY int
	gravity                            float64
	fontName                           string
}

var singer *Config

// Get returns a pointer to the singleton Config
func Get() *Config {
	if singer == nil {
		singer = &Config{
			16, 16, 10, 8, -0.25,
			"ponderosa.ttf",
		}
	}
	return singer
}

// ScreenHeight returns the calculated viewport height based on config
func (c *Config) ScreenHeight() int { return c.TileDimY * c.TilesY }

// ScreenWidth returns the calculated viewport width based on config
func (c *Config) ScreenWidth() int { return c.TileDimX * c.TilesX }

func (c *Config) scrollBoundLowerX() int {
	return (c.ScreenWidth() / 2) - (c.ScreenWidth() / 8) - (c.TileDimX / 2)
}
func (c *Config) scrollBoundUpperX() int {
	return (c.ScreenWidth() / 2) + (c.ScreenWidth() / 8) - (c.TileDimX / 2)
}
func (c *Config) scrollBoundLowerY() int {
	return (c.ScreenHeight() / 2) - (c.ScreenHeight() / 8) - (c.TileDimY / 2)
}
func (c *Config) scrollBoundUpperY() int {
	return (c.ScreenHeight() / 2) + (c.ScreenHeight() / 8) - (c.TileDimY / 2)
}

// Button configuration
const (
	ConfirmKey = iota
	CancelKey
	JumpKey
	DashKey
	UpKey
	DownKey
	LeftKey
	RightKey
)

func (c *Config) KeyUp() ebiten.Key      { return c.buttonSetting(UpKey) }
func (c *Config) KeyDown() ebiten.Key    { return c.buttonSetting(DownKey) }
func (c *Config) KeyLeft() ebiten.Key    { return c.buttonSetting(LeftKey) }
func (c *Config) KeyRight() ebiten.Key   { return c.buttonSetting(RightKey) }
func (c *Config) KeyConfirm() ebiten.Key { return c.buttonSetting(ConfirmKey) }
func (c *Config) KeyCancel() ebiten.Key  { return c.buttonSetting(CancelKey) }
func (c *Config) KeyJump() ebiten.Key    { return c.buttonSetting(JumpKey) }
func (c *Config) KeyDash() ebiten.Key    { return c.buttonSetting(DashKey) }

// ButtonSetting takes in a button function and returns what key it maps to.
func (c *Config) buttonSetting(k int) ebiten.Key {
	switch k {
	case ConfirmKey:
		return ebiten.KeyZ
	case CancelKey:
		return ebiten.KeyX
	case JumpKey:
		return ebiten.KeySpace
	case DashKey:
		return ebiten.KeyShift
	case UpKey:
		return ebiten.KeyUp
	case DownKey:
		return ebiten.KeyDown
	case LeftKey:
		return ebiten.KeyLeft
	case RightKey:
		return ebiten.KeyRight
	default:
		return ebiten.KeySpace
	}
}

// ScrollBoundaries are the U, R, D, L values indicating how far the character
// can walk in screen coordinates before scrolling can begin
func (c *Config) ScrollBoundaries() (int, int, int, int) {
	return c.scrollBoundLowerY(),
		c.scrollBoundUpperX(),
		c.scrollBoundUpperY(),
		c.scrollBoundLowerX()
}

// Gravity - the gravity coefficient in game
func (c *Config) Gravity() float64 {
	return c.gravity
}

func (c *Config) Font() string {
	return c.fontName
}

func (c *Config) TextSpeed() int {
	return 1
}

func (c *Config) WindowColor() color.Color {
	return color.Black
}
