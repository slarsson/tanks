package game

import (
	"encoding/json"
	"fmt"
)

type Map struct {
	Obstacles  []*Polygon
	Boundaries [4]float32 // maxX, minX, maxY, minY
}

type MapDataBlock struct {
	Name   string      `json:"name"`
	Coords [][]float32 `json:"coords"` // clockwise rotation !!
}

type MapData struct {
	Name   string         `json:"name"`
	Blocks []MapDataBlock `json:"blocks"`
}

func NewMap() *Map {
	data := []byte(`{
		"name": "Biltema",
		"blocks": [
			{
				"name": "wall",
				"coords": [[0, 16, 0], [10, 16, 0], [10, 15, 0], [0, 15, 0]]
			},
			{
				"name": "house1",
				"coords": [[10, 10, 0], [20, 10, 0], [20, 0, 0], [10, 0, 0]]
			}
		]
	}`)

	var jsonData MapData
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		panic(err)
	}
	//fmt.Println(jsonData)

	x := []*Polygon{}

	for _, arr := range jsonData.Blocks {
		poly := &Polygon{}
		for _, point := range arr.Coords {
			poly.Add(point[0], point[1], 0)
			//poly = append(*poly, &Vector3{X: point[0], Y: point[1], Z: 0})
		}
		x = append(x, poly)
	}
	fmt.Println(x[0])

	// test := &Polygon{
	// 	&Vector3{X: 0, Y: 16, Z: 0},
	// 	&Vector3{X: 10, Y: 16, Z: 0},
	// 	&Vector3{X: 10, Y: 15, Z: 0},
	// 	&Vector3{X: 0, Y: 15, Z: 0},
	// }

	return &Map{
		Obstacles: x,
		//Obstacles:  []*Polygon{test},
		Boundaries: [4]float32{50, -50, 50, -50},
	}
}

func (m *Map) OutOfBounds(point *Vector3) bool {
	if point.X > m.Boundaries[0] || point.X < m.Boundaries[1] || point.Y > m.Boundaries[2] || point.Y < m.Boundaries[3] {
		fmt.Println("out of map")
		return true
	}
	return false
}
