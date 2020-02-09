package scene

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Data is unmarshaled from a room json file
type Data struct {
	Name   string       `json:"name"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
	Actors []*ActorData `json:"actors"`
}

// ActorData false
type ActorData struct {
	Name     string        `json:"name"`
	Kind     string        `json:"kind"`
	Sprite   *spriteData   `json:"sprite"`
	Collider *colliderData `json:"collider"`
	OffsetX  int           `json:"offsetX"`
	OffsetY  int           `json:"offsetY"`
}

type spriteData struct {
	*TileSpriteData
	*ShapeSpriteData
	Kind string `json:"kind"`
}

// TileSpriteData false
type TileSpriteData struct {
	Rows  int    `json:"rows"`
	Cols  int    `json:"cols"`
	Sheet string `json:"sheet"`
	Tiles []int  `json:"tiles"`
}

// ShapeSpriteData false
type ShapeSpriteData struct {
	// TODO: implement a shape sprite data
}

type colliderData struct {
	*BlockColliderData
	*TriangleColliderData
	Kind     string `json:"kind"`
	Blocking bool   `json:"blocking"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Z        int    `json:"z"`
	D        int    `json:"d"`
	Name     string `json:"name"`
}

// BlockColliderData false
type BlockColliderData struct {
	W int `json:"w"`
	H int `json:"h"`
}

// TriangleColliderData false
type TriangleColliderData struct {
	Rx2  int `json:"rx2"`
	Ry2  int `json:"ry2"`
	Rx3  int `json:"rx3"`
	Ry3  int `json:"ry3"`
	Axis int `json:"axis"`
}

// FromJSON creates a new Data struct from the bytes of json file
func FromJSON(source string) *Data {
	var out Data
	body, ferr := ioutil.ReadFile(source)
	if ferr != nil {
		fmt.Printf("error reading file")
		panic(ferr)
	}
	if err := json.Unmarshal(body, &out); err != nil {
		fmt.Printf("error unmarshaling json")
		panic(err)
	}

	return &out
}
