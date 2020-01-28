package colliders

import (
	"math"

	"enewey.com/golang-game/types"
	"enewey.com/golang-game/utils"
	"github.com/enewey/resolv/resolv"
)

// Triangle - 3D triangular prism along the z-axis
type Triangle struct {
	BaseCollider
	// r variables mean they are RELATIVE to the BaseCollider x,y,z
	rx2, ry2, rx3, ry3, d, axis int
}

// Axis constants, to help build triangle prisms.
const (
	XAxis = iota
	YAxis
	ZAxis
)

// NewTriangle - creates a new 3D triangular prism, where the prism runs along
// the specified axis. The named variables (x, y, z etc) are all regular XYZ
// When the XAxis is specified: rx2/rx3 will map to the Z axis
// When the YAxis is specified: ry2/ry3 will map to the Z axis
func NewTriangle(x, y, z, rx2, ry2, rx3, ry3, d, axis int, blocking, reactive bool, name string) Collider {
	tri := &Triangle{rx2: rx2, ry2: ry2, rx3: rx3, ry3: ry3, d: d, axis: axis}
	tri.x, tri.y, tri.z = x, y, z
	tri.name = name
	tri.ref = -1
	tri.bodyType = &BodyType{blocking: blocking, reactive: reactive}

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

	return tri
}

// ZDepth for Triangles is the z span at the given point
func (b *Triangle) ZDepth(x, y int) int {
	switch b.axis {
	case ZAxis:
		return b.d
	case XAxis:
		var pts []*types.Point
		var mainp *types.Point
		p1 := types.NewPoint(b.z, b.y)
		p2 := types.NewPoint(b.z+b.rx2, b.y+b.ry2)
		p3 := types.NewPoint(b.z+b.rx3, b.y+b.ry3)
		if p1.X > p2.X && p1.X > p3.X {
			pts = []*types.Point{p2, p3}
			mainp = p1
		} else if p2.X > p1.X && p2.X > p3.X {
			pts = []*types.Point{p1, p3}
			mainp = p2
		} else {
			pts = []*types.Point{p1, p2}
			mainp = p3
		}

		// edge case
		if y == mainp.Y {
			return mainp.X
		}

		var ret int
		for _, v := range pts {
			if (y < v.Y && y < mainp.Y) || (y > v.Y && y > mainp.Y) {
				continue
			}
			slope := float64(v.Y-mainp.Y) / float64(v.X-mainp.X)
			b := float64(v.Y) - (slope * float64(v.X))
			ret = utils.Max(int(math.Abs((float64(y)-b)/slope)), ret)
		}
		return ret

	case YAxis:
		var pts []*types.Point
		var mainp *types.Point
		p1 := types.NewPoint(b.x, b.z)
		p2 := types.NewPoint(b.x+b.rx2, b.z+b.ry2)
		p3 := types.NewPoint(b.x+b.rx3, b.z+b.ry3)
		if p1.Y > p2.Y && p1.Y > p3.Y {
			pts = []*types.Point{p2, p3}
			mainp = p1
		} else if p2.Y > p1.Y && p2.Y > p3.Y {
			pts = []*types.Point{p1, p3}
			mainp = p2
		} else {
			pts = []*types.Point{p1, p2}
			mainp = p3
		}

		// edge case
		if x == mainp.X {
			return mainp.Y
		}

		var ret int
		for _, v := range pts {
			if (x < v.X && x < mainp.X) || (x > v.X && x > mainp.X) {
				continue
			}
			slope := float64(v.Y-mainp.Y) / float64(v.X-mainp.X)
			b := float64(v.Y) - (slope * float64(v.X))
			ret = utils.Max(int(math.Abs((float64(x)-b)/slope)), ret)
		}
		return ret
	}
	return b.d
}

// YDepth for Triangles is the YSpan at a given x/z point
func (b *Triangle) YDepth(x, z int) int {
	// TODO: parse that fucking math above, jesus how did i do that
	return b.d
}

// XDepth for Triangles is the XSpan at a given y/z point
func (b *Triangle) XDepth(y, z int) int {
	// TODO: parse that fucking math above... god help me
	return b.d
}
