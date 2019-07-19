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
	"enewey.com/golang-game/utils"
)

// Room - encapsulates what is parsed from a .room file.
type Room struct {
	layers []*Layer
	walls  colliders.Colliders
}

// Layers woo
func (room *Room) Layers() []*Layer { return room.layers }

// Colliders woo
func (room *Room) Colliders() colliders.Colliders { return room.walls }

// NewRoomFromFile woo
func NewRoomFromFile(source string) *Room {
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	layers, blocks := parseRoomFile(r)
	return &Room{layers, blocks}
}

func parseRoomFile(r *csv.Reader) ([]*Layer, colliders.Colliders) {
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
	return retLyr, retBlk
}

func parseBlock(strs []string) *colliders.Block {
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
		y = utils.Flint(yy * 16)
		x = utils.Flint(xx * 16)
		z = utils.Flint(zz * 16)
		w = utils.Flint(ww * 16)
		h = utils.Flint(hh * 16)
		d = utils.Flint(dd * 16)
	}

	return colliders.NewBlock(x, y, z, w, h, d, name).(*colliders.Block)
}

func parseTriangle(strs []string) *colliders.Triangle {
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
		ry2 = utils.Flint(yy2 * 16)
		rx2 = utils.Flint(xx2 * 16)
		ry3 = utils.Flint(yy3 * 16)
		rx3 = utils.Flint(xx3 * 16)
		x = utils.Flint(xx * 16)
		y = utils.Flint(yy * 16)
		z = utils.Flint(zz * 16)
		d = utils.Flint(dd * 16)
	}

	fmt.Printf("triangle created %d %d %d %d %d %d %d %d %s\n", x, y, z, rx2, ry2, rx3, ry3, d, name)

	return colliders.
		NewTriangle(x, y, z, rx2, ry2, rx3, ry3, d, axis, name).(*colliders.Triangle)
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
