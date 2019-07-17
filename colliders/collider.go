package colliders

import (
	"enewey.com/golang-game/utils"
	"github.com/enewey/resolv/resolv"
)

// Collider - BaseZ gets the "root" Z level of the collider.
//            Depth is how many Z levels it spans.
//            Need 3 shapes for 3D collision (see the PreventCollisions method)
type Collider interface {
	X() int
	Y() int
	Z() int
	ZDepth(int, int) int
	Pos() (int, int, int)
	SetPos(int, int, int)
	Translate(int, int, int)
	XYShape() resolv.Shape
	XZShape() resolv.Shape
	ZYShape() resolv.Shape
	Name() string
}

// BaseCollider is an anonymous struct included in each Collider
type BaseCollider struct {
	xyshape resolv.Shape
	xzshape resolv.Shape
	zyshape resolv.Shape
	name    string
	x, y, z int
}

// X returns the root x position of this Collider
func (b *BaseCollider) X() int { return b.x }

// Y returns the root y position of this Collider
func (b *BaseCollider) Y() int { return b.y }

// Z returns the root z position of this Collider
func (b *BaseCollider) Z() int { return b.z }

// XYShape returns the shape of this collider on the XY plane
func (b *BaseCollider) XYShape() resolv.Shape { return b.xyshape }

// XZShape returns the shape of this collider on the XZ plane
func (b *BaseCollider) XZShape() resolv.Shape { return b.xzshape }

// ZYShape returns the shape of this collider on the YZ plane
func (b *BaseCollider) ZYShape() resolv.Shape { return b.zyshape }

// Name returns the name of this collider. Used primarily for debugging.
func (b *BaseCollider) Name() string { return b.name }

// Pos - get the x,y,z position
func (b *BaseCollider) Pos() (int, int, int) {
	return b.x, b.y, b.z
}

//x, y, z, w, h, d int
func (b *BaseCollider) setX(x int) {
	b.x = x
	b.xyshape.SetXY(int32(x), int32(b.y))
	b.xzshape.SetXY(int32(x), int32(b.z))
}

func (b *BaseCollider) setY(y int) {
	b.y = y
	b.xyshape.SetXY(int32(b.x), int32(y))
	b.zyshape.SetXY(int32(b.z), int32(y))
}

func (b *BaseCollider) setZ(z int) {
	b.z = z
	b.xzshape.SetXY(int32(b.x), int32(z))
	b.zyshape.SetXY(int32(z), int32(b.y))
}

// SetPos - set the x,y,z position of this collider
func (b *BaseCollider) SetPos(x, y, z int) {
	b.setX(x)
	b.setY(y)
	b.setZ(z)
}

// Translate - move the x,y,z position of this collider by a delta
func (b *BaseCollider) Translate(dx, dy, dz int) {
	cx, cy, cz := b.Pos()
	b.setX(dx + int(cx))
	b.setY(dy + int(cy))
	b.setZ(dz + int(cz))
}

// Colliders woo
type Colliders []Collider

func (cs Colliders) getCollidingXY(subject Collider) Colliders {
	var ret = make(Colliders, len(cs))
	i := 0
	for _, v := range cs {
		if v.XYShape().IsColliding(subject.XYShape()) {
			ret[i] = v
			i++
		}
	}
	return ret[:i]
}

func (cs Colliders) getCollidingXZ(subject Collider) Colliders {
	var ret = make(Colliders, len(cs))
	i := 0
	for _, v := range cs {
		if v.XZShape().IsColliding(subject.XZShape()) {
			ret[i] = v
			i++
		}
	}
	return ret[:i]
}

func (cs Colliders) getCollidingZY(subject Collider) Colliders {
	var ret = make(Colliders, len(cs))
	i := 0
	for _, v := range cs {
		if v.ZYShape().IsColliding(subject.ZYShape()) {
			ret[i] = v
			i++
		}
	}
	return ret[:i]
}

// FindFloor woo
func (cs Colliders) FindFloor(subject Collider) int {
	sx, sy, sz := subject.Pos()
	colls := cs.getCollidingXY(subject)
	var floorZ = -99
	for _, v := range colls {
		z := v.Z() + v.ZDepth(sx, sy)
		if z > floorZ && z <= sz {
			floorZ = z
		}
	}
	return floorZ
}

// PreventCollision - checks if the subject would collide against the provided colliders.
//		if a collision would occur, translates the subject collider to prevent the collision.
//	Returns three booleans: hitGround, hitCeiling, and hitWall.
func (cs Colliders) PreventCollision(dx, dy, dz int, subject Collider) (bool, bool, bool) {
	var hitGround, hitCeiling, hitWall bool
	var ax, ay, az = dx, dy, dz

	// resolve on Y axis
	xzfcoll := cs.getCollidingXZ(subject)

	for _, v := range xzfcoll {
		resXY := resolv.Resolve(subject.XYShape(), v.XYShape(), 0, int32(dy))
		resZY := resolv.Resolve(subject.ZYShape(), v.ZYShape(), 0, int32(dy))
		// z-collision occurred only if *both* shapes collide
		if resXY.Colliding() && resZY.Colliding() {
			ay = utils.Min(int(resXY.ResolveY), int(resZY.ResolveY))
			hitWall = true
			break
		}
	}
	subject.Translate(0, ay, 0)

	// resolve on X axis
	zyfcoll := cs.getCollidingZY(subject)

	for _, v := range zyfcoll {
		resXY := resolv.Resolve(subject.XYShape(), v.XYShape(), int32(dx), 0)
		resXZ := resolv.Resolve(subject.XZShape(), v.XZShape(), int32(dx), 0)
		// z-collision occurred only if *both* shapes collide
		if resXY.Colliding() && resXZ.Colliding() {
			ax = utils.Min(int(resXY.ResolveX), int(resXZ.ResolveX))
			hitWall = true
			break
		}
	}
	subject.Translate(ax, 0, 0)

	// resolve on z axis
	xyfcoll := cs.getCollidingXY(subject)

	for _, v := range xyfcoll {
		resXZ := resolv.Resolve(subject.XZShape(), v.XZShape(), 0, int32(dz))
		resZY := resolv.Resolve(subject.ZYShape(), v.ZYShape(), int32(dz), 0)
		// z-collision occurred only if *both* shapes collide
		if resXZ.Colliding() && resZY.Colliding() {
			az = utils.Min(int(resXZ.ResolveY), int(resZY.ResolveX))
			hitGround = dz < 0
			hitCeiling = dz > 0
			break
		}
	}

	subject.Translate(0, 0, az)
	return hitGround, hitCeiling, hitWall
}
