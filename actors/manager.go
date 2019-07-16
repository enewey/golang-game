package actors

import (
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/types"
)

type Manager struct {
	actors map[int]*Actor // actor 0 is always the player-controller actor
	actions Actions
}

func (m *Manager) Actors() map[int]*Actor { return m.actors }
func (m *Manager) Actions() Actions { return m.actions }

func (m *Manager) Act(df types.Frame) {

}

func (m *Manager) Colliders() colliders.Colliders {
	ret := make(colliders.Colliders, len(m.actors))
	var i int
	for _, v := range m.actors {
		ret[i] = v.Collider()
		i++
	}

	return ret
}
