package room

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"enewey.com/golang-game/collider"
)

// Room woo
type Room struct {
	layers []*Layer
	blocks collider.Colliders
}

// Layers woo
func (room *Room) Layers() []*Layer { return room.layers }

// Colliders woo
func (room *Room) Colliders() collider.Colliders { return room.blocks }

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

func parseRoomFile(r *csv.Reader) ([]*Layer, collider.Colliders) {
	var retLyr []*Layer
	var retBlk = make(collider.Colliders, 0)
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
			retBlk = append(retBlk, parseBlock(line[1:]))
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

func parseBlock(strs []string) *collider.Collider {
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
			y = int(yy * 16)
			break
		case "x":
			xx, err = strconv.ParseFloat(sp[1], 64)
			x = int(xx * 16)
			break
		case "z":
			zz, err = strconv.ParseFloat(sp[1], 64)
			z = int(zz * 16)
			break
		case "w":
			ww, err = strconv.ParseFloat(sp[1], 64)
			w = int(ww * 16)
			break
		case "h":
			hh, err = strconv.ParseFloat(sp[1], 64)
			h = int(hh * 16)
			break
		case "d":
			dd, err = strconv.ParseFloat(sp[1], 64)
			d = int(dd * 16)
			break
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	return collider.NewBlock(x, y, z, w, h, d, name)
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
