package main

import (
	"fmt"
	"math"
)

type PositionState struct {
	Pivot *Vector3
	Edges *Edges
}

type Edges struct {
	A *Vector3
	B *Vector3
	C *Vector3
	D *Vector3
}

type Polygon []*Vector3

type MTV struct {
	Vector    *Vector3
	Magnitude float32
}

// = Polygon{
// 	&Vector3{X: -1.5, Y: 3, Z: 0},
// 	&Vector3{X: 1.5, Y: 3, Z: 0},
// 	&Vector3{X: 1.5, Y: -3, Z: 0},
// 	&Vector3{X: -1.5, Y: -3, Z: 0},
// }

func newPoygon() *Polygon {
	return &Polygon{
		&Vector3{X: -1.5, Y: 3, Z: 0},
		&Vector3{X: 1.5, Y: 3, Z: 0},
		&Vector3{X: 1.5, Y: -3, Z: 0},
		&Vector3{X: -1.5, Y: -3, Z: 0},
	}
}

func (p *Polygon) rotate(rot float32, point *Vector3) {
	sin := float32(math.Sin(float64(rot)))
	cos := float32(math.Cos(float64(rot)))

	swag := newPoygon() // fix this plz

	for i, v := range *p {
		v.X = (*swag)[i].X*cos - (*swag)[i].Y*sin
		v.X += point.X
		v.Y = (*swag)[i].X*sin + (*swag)[i].Y*cos
		v.Y += point.Y
	}
}

func (p *Polygon) print() {

	for _, v := range *p {
		fmt.Printf("x: %f, y: %f \n", v.X, v.Y)
	}

	// 	fmt.Printf("x: %f, y: %f \n", v.X, v.Y)
	// }
	fmt.Println("===========")
}

// func (p Polygon) maxminX() (float32, float32) {
// 	min := p[0].X
// 	max := p[0].X

// 	for _, v := range p[1:] {
// 		if v.X < min {
// 			min = v.X
// 		}
// 		if v.X > max {
// 			max = v.X
// 		}
// 	}

// 	return max, min
// }

// func (p Polygon) maxminY() (float32, float32) {
// 	min := p[0].Y
// 	max := p[0].Y

// 	for _, v := range p[1:] {
// 		if v.Y < min {
// 			min = v.Y
// 		}
// 		if v.Y > max {
// 			max = v.Y
// 		}
// 	}

// 	return max, min
// }

func (p Polygon) projectx(n *Vector3) (float32, float32) {
	var min, max, project float32
	for i := 0; i < len(p); i++ {
		project = n.X*p[i].X + n.Y*p[i].Y

		if i == 0 {
			min = project
			max = project
			continue
		}

		if project < min {
			min = project
		}

		if project > max {
			max = project
		}
	}

	return min, max
}

func (p Polygon) test(x *Polygon) (bool, *MTV) {

	var n Vector3

	overlap := math.MaxFloat32
	var smallest Vector3

	s := len(p) - 1
	for i := 0; i <= s; i++ {
		if i == s {
			n = Vector3{
				X: p[s].Y - p[0].Y, // normals pointing outwards
				Y: p[0].X - p[s].X,
			}
		} else {
			n = Vector3{
				X: p[i].Y - p[i+1].Y,
				Y: p[i+1].X - p[i].X,
			}
		}
		n.norm()

		minA, maxA := p.projectx(&n)
		minB, maxB := (*x).projectx(&n)

		if maxA < minB || maxB < minA {
			return false, nil
		}

		//o := math.Abs(math.Max(float64(minA), float64(maxB)) - math.Min(float64(maxA), float64(maxB)))
		o := math.Abs(math.Min(float64(maxA), float64(maxB)) - math.Max(float64(minA), float64(minB)))
		if o > 0 && o < overlap {
			overlap = o
			smallest = n
		}
	}

	s = len(*x) - 1
	for i := 0; i <= s; i++ {
		if i == s {
			n = Vector3{
				X: (*x)[s].Y - (*x)[0].Y, // normals pointing outwards
				Y: (*x)[0].X - (*x)[s].X,
			}
		} else {
			n = Vector3{
				X: (*x)[i].Y - (*x)[i+1].Y,
				Y: (*x)[i+1].X - (*x)[i].X,
			}
		}
		n.norm()

		minA, maxA := p.projectx(&n)
		minB, maxB := (*x).projectx(&n)

		if maxA < minB || maxB < minA {
			return false, nil
		}

		//o := math.Abs(math.Max(float64(minA), float64(maxB)) - math.Min(float64(maxA), float64(maxB)))
		o := math.Abs(math.Min(float64(maxA), float64(maxB)) - math.Max(float64(minA), float64(minB)))
		if o > 0 && o < overlap {
			overlap = o
			smallest = n
		}
	}

	smallest.norm()

	// smallest.X *= -1
	// smallest.Y *= -1

	fmt.Println("MAG:", overlap)

	return true, &MTV{Vector: &smallest, Magnitude: 1.01 * float32(overlap)}
}
