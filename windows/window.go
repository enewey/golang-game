package windows

// Window represents a box to contain text or images that appears in a separate
// context from actors.
type Window interface {
	Draw()
	HandleInput()
}
