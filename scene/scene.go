package scene

import (
	"fmt"

	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/events"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/utils"
	"enewey.com/golang-game/windows"

	"github.com/hajimehoshi/ebiten"
)

// Scene -	coordinates window, actor, and room entities.
// 			processes inputs, delegates queued events, triggers actions, and
//			resolves collisions.
type Scene struct {
	WindowM          *windows.Manager
	ActorM           *actors.Manager
	width, height    int
	offsetX, offsetY int // room rendering offsets
}

var cfg *config.Config

func init() {
	cfg = config.Get()
}

// New w
func New(player actors.Actor, dataFile string) *Scene {
	wmgr := windows.NewManager()
	mgr := actors.NewManager()
	mgr.SetPlayer(player)

	room := createRoom(FromJSON(dataFile))
	boundaries := NewBoundaries(room.Width, room.Height)
	for _, bound := range boundaries {
		mgr.AddActor(bound)
	}

	for _, actor := range room.actors {
		mgr.AddActor(actor)
	}

	px, py, pz := player.Pos()
	ox, oy := getScrollOffset(
		room.Width*cfg.TileDimX,
		room.Height*cfg.TileDimY,
		0, 0,
		px, py, pz)
	return &Scene{wmgr, mgr, room.Width, room.Height, ox, oy}
}

// func roomToActors(rm *room.Room, tiles *sprites.Spritesheet, px, dimX int) []actors.Actor {
// 	// the current room situation is such that tiles on an even priority are floors,
// 	// and tiles on an odd priority are walls.
// 	// Floor xyz calculation:
// 	//		x = column * px
// 	//		y = (row + (priority/2)) * px
// 	//		z = (priority/2) * px
// 	// Wall xyz calculation:
// 	//		x = column * px
// 	//		y = (row + 1 + (priority/2)) * px
// 	//		z = (priority/2) * px
// 	//
// 	// note the only difference is adding 1 to the Y when dealing with a wall
// 	// which only happens on an odd-numbered priority
// 	ret := []actors.Actor{}

// 	for lyrNum, lyr := range rm.Layers() {
// 		yFactor := lyr.Priority() % 2
// 		isWalls := yFactor == 1

// 		for i, tile := range lyr.Tiles() {
// 			if tile == 0 {
// 				continue
// 			}
// 			r := i / dimX
// 			c := i % dimX

// 			x := c * px
// 			y := (r + yFactor + (lyr.Priority() / 2)) * px
// 			z := (lyr.Priority() / 2) * px

// 			if isWalls {
// 				ret = append(ret, actors.NewSpriteActor(
// 					"wall",
// 					sprites.NewStaticSpritemap(tiles.GetSprite(tile)),
// 					colliders.NewBlock(x, y, z, px, 0, px, false, fmt.Sprintf("wall-%d-%d", lyrNum, i)),
// 					0, -px,
// 				))
// 			} else {
// 				ret = append(ret, actors.NewSpriteActor(
// 					"floor",
// 					sprites.NewStaticSpritemap(tiles.GetSprite(tile)),
// 					colliders.NewBlock(x, y, z, px, px, 0, false, fmt.Sprintf("floor-%d-%d", lyrNum, i)),
// 					0, 0,
// 				))
// 			}
// 		}
// 	}

// 	return ret
// }

// AddActor adds an actor to the scene
func (s *Scene) AddActor(actor actors.Actor) {
	s.ActorM.AddActor(actor)
}

// Update - main update loop
func (s *Scene) Update(df types.Frame) {
	// first process inputs
	state := input.State().Tick(df)

	// windows take priority over actors
	if !s.WindowM.HandleInput(state) {
		s.ActorM.HandleInput(state)
	}

	// then process/delegate events
	s.processEvents()

	// then call the manager act() functions
	if !s.WindowM.Act(df) {
		// actors only get to act if window manager doesnt declare focus
		s.ActorM.Act(df)
		s.ActorM.ResolveCollisions()
	}

	//At the end of it, get the player's position and adjust the scroll offset
	px, py, pz := s.ActorM.GetPlayer().Pos()
	s.offsetX, s.offsetY = getScrollOffset(
		s.width*cfg.TileDimX,
		s.height*cfg.TileDimY,
		s.offsetX, s.offsetY,
		px, py, pz)
}

func (s *Scene) processEvents() {
	for events.HasNext() {
		ev := events.Read()
		fmt.Printf("processing events %d :: ", ev.Code())
		switch ev.Scope() {
		case events.Actor:
			s.ActorM.Actions().Add(actors.InterpretEvent(ev))
		case events.Window:
			s.WindowM.AddWindow(windows.InterpretEvent(ev))
		default:
			fmt.Printf("unknown event scope %d\n", ev.Scope())
			continue
		}
	}
}

func getScrollOffset(w, h, ox, oy, px, py, pz int) (int, int) {
	var retx, rety = ox, oy
	bu, br, bd, bl := cfg.ScrollBoundaries() // up, right, down, left

	// scroll limits based on screen width and actual room width
	maxX := w - cfg.ScreenWidth()
	maxY := h - cfg.ScreenHeight()

	// x-scroll adjustment
	if px-ox > br {
		retx = utils.Min(px-br, maxX)
	} else if px-ox < bl {
		retx = utils.Max(px-bl, 0)
	}

	// y-scroll adjustment
	y := py - pz // visual y coordinate the sum of actual y and negative z
	if y-oy > bd {
		rety = utils.Min(y-bd, maxY)
	} else if y-oy < bu {
		rety = utils.Max(y-bu, 0)
	}

	return retx, rety
}

// Render - called by main render loop
func (s *Scene) Render(img *ebiten.Image) *ebiten.Image {
	s.ActorM.Render(img, s.offsetX, s.offsetY)
	s.WindowM.Render(img, s.offsetX, s.offsetY)
	// windows render on TOP.. i.e. AFTER

	return img
}
