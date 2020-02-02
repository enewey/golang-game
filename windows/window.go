package windows

import (
	"fmt"
	"image/color"

	"enewey.com/golang-game/cache"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/types"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

// Window represents a box to contain text or images that appears in a separate
// context from actors. Unlike Actors, windows do not rely on Actions to progress
// their gamestate. They receive input and are stateful.
type Window interface {
	Draw(*ebiten.Image)
	HandleInput(input.Input) bool
	Act(types.Frame)
	IsDisposed() bool
}

// BaseWindow is a
type BaseWindow struct {
	skin       *sprites.SpriteRect
	x, y, w, h int
	disposed   bool
}

// NewBlankWindow returns a blank window.
func NewBlankWindow(x, y, w, h int, c color.Color) *BaseWindow {
	sprite := sprites.NewSpriteRect(w, h)
	sprite.Sprite.Img().Fill(c)
	return &BaseWindow{sprite, x, y, w, h, false}
}

// IsDisposed tells whether this window should be disposed or not.
func (w *BaseWindow) IsDisposed() bool { return w.disposed }

func (w *BaseWindow) dispose() { w.disposed = true }

// Draw does a draw
// func (w *BaseWindow) Draw(df types.Frame, img *ebiten.Image) {
// 	w.skin.Sprite.Draw(w.x, w.y, img)
// }

// MessageWindow is a window containing a text message
type MessageWindow struct {
	*BaseWindow
	message     string
	currMessage string
	speed       types.Frame
	elapsed     types.Frame
	end         types.Frame
}

// NewMessageWindow returns a new message window, where the speed is how many frames
// it takes between each letter
func NewMessageWindow(x, y, w, h int, c color.Color, message string, speed types.Frame) *MessageWindow {
	end := len(message) * speed
	return &MessageWindow{NewBlankWindow(x, y, w, h, c), message, "", speed, 0, end}
}

// Act ticks up the window elapsed frames and shows more of the message
func (w *MessageWindow) Act(df types.Frame) {
	if w.elapsed < w.end {
		w.elapsed += df
		substr := int(w.elapsed / w.speed)
		w.currMessage = w.message[0:substr]
	} else {
		w.currMessage = w.message
	}
}

// Draw does a draw
func (w *MessageWindow) Draw(img *ebiten.Image) {

	fmt.Printf(w.currMessage)
	w.skin.Sprite.Draw(w.x, w.y, img)

	font := cache.Get().LoadFont(config.Get().Font())
	text.Draw(img, w.currMessage, font, w.x, w.y, color.White)
}

// HandleInput - so long as the message window is active, it will consume input.
func (w *MessageWindow) HandleInput(state input.Input) bool {
	cfg := config.Get()

	if state[cfg.KeyConfirm()].Pressed() {
		if w.elapsed >= w.end {
			w.dispose()
		} else {
			w.elapsed = w.end
		}
	}

	if state[cfg.KeyCancel()].Pressed() {
		w.elapsed = w.end
		w.dispose()
	}
	return true
}
