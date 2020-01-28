package actors

import (
	"math"
	"sort"

	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/utils"
	"github.com/hajimehoshi/ebiten"
)

// Manager - manages a group of actors (all actors in a scene)
type Manager struct {
	actors       map[int]Actor // actor 0 is always the player-controller actor
	sortedActors []Actor
	actions      Actions
	hooks        *Hooks
}

// NewManager create a new actor manager
func NewManager() *Manager {
	return &Manager{
		make(map[int]Actor),
		[]Actor{},
		make([]Action, 5),
		&Hooks{[]PostCollisionHook{}},
	}
}

// Actors w
func (m *Manager) Actors() map[int]Actor { return m.actors }

// Actions w
func (m *Manager) Actions() Actions { return m.actions }

// Act - process all queued actions
func (m *Manager) Act(df types.Frame) {
	i := 0
	for i < len(m.actions) {
		action := m.actions[i]
		if action != nil {
			// when an action returns true, it is done processing
			if action.Process(df) {
				m.actions[i] = nil
			}
		}
		i++
	}
}

// AddHook - add a hook to be processed by the manager
func (m *Manager) AddHook(hook Hook) {
	hook.SetManager(m)
	m.hooks.AddHook(hook)
}

// SetPlayer - set the player-controlled actor
func (m *Manager) SetPlayer(a Actor) {
	a.SetID(0)
	m.setActor(0, a)
}

// GetPlayer returns a pointer to the actor whom is controlled by the player.
func (m *Manager) GetPlayer() Actor {
	return m.actors[0]
}

// AddActor - add a new actor to the manager
func (m *Manager) AddActor(a Actor) {
	a.SetID(len(m.actors) + 1)
	m.setActor(a.ID(), a)
}

func (m *Manager) setActor(id int, a Actor) {
	m.actors[id] = a
	a.Collider().SetRef(id)
	m.sortedActors = append(m.sortedActors, a)
}

