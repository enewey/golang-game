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

// Actor woo
type Actor struct {
	id        int // unique identifier
	direction int
	category  string // denotes the "type" of actor

	spritemap *sprites.Spritemap
	shadow    *sprites.Sprite
	collider  colliders.Collider

	vx, vy, vz float64
	shadowZ    int // shadow z-position

	onGround   bool
	controlled bool
	dashed     bool
}

// NewActor create a new actor
func NewActor(category string, sprite *sprites.Spritemap, shadow *sprites.Sprite,
	collider colliders.Collider) *Actor {

	return &Actor{
		-1, types.Down,
		category,
		sprite,
		shadow,
		collider,
		0, 0, 0, 0, false, false, false,
	}
}

// Direction - gets the last calculated direction for this actor
func (a *Actor) Direction() int {
	return a.direction
}

// FacingVertical returns true if the actor's direction is Up or Down
func (a *Actor) FacingVertical() bool {
	return (a.direction == types.Up || a.direction == types.Down)
}

// FacingHorizontal returns true if the actor's direction is Left or Right
func (a *Actor) FacingHorizontal() bool {
	return (a.direction == types.Left || a.direction == types.Right)
}

// Orthogonal returns true if the hero is facing Up, Down, Left or Right
func (a *Actor) Orthogonal() bool {
	return a.FacingVertical() || a.FacingHorizontal()
}

// FacingDiagonal returns true if the actor's direction is diagonal
func (a *Actor) FacingDiagonal() bool {
	return (a.direction == types.UpRight || a.direction == types.UpLeft ||
		a.direction == types.DownRight || a.direction == types.DownLeft)
}

// Dashed - get the "dashed" state -- set by the dash action.
func (a *Actor) Dashed() bool {
	return a.dashed
}

// CalcDirection - resolves the actor's direciton based on its current velocity.
func (a *Actor) CalcDirection() int {
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
func (a *Actor) OnGround() bool {
	return a.onGround
}

// Controlled - this actor is being controlled by actions and cannot respond to
// input
func (a *Actor) Controlled() bool {
	return a.controlled
}

// Pos woo
func (a *Actor) Pos() (int, int, int) {
	return a.collider.Pos()
}

// Bottom returns the "bottom" of the actor's graphic position.
// TODO: fuck
func (a *Actor) Bottom() int {
	return a.collider.Y() + 8
}

// Vel - get the actor velocity, which is how many pixels the actor will attempt
//		 to move each frame update
func (a *Actor) Vel() (float64, float64, float64) {
	return a.vx, a.vy, a.vz
}

// SetVel woo
func (a *Actor) SetVel(x, y, z float64) {
	a.vx, a.vy, a.vz = x, y, z
}

// SetVelX w
func (a *Actor) SetVelX(x float64) { a.vx = x }

// SetVelY y
func (a *Actor) SetVelY(y float64) { a.vy = y }

// SetVelZ z
func (a *Actor) SetVelZ(z float64) { a.vz = z }

// Collider woo
func (a *Actor) Collider() colliders.Collider {
	return a.collider
}

// Sprite woo
func (a *Actor) Sprite() *sprites.Sprite {
	return (*a.spritemap)[a.direction]
}

func (a *Actor) draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, z := a.Pos()
	return a.spritemap.DrawSprite(a.direction, x-4+offsetX, y-z-8+offsetY, img)
}

func (a *Actor) drawShadow(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, _ := a.Pos()
	return a.shadow.DrawSprite(x-4+offsetX, y-a.shadowZ-8+offsetY, img)
}
