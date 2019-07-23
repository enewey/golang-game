package config

// Config provides global configurations for the game, primarily related to graphics.
type Config struct {
	TileDimX, TileDimY, TilesX, TilesY int
}

var singer *Config

// Get returns a pointer to the singleton Config
func Get() *Config {
	if singer == nil {
		singer = &Config{
			16, 16, 10, 8,
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

// ScrollBoundaries are the U, R, D, L values indicating how far the character
// can walk in screen coordinates before scrolling can begin
func (c *Config) ScrollBoundaries() (int, int, int, int) {
	return c.scrollBoundLowerY(),
		c.scrollBoundUpperX(),
		c.scrollBoundUpperY(),
		c.scrollBoundLowerX()
}
