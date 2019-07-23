package actors

import (
	"math"

	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/utils"
	"github.com/hajimehoshi/ebiten"
)

// Manager - manages a group of actors (all actors in a scene)
type Manager struct {
	actors  map[int]*Actor // actor 0 is always the player-controller actor
	actions Actions

	playerDrawn bool
}

// NewManager create a new actor manager
func NewManager() *Manager {
	return &Manager{
		make(map[int]*Actor),
		make([]Action, 5),
		false,
	}
}

// Actors w
func (m *Manager) Actors() map[int]*Actor { return m.actors }

// Actions w
func (m *Manager) Actions() Actions { return m.actions }

// Act w
func (m *Manager) Act(df types.Frame) {
	i := 0
	for i < len(m.actions) {
		action := m.actions[i]
		if action != nil {
			if action.Process(df) {
				m.actions[i] = nil
			}
		}
		i++
	}
}

// SetPlayer - set the player-controlled actor
func (m *Manager) SetPlayer(a *Actor) {
	a.id = 0
	m.setActor(0, a)
}

// GetPlayer returns a pointer to the actor whom is controlled by the player.
func (m *Manager) GetPlayer() *Actor {
	return m.actors[0]
}

// AddActor - add a new actor to the manager
func (m *Manager) AddActor(a *Actor) {
	a.id = len(m.actors) + 1
	m.setActor(a.id, a)
}

func (m *Manager) setActor(id int, a *Actor) {
	m.actors[id] = a
}

// HandleInput - returns "true" if input is captured, disallowing any other
// 				 manager from handling the input.
func (m *Manager) HandleInput(state input.Input) bool {
	player := m.actors[0]

	if player.Controlled() {
		return false
	}
	if player.OnGround() {
		player.SetVel(0, 0, player.vz)
	}
	var dx, dy float64
	if state[ebiten.KeyUp].Pressed() {
		dy--
	}
	if state[ebiten.KeyDown].Pressed() {
		dy++
	}
	if state[ebiten.KeyLeft].Pressed() {
		dx--
	}
	if state[ebiten.KeyRight].Pressed() {
		dx++
	}
	player.vx, player.vy = dx, dy
	player.CalcDirection()

	if state[ebiten.KeySpace].JustPressed() && player.OnGround() {
		action := NewJumpAction(player, 4.0)
		m.actions.Add(action)
	}

	if state[ebiten.KeyShift].JustPressed() && !player.Dashed() && player.OnGround() {
		vx, vy := utils.Itof(DirToVec(player.direction))
		action := NewDashAction(player, vx*2.5, vy*2.5)
		m.actions.Add(action)
	}
	return true
}

// ResolveCollisions - every actor being managed will check collision against
//		the provided Colliders.
// 		Also alters velocity of actors in the air for gravity.
func (m *Manager) ResolveCollisions(scoll colliders.Colliders) {
	for _, v := range m.actors {
		// resolve the actor's direction
		v.CalcDirection()
		dx, dy, dz := v.Vel()
		dz = math.Max(dz, -6)

		hitG, hitC, hitW, ax, ay, _ :=
			scoll.PreventCollision(int(dx), int(dy), int(dz), v.Collider())

		// traversing up or down a slope
		if v.onGround { // going up
			if hitW {
				vx, vy := DirToVec(v.Direction())
				if !scoll.WouldCollide(vx-ax, vy-ay, 1, v.Collider()) &&
					scoll.WouldCollide(vx-ax, vy-ay, 0, v.Collider()) {
					v.Collider().Translate(vx-ax, vy-ay, 1)
					hitW = false
					hitG = true
				}
			} else { // going down
				if !scoll.WouldCollide(0, 0, -1, v.Collider()) &&
					scoll.WouldCollide(0, 0, -2, v.Collider()) {
					v.Collider().Translate(0, 0, -1)
					hitG = true
				}
			}
		}

		// glancing collisions - collisions where only one pixel is the
		// difference, just force the actor to the side to avoid the collision.
		if hitW && v.Orthogonal() {
			switch v.Direction() {
			case Left:
				if !scoll.WouldCollide(-1, 1, 0, v.Collider()) {
					v.Collider().Translate(-1, 1, 0)
				} else if !scoll.WouldCollide(-1, -1, 0, v.Collider()) {
					v.Collider().Translate(-1, -1, 0)
				}
				break
			case Right:
				if !scoll.WouldCollide(1, 1, 0, v.Collider()) {
					v.Collider().Translate(1, 1, 0)
				} else if !scoll.WouldCollide(1, -1, 0, v.Collider()) {
					v.Collider().Translate(1, -1, 0)
				}
				break
			case Up:
				if !scoll.WouldCollide(-1, -1, 0, v.Collider()) {
					v.Collider().Translate(-1, -1, 0)
				} else if !scoll.WouldCollide(1, -1, 0, v.Collider()) {
					v.Collider().Translate(1, -1, 0)
				}
				break
			case Down:
				if !scoll.WouldCollide(1, 1, 0, v.Collider()) {
					v.Collider().Translate(1, 1, 0)
				} else if !scoll.WouldCollide(-1, 1, 0, v.Collider()) {
					v.Collider().Translate(-1, 1, 0)
				}
				break
			}
		}

		if hitG {
			v.onGround = true
			v.vz = 0
		} else if math.Abs(dz) >= 1 {
			v.onGround = false
		}

		// if the actor is in a "dashed" state,
		// make sure it gets cleared when the actor hits the ground
		if v.dashed {
			v.dashed = !v.onGround
		}

		if hitC {
			v.vz = 0
		} else if !hitG {
			v.vz -= 0.3
		}

		v.shadowZ = scoll.FindFloor(v.Collider())
	}
}

// Render - draw the actors given a priority and row
func (m *Manager) Render(img *ebiten.Image, layer, row, ox, oy int) *ebiten.Image {
	for _, actor := range m.actors {
		sx, sy, sz := actor.Pos()
		sd := actor.Collider().ZDepth(sx, sy)
		charPr := int(math.Round(float64(sz+8) / 8))
		shadowPr := int(math.Floor(float64(actor.shadowZ+sd) / 8))
		charRow := int(math.Round(float64(sy+8) / 16))

		if shadowPr == layer && utils.Max(charRow-layer, 0) == row {
			actor.drawShadow(img, -ox, -oy)
		}
		if charPr == layer && utils.Max(charRow-layer, 0) == row {
			// fmt.Printf("drawing actor, prioritys %d %d rows %d %d\n", layer, charPr, row, charRow)
			actor.draw(img, -ox, -oy)
		}
	}
	return img
}
