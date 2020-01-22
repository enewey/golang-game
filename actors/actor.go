package actors

import (
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
)

// DirToVec converts a Direction to a 2D vector
func DirToVec(d int) (int, int) {
	switch d {
	case types.Up:
		return 0, -1
	case types.UpRight:
		return 1, -1
	case types.Right:
		return 1, 0
	case types.DownRight:
		return 1, 1
	case types.Down:
		return 0, 1
	case types.DownLeft:
		return -1, 1
	case types.Left:
		return -1, 0
	case types.UpLeft:
		return -1, -1
	}
	return 0, 0
}

// CanMove is an interface for entities which can be moved and/or controlled.
type CanMove interface {
	Direction() int
	FacingVertical() bool
	FacingHorizontal() bool
	Orthogonal() bool
	FacingDiagonal() bool
	Dashed() bool
	SetDashed(bool)
	CalcDirection() int
	OnGround() bool
	SetOnGround(bool)
	Controlled() bool
	SetControlled(bool)
	Vel() (float64, float64, float64)
	SetVel(x, y, z float64)
	SetVelX(x float64)
	SetVelY(y float64)
	SetVelZ(z float64)
	Collider() colliders.Collider
}

// Drawable is an interface for entities which can be drawn on the screen
type Drawable interface {
	DrawOffset() (int, int)
	DrawPos() (int, int)
	Sprite() *sprites.Sprite
	draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image
}

// Actor interface
type Actor interface {
	ID() int
	SetID(int)
	Pos() (int, int, int)
	Collider() colliders.Collider
	CanCollide() bool
	Category() string
	IsBehind(Actor) bool
}

type baseActor struct {
	id       int
	category string // denotes the "type" of actor

	spritemap sprites.Spritemap
	collider  colliders.Collider
}

// ID - unique id for actor
func (a *baseActor) ID() int { return a.id }

// SetID - set the unique ID
func (a *baseActor) SetID(id int) { a.id = id }

// Pos - returns an x,y,z tuple of the actor position
func (a *baseActor) Pos() (int, int, int) { return a.collider.Pos() }

// Collider - returns the raw collider for the actor
func (a *baseActor) Collider() colliders.Collider { return a.collider }

func (a *baseActor) CanCollide() bool { return false }

// Category - returns the designated category metadata of the actor
func (a *baseActor) Category() string { return a.category }

// SpriteActor s
type SpriteActor struct {
	baseActor
	// drawn offset
	ox, oy int
}

// NewSpriteActor s
func NewSpriteActor(
	category string,
	sprite sprites.Spritemap,
	collider colliders.Collider,
	ox, oy int,
) *SpriteActor {

	return &SpriteActor{
		baseActor{-1, category, sprite, collider},
		ox, oy,
	}
}

// Sprite woo
func (a *SpriteActor) Sprite() *sprites.Sprite {
	return a.spritemap.Sprite(0)
}

// DrawPos - returns the position this actor should be drawn in world space
func (a *SpriteActor) DrawPos() (int, int) {
	x, y, z := a.Pos()
	return (x + a.ox), (y - z + a.oy)
}

func (a *SpriteActor) draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, z := a.Pos()
	return a.spritemap.Sprite(0).Draw(x+a.ox+offsetX, y-z+a.oy+offsetY, img)
}

// DrawOffset s
func (a *SpriteActor) DrawOffset() (int, int) { return a.ox, a.oy }

// IsBehind s
func (a *SpriteActor) IsBehind(b Actor) bool {
	defaultToID := func(compVal int) bool {
		if compVal == 0 {
			return a.ID() < b.ID()
		}
		return compVal < 0
	}

	ax, ay, az := a.Pos()
	bx, by, bz := b.Pos()
	ad := a.Collider().ZDepth(ax, ay)
	alen := a.Collider().YDepth(ax, az)
	bd := b.Collider().ZDepth(bx, by)
	blen := b.Collider().YDepth(bx, bz)

	yIntersects := (ay < by+blen) && (by < ay+alen)
	zIntersects := (az < bz+bd) && (bz < az+ad)

	// if the two actors intersect on the Z or Y plane,
	// sort by the higher Y or Z position
	// (i.e. sort by the non-intersecting plane)
	if yIntersects && !zIntersects {
		return defaultToID((az + ad) - (bz + bd))
	} else if zIntersects && !yIntersects {
		return defaultToID((ay + alen) - (by + blen))
	}
	return defaultToID((az + ad + ay + alen) - (bz + bd + by + blen))
}

// StaticActor - actor that has collision.
type StaticActor struct {
	SpriteActor
}

// NewStaticActor s
func NewStaticActor(
	category string,
	sprite sprites.Spritemap,
	collider colliders.Collider,
	ox, oy int,
) *StaticActor {

	return &StaticActor{
		*NewSpriteActor(category, sprite, collider, ox, oy),
	}
}

