package room

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Data is unmarshaled from a room json file
type Data struct {
	Name          string
	Width, Height int
	Actors        []*ActorData
}

// ActorData false
type ActorData struct {
	Name     string
	Sprite   *spriteData
	Collider *colliderData
}

type spriteData struct {
	*TileSpriteData
	*ShapeSpriteData
}

// TileSpriteData false
type TileSpriteData struct {
	Kind       string
	Rows, Cols int
	Sheet      string
	Tiles      []int
}

// ShapeSpriteData false
type ShapeSpriteData struct {
	Kind string
	// TODO: implement a shape sprite data
}

type colliderData struct {
	*BlockColliderData
	*TriangleColliderData
}

// BlockColliderData false
type BlockColliderData struct {
	Kind             string
	X, Y, Z, W, H, D int
	Blocking         bool
	Name             string
}

// TriangleColliderData false
type TriangleColliderData struct {
	Kind                                 string
	X, Y, Z, Rx2, Ry2, Rx3, Ry3, D, Axis int
	Blocking                             bool
	Name                                 string
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
