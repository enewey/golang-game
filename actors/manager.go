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
}

// NewManager create a new actor manager
func NewManager() *Manager {
	return &Manager{
		make(map[int]*Actor),
		make([]Action, 5),
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
	var dx, dy int
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
	if dx != 0 || dy != 0 {
		action := NewMoveByAction(m.actors[0], dx, dy, 0, 1)
		m.actions = append(m.actions, action)
		// ev := NewMoveByActorEvent(0, -1, dx, dy, 0, 1)
		// events.Hub().ActorEvents().Enqueue(ev)
	}

	if state[ebiten.KeySpace].Pressed() {
		action := NewJumpAction(m.actors[0], 4.0)
		m.actions = append(m.actions, action)
		// ev := NewJumpActorEvent(0, -1, 4.0)
		// events.Hub().ActorEvents().Enqueue(ev)
	}
	return false
}

// ResolveCollisions - every actor being managed will check collision against
//		the provided Colliders.
// 		Also alters velocity of actors in the air for gravity.
func (m *Manager) ResolveCollisions(scoll colliders.Colliders) {
	for _, v := range m.actors {
		dx, dy, dz := v.Vel()
		dz = math.Max(dz, -6)
		var ax, ay, az int = int(dx), int(dy), int(dz)
		var hitGround, hitCeiling, xResolved, yResolved bool
		var unresolved = true
		for unresolved {
			ax, ay, az, hitGround, hitCeiling, xResolved, yResolved =
				colliders.ResolveCollision(ax, ay, az, v.Collider(), scoll)
			unresolved = xResolved || yResolved
			if xResolved {
				// fmt.Printf("resolved x %d - ", ax)
				v.Collider().Translate(ax, 0, 0)
				v.vx = 0
				ax = 0
			} else if yResolved {
				// fmt.Printf("resolved y %d - ", ay)
				v.Collider().Translate(0, ay, 0)
				v.vy = 0
				ay = 0
			} else {
				// fmt.Printf("resolved all three %d %d %d\n", ax, ay, az)
				v.Collider().Translate(ax, ay, az)
			}
		}
		// fmt.Printf("dz %d az %d\n fallV %f\n", dz, az, fallV)
		if hitGround {
			v.onGround = true
			v.vz = 0
		} else if az != 0 {
			v.onGround = false
		}
		if v.onGround && dz < -1 && az >= 0 {
			v.vz = 0
		} else if hitCeiling && dz > 0 && az <= 0 {
			v.vz = 0
		} else {
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