// CanCollide - for static actors, yes.
func (a *StaticActor) CanCollide() bool { return true }

// CharActor woo
type CharActor struct {
	SpriteActor

	direction int
	shadow    *sprites.Sprite

	vx, vy, vz float64
	shadowZ    int // shadow z-position

	onGround   bool
	controlled bool
	dashed     bool
}

// NewCharActor create a new char actor
func NewCharActor(
	category string,
	sprite sprites.Spritemap,
	shadow *sprites.Sprite,
	collider colliders.Collider,
	ox, oy int,
) Actor {

	return &CharActor{
		*NewSpriteActor(category, sprite, collider, ox, oy),
		types.Down,
		shadow,
		0, 0, 0, 0, false, false, false,
	}
}

// Direction - gets the last calculated direction for this actor
func (a *CharActor) Direction() int { return a.direction }

// FacingVertical returns true if the actor's direction is Up or Down
func (a *CharActor) FacingVertical() bool {
	return (a.direction == types.Up || a.direction == types.Down)
}

// FacingHorizontal returns true if the actor's direction is Left or Right
func (a *CharActor) FacingHorizontal() bool {
	return (a.direction == types.Left || a.direction == types.Right)
}

// Orthogonal returns true if the hero is facing Up, Down, Left or Right
func (a *CharActor) Orthogonal() bool {
	return a.FacingVertical() || a.FacingHorizontal()
}

// FacingDiagonal returns true if the actor's direction is diagonal
func (a *CharActor) FacingDiagonal() bool {
	return (a.direction == types.UpRight || a.direction == types.UpLeft ||
		a.direction == types.DownRight || a.direction == types.DownLeft)
}

// Dashed - get the "dashed" state -- set by the dash action.
func (a *CharActor) Dashed() bool { return a.dashed }

// SetDashed woow
func (a *CharActor) SetDashed(b bool) { a.dashed = b }

// CalcDirection - resolves the actor's direciton based on its current velocity.
func (a *CharActor) CalcDirection() int {
	if a.vx < 0 && a.vy < 0 {
		a.direction = types.UpLeft
	} else if a.vx > 0 && a.vy < 0 {
		a.direction = types.UpRight
	} else if a.vx > 0 && a.vy > 0 {
		a.direction = types.DownRight
	} else if a.vx < 0 && a.vy > 0 {
		a.direction = types.DownLeft
	} else if a.vx == 0 && a.vy > 0 {
		a.direction = types.Down
	} else if a.vx == 0 && a.vy < 0 {
		a.direction = types.Up
	} else if a.vx > 0 && a.vy == 0 {
		a.direction = types.Right
	} else if a.vx < 0 && a.vy == 0 {
		a.direction = types.Left
	}
	return a.direction
}

// OnGround woo
func (a *CharActor) OnGround() bool { return a.onGround }

// SetOnGround woo
func (a *CharActor) SetOnGround(b bool) { a.onGround = b }

// IsStatic - a static actor does not move, and does not need collision checks
func (a *CharActor) IsStatic() bool {
	return a.category == "static"
}

// Controlled - this actor is being controlled by actions and cannot respond to input
func (a *CharActor) Controlled() bool { return a.controlled }

// SetControlled - this actor is being controlled by actions and cannot respond to input
func (a *CharActor) SetControlled(b bool) { a.controlled = b }

// Vel - get the actor velocity, which is how many pixels the actor will attempt
//		 to move each frame update
func (a *CharActor) Vel() (float64, float64, float64) {
	return a.vx, a.vy, a.vz
}

// SetVel woo
func (a *CharActor) SetVel(x, y, z float64) {
	a.vx, a.vy, a.vz = x, y, z
}

// SetVelX w
func (a *CharActor) SetVelX(x float64) { a.vx = x }

// SetVelY y
func (a *CharActor) SetVelY(y float64) { a.vy = y }

// SetVelZ z
func (a *CharActor) SetVelZ(z float64) { a.vz = z }

// Sprite woo
func (a *CharActor) Sprite() *sprites.Sprite {
	return a.spritemap.Sprite(a.direction)
}

// DrawPos - returns the position this actor should be drawn in world space
func (a *CharActor) DrawPos() (int, int) {
	x, y, z := a.Pos()
	return (x + a.ox), (y - z + a.oy)
}

func (a *CharActor) draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, z := a.Pos()
	return a.spritemap.Sprite(a.direction).Draw(x+a.ox+offsetX, y-z+a.oy+offsetY, img)
}

func (a *CharActor) drawShadow(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, _ := a.Pos()
	return a.shadow.Draw(x-4+offsetX, y-a.shadowZ-8+offsetY, img)
}

// DrawOffset s
func (a *CharActor) DrawOffset() (int, int) { return a.ox, a.oy }

// CanCollide - for char actors, yes.
func (a *CharActor) CanCollide() bool { return true }
