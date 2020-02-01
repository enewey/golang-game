package colliders

import (
	"enewey.com/golang-game/events"
	"github.com/enewey/resolv/resolv"
)

// Block - 3D rectangle
type Block struct {
	BaseCollider
	w, h, d int
}

// Width - x span
func (b *Block) Width() int { return b.w }

// Height - y span
func (b *Block) Height() int { return b.h }

// XDepth for Blocks, x span constant at any y/z
func (b *Block) XDepth(y, z int) int { return b.w }

// ZDepth for Blocks, z span is constant at any x/y
func (b *Block) ZDepth(x, y int) int { return b.d }

// YDepth for Blocks, y span is constant at any x/z
func (b *Block) YDepth(x, y int) int { return b.h }

// Center returns a point in the center of the rectangle
func (b *Block) Center() (int, int, int) {
	return b.x + (b.w / 2), b.y + (b.h / 2), b.z + (b.d / 2)
}

// NewBlock - creates a new 3D rectangle collider.
func NewBlock(x, y, z, w, h, d int, blocking, reactive bool, name string) Collider {
	b := &Block{w: w, h: h, d: d}
	b.x, b.y, b.z = x, y, z
	b.name = name
	b.xyshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(w), int32(h)))
	b.xzshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(z), int32(w), int32(d)))
	b.zyshape = resolv.Shape(resolv.NewRectangle(int32(z), int32(y), int32(d), int32(h)))
	b.ref = -1
	b.bodyType = &BodyType{blocking: blocking, reactive: reactive}
	b.reactionHub = events.NewReactionHub()
	return b
}
