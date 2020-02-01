package room

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"enewey.com/golang-game/colliders"
	"enewey.com/golang-game/config"
	"enewey.com/golang-game/utils"
)

// Room - encapsulates what is parsed from a .room file.
type Room struct {
	layers        []*Layer
	walls         colliders.Colliders
	width, height int
}

// Layers woo
func (room *Room) Layers() []*Layer { return room.layers }

// Colliders woo
func (room *Room) Colliders() colliders.Colliders { return room.walls }

// Width is the room width in tiles
func (room *Room) Width() int { return room.width }

// Height is the room height in tiles
func (room *Room) Height() int { return room.height }

// NewRoomFromFile woo
func NewRoomFromFile(source string) *Room {
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	layers, blocks, w, h := parseRoomFile(r)
	return &Room{layers, blocks, w, h}
}

func parseRoomFile(r *csv.Reader) ([]*Layer, colliders.Colliders, int, int) {
	var retLyr []*Layer
	var retBlk = make(colliders.Colliders, 0)
	var curr []string
	var priority, width, height int
	var appender = func() []*Layer {
		return append(retLyr, NewLayer(
			parseLayer(padArray(curr, width*height)),
			priority,
		))
	}

	for {
		line, err := r.Read()
		if err == io.EOF {
			retLyr = appender()
			break
		}

		if len(line) == 1 && strings.TrimSpace(line[0]) == "" {
			continue
		} else if strings.HasPrefix(line[0], "collider") {
			sp := strings.Split(line[0], " ")
			switch sp[1] {
			case "triangle":
				retBlk = append(retBlk, parseTriangle(line[1:]))
				break
			case "block":
				retBlk = append(retBlk, parseBlock(line[1:]))
				break
			}
		} else if strings.HasPrefix(line[0], "width") {
			width, _ = strconv.Atoi(strings.Split(line[0], " ")[1])
		} else if strings.HasPrefix(line[0], "height") {
			height, _ = strconv.Atoi(strings.Split(line[0], " ")[1])
		} else if strings.HasPrefix(line[0], "layer") {
			if curr != nil {
				retLyr = appender()
			}
			priority, err = strconv.Atoi(strings.Split(line[0], " ")[1])
			curr = []string{}
		} else {
			curr = append(curr, padArray(line, width)...)
		}
	}
	sort.Sort(ByPriority(retLyr))
	return retLyr, retBlk, width, height
}

func parseBlock(strs []string) colliders.Collider {
	var y, x, z, w, h, d int
	var yy, xx, zz, ww, hh, dd float64
	var name string
	var err error
	for _, v := range strs {
		sp := strings.Split(v, " ")
		switch sp[0] {
		case "name":
			name = sp[1]
		case "y":
			yy, err = strconv.ParseFloat(sp[1], 64)
			break
		case "x":
			xx, err = strconv.ParseFloat(sp[1], 64)
			break
		case "z":
			zz, err = strconv.ParseFloat(sp[1], 64)
			break
		case "w":
			ww, err = strconv.ParseFloat(sp[1], 64)
			break
		case "h":
			hh, err = strconv.ParseFloat(sp[1], 64)
			break
		case "d":
			dd, err = strconv.ParseFloat(sp[1], 64)
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		tdx := float64(config.Get().TileDimX)
		y = utils.Flint(yy * tdx)
		x = utils.Flint(xx * tdx)
		z = utils.Flint(zz * tdx)
		w = utils.Flint(ww * tdx)
		h = utils.Flint(hh * tdx)
		d = utils.Flint(dd * tdx)
	}

	return colliders.NewBlock(x, y, z, w, h, d, true, name)
}

func parseTriangle(strs []string) colliders.Collider {
	var rx2, ry2, rx3, ry3, x, y, z, d, axis int
	var xx2, yy2, xx3, yy3, xx, yy, zz, dd float64
	var name string
	var err error
	for _, v := range strs {
		sp := strings.Split(v, " ")
		switch sp[0] {
		case "name":
			name = sp[1]
		case "axis":
			if sp[1] == "x" {
				axis = colliders.XAxis
			} else if sp[1] == "y" {
				axis = colliders.YAxis
			} else {
				axis = colliders.ZAxis
			}
		case "ry2":
			yy2, err = strconv.ParseFloat(sp[1], 64)
			break
		case "rx2":
			xx2, err = strconv.ParseFloat(sp[1], 64)
			break
		case "ry3":
			yy3, err = strconv.ParseFloat(sp[1], 64)
			break
		case "rx3":
			xx3, err = strconv.ParseFloat(sp[1], 64)
			break
		case "x":
			xx, err = strconv.ParseFloat(sp[1], 64)
			break
		case "y":
			yy, err = strconv.ParseFloat(sp[1], 64)
			break
		case "z":
			zz, err = strconv.ParseFloat(sp[1], 64)
			break
		case "d":
			dd, err = strconv.ParseFloat(sp[1], 64)
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		tdx := float64(config.Get().TileDimX)
		ry2 = utils.Flint(yy2 * tdx)
		rx2 = utils.Flint(xx2 * tdx)
		ry3 = utils.Flint(yy3 * tdx)
		rx3 = utils.Flint(xx3 * tdx)
		x = utils.Flint(xx * tdx)
		y = utils.Flint(yy * tdx)
		z = utils.Flint(zz * tdx)
		d = utils.Flint(dd * tdx)
	}

	fmt.Printf("triangle created %d %d %d %d %d %d %d %d %s\n", x, y, z, rx2, ry2, rx3, ry3, d, name)

	return colliders.
		NewTriangle(x, y, z, rx2, ry2, rx3, ry3, d, axis, true, name)
}

func parseLayer(strs []string) []int {
	var ret = make([]int, len(strs))
	for i, v := range strs {
		ret[i], _ = strconv.Atoi(strings.TrimSpace(v))
	}

	return ret
}

func padArray(arr []string, size int) []string {
	ret := arr[:]
	if len(ret) < size {
		ret = append(ret, make([]string, size-len(ret))...)
		for i := len(ret); i < size; i++ {
			ret[i] = "0"
		}
	}
	return ret
}
