package input

import (
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
)

var keys = []ebiten.Key{
	ebiten.KeyRight, ebiten.KeyLeft, ebiten.KeyDown, ebiten.KeyUp,
	ebiten.KeySpace, ebiten.KeyTab, ebiten.KeyShift,
	ebiten.KeyZ, ebiten.KeyX,
}

// Input - map of ebiten Keys to how many frames they have been held down
type Input map[ebiten.Key]*KeyState

var state Input

// State - get the state of all related inputs (as defined in the array at the top).
func State() Input {
	if state == nil {
		state = make(map[ebiten.Key]*KeyState)
		for _, k := range keys {
			if state[k] == nil {
				state[k] = &KeyState{k, 0}
			}
		}
	}
	return state
}

// Tick - meant to be called every frame. df = delta frames (since last tick).
func (in Input) Tick(df types.Frame) Input {
	for _, v := range in {
		v.frames = v.CalcPress(df)
	}

	return in
}

// KeyState - represents the state of an individual input key.
type KeyState struct {
	key    ebiten.Key
	frames types.Frame
}

// Pressed indicates whether this button is pressed down.
func (k *KeyState) Pressed() bool { return ebiten.IsKeyPressed(k.key) }

// JustPressed indicates whether this button was pressed on the current frame.
func (k *KeyState) JustPressed() bool { return k.frames == 1 }

// Frames tells how many frames this button has been pressed.
func (k *KeyState) Frames() types.Frame { return k.frames }

// PressedUnder returns true if the button has been pressed for less than the number of frames.
func (k *KeyState) PressedUnder(f int) bool { return k.frames < f }

// PressedOver returns true if the button has been pressed for more than the number of frames.
func (k *KeyState) PressedOver(f int) bool { return k.frames > f }

// PressedWindow returns true if the button has been pressed for a number of frames greater than f but less than g
func (k *KeyState) PressedWindow(f, g int) bool { return k.frames > f && k.frames < g }

// CalcPress - accumulates the frames this key has been pressed.
func (k *KeyState) CalcPress(df types.Frame) types.Frame {
	if !k.Pressed() {
		return 0
	}
	return k.frames + df
}
