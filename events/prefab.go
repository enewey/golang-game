package events

// NewMessageWindowEvent - event for a window message
func NewMessageWindowEvent(x, y, w, h int, msg string) *Event {
	return &Event{2, 0, []interface{}{x, y, w, h, msg}}
}
