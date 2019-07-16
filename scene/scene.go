package scene

import (
	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/room"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/sprites"

	"github.com/hajimehoshi/ebiten"
)

// Scene -	coordinates window, actor, and room entities.
// 			processes inputs, delegates queued events, triggers actions, and
//			resolves collisions.
type Scene struct {
	actorM *actors.Manager
	room   *room.Room
	tiles  *sprites.Spritesheet
}

// New w
func New(
	player *actors.Actor,
	room *room.Room,
	tiles *sprites.Spritesheet,
) *Scene {
	mgr := actors.NewManager()
	mgr.SetPlayer(player)
	return &Scene{mgr, room, tiles}
}

// Update w
func (s *Scene) Update(df types.Frame) {
	// first process inputs
	state := input.State().Tick(df)
	blocked := false
	if !blocked {
		blocked = s.actorM.HandleInput(state)
	}

	// then process/delegate events
	s.processEvents()

	// then call the manager act() functions
	s.act(df)

	// resolve collisions of actor against room based on staged actor velocities
	s.resolveCollisions()

	// allow actor manager to resolve collisions between actors
	// (which may generate more events)
}

func (s *Scene) processEvents() {
	// queue := events.Hub().ActorEvents()
	// for queue.HasNext() {
	// 	ev := queue.Read()
	// 	ev.Process(s.actorM)
	// }
}

func (s *Scene) act(df int) {
	s.actorM.Act(df)
}

func (s *Scene) resolveCollisions() {
	s.actorM.ResolveCollisions(s.room.Colliders())
}

// Render - called by main render loop
func (s *Scene) Render(img *ebiten.Image) *ebiten.Image {
	for row := 0; row < 8; row++ {
		for pr := 0; pr < 10; pr++ {
			s.RenderRow(img, pr, row)
		}
	}

	return img
}

// RenderRow - render a single row
func (s *Scene) RenderRow(img *ebiten.Image, pr, row int) *ebiten.Image {
	for _, layer := range s.room.Layers() {
		if pr != layer.Priority() {
			continue
		}
		
		mapTiles := layer.TilesRow(row, 10)
		for col := 0; col < len(mapTiles); col++ {
			tile := s.tiles.GetSprite(mapTiles[col])
			tile.DrawSprite(col*16, row*16, img)
		}
	}

	s.actorM.Render(img, pr, row)

	return img
}
