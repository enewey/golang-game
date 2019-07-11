package collider

import (
	"github.com/SolarLune/resolv/resolv"
)

// Collider - BaseZ gets the "root" Z level of the collider.
//            Depth is how many Z levels it spans.
type Collider struct {
	shape resolv.Shape
	zh    int
	zd    int
}

// Colliders woo
type Colliders []*Collider

// Shape woo
func (b *Collider) Shape() resolv.Shape { return b.shape }

// BaseZ woo
func (b *Collider) BaseZ() int { return b.zh }

// DepthZ woo
func (b *Collider) DepthZ() int { return b.zd }

// NewBlock woo
func NewBlock(x, y, w, h, zh, zd int) *Collider {
	return &Collider{
		resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(w), int32(h))),
		zh,
		zd,
	}
}

// GetGroup woo
func (cs Colliders) GetGroup(z int, tag string) *resolv.Space {
	ret := resolv.NewSpace()
	ret.AddTags(tag)
	for _, v := range cs {
		if v.zh <= z && (v.zh+v.zd) >= z {
			ret.Add(v.shape)
		}
	}
	return ret
}
