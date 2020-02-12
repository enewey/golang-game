package events

import (
	"enewey.com/golang-game/config"
)

// NewMessageReaction creates a reaction that produces one or more default messages windows.
func NewMessageReaction(messages []string) Reaction {
	cfg := config.Get()
	ret := []*Event{}
	for _, v := range messages {
		ret = append(ret, NewMessageWindowEvent(0, (cfg.ScreenHeight()*2)/3,
			cfg.ScreenWidth(), (cfg.ScreenHeight()/3)+1, v))
	}
	return NewReaction(func(...interface{}) { EnqueueAll(ret) })
}

// NewMessageWindowEvent - event for a window message
func NewMessageWindowEvent(x, y, w, h int, msg string) *Event {
	return &Event{2, 0, []interface{}{x, y, w, h, msg}}
}