// HandleInput - returns "true" if input is captured, disallowing any other
// 				 manager from handling the input.
func (m *Manager) HandleInput(state input.Input) bool {
	playerActor := m.actors[0]
	player := playerActor.(CanMove)

	if player.Controlled() {
		return false
	}
	if player.OnGround() {
		_, _, vz := player.Vel()
		player.SetVel(0, 0, vz)
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
	player.SetVelX(dx)
	player.SetVelY(dy)
	player.CalcDirection()

	if state[ebiten.KeySpace].JustPressed() && player.OnGround() {
		action := NewJumpAction(playerActor, 4.0)
		m.actions.Add(action)
	}

	if state[ebiten.KeyShift].JustPressed() && !player.Dashed() && player.OnGround() {
		vx, vy := utils.Normalize2(utils.Itof(DirToVec(player.Direction())))
		action := NewDashAction(playerActor, vx*2.5, vy*2.5, 0.0)
		m.actions.Add(action)
	}
	return true
}

// ResolveCollisions - every CanMove actor being managed will check collision against
//		the provided Colliders.
// 		Also alters velocity of actors in the air for gravity.
func (m *Manager) ResolveCollisions(scoll colliders.Colliders) {
	var mcolls colliders.Colliders = scoll[:]
	for _, ac := range m.actors {
		if ac.CanCollide() {
			mcolls = append(mcolls, ac.Collider())
		}
	}

	for _, ac := range m.actors {
		if _, ok := ac.(CanMove); !ok {
			continue
		}
		subject := ac.(CanMove)

		// Exclude the subject actor
		colliderCtx := mcolls.Remove(subject.Collider())

		dx, dy, dz := subject.Vel()

		// First, run subject against colliders with custom behavior (reactive colliders)
		reactors := colliderCtx.GetReactive()
		for _, r := range reactors.GetColliding(int(dx), int(dy), int(dz), subject.Collider()) {
			r.Reaction()(ac, m.actors[r.Ref()])
		}

		dx, dy, dz = subject.Vel()

		// Second, check collision against blocking colliders and prevent the collisions.
		handleBlockingCollisions(dx, dy, dz, subject, colliderCtx.GetBlocking())

		for _, hook := range m.hooks.PostCollision {
			hook.Tap(colliderCtx)
		}
	}
}

func handleBlockingCollisions(dx, dy, dz float64, v CanMove, colliderCtx colliders.Colliders) {
	// resolve the actor's direction
	v.CalcDirection()
	dz = math.Max(dz, -6)

	hitG, hitC, hitW, ax, ay, _ :=
		colliderCtx.PreventCollision(int(dx), int(dy), int(dz), v.Collider())

	// traversing up or down a slope
	if v.OnGround() { // going up
		if hitW {
			vx, vy := DirToVec(v.Direction())
			if !colliderCtx.WouldCollide(vx-ax, vy-ay, 1, v.Collider()) &&
				colliderCtx.WouldCollide(vx-ax, vy-ay, 0, v.Collider()) {
				v.Collider().Translate(vx-ax, vy-ay, 1)
				hitW = false
				hitG = true
			}
		} else { // going down
			if !colliderCtx.WouldCollide(0, 0, -1, v.Collider()) &&
				colliderCtx.WouldCollide(0, 0, -2, v.Collider()) {
				v.Collider().Translate(0, 0, -1)
				hitG = true
			}
		}
	}

	// glancing collisions - collisions where only one pixel is the
	// difference, just force the actor to the side to avoid the collision.
	if hitW && v.Orthogonal() {
		switch v.Direction() {
		case types.Left:
			if !colliderCtx.WouldCollide(-1, 1, 0, v.Collider()) {
				v.Collider().Translate(-1, 1, 0)
			} else if !colliderCtx.WouldCollide(-1, -1, 0, v.Collider()) {
				v.Collider().Translate(-1, -1, 0)
			}
			break
		case types.Right:
			if !colliderCtx.WouldCollide(1, 1, 0, v.Collider()) {
				v.Collider().Translate(1, 1, 0)
			} else if !colliderCtx.WouldCollide(1, -1, 0, v.Collider()) {
				v.Collider().Translate(1, -1, 0)
			}
			break
		case types.Up:
			if !colliderCtx.WouldCollide(-1, -1, 0, v.Collider()) {
				v.Collider().Translate(-1, -1, 0)
			} else if !colliderCtx.WouldCollide(1, -1, 0, v.Collider()) {
				v.Collider().Translate(1, -1, 0)
			}
			break
		case types.Down:
			if !colliderCtx.WouldCollide(1, 1, 0, v.Collider()) {
				v.Collider().Translate(1, 1, 0)
			} else if !colliderCtx.WouldCollide(-1, 1, 0, v.Collider()) {
				v.Collider().Translate(-1, 1, 0)
			}
			break
		}
	}

	if !colliderCtx.WouldCollide(0, 0, -1, v.Collider()) {
		v.SetOnGround(false)
	}

	if hitG {
		v.SetOnGround(true)
		v.SetVelZ(0)
	} else if math.Abs(dz) >= 1 {
		v.SetOnGround(false)
	}

	// if the actor is in a "dashed" state,
	// make sure it gets cleared when the actor hits the ground
	if v.Dashed() {
		v.SetDashed(!v.OnGround())
	}

	if hitC {
		v.SetVelZ(0)
	} else if !hitG && !v.OnGround() {
		_, _, vz := v.Vel()
		v.SetVelZ(math.Max(vz+config.Get().Gravity(), -6.0))
	}

	// v.shadowZ = scoll.FindFloor(v.Collider())
}

// Render - draw the actors given a priority and row
func (m *Manager) Render(img *ebiten.Image, ox, oy int) *ebiten.Image {
	m.drawSort()
	for _, actor := range m.sortedActors {
		drawable := actor.(Drawable)
		drawable.draw(img, -ox, -oy)
	}
	return img
}

func (m *Manager) drawSort() {
	sort.Slice(m.sortedActors, func(i, j int) bool {
		return m.sortedActors[i].IsBehind(m.sortedActors[j])
	})
}
