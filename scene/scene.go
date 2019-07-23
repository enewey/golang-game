package scene

import (
	"enewey.com/golang-game/actors"
	"enewey.com/golang-game/input"
	"enewey.com/golang-game/room"
	"enewey.com/golang-game/sprites"
	"enewey.com/golang-game/types"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/utils"

	"github.com/hajimehoshi/ebiten"
)

// Scene -	coordinates window, actor, and room entities.
// 			processes inputs, delegates queued events, triggers actions, and
//			resolves collisions.
type Scene struct {
	actorM *actors.Manager
	room   *room.Room
	tiles  *sprites.Spritesheet
	offsetX, offsetY int  // room rendering offsets
}

var cfg *config.Config
func init() {
	cfg = config.Get()
}

// New w
func New(
	player *actors.Actor,
	room *room.Room,
	tiles *sprites.Spritesheet,
) *Scene {
	mgr := actors.NewManager()
	mgr.SetPlayer(player)
	px,py,pz := player.Pos()
	ox,oy := getScrollOffset(
		room.Width()*cfg.TileDimX, 
		room.Height()*cfg.TileDimY,
		0, 0,
		px,py,pz)
	return &Scene{mgr, room, tiles, ox, oy}
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

	//At the end of it, get the player's position and adjust the scroll offset
	px,py,pz := s.actorM.GetPlayer().Pos()
	s.offsetX, s.offsetY = getScrollOffset(
		s.room.Width()*cfg.TileDimX, 
		s.room.Height()*cfg.TileDimY,
		s.offsetX, s.offsetY,
		px, py, pz)
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
	y := py-pz // visual y coordinate the sum of actual y and negative z
	if y - oy > bd {
		rety = utils.Min(y-bd, maxY)
	} else if y - oy < bu {
		rety = utils.Max(y-bu, 0)
	}

	return retx,rety
}

// Render - called by main render loop
func (s *Scene) Render(img *ebiten.Image) *ebiten.Image {
	rowOffset := s.offsetY / 16
	rowmax := rowOffset+cfg.TilesY+1
	if rowmax > s.room.Height() {
		rowmax--
	}
	for pr := 0; pr < 10; pr++ {
		for row := rowOffset; row < rowmax; row++ {
			s.RenderRow(img, pr, row)
		}
	}

	return img
}

// RenderRow - render a single row
func (s *Scene) RenderRow(img *ebiten.Image, pr, row int) *ebiten.Image {
	colOffset := s.offsetX / 16
	rowOffset := s.offsetY / 16
	xPixelOffset := s.offsetX % 16
	yPixelOffset := s.offsetY % 16
	for _, layer := range s.room.Layers() {
		if pr != layer.Priority() {
			continue
		}

		mapTiles := layer.TilesRow(row, s.room.Width())
		if colOffset+cfg.TilesX < len(mapTiles) {
			mapTiles = mapTiles[colOffset:colOffset+cfg.TilesX+1]
		} else {
			mapTiles = mapTiles[colOffset:colOffset+cfg.TilesX]
		}
		for col := 0; col < len(mapTiles); col++ {
			tile := s.tiles.GetSprite(mapTiles[col])
			tile.DrawSprite(
				(col*cfg.TileDimX)-xPixelOffset,
				((row - rowOffset)*cfg.TileDimY)-yPixelOffset,
				img)
		}
	}

	s.actorM.Render(img, pr, row, s.offsetX, s.offsetY)

	return img
}
