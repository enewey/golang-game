package room

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// Room woo
type Room struct {
	layers []*Layer
}

// NewRoomFromFile woo
func NewRoomFromFile(source string) *Room {
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	var layers []*Layer
	var priority int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		layers = append(layers, NewLayer(parseLayer(record), priority))
	}

	return &Room{layers}
}

func parseLayer(strs []string) []int {
	var ret = make([]int, len(strs))
	for i, v := range strs {
		ret[i], _ = strconv.Atoi(v)
	}

	return ret
}

// Layers woo
func (room *Room) Layers() []*Layer { return room.layers }
