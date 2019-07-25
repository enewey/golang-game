package colliders

import "github.com/enewey/resolv/resolv"

// Block - 3D rectangle
type Block struct {
	BaseCollider
	w, h, d int
}

// Width - x span
func (b *Block) Width() int { return b.w }

// Height - y span
func (b *Block) Height() int { return b.h }

// ZDepth for Blocks, z span is constant at any x/y
func (b *Block) ZDepth(x, y int) int { return b.d }

// NewBlock - creates a new 3D rectangle collider.
func NewBlock(x, y, z, w, h, d int, name string) Collider {
	b := &Block{w: w, h: h, d: d}
	b.x, b.y, b.z = x, y, z
	b.name = name
	b.xyshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(w), int32(h)))
	b.xzshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(z), int32(w), int32(d)))
	b.zyshape = resolv.Shape(resolv.NewRectangle(int32(z), int32(y), int32(d), int32(h)))
	return b
}
