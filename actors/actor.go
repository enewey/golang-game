package actors

import (
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/sprites"
	"github.com/hajimehoshi/ebiten"
)

// Directions for the actor 'direction' property
const (
	Up = iota
	UpRight
	Right
	DownRight
	Down
	DownLeft
	Left
	UpLeft
)

// DirToVec converts a Direction to a 2D vector
func DirToVec(d int) (int, int) {
	switch d {
	case Up:
		return 0, -1
	case UpRight:
		return 1, -1
	case Right:
		return 1, 0
	case DownRight:
		return 1, 1
	case Down:
		return 0, 1
	case DownLeft:
		return -1, 1
	case Left:
		return -1, 0
	case UpLeft:
		return -1, -1
	}
	return 0, 0
}

// Actor woo
type Actor struct {
	id        int // unique identifier
	direction int
	category  string // denotes the "type" of actor

	sprite   *sprites.Sprite
	shadow   *sprites.Sprite
	collider colliders.Collider

	vx, vy, vz float64
	shadowZ    int // shadow z-position

	onGround   bool
	controlled bool
	dashed     bool
}

// NewActor create a new actor
func NewActor(category string, sprite, shadow *sprites.Sprite,
	collider colliders.Collider) *Actor {

	return &Actor{
		-1, Down,
		category, sprite, shadow, collider,
		0, 0, 0, 0, false, false, false,
	}
}

// Direction - gets the last calculated direction for this actor
func (a *Actor) Direction() int {
	return a.direction
}

// FacingVertical returns true if the actor's direction is Up or Down
func (a *Actor) FacingVertical() bool {
	return (a.direction == Up || a.direction == Down)
}

// FacingHorizontal returns true if the actor's direction is Left or Right
func (a *Actor) FacingHorizontal() bool {
	return (a.direction == Left || a.direction == Right)
}

// Orthogonal returns true if the hero is facing Up, Down, Left or Right
func (a *Actor) Orthogonal() bool {
	return a.FacingVertical() || a.FacingHorizontal()
}

// FacingDiagonal returns true if the actor's direction is diagonal
func (a *Actor) FacingDiagonal() bool {
	return (a.direction == UpRight || a.direction == UpLeft ||
		a.direction == DownRight || a.direction == DownLeft)
}

// Dashed - get the "dashed" state -- set by the dash action.
func (a *Actor) Dashed() bool {
	return a.dashed
}

// CalcDirection - resolves the actor's direciton based on its current velocity.
func (a *Actor) CalcDirection() int {
	if a.vx < 0 && a.vy < 0 {
		a.direction = UpLeft
	} else if a.vx > 0 && a.vy < 0 {
		a.direction = UpRight
	} else if a.vx > 0 && a.vy > 0 {
		a.direction = DownRight
	} else if a.vx < 0 && a.vy > 0 {
		a.direction = DownLeft
	} else if a.vx == 0 && a.vy > 0 {
		a.direction = Down
	} else if a.vx == 0 && a.vy < 0 {
		a.direction = Up
	} else if a.vx > 0 && a.vy == 0 {
		a.direction = Right
	} else if a.vx < 0 && a.vy == 0 {
		a.direction = Left
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
	return a.sprite
}

func (a *Actor) draw(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, z := a.Pos()
	return a.sprite.DrawSprite(x-4+offsetX, y-z-8+offsetY, img)
}

func (a *Actor) drawShadow(img *ebiten.Image, offsetX, offsetY int) *ebiten.Image {
	x, y, _ := a.Pos()
	return a.shadow.DrawSprite(x-4+offsetX, y-a.shadowZ-8+offsetY, img)
}
