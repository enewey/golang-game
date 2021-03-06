package colliders

import (
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/utils"
	"github.com/enewey/resolv/resolv"
)

// BodyType - set of flags to describe the physical body of the collider
type BodyType struct {
	blocking bool
}

// Collider - BaseZ gets the "root" Z level of the collider.
//            Depth is how many Z levels it spans.
//            Need 3 shapes for 3D collision (see the PreventCollisions method)
type Collider interface {
	X() int
	Y() int
	Z() int
	XDepth(int, int) int
	YDepth(int, int) int
	ZDepth(int, int) int
	Pos() (int, int, int)
	Center() (int, int, int)
	SetPos(int, int, int)
	Translate(int, int, int)
	Copy() Collider
	XYShape() resolv.Shape
	XZShape() resolv.Shape
	ZYShape() resolv.Shape
	Name() string
	Ref() int
	SetRef(int)
	IsBlocking() bool
	IsReactive(int) bool
	Reactions() *events.ReactionHub
}

// BaseCollider is an anonymous struct included in each Collider
type BaseCollider struct {
	xyshape     resolv.Shape
	xzshape     resolv.Shape
	zyshape     resolv.Shape
	name        string
	x, y, z     int
	ref         int
	bodyType    *BodyType
	reactionHub *events.ReactionHub
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

// Ref - reference identifier for this collider; typically matches an actor ID
func (b *BaseCollider) Ref() int {
	return b.ref
}

// SetRef - set the reference identifier for this collider. Used when registering an actor with an ID.
func (b *BaseCollider) SetRef(ref int) {
	b.ref = ref
}

// IsBlocking - tells whether this is a blocking collider or not.
func (b *BaseCollider) IsBlocking() bool {
	return b.bodyType.blocking
}

// IsReactive - indicates the collision behavior for this collider is custom.
func (b *BaseCollider) IsReactive(T int) bool {
	return b.Reactions().HasReactions(T)
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

// Reactions - retrieves the ReactionMux for this collider
func (b *BaseCollider) Reactions() *events.ReactionHub {
	return b.reactionHub
}

//
// ==========================================================
// ======== Colliders =======================================
// ==========================================================
//

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

// GetColliding - get all the colliders that would collide with the subject
// if the subject were to move by the specified deltas
func (cs Colliders) GetColliding(dx, dy, dz int, subject Collider) Colliders {
	var ret = make(Colliders, len(cs))
	i := 0
	for _, v := range cs {
		if subject.XYShape().WouldBeColliding(v.XYShape(), int32(dx), int32(dy)) &&
			subject.XZShape().WouldBeColliding(v.XZShape(), int32(dx), int32(dz)) &&
			subject.ZYShape().WouldBeColliding(v.ZYShape(), int32(dz), int32(dy)) {
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

// WouldCollide tests in what planes the subject WOULD collide if it were moved
// by the provided deltas.
func (cs Colliders) WouldCollide(dx, dy, dz int, subject Collider) bool {
	ret := false
	for _, v := range cs {
		if subject.XYShape().WouldBeColliding(v.XYShape(), int32(dx), int32(dy)) &&
			subject.XZShape().WouldBeColliding(v.XZShape(), int32(dx), int32(dz)) &&
			subject.ZYShape().WouldBeColliding(v.ZYShape(), int32(dz), int32(dy)) {
			ret = true
			break
		}
	}

	return ret
}

// TestXCollision - checks if a movement in the X direction for the subject would
// collide into the colliders. Returns the resolved dx that is safe, and whether
// or not the movement actually would collide. Does not translate the subject
// collider.
func (cs Colliders) TestXCollision(dx int, subject Collider) (int, bool) {
	// resolve on X axis
	zyfcoll := cs.getCollidingZY(subject)

	for _, v := range zyfcoll {
		resXY := resolv.Resolve(subject.XYShape(), v.XYShape(), int32(dx), 0)
		resXZ := resolv.Resolve(subject.XZShape(), v.XZShape(), int32(dx), 0)
		// x-collision occurred only if *both* shapes collide
		if resXY.Colliding() && resXZ.Colliding() {
			return utils.Min(int(resXY.ResolveX), int(resXZ.ResolveX)), true
		}
	}
	return dx, false
}

// TestYCollision returns the resolved y, whether a collision happened, and.
func (cs Colliders) TestYCollision(dy int, subject Collider) (int, bool) {
	// resolve on Y axis
	xzfcoll := cs.getCollidingXZ(subject)

	for _, v := range xzfcoll {
		resXY := resolv.Resolve(subject.XYShape(), v.XYShape(), 0, int32(dy))
		resZY := resolv.Resolve(subject.ZYShape(), v.ZYShape(), 0, int32(dy))
		// y-collision occurred only if *both* shapes collide
		if resXY.Colliding() && resZY.Colliding() {
			return utils.Min(int(resXY.ResolveY), int(resZY.ResolveY)), true
		}
	}
	return dy, false
}

// TestZCollision returns the resolved z, and if the collision happened. The
// two booleans represent hitGround and hitCeiling, respectively.
func (cs Colliders) TestZCollision(dz int, subject Collider) (int, bool, bool) {
	var rz = dz
	var hg, hc bool
	xyfcoll := cs.getCollidingXY(subject)

	for _, v := range xyfcoll {
		resXZ := resolv.Resolve(subject.XZShape(), v.XZShape(), 0, int32(dz))
		resZY := resolv.Resolve(subject.ZYShape(), v.ZYShape(), int32(dz), 0)
		// z-collision occurred only if *both* shapes collide
		if resXZ.Colliding() && resZY.Colliding() {
			// things get weird if more than one collision occurs, so keep track
			// and use the collision that resolves the smallest delta.
			z := utils.Min(int(resXZ.ResolveY), int(resZY.ResolveX))
			if utils.Abs(rz) > utils.Abs(z) {
				rz = z
				hg = dz < 0
				hc = dz > 0
			}

		}
	}
	return rz, hg, hc
}

// PreventCollision - checks if the subject would collide against the provided colliders.
//		if a collision would occur, translates the subject collider to prevent the collision.
//	Returns three booleans: hitGround, hitCeiling, and hitWall, and three ints,
//	x,y,z that describe how far the actor moved.
func (cs Colliders) PreventCollision(dx, dy, dz int, subject Collider) (bool, bool, bool, int, int, int) {
	var hitGround, hitCeiling, hitWallX, hitWallY bool
	var ax, ay, az = dx, dy, dz

	az, hitGround, hitCeiling = cs.TestZCollision(az, subject)
	subject.Translate(0, 0, az)

	ax, hitWallX = cs.TestXCollision(ax, subject)
	subject.Translate(ax, 0, 0)

	ay, hitWallY = cs.TestYCollision(ay, subject)
	subject.Translate(0, ay, 0)

	hitWall := hitWallX || hitWallY

	return hitGround, hitCeiling, hitWall, ax, ay, az
}

// ShoveCollision - shove all of the colliders aside in favor of the subject movement.
func (cs Colliders) ShoveCollision(dx, dy, dz int, subject Collider) Colliders {
	clone := subject.Copy()
	x, y, z := subject.Pos()
	clone.SetPos(dx+x, dy+y, dz+z)

	var sub Colliders = []Collider{clone}
	// shoved := cs.GetColliding(dx, dy, dz, subject)

	for _, v := range cs {
		// this is some hacky shit.
		// we translate the "shoved" collider by the movement of the heavy collider
		// and then attempt to move it back via prevent collision.
		// this is because Resolv does not try to prevent a collision
		// if the deltas are zero.
		// This will undoubtedly create some wacky edge cases, so need to do some
		// diligent testing around this stuff.
		vx, vy, vz := v.Pos()
		v.SetPos(vx+dx, vy+dy, vz+dz)

		var ax, ay, az = -dx, -dy, -dz

		sub.PreventCollision(ax, ay, az, v)
	}
	return cs
}

// Filter - filter func for Colliders.
// Returns a new slice where all the colliders return true for the test function.
func (cs Colliders) Filter(test func(Collider, int) bool) Colliders {
	ret := make([]Collider, len(cs))
	it := 0
	for i, v := range cs {
		if test(v, i) {
			ret[it] = v
			it++
		}
	}
	return ret[:it]
}

// ExcludeByCollider - return a new array excluding the given collider
func (cs Colliders) ExcludeByCollider(test Collider) Colliders {
	return cs.Filter(func(c Collider, i int) bool {
		return c != test
	})
}

// GetBlocking - returns a new slice of colliders which are blocking
func (cs Colliders) GetBlocking() Colliders {
	return cs.Filter(func(c Collider, i int) bool {
		return c.IsBlocking()
	})
}

// GetReactive - returns a new slice of colliders which are reactive
func (cs Colliders) GetReactive(T int) Colliders {
	return cs.Filter(func(c Collider, i int) bool {
		return c.IsReactive(T)
	})
}

// ExcludeByRef - returns a new slice of colliders without the selected ref
func (cs Colliders) ExcludeByRef(ref int) Colliders {
	return cs.Filter(func(c Collider, i int) bool {
		return c.Ref() != ref
	})
}
