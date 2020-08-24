package actors

import (
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
)

// DirToVec converts a Direction to a 2D vector
func DirToVec(d types.Direction) (int, int) {
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

// VecToDir - transform a 2D vector to a Direction
func VecToDir(vx, vy float64, def types.Direction) types.Direction {
	if vx < 0 && vy < 0 {
		return types.UpLeft
	} else if vx > 0 && vy < 0 {
		return types.UpRight
	} else if vx > 0 && vy > 0 {
		return types.DownRight
	} else if vx < 0 && vy > 0 {
		return types.DownLeft
	} else if vx == 0 && vy > 0 {
		return types.Down
	} else if vx == 0 && vy < 0 {
		return types.Up
	} else if vx > 0 && vy == 0 {
		return types.Right
	} else if vx < 0 && vy == 0 {
		return types.Left
	}
	return def
}

// CanMove is an interface for entities which can be moved and/or controlled.
type CanMove interface {
	OnGround() bool
	SetOnGround(bool)
	SubPos() (float64, float64, float64)
	SetSubPos(float64, float64, float64)
	MoveDelta(float64, float64, float64) (float64, float64, float64)
	Vel() (float64, float64, float64)
	SetVel(x, y, z float64)
	SetVelX(x float64)
	SetVelY(y float64)
	SetVelZ(z float64)
	Collider() colliders.Collider
	Weight() int

	Direction() types.Direction
	FacingVertical() bool
	FacingHorizontal() bool
	Orthogonal() bool
	FacingDiagonal() bool
	CalcDirection() types.Direction
}

var _ CanMove = &MovingActor{}
var _ CanMove = &CharActor{}

// Controllable w
type Controllable interface {
	Controlled() bool
	SetControlled(bool)
}

var _ Controllable = &CharActor{}

// CanDash w
type CanDash interface {
	Dashed() bool
	SetDashed(bool)
}

var _ CanDash = &CharActor{}

// Drawable is an interface for entities which can be drawn on the screen
type Drawable interface {
	DrawOffset() (int, int)
	DrawPos() (int, int)
	Sprite() *sprites.Sprite
	draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image
}

var _ Drawable = &CharActor{}
var _ Drawable = &SpriteActor{}

// Actor interface
type Actor interface {
	ID() int
	SetID(int)
	Pos() (int, int, int)
	SetPos(int, int, int)
	Collider() colliders.Collider
	CanCollide() bool
	Category() string
	IsBehind(Actor) bool
}

var _ Actor = &baseActor{}
var _ Actor = &SpriteActor{}
var _ Actor = &StaticActor{}
var _ Actor = &MovingActor{}
var _ Actor = &CharActor{}

type baseActor struct {
	id       int
	category string // denotes the "type" of actor

	collider colliders.Collider
}

// ID - unique id for actor
func (a *baseActor) ID() int { return a.id }

// SetID - set the unique ID
func (a *baseActor) SetID(id int) { a.id = id }

// Pos - returns an x,y,z tuple of the actor position
func (a *baseActor) Pos() (int, int, int) { return a.collider.Pos() }

// SetPos - sets the position of the actor
func (a *baseActor) SetPos(x, y, z int) {
	a.collider.SetPos(x, y, z)
}

// Collider - returns the raw collider for the actor
func (a *baseActor) Collider() colliders.Collider { return a.collider }

// CanCollide tells whether this actor can resolve collisions.
func (a *baseActor) CanCollide() bool { return false }

// Category - returns the designated category metadata of the actor
func (a *baseActor) Category() string { return a.category }

func (a *baseActor) IsBehind(b Actor) bool { return true }

// InvisibleActor is an actor with no sprite, but does have collision.
type InvisibleActor struct {
	baseActor
}

// NewInvisibleActor returns a new invisible actor with the provided collider.
func NewInvisibleActor(category string, collider colliders.Collider) *InvisibleActor {
	return &InvisibleActor{baseActor{-1, category, collider}}
}

// CanCollide tells whether this actor can resolve collisions.
func (a *InvisibleActor) CanCollide() bool { return true }

// SpriteActor is an actor that has a sprite.
type SpriteActor struct {
	baseActor
	// drawn offset
	spritemap sprites.Spritemap
	ox, oy    int
}

// NewSpriteActor returns a new SpriteActor, an actor with collision and a sprite.
func NewSpriteActor(
	category string,
	sprite sprites.Spritemap,
	collider colliders.Collider,
	ox, oy int,
) *SpriteActor {
	return &SpriteActor{
		baseActor{-1, category, collider},
		sprite,
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

// IsBehind tests whether this actor is "behind" another actor, which is used to determine drawing order.
func (a *SpriteActor) IsBehind(b Actor) bool {
	defaultToID := func(args ...int) bool {
		for _, v := range args {
			if v != 0 {
				return v < 0
			}
		}
		return a.ID() < b.ID()
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
		return defaultToID((az+ad)-(bz+bd), (ay+alen)-(by+blen), (az+ad+ay+alen)-(bz+bd+by+blen))
	} else if zIntersects && !yIntersects {
		return defaultToID((ay+alen)-(by+blen), (az+ad)-(bz+bd), (az+ad+ay+alen)-(bz+bd+by+blen))
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

// MovingActor is like a static actor, but can move.
// vx/vy/vz is the actor's velocity
// subx/suby/subz is the actor's sub-pixel position
type MovingActor struct {
	StaticActor
	vx, vy, vz       float64
	subx, suby, subz float64
	direction        types.Direction
	weight           int
	onGround         bool
}

// NewMovingActor creates a new MovingActor, which is like a static actor that can move.
func NewMovingActor(
	category string,
	sprite sprites.Spritemap,
	collider colliders.Collider,
	ox, oy int,
	weight int,
	onGround bool,
) Actor {
	return &MovingActor{
		*NewStaticActor(category, sprite, collider, ox, oy),
		0, 0, 0,
		0, 0, 0,
		types.Down,
		weight,
		true,
	}
}

// Vel - get the actor velocity, which is how many pixels the actor will attempt
//		 to move each frame update
func (a *MovingActor) Vel() (float64, float64, float64) {
	return a.vx, a.vy, a.vz
}

// MoveDelta - gets a calculated movement delta using the provided velocity and
// the actor's sub-pixel position.
func (a *MovingActor) MoveDelta(vx, vy, vz float64) (float64, float64, float64) {
	sx, sy, sz := a.SubPos()
	return sx + vx, sy + vy, sz + vz
}

// SetVel woo
func (a *MovingActor) SetVel(x, y, z float64) {
	a.vx, a.vy, a.vz = x, y, z
}

// SetVelX w
func (a *MovingActor) SetVelX(x float64) { a.vx = x }

// SetVelY y
func (a *MovingActor) SetVelY(y float64) { a.vy = y }

// SetVelZ z
func (a *MovingActor) SetVelZ(z float64) { a.vz = z }

// SubPos - gets the sub-pixel offset of the actor's position
func (a *MovingActor) SubPos() (float64, float64, float64) {
	return a.subx, a.suby, a.subz
}

// SetSubPos -- set the actor's subpixel position
func (a *MovingActor) SetSubPos(cx, cy, cz float64) {
	a.subx, a.suby, a.subz = cx, cy, cz
}

// Direction - gets the last calculated direction for this actor
func (a *MovingActor) Direction() types.Direction { return a.direction }

// FacingVertical returns true if the actor's direction is Up or Down
func (a *MovingActor) FacingVertical() bool {
	return (a.direction == types.Up || a.direction == types.Down)
}

// FacingHorizontal returns true if the actor's direction is Left or Right
func (a *MovingActor) FacingHorizontal() bool {
	return (a.direction == types.Left || a.direction == types.Right)
}

// Orthogonal returns true if the hero is facing Up, Down, Left or Right
func (a *MovingActor) Orthogonal() bool {
	return a.FacingVertical() || a.FacingHorizontal()
}

// FacingDiagonal returns true if the actor's direction is diagonal
func (a *MovingActor) FacingDiagonal() bool {
	return (a.direction == types.UpRight || a.direction == types.UpLeft ||
		a.direction == types.DownRight || a.direction == types.DownLeft)
}

// CalcDirection - resolves the actor's direciton based on its current velocity.
func (a *MovingActor) CalcDirection() types.Direction {
	a.direction = VecToDir(a.vx, a.vy, a.direction)
	return a.direction
}

// OnGround woo
func (a *MovingActor) OnGround() bool { return a.onGround }

// SetOnGround woo
func (a *MovingActor) SetOnGround(b bool) { a.onGround = b }

// Weight is the priority of this actor in terms of blocking; heavier actors will push around lighter actors.
func (a *MovingActor) Weight() int { return a.weight }

// CharActor woo
type CharActor struct {
	MovingActor

	controlled bool
	dashed     bool
}

// NewCharActor create a new char actor
func NewCharActor(
	category string,
	sprite sprites.Spritemap,
	collider colliders.Collider,
	ox, oy int,
	weight int,
) Actor {
	return &CharActor{
		*NewMovingActor(category, sprite, collider, ox, oy, weight, true).(*MovingActor),
		false, false,
	}
}

// Dashed - get the "dashed" state -- set by the dash action.
func (a *CharActor) Dashed() bool { return a.dashed }

// SetDashed woow
func (a *CharActor) SetDashed(b bool) { a.dashed = b }

// Controlled - this actor is being controlled by actions and cannot respond to input
func (a *CharActor) Controlled() bool { return a.controlled }

// SetControlled - this actor is being controlled by actions and cannot respond to input
func (a *CharActor) SetControlled(b bool) { a.controlled = b }

// Sprite woo
func (a *CharActor) Sprite() *sprites.Sprite {
	return a.spritemap.Sprite(int(a.direction))
}

// DrawPos - returns the position this actor should be drawn in world space
func (a *CharActor) DrawPos() (int, int) {
	x, y, z := a.Pos()
	return (x + a.ox), (y - z + a.oy)
}

func (a *CharActor) draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y := a.DrawPos()
	return a.spritemap.Sprite(int(a.direction)).Draw(x+offsetX, y+offsetY, img)
}

// DrawOffset s
func (a *CharActor) DrawOffset() (int, int) { return a.ox, a.oy }

// CanCollide - for char actors, yes.
func (a *CharActor) CanCollide() bool { return true }
