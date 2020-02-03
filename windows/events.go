package windows

import (
	"fmt"

	"enewey.com/golang-game/config"
	"enewey.com/golang-game/events"
)

// Window types as Event codes
const (
	Message = iota
)

// InterpretEvent translates an event into a window
func InterpretEvent(ev *events.Event) Window {
	cfg := config.Get()
	p := ev.Payload()
	fmt.Printf("interpreting window event %d :: ", ev.Code())
	switch ev.Code() {
	case Message:
		return messageWindowEvent(p)
	default:
		fmt.Printf("unknown window event code %d\n", ev.Code())
		return NewMessageWindow(0, 0, 100, 100, cfg.WindowColor(), "", cfg.TextSpeed())
	}
}

func messageWindowEvent(p []interface{}) *MessageWindow {
	cfg := config.Get()
	x, y, w, h := p[0].(int), p[1].(int), p[2].(int), p[3].(int)
	msg := p[4].(string)
	fmt.Printf("message window interpreted %d %d %d %d %s", x, y, w, h, msg)
	return NewMessageWindow(x, y, w, h, cfg.WindowColor(), msg, cfg.TextSpeed())
}
