package colliders

import (
	"fmt"
	"github.com/enewey/resolv/resolv"
	"enewey.com/golang-game/utils"
)

// Collider - BaseZ gets the "root" Z level of the collider.
//            Depth is how many Z levels it spans.
//            Need 3 shapes for 3D collision (see the PreventCollisions method)
type Collider interface {
	X() int
	Y() int
	Z() int
	ZDepth(int,int) int
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
	xyshape          resolv.Shape
	xzshape          resolv.Shape
	zyshape          resolv.Shape
	name             string
	x,y,z int
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

// Block - 3D rectangle
type Block struct {
	BaseCollider
	w,h,d int
}

// Width - x span
func (b *Block) Width() int { return b.w }

// Height - y span
func (b *Block) Height() int { return b.h }

// ZDepth - z span is constant at any x/y
func (b *Block) ZDepth(x,y int) int { return b.d }

// NewBlock - creates a new 3D rectangle collider.
func NewBlock(x, y, z, w, h, d int, name string) Collider {
	b := &Block{w: w, h: h, d: d}
	b.x,b.y,b.z = x,y,z
	b.name = name
	b.xyshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(w), int32(h)))
	b.xzshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(z), int32(w), int32(d)))
	b.zyshape = resolv.Shape(resolv.NewRectangle(int32(z), int32(y), int32(d), int32(h)))
	return b
}

// Triangle - 3D triangular prism along the z-axis
type Triangle struct {
	BaseCollider
	// r variables mean they are RELATIVE to the BaseCollider x,y,z
	rx2,ry2,rx3,ry3,d,axis int
}

// // X2 gets the X-coord for the 2nd point of this triangle
// func (b *Triangle) X2() int { return b.x2 }
// // Y2 gets the Y-coord for the 2nd point of this triangle
// func (b *Triangle) Y2() int { return b.y2 }
// // X3 gets the X-coord for the 3rd point of this triangle
// func (b *Triangle) X3() int { return b.x3 }
// // Y3 gets the Y-coord for the 3rd point of this triangle
// func (b *Triangle) Y3() int { return b.y3 }

// ZDepth - z span at the given point
func (b *Triangle) ZDepth(x,y int) int {
	// switch (b.axis) {
	// case ZAxis:
	// 	return b.d
	// case XAxis:

	// case YAxis:
	// }
	return b.d
}

// Axis constants, to help build triangle prisms.
const (
	XAxis = iota
	YAxis
	ZAxis
)

// NewTriangle - creates a new 3D triangular prism, where the prism runs along 
// the specified axis. The named variables (x1, y1, z etc) are named as if the
// prism runs "d" (depth) units along the Z axis.
// When the XAxis is specified: x coordinates will map to the Z axis
// When the YAxis is specified: y coordinates will map to the Z axis
func NewTriangle(x, y, z, rx2, ry2, rx3, ry3, d, axis int, name string) Collider {
	tri := &Triangle{rx2: rx2, ry2: ry2, rx3: rx3, ry3: ry3, d: d, axis: axis}
	tri.x, tri.y, tri.z = x, y, z
	tri.name = name

	switch axis {
	case XAxis:
		//fmt.Printf("parsed x triangle")
		w := d
		h := utils.Max(y, y+ry2, y+ry3) - utils.Min(y, y+ry2, y+ry3)
		tri.zyshape = resolv.NewTriangle(
			int32(z), int32(y),
			int32(z+rx2), int32(y+ry2),
			int32(z+rx3), int32(y+ry3))
		tri.xzshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(z), int32(w), int32(d)))
		tri.xyshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(d), int32(h)))
	case YAxis:
		w := utils.Max(x, x+rx2, x+rx3) - utils.Min(x, x+rx2, x+rx3)
		h := d
		tri.xzshape = resolv.NewTriangle(
			int32(x), int32(z), 
			int32(x+rx2), int32(z+ry2), 
			int32(x+rx3), int32(z+ry3))
		tri.xyshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(y), int32(w), int32(d)))
		tri.zyshape = resolv.Shape(resolv.NewRectangle(int32(z), int32(y), int32(d), int32(h)))
	case ZAxis:
		w := utils.Max(x, x+rx2, x+rx3) - utils.Min(x, x+rx2, x+rx3)
		h := utils.Max(y, y+ry2, y+ry3) - utils.Min(y, y+ry2, y+ry3)
		tri.xyshape = resolv.NewTriangle(
			int32(x), int32(y), 
			int32(x+rx2), int32(y+ry2), 
			int32(x+rx3), int32(y+ry3))
		tri.xzshape = resolv.Shape(resolv.NewRectangle(int32(x), int32(z), int32(w), int32(d)))
		tri.zyshape = resolv.Shape(resolv.NewRectangle(int32(z), int32(y), int32(d), int32(h)))
	}

	if name == "platform_chunk4" {
		fmt.Printf("triangle made %s axis %d xy %v xz %v zy %v\n", tri.name, tri.axis, tri.xyshape, tri.xzshape, tri.zyshape )
	}
	

	return tri
}

