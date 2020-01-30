package scene

import (
	"fmt"

	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/utils"
)

// NewTrampoline returns a trampoline actor
func NewTrampoline(x, y, z int, sprite sprites.Spritemap) actors.Actor {
	rock := actors.NewStaticActor(
		"wall",
		sprite,
		colliders.NewBlock(x, y, z, 12, 8, 8, true, true, fmt.Sprintf("manual-trampoline")),
		-2, -8,
	)
	reaction := events.NewReaction(func(args ...interface{}) {
		subject := args[0].(actors.CanMove)
		object := args[1].(actors.Actor)

		vx, vy, vz := subject.Vel()
		_, _, sz := subject.Collider().Pos()
		ox, oy, oz := object.Collider().Pos()
		od := object.Collider().ZDepth(ox, oy)
		if sz >= oz+od && vz < 0 {
			subject.SetVel(vx, vy, 0)
			subject.SetOnGround(false)
			events.Enqueue(events.New(1, actors.Dash, []interface{}{subject, 0.0, 0.0, (vz * -1.0)}))
		}
	})
	rock.Collider().SetReaction(reaction)
	return rock
}

// CreateShadow creates a shadow actor for a given drawable actor
func CreateShadow(subject actors.Actor, shadowSprite *sprites.Sprite) (actors.Actor, actors.PostCollisionHook) {
	x, y, z := subject.Pos()
	w := subject.Collider().XDepth(y, z)
	h := subject.Collider().YDepth(x, z)
	collider := colliders.NewBlock(x, y, z, w, h, 9, false, false, fmt.Sprintf("%s-shadow", subject.Collider().Name()))
	ox, oy := subject.(actors.Drawable).DrawOffset()
	shadow := actors.NewStaticActor("shadow", sprites.NewStaticSpritemap(shadowSprite), collider, ox, oy+4)

	hook := actors.NewShadowHook(shadow, subject)
	return shadow, hook
}

// NewPushBlock - a block that can be pushed!
func NewPushBlock(x, y, z int, sprite sprites.Spritemap) actors.Actor {
	block := actors.NewMovingActor(
		"block",
		sprite,
		colliders.NewBlock(x, y, z, 16, 16, 15, true, true, "push-block-boi"),
		0, -16, true,
	)
	reaction := events.NewAfterConsecutiveReaction(
		func(args ...interface{}) {
			fmt.Printf("reaction triggered\n")
			subject := args[0].(actors.CanMove)
			object := args[1].(actors.CanMove)

			x1, y1, z1 := object.Collider().Center()
			x2, y2, z2 := subject.Collider().Center()

			dx, dy, _ := utils.DominantAxis(utils.Cast(
				float64(x1), float64(y1), float64(z1),
				float64(x2), float64(y2), float64(z2),
			))
			events.Enqueue(
				events.New(
					events.Actor, actors.MoveBy, []interface{}{object, int(dx * 16), int(dy * 16), 0, 32},
				),
			)
		},
		func(args ...interface{}) bool {
			// subject := args[1].(actors.CanMove)
			object := args[1].(actors.CanMove)
			vx, vy, vz := object.Vel()
			// x, y, z := object.Pos()
			return (vx == 0 && vy == 0 && vz == 0)
		},
		30,
		120,
	)
	block.Collider().SetReaction(reaction)
	return block
}
