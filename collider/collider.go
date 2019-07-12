package collider

import (
	"fmt"
	"math"

	"github.com/SolarLune/resolv/resolv"
)

// Collider - BaseZ gets the "root" Z level of the collider.
//            Depth is how many Z levels it spans.
type Collider struct {
	xyshape          resolv.Shape
	xzshape          resolv.Shape
	zyshape          resolv.Shape
	x, y, z, w, h, d int
	name             string
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
func NewBlock(x, y, z, w, h, d int, name string) *Collider {
	return &Collider{
		x: x, y: y, z: z,
		w: w, h: h, d: d,
		name:    name,
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
	var hitGround bool
	// _, xzgroup, zygroup := cache.getPlanes("walls")

	for _, v := range cache.colliders {
		resXZ := resolv.Resolve(subject.xzshape, v.xzshape, 0, int32(dz))
		resZY := resolv.Resolve(subject.zyshape, v.zyshape, int32(dz), 0)
		if resXZ.Colliding() && resZY.Colliding() {
			if math.Abs(float64(resXZ.ResolveY)) < math.Abs(float64(resZY.ResolveX)) {
				rz = int(resXZ.ResolveY)
			} else {
				rz = int(resZY.ResolveX)
			}
			hitGround = true
			break
		}
	}
	if !hitGround {
		rz = dz
	}

	filterColls := colliderFilter(
		subject.z,
		subject.z+subject.d,
		cache.colliders,
		filterByZRange)
	xygroup := filterColls.getXYGroup("walls")

	fmt.Printf("before %d after %d\n", len(cache.colliders), len(filterColls))
	resX := xygroup.Resolve(subject.xyshape, int32(dx), 0)
	if resX.Colliding() {
		rx = int(resX.ResolveX)
	} else {
		rx = dx
	}

	resY := xygroup.Resolve(subject.xyshape, 0, int32(dy))
	if resY.Colliding() {
		ry = int(resY.ResolveY)
	} else {
		ry = dy
	}

	return rx, ry, rz, hitGround
}

func filterByZRange(zmin, zmax int, collider *Collider) bool {
	return collider.z+collider.d > zmin && collider.z < zmax
}

func colliderFilter(zmin, zmax int, arr Colliders, f func(int, int, *Collider) bool) Colliders {
	var ret Colliders
	for _, v := range arr {
		if f(zmin, zmax, v) {
			ret = append(ret, v)
		}
	}
	return ret
}
