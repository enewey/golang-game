package colliders

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

// Pos - get the x,y,z position
func (b *Collider) Pos() (int, int, int) {
	return b.x, b.y, b.z
}

// Width - x span
func (b *Collider) Width() int { return b.w }

// Height - y span
func (b *Collider) Height() int { return b.h }

// Depth - z span
func (b *Collider) Depth() int { return b.d }

// SetPos - set the x,y,z position of this collider
func (b *Collider) SetPos(x, y, z int) {
	b.setX(x)
	b.setY(y)
	b.setZ(z)
}

// Translate - move the x,y,z position of this collider by a delta
func (b *Collider) Translate(dx, dy, dz int) {
	cx, cy, cz := b.Pos()
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

func (cs Colliders) getCollidingXY(subject *Collider) Colliders {
	var ret Colliders
	for _, v := range cs {
		if v.xyshape.IsColliding(subject.xyshape) {
			ret = append(ret, v)
		}
	}
	return ret
}

// FindFloor woo
func (cs Colliders) FindFloor(subject *Collider) int {
	_, _, sz := subject.Pos()
	colls := cs.getCollidingXY(subject)
	var floorZ = -99
	for _, v := range colls {
		z := v.z + v.d
		if z > floorZ && z <= sz {
			floorZ = z
		}
	}
	return floorZ
}

// ResolveCollision woo
func ResolveCollision(dx, dy, dz int, subject *Collider, colliders Colliders) (int, int, int, bool, bool, bool, bool) {
	var rx, ry, rz int = dx, dy, dz
	var hitGround, hitCeiling bool

	// to resolve the XY collision, filter out the colliders that are NOT in the
	// range of Z that we care about.
	filterColls := colliderFilter(
		subject.z,
		subject.z+subject.d,
		colliders,
		filterByZRange)
	xygroup := filterColls.getXYGroup("walls")

	// now that we have our group of XY shapes we care about, resolve the deltas
	resX := xygroup.Resolve(subject.xyshape, int32(dx), 0)
	if resX.Colliding() {
		rx = int(resX.ResolveX)
		return rx, ry, rz, hitGround, hitCeiling, true, false
	}

	resY := xygroup.Resolve(subject.xyshape, 0, int32(dy))
	if resY.Colliding() {
		ry = int(resY.ResolveY)
		return rx, ry, rz, hitGround, hitCeiling, false, true
	}

	for _, v := range colliders {
		resXZ := resolv.Resolve(subject.xzshape, v.xzshape, 0, int32(dz))
		resZY := resolv.Resolve(subject.zyshape, v.zyshape, int32(dz), 0)
		// z-collision occurred only if *both* shapes collide
		if resXZ.Colliding() && resZY.Colliding() {
			rz = int(resXZ.ResolveY)
			hitGround = dz < 0
			hitCeiling = dz > 0
			break
		}
	}

	return rx, ry, rz, hitGround, hitCeiling, false, false
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