package sprites

import (
	"enewey.com/golang-game/types"
)

// Spritemap allows mapping an int to sprites -- handy for use with the Direction iota
type Spritemap interface {
	Sprite(id int) *Sprite
}

// CharaMap c
type CharaMap map[int]*Sprite

// NewCharaSpritemap returns a new 4 directional spritemap for an actor
func NewCharaSpritemap(d, r, u, l *Sprite) *CharaMap {
	return &CharaMap{
		int(types.Up):        u,
		int(types.Down):      d,
		int(types.Right):     r,
		int(types.Left):      l,
		int(types.UpRight):   u,
		int(types.UpLeft):    u,
		int(types.DownRight): d,
		int(types.DownLeft):  d,
	}
}

// Sprite returns the Sprite with the given ID
func (sm *CharaMap) Sprite(id int) *Sprite {
	return (*sm)[id]
}

// StaticMap s
type StaticMap struct {
	s *Sprite
}

// NewStaticSpritemap returns a new single-sprite map.
func NewStaticSpritemap(s *Sprite) *StaticMap {
	return &StaticMap{s}
}

// Sprite for a StaticMap always returns the same sprite
func (sm StaticMap) Sprite(id int) *Sprite {
	return sm.s
}
