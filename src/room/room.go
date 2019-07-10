package room

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Room woo
type Room struct {
	layers []*Layer
}

// Layers woo
func (room *Room) Layers() []*Layer { return room.layers }

// NewRoomFromFile woo
func NewRoomFromFile(source string) *Room {
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	return &Room{parseRoomFile(r)}
}

func parseRoomFile(r *csv.Reader) []*Layer {
	var ret []*Layer
	var curr []string
	var priority, width, height int
	var appender = func() []*Layer {
		return append(ret, NewLayer(
			parseLayer(padArray(curr, width*height)),
			priority))
	}

	for {
		line, err := r.Read()
		if err == io.EOF {
			ret = appender()
			break
		}

		if strings.HasPrefix(line[0], "width") {
			width, _ = strconv.Atoi(strings.Split(line[0], " ")[1])
		} else if strings.HasPrefix(line[0], "height") {
			height, _ = strconv.Atoi(strings.Split(line[0], " ")[1])
		} else if strings.HasPrefix(line[0], "layer") {
			if curr != nil {
				ret = appender()
			}
			priority, err = strconv.Atoi(strings.Split(line[0], " ")[1])
			curr = []string{}
		} else {
			curr = append(curr, padArray(line, width)...)
		}
	}
	sort.Sort(ByPriority(ret))
	return ret
}

func parseLayer(strs []string) []int {
	var ret = make([]int, len(strs))
	for i, v := range strs {
		ret[i], _ = strconv.Atoi(v)
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
