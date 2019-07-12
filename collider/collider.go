package collider

import (
	"github.com/SolarLune/resolv/resolv"
)

// Collider - BaseZ gets the "root" Z level of the collider.
//            Depth is how many Z levels it spans.
type Collider struct {
	xyshape          resolv.Shape
	xzshape          resolv.Shape
	zyshape          resolv.Shape
	x, y, z, w, h, d int
}

func (b *Collider) setX(x int) {
	b.x = x
	b.xyshape.SetXY(int32(x), int32(b.y))
	b.xzshape.SetXY(int32(x), int32(b.z))
}

func (b *Collider) setY(y int) {
	b.y = y
	b.xyshape.SetXY(int32(b.x), int32(y))
	b.zyshape.SetXY(int32(b.z), int32(y))
}

func (b *Collider) setZ(z int) {
	b.z = z
	b.xzshape.SetXY(int32(b.x), int32(z))
	b.zyshape.SetXY(int32(z), int32(b.y))
}

// GetPos - get the x,y,z position
func (b *Collider) GetPos() (int, int, int) {
	return b.x, b.y, b.z
}

// SetPos - set the x,y,z position of this collider
func (b *Collider) SetPos(x, y, z int) {
	b.setX(x)
	b.setY(y)
	b.setZ(z)
}

// Translate - move the x,y,z position of this collider by a delta
func (b *Collider) Translate(dx, dy, dz int) {
	cx, cy := b.xyshape.GetXY()
	_, cz := b.xzshape.GetXY()
	b.setX(dx + int(cx))
	b.setY(dy + int(cy))
	b.setZ(dz + int(cz))
}

// NewBlock woo
func NewBlock(x, y, z, w, h, d int) *Collider {
	return &Collider{
		x: x, y: y, z: z,
		w: w, h: h, d: d,
		xyshape: resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(w), int32(h))),
		xzshape: resolv.Shape(resolv.NewRectangle(int32(x), int32(z), int32(w), int32(d))),
		zyshape: resolv.Shape(resolv.NewRectangle(int32(z), int32(y), int32(d), int32(h))),
	}
}

// Colliders woo
type Colliders []*Collider

// GetXYGroup woo
func (cs Colliders) getXYGroup(tag string) *resolv.Space {
	ret := resolv.NewSpace()
	ret.AddTags(tag)
	for _, b := range cs {
		ret.Add(b.xyshape)
	}
	return ret
}

// GetXZGroup woo
func (cs Colliders) getXZGroup(tag string) *resolv.Space {
	ret := resolv.NewSpace()
	ret.AddTags(tag)
	for _, b := range cs {
		ret.Add(b.xzshape)
	}
	return ret
}

// GetZYGroup woo
func (cs Colliders) getZYGroup(tag string) *resolv.Space {
	ret := resolv.NewSpace()
	ret.AddTags(tag)
	for _, b := range cs {
		ret.Add(b.zyshape)
	}
	return ret
}

const (
	exy = iota
	exz
	ezy
)

// SpaceCache woo
type SpaceCache struct {
	colliders Colliders
	cache     map[int]*resolv.Space
}

// NewSpaceCache woo
func NewSpaceCache(colls Colliders) *SpaceCache {
	return &SpaceCache{colls, make(map[int]*resolv.Space)}
}

func (c *SpaceCache) getPlanes(tag string) (*resolv.Space, *resolv.Space, *resolv.Space) {
	if c.cache[exy] == nil {
		c.cache[exy] = c.colliders.getXYGroup(tag)
	}
	if c.cache[exz] == nil {
		c.cache[exz] = c.colliders.getXZGroup(tag)
	}
	if c.cache[ezy] == nil {
		c.cache[ezy] = c.colliders.getZYGroup(tag)
	}

	return c.cache[exy], c.cache[exz], c.cache[ezy]
}

// ResolveCollision woo
func ResolveCollision(dx, dy, dz int, subject *Collider, cache *SpaceCache) (int, int, int, bool) {
	var rx, ry, rz int
	var onGround bool
	_, xzgroup, zygroup := cache.getPlanes("walls")

	resX := xzgroup.Resolve(subject.xzshape, int32(dx), 0)
	if resX.Colliding() {
		rx = int(resX.ResolveX)
	} else {
		rx = dx
	}

	resY := zygroup.Resolve(subject.zyshape, 0, int32(dy))
	if resY.Colliding() {
		ry = int(resY.ResolveY)
	} else {
		ry = dy
	}

	resZ := xzgroup.Resolve(subject.xzshape, 0, int32(dz))
	if resZ.Colliding() {
		rz = int(resZ.ResolveY)
		onGround = true
	} else {
		rz = dz
		onGround = false
	}

	return rx, ry, rz, onGround
}
