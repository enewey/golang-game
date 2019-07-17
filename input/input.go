package input

import (
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var keys = []ebiten.Key{
	ebiten.KeyRight, ebiten.KeyLeft, ebiten.KeyDown, ebiten.KeyUp,
	ebiten.KeySpace, ebiten.KeyTab, ebiten.KeyShift,
}

// KeyState w
type KeyState struct {
	key    ebiten.Key
	frames types.Frame
}

// Pressed w
func (k *KeyState) Pressed() bool { return ebiten.IsKeyPressed(k.key) }

// JustPressed w
func (k *KeyState) JustPressed() bool { return inpututil.IsKeyJustPressed(k.key) }

// Frames w
func (k *KeyState) Frames() types.Frame { return k.frames }

// CalcPress - accumulates the frames this key has been pressed.
func (k *KeyState) CalcPress(df types.Frame) types.Frame {
	if !k.Pressed() {
		return 0
	} else if k.JustPressed() {
		return df
	}
	return k.frames + df
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
