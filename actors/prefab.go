package actors

import (
	"fmt"

	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/sprites"
)

// NewTrampoline returns a trampoline actor
func NewTrampoline(x, y, z int, sprite sprites.Spritemap) Actor {
	rock := NewStaticActor(
		"wall",
		sprite,
		colliders.NewBlock(x, y, z, 12, 8, 8, true, true, fmt.Sprintf("manual-trampoline")),
		-2, -8,
	)
	reaction := func(args ...interface{}) {
		subject := args[0].(CanMove)
		object := args[1].(Actor)

		_, _, vz := subject.Vel()
		_, _, sz := subject.Collider().Pos()
		ox, oy, oz := object.Collider().Pos()
		od := object.Collider().ZDepth(ox, oy)
		if sz >= oz+od && vz < 0 {
			events.Enqueue(events.New(1, 3, []interface{}{subject, 0.0, 0.0, 2.5}))
		}
	}
	rock.Collider().SetReaction(reaction)
	return rock
}
