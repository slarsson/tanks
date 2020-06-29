package game

type Map struct {
	Obstacles []*Polygon
}

func NewMap() *Map {
	test := &Polygon{
		&Vector3{X: 0, Y: 16, Z: 0},
		&Vector3{X: 10, Y: 16, Z: 0},
		&Vector3{X: 10, Y: 15, Z: 0},
		&Vector3{X: 0, Y: 15, Z: 0},
	}

	return &Map{
		Obstacles: []*Polygon{test},
	}
}
