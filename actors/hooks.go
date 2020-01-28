package actors

import "enewey.com/golang-game/colliders"

// Hooks w
type Hooks struct {
	PostCollision []PostCollisionHook
}

// AddHook - add a hook to the Hooks structure
func (hs *Hooks) AddHook(hook Hook) {
	if h, ok := hook.(PostCollisionHook); ok {
		hs.PostCollision = append(hs.PostCollision, h)
	}
}

// PostCollisionHook - TODO: rename this
type PostCollisionHook interface {
	SetManager(*Manager)
	Tap(colliders.Colliders)
}

// Hook - The base Hook interface
type Hook interface {
	SetManager(*Manager)
}

type baseHook struct {
	manager *Manager
}

func (h *baseHook) SetManager(m *Manager) {
	h.manager = m
}

// ShadowHook - hook for making on actor act as another actor's shadow.
type ShadowHook struct {
	*baseHook
	shadow, subject Actor
}

// NewShadowHook creates a new shadow hook
func NewShadowHook(shadow, subject Actor) *ShadowHook {
	return &ShadowHook{&baseHook{}, shadow, subject}
}

// Tap - queues a change position action
func (h *ShadowHook) Tap(colls colliders.Colliders) {
	x, y, _ := h.subject.Pos()
	floor := colls.GetBlocking().FindFloor(h.subject.Collider())
	h.shadow.Collider().SetPos(x, y, floor+4)
}
