package actors

import (
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/sprites"
	"github.com/hajimehoshi/ebiten"
)

// Actor woo
type Actor struct {
	id       int    // unique identifier
	category string // denotes the "type" of actor
	sprite   *sprites.Sprite
	shadow   *sprites.Sprite
	collider *colliders.Collider

	vx, vy, vz float64
	shadowZ    int // shadow z-position
	onGround   bool
	controlled bool
}

// NewActor create a new actor
func NewActor(category string, sprite, shadow *sprites.Sprite,
	collider *colliders.Collider) *Actor {

	return &Actor{
		-1,
		category, sprite, shadow, collider,
		0, 0, 0, 0, false, false,
	}
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

// Dims woo
func (a *Actor) Dims() (int, int, int) {
	return a.collider.Width(), a.collider.Height(), a.collider.Depth()
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
func (a *Actor) Collider() *colliders.Collider {
	return a.collider
}

// Sprite woo
func (a *Actor) Sprite() *sprites.Sprite {
	return a.sprite
}

func (a *Actor) draw(img *ebiten.Image) *ebiten.Image {
	x, y, z := a.Pos()
	return a.sprite.DrawSprite(x-4, y-z-8, img)
}

func (a *Actor) drawShadow(img *ebiten.Image) *ebiten.Image {
	x, y, _ := a.Pos()
	return a.shadow.DrawSprite(x-4, y-a.shadowZ-8, img)
}
