package game

import "fmt"

type Map struct {
	Obstacles  []*Polygon
	Boundaries [4]float32 // maxX, minX, maxY, minY
}

func NewMap() *Map {
	test := &Polygon{
		&Vector3{X: 0, Y: 16, Z: 0},
		&Vector3{X: 10, Y: 16, Z: 0},
		&Vector3{X: 10, Y: 15, Z: 0},
		&Vector3{X: 0, Y: 15, Z: 0},
	}

	return &Map{
		Obstacles:  []*Polygon{test},
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
