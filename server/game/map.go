package game

import (
	"encoding/json"
	"math/rand"
	"time"
)

type Map struct {
	Boundaries [4]float32 // maxX, minX, maxY, minY
	Spawns     []Point2
	Obstacles  []Obstacle
}

type Obstacle struct {
	Polygon  *Polygon
	Centroid *Vector3
	Radius   float32
}

type MapData struct {
	Name       string           `json:"name"`
	Boundaries [4]float32       `json:"boundaries"`
	Spawns     []Point2         `json:"spawns"`
	Crane      Point2           `json:"crane"`
	Containers []ContainerGroup `json:"containers"`
}

type Point2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type ContainerGroup struct {
	Position Vector3 `json:"position"`
	Rotation float32 `json:"rotation"`
	Total    int     `json:"total"`
	Bottom   int     `json:"bottom"`
}

func NewMap() *Map {
	// TODO: map.json (in ./client/assets) should be synced automatically..
	data := []byte(`{
		"name": "Port of Nrkp",
		"boundaries": [50, -50, 50, -50],
		"spawns": [
			{
				"x": 0, 
				"y":0
			}, 
			{
				"x": 20, 
				"y": -20
			}
		],
		"crane": {
			"x": -25,
			"y": 0
		},
		"containers": [
			{
				"position": {
					"x": 40,
					"y": 0,
					"z": 0
				},
				"rotation": 0,
				"total": 25,
				"bottom": 10
			},
			{
				"position": {
					"x": 20,
					"y": 7,
					"z": 0
				},
				"rotation": 0,
				"total": 16,
				"bottom": 8
			},
			{
				"position": {
					"x": -25,
					"y": 0,
					"z": 0
				},
				"rotation": 1.57,
				"total": 5,
				"bottom": 5
			},
			{
				"position": {
					"x": -25,
					"y": -9,
					"z": 0
				},
				"rotation": 1.57,
				"total": 5,
				"bottom": 5
			},
			{
				"position": {
					"x": -25,
					"y": -35,
					"z": 0
				},
				"rotation": 1.57,
				"total": 5,
				"bottom": 5
			},
			{
				"position": {
					"x": -25,
					"y": 25,
					"z": 0
				},
				"rotation": 1.3,
				"total": 1,
				"bottom": 1
			},
			{
				"position": {
					"x": 0,
					"y": 30,
					"z": 0
				},
				"rotation": 0,
				"total": 1,
				"bottom": 1
			},
			{
				"position": {
					"x": 30,
					"y": -40,
					"z": 0
				},
				"rotation": -1.2,
				"total": 5,
				"bottom": 3
			},
			{
				"position": {
					"x": 0,
					"y": -40,
					"z": 0
				},
				"rotation": 1.2,
				"total": 1,
				"bottom": 1
			},
			{
				"position": {
					"x": 5,
					"y": -32,
					"z": 0
				},
				"rotation": 1.2,
				"total": 1,
				"bottom": 1
			}
		]
	}`)

	var manifest MapData
	err := json.Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	obstacles := []Obstacle{}
	for _, val := range manifest.Containers {

		// TODO: hmm.. should not hardcode the dimensions
		x := 0.5 * float32(8)
		y := 0.5 * float32(val.Bottom) * 3.75

		poly := &Polygon{
			&Vector3{
				X: val.Position.X - x,
				Y: val.Position.Y + y,
				Z: 0,
			},
			&Vector3{
				X: val.Position.X + x,
				Y: val.Position.Y + y,
				Z: 0,
			},
			&Vector3{
				X: val.Position.X + x,
				Y: val.Position.Y - y,
				Z: 0,
			},
			&Vector3{
				X: val.Position.X - x,
				Y: val.Position.Y - y,
				Z: 0,
			},
		}
		poly.Rotate(val.Rotation, &val.Position)

		obstacles = append(obstacles, Obstacle{
			Polygon:  poly,
			Centroid: val.Position.Clone(),
			Radius:   poly.FindRadius(&val.Position),
		})
	}

	// test
	craneWidth := 0.5 * float32(1)
	craneHeight := 0.5 * float32(14.75)
	//cranePosition := Vector3{X: manifest.Crane.X, Y: manifest.Crane.Y, Z: 0}

	pos1 := Vector3{
		X: manifest.Crane.X + 0.5*21.75 - 0.5,
		Y: manifest.Crane.Y,
		Z: 0,
	}

	pos2 := Vector3{
		X: manifest.Crane.X - 0.5*21.75 + 0.5,
		Y: manifest.Crane.Y,
		Z: 0,
	}

	poly1 := &Polygon{
		&Vector3{
			X: pos1.X - craneWidth,
			Y: pos1.Y + craneHeight,
			Z: 0,
		},
		&Vector3{
			X: pos1.X + craneWidth,
			Y: pos1.Y + craneHeight,
			Z: 0,
		},
		&Vector3{
			X: pos1.X + craneWidth,
			Y: pos1.Y - craneHeight,
			Z: 0,
		},
		&Vector3{
			X: pos1.X - craneWidth,
			Y: pos1.Y - craneHeight,
			Z: 0,
		},
	}

	poly2 := &Polygon{
		&Vector3{
			X: pos2.X - craneWidth,
			Y: pos2.Y + craneHeight,
			Z: 0,
		},
		&Vector3{
			X: pos2.X + craneWidth,
			Y: pos2.Y + craneHeight,
			Z: 0,
		},
		&Vector3{
			X: pos2.X + craneWidth,
			Y: pos2.Y - craneHeight,
			Z: 0,
		},
		&Vector3{
			X: pos2.X - craneWidth,
			Y: pos2.Y - craneHeight,
			Z: 0,
		},
	}

	obstacles = append(obstacles, Obstacle{
		Polygon:  poly1,
		Centroid: pos1.Clone(),
		Radius:   poly1.FindRadius(&pos1),
	})

	obstacles = append(obstacles, Obstacle{
		Polygon:  poly2,
		Centroid: pos2.Clone(),
		Radius:   poly2.FindRadius(&pos2),
	})

	// obstacles = append(obstacles, Obstacle{
	// 	Polygon:  poly2,
	// 	Centroid: p.Clone(),
	// 	Radius:   poly.FindRadius(&p),
	// })

	// p := Vector3{X: 0.5*21.75 - 0.5, Y: 0, Z: 0}
	// p.X -= 25
	// x := 0.5 * float32(1)
	// y := 0.5 * float32(14.75)

	// poly := &Polygon{
	// 	&Vector3{
	// 		X: p.X - x,
	// 		Y: p.Y + y,
	// 		Z: 0,
	// 	},
	// 	&Vector3{
	// 		X: p.X + x,
	// 		Y: p.Y + y,
	// 		Z: 0,
	// 	},
	// 	&Vector3{
	// 		X: p.X + x,
	// 		Y: p.Y - y,
	// 		Z: 0,
	// 	},
	// 	&Vector3{
	// 		X: p.X - x,
	// 		Y: p.Y - y,
	// 		Z: 0,
	// 	},
	// }
	// obstacles = append(obstacles, Obstacle{
	// 	Polygon:  poly,
	// 	Centroid: p.Clone(),
	// 	Radius:   poly.FindRadius(&p),
	// })

	// p = Vector3{X: -0.5*21.75 + 0.5, Y: 0, Z: 0}
	// p.X -= 25
	// x = 0.5 * float32(1)
	// y = 0.5 * float32(14.75)

	// poly2 := &Polygon{
	// 	&Vector3{
	// 		X: p.X - x,
	// 		Y: p.Y + y,
	// 		Z: 0,
	// 	},
	// 	&Vector3{
	// 		X: p.X + x,
	// 		Y: p.Y + y,
	// 		Z: 0,
	// 	},
	// 	&Vector3{
	// 		X: p.X + x,
	// 		Y: p.Y - y,
	// 		Z: 0,
	// 	},
	// 	&Vector3{
	// 		X: p.X - x,
	// 		Y: p.Y - y,
	// 		Z: 0,
	// 	},
	// }
	// obstacles = append(obstacles, Obstacle{
	// 	Polygon:  poly2,
	// 	Centroid: p.Clone(),
	// 	Radius:   poly.FindRadius(&p),
	// })

	return &Map{
		Boundaries: manifest.Boundaries,
		Spawns:     manifest.Spawns,
		Obstacles:  obstacles,
	}
}

func (m *Map) OutOfBounds(point *Vector3) bool {
	return point.X > m.Boundaries[0] || point.X < m.Boundaries[1] || point.Y > m.Boundaries[2] || point.Y < m.Boundaries[3]
}

func (m Map) OffsetMap(point *Vector3, distance float32) bool {
	return point.X > (m.Boundaries[0]+distance) || point.X < (m.Boundaries[1]-distance) || point.Y > (m.Boundaries[2]+distance) || point.Y < (m.Boundaries[3]-distance)
}

func (m Map) RandomSpawn() (float32, float32) {
	n := len(m.Spawns)
	if n == 0 {
		return 0, 0
	}

	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
	return m.Spawns[idx].X, m.Spawns[idx].Y
}
