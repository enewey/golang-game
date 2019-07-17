package actors

import (
	"math"

	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/types"
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
	// action := NewMoveByAction(m.actors[0], dx, dy, 0, 1)
	// m.actions.Add(action)
	// ev := NewMoveByActorEvent(0, -1, dx, dy, 0, 1)
	// events.Hub().ActorEvents().Enqueue(ev)

	if state[ebiten.KeySpace].JustPressed() && player.OnGround() {
		action := NewJumpAction(player, 4.0)
		m.actions.Add(action)
		// ev := NewJumpActorEvent(0, -1, 4.0)
		// events.Hub().ActorEvents().Enqueue(ev)
	}

	if state[ebiten.KeyShift].JustPressed() && !player.Dashed() && player.OnGround() {
		var vx, vy float64
		switch player.direction {
		case Up:
			vy = -1.0
			break
		case Down:
			vy = 1.0
			break
		case Right:
			vx = 1.0
			break
		case Left:
			vx = -1.0
			break
		case UpRight:
			vx, vy = 1.0, -1.0
			break
		case UpLeft:
			vx, vy = -1.0, -1.0
			break
		case DownRight:
			vx, vy = 1.0, 1.0
			break
		case DownLeft:
			vx, vy = -1.0, 1.0
			break
		}
		action := NewDashAction(player, vx*2.5, vy*2.5)
		m.actions.Add(action)
		// ev := NewJumpActorEvent(0, -1, 4.0)
		// events.Hub().ActorEvents().Enqueue(ev)
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

		hitG, hitC, _ :=
			scoll.PreventCollision(int(dx), int(dy), int(dz), v.Collider())

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
func (m *Manager) Render(img *ebiten.Image, layer, row int) *ebiten.Image {
	for _, actor := range m.actors {
		_, sy, sz := actor.Pos()
		sd := actor.Collider().Depth()
		charPr := int(math.Ceil(float64(sz+1) / 8))
		shadowPr := int(math.Floor(float64(actor.shadowZ+sd) / 8))
		charRow := int(math.Ceil(float64(sy) / 16))

		if shadowPr == layer && max(charRow-layer, 0) == row {
			actor.drawShadow(img)
		}
		if charPr == layer && max(charRow-layer, 0) == row {
			actor.draw(img)
		}
	}
	return img
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
