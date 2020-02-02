package windows

import (
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
)

// Manager manages all active windows.
// Windows are arranged as a stack, where the last window in the windows array
// is the window with focus.
type Manager struct {
	windows []Window
}

// NewManager - create a new window manager
func NewManager() *Manager {
	return &Manager{[]Window{}}
}

// AddWindow adds a new window to the window stack
func (m *Manager) AddWindow(win Window) {
	m.windows = append(m.windows, win)
}

// HandleInput - handles the input state. Returns true if input is consumed.
func (m *Manager) HandleInput(state input.Input) bool {
	if len(m.windows) == 0 {
		return false
	}
	// iterate over windows in reverse; last window added should have focus
	for i := range m.windows {
		pos := len(m.windows) - i - 1
		v := m.windows[pos]
		if v.HandleInput(state) {
			return true
		}
	}
	return false
}

// Act w
func (m *Manager) Act(df types.Frame) bool {
	// if there are no windows, abort and return false to signal there is no focus
	if len(m.windows) == 0 {
		return false
	}
	// perform actions on each window at the top of the window stack.
	// if the window is disposed, get rid of it, and go to the next.
	for {
		win := m.windows[len(m.windows)-1]
		win.Act(df)
		if win.IsDisposed() {
			m.windows = m.windows[:len(m.windows)-1]
			if len(m.windows) == 0 {
				break
			}
		} else {
			break
		}
	}
	return true
}

// Render - draw each window in the window stack in order
// (window with focus will be drawn last, i.e. on top)
func (m *Manager) Render(img *ebiten.Image, ox, oy int) *ebiten.Image {
	for _, v := range m.windows {
		v.Draw(img, ox, oy)
	}
	return img
}
