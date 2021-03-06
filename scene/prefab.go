package scene

import (
	"fmt"

	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/utils"
)

// NewBoundaries creates boundaries for a room.
func NewBoundaries(roomWidth, roomHeight int) []actors.Actor {
	dimx, dimy := config.Get().TileDimX, config.Get().TileDimY
	width, height := dimx*roomWidth, dimy*roomHeight
	//north
	north := actors.NewInvisibleActor("boundary", colliders.NewBlock(-dimx, -dimy, 0, width+(dimx*2), dimy, 99999, true, "north_boundary"))
	//east
	east := actors.NewInvisibleActor("boundary", colliders.NewBlock(width, -dimy, 0, dimx, height+(dimy*2), 99999, true, "east_boundary"))
	//south
	south := actors.NewInvisibleActor("boundary", colliders.NewBlock(-dimx, height, 0, width+(dimx*2), dimy, 99999, true, "south_boundary"))
	//west
	west := actors.NewInvisibleActor("boundary", colliders.NewBlock(-dimx, -dimy, 0, dimx, height+(dimy*2), 99999, true, "west_boundary"))
	return []actors.Actor{north, east, south, west}
}

// NewTrampoline returns a trampoline actor
func NewTrampoline(x, y, z int, sprite sprites.Spritemap) actors.Actor {
	rock := actors.NewStaticActor(
		"wall",
		sprite,
		colliders.NewBlock(x, y, z, 12, 8, 8, true, fmt.Sprintf("manual-trampoline")),
		-2, -8,
	)
	reaction := events.NewReaction(func(args ...interface{}) {
		subject := args[0].(actors.CanMove)
		object := args[1].(actors.Actor)

		vx, vy, vz := subject.Vel()
		_, _, sz := subject.Collider().Pos()
		ox, oy, oz := object.Collider().Pos()
		od := object.Collider().ZDepth(ox, oy)

		var upward float64
		pressed := input.State()[config.Get().KeyJump()]
		if pressed.PressedWindow(0, 24) {
			upward = 5.0
		} else {
			upward = 3.3
		}

		if sz >= oz+od && vz < 0 {
			subject.SetVel(vx, vy, 0)
			subject.SetOnGround(false)
			events.Enqueue(events.New(1, actors.DashActionType, []interface{}{subject, 0.0, 0.0, upward}))
		}
	})
	rock.Collider().Reactions().Push(events.ReactionOnCollision, reaction)
	return rock
}

// CreateShadow creates a shadow actor for a given drawable actor
func CreateShadow(subject actors.Actor, shadowSprite *sprites.Sprite) (actors.Actor, actors.PostCollisionHook) {
	x, y, z := subject.Pos()
	w := subject.Collider().XDepth(y, z)
	h := subject.Collider().YDepth(x, z)
	collider := colliders.NewBlock(x, y, z, w, h, 9, false, fmt.Sprintf("%s-shadow", subject.Collider().Name()))
	ox, oy := subject.(actors.Drawable).DrawOffset()
	shadow := actors.NewStaticActor("shadow", sprites.NewStaticSpritemap(shadowSprite), collider, ox, oy+4)

	hook := actors.NewShadowHook(shadow, subject)
	return shadow, hook
}

// NewPushBlock - a block that can be pushed!
func NewPushBlock(x, y, z int, name string, sprite sprites.Spritemap) actors.Actor {
	block := actors.NewMovingActor(
		"block",
		sprite,
		colliders.NewBlock(x, y, z, 16, 16, 15, true, name),
		0, -16, 10, true,
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
					events.Actor, actors.MoveByActionType, []interface{}{object, dx * 16.0, dy * 16.0, 0.0, 32},
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
	block.Collider().Reactions().Push(events.ReactionOnCollision, reaction)
	// interaction := events.NewReaction(func(args ...interface{}) {
	// 	events.Enqueue(
	// 		events.NewMessageWindowEvent(0, (config.Get().ScreenHeight()*2)/3,
	// 			config.Get().ScreenWidth(), (config.Get().ScreenHeight()/3)+1,
	// 			"Hello! This is a test of drawing text\non a message window.\nNeato!~"),
	// 	)
	// })
	interaction := events.NewMessageReaction([]string{"Hello! This is a test of drawing text\non a message window.\nNeato!~", "This is a second message!"})
	block.Collider().Reactions().Push(events.ReactionOnInteraction, interaction)
	return block
}
