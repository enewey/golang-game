package actors

import (
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/sprites"
)

type Actor struct {
	id       int    // unique identifier
	category string // denotes the "type" of actor
	sprite   *sprites.Sprite
	collider *colliders.Collider

	vx, vy, vz float64
	onGround   bool
}

func (a *Actor) OnGround() bool {
	return a.onGround
}

func (a *Actor) Pos() (int, int, int) {
	return a.collider.Pos()
}

func (a *Actor) Dims() (int, int, int) {
	return a.collider.Width(), a.collider.Height(), a.collider.Depth()
}

// Vel - get the actor velocity, which is how many pixels the actor will attempt
//		 to move each frame update
func (a *Actor) Vel() (float64, float64, float64) {
	return a.vx, a.vy, a.vz
}

func (a *Actor) SetVel(x, y, z float64) {
	a.vx, a.vy, a.vz = x, y, z
}

func (a *Actor) Collider() *colliders.Collider {
	return a.collider
}

func (a *Actor) Sprite() *sprites.Sprite {
	return a.sprite
}
