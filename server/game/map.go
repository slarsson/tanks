package game

import (
	"encoding/json"
	"fmt"
)

type Map struct {
	Obstacles  []*Polygon
	Boundaries [4]float32 // maxX, minX, maxY, minY
}

type MapData struct {
	Name       string           `json:"name"`
	Boundaries [4]float32       `json:"boundaries"`
	Containers []ContainerGroup `json:"containers"`
}

type ContainerGroup struct {
	Position Vector3 `json:"position"`
	Total    int     `json:"total"`
	Bottom   int     `json:"bottom"`
}

func NewMap() *Map {
	data := []byte(`{
		"name": "Port of Nrkp",
		"boundaries": [50, -50, 50, -50], 
		"containers": [
			{
				"position": {
					"x": 40,
					"y": 0,
					"z": 0
				},
				"total": 25,
				"bottom": 10
			},
			{
				"position": {
					"x": -20,
					"y": 15,
					"z": 0
				},
				"total": 5,
				"bottom": 5
			},
			{
				"position": {
					"x": 0,
					"y": 30,
					"z": 0
				},
				"total": 2,
				"bottom": 1
			}
		]
	}`)

	var manifest MapData
	err := json.Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}
	fmt.Println(manifest.Boundaries)

	obstacles := []*Polygon{}

	for _, val := range manifest.Containers {
		//fmt.Println(val)
		//poly := &Polygon{}

		yy := 0.5 * float32(val.Bottom) * 3.75
		xx := 0.5 * float32(8)
		//fmt.Println("size1:", s1, s2)

		c1 := &Vector3{
			X: val.Position.X - xx,
			Y: val.Position.Y + yy,
			Z: 0,
		}

		c2 := &Vector3{
			X: val.Position.X + xx,
			Y: val.Position.Y + yy,
			Z: 0,
		}

		c3 := &Vector3{
			X: val.Position.X + xx,
			Y: val.Position.Y - yy,
			Z: 0,
		}

		c4 := &Vector3{
			X: val.Position.X - xx,
			Y: val.Position.Y - yy,
			Z: 0,
		}

		obstacles = append(obstacles, &Polygon{c1, c2, c3, c4})
	}

	// x := []*Polygon{}

	// for _, arr := range jsonData.Blocks {
	// 	poly := &Polygon{}
	// 	for _, point := range arr.Coords {
	// 		poly.Add(point[0], point[1], 0)
	// 		//poly = append(*poly, &Vector3{X: point[0], Y: point[1], Z: 0})
	// 	}
	// 	x = append(x, poly)
	// }
	// fmt.Println(x[0])

	// test := &Polygon{
	// 	&Vector3{X: 0, Y: 16, Z: 0},
	// 	&Vector3{X: 10, Y: 16, Z: 0},
	// 	&Vector3{X: 10, Y: 15, Z: 0},
	// 	&Vector3{X: 0, Y: 15, Z: 0},
	// }

	return &Map{
		Obstacles: obstacles,
		//Boundaries: [4]float32{50, -50, 50, -50},
		Boundaries: manifest.Boundaries,
	}
}

func (m *Map) OutOfBounds(point *Vector3) bool {
	if point.X > m.Boundaries[0] || point.X < m.Boundaries[1] || point.Y > m.Boundaries[2] || point.Y < m.Boundaries[3] {
		fmt.Println("out of map")
		return true
	}
	return false
}