// Colliders woo
type Colliders []Collider

// // GetXYGroup woo
// func (cs Colliders) getXYGroup(tag string) *resolv.Space {
// 	ret := resolv.NewSpace()
// 	ret.AddTags(tag)
// 	for _, b := range cs {
// 		ret.Add(b.XYShape())
// 	}
// 	return ret
// }

// // GetXZGroup woo
// func (cs Colliders) getXZGroup(tag string) *resolv.Space {
// 	ret := resolv.NewSpace()
// 	ret.AddTags(tag)
// 	for _, b := range cs {
// 		ret.Add(b.XZShape())
// 	}
// 	return ret
// }

// // GetZYGroup woo
// func (cs Colliders) getZYGroup(tag string) *resolv.Space {
// 	ret := resolv.NewSpace()
// 	ret.AddTags(tag)
// 	for _, b := range cs {
// 		ret.Add(b.ZYShape())
// 	}
// 	return ret
// }

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

	// // to resolve the XY collision, filter out the colliders that are NOT in the
	// // range of Z that we care about.
	// filterColls := colliderFilter(
	// 	subject.Z(),
	// 	subject.Z()+subject.Depth(),
	// 	cs,
	// 	filterByZRange)
	// xygroup := filterColls.getXYGroup("walls")

	// // now that we have our group of XY shapes we care about, resolve the deltas
	// // do X and Y as two individual checks to prevent stupid crap like jumping
	// // into corners and falling through the floor.
	// resX := xygroup.Resolve(subject.XYShape(), int32(dx), 0)
	// if resX.Colliding() {
	// 	subject.Translate(int(resX.ResolveX), 0, 0)
	// 	ax = 0
	// 	hitWall = true
	// }
	// if ax != 0 {
	// 	subject.Translate(ax, 0, 0)
	// }

	// resY := xygroup.Resolve(subject.XYShape(), 0, int32(dy))
	// if resY.Colliding() {
	// 	subject.Translate(0, int(resY.ResolveY), 0)
	// 	ay = 0
	// 	hitWall = true
	// }
	// if ay != 0 {
	// 	subject.Translate(0, ay, 0)
	// }

	// Now for Z collisions, imagine the XZ plane (camera facing down the Y axis
	// where the horizon is the X axis)	and the ZY plane (camera facing up the
	// X axis, where the horizon is the Z axis). If a collider's XZ and ZY
	// planes *both* collide with the subject, then there is a Z collision.
	//
	// This is a roundabout way to fit this square peg into a round hole
	// (the resolv library is only meant for 2D, not 3D, collision) but it works
	// pretty nicely.
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

// func filterByZRange(zmin, zmax int, collider Collider) bool {
// 	return collider.Z()+collider.Depth() > zmin && collider.Z() < zmax
// }

// func colliderFilter(min, max int, arr Colliders, f func(int, int, Collider) bool) Colliders {
// 	var ret Colliders
// 	for _, v := range arr {
// 		if f(min, max, v) {
// 			ret = append(ret, v)
// 		}
// 	}
// 	return ret
// }
