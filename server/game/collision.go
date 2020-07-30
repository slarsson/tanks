package game

import (
	"fmt"
	"math"
)

// Polygon contains a slice of vectors that builds the object
type Polygon []*Vector3

// MTV (Minimum Translation Vector)
type MTV struct {
	Vector    *Vector3
	Magnitude float32
	Test      int
}

func (p *Polygon) Add(x float32, y float32, z float32) {
	*p = append(*p, &Vector3{X: x, Y: y, Z: z})
}

// NewTankHullPolygon creates a new Polygon representing the hull of the tank
func NewTankHullPolygon() *Polygon {
	return &Polygon{
		&Vector3{X: -1.5, Y: 3, Z: 0},
		&Vector3{X: 1.5, Y: 3, Z: 0},
		&Vector3{X: 1.5, Y: -3, Z: 0},
		&Vector3{X: -1.5, Y: -3, Z: 0},
	}
}

// // Rotate the vertex of a Polygon (using the start values)
// //
// // TODO: this is bad, make if better..
// func (p *Polygon) Rotate(rot float32, point *Vector3) {
// 	polygonType := NewTankHullPolygon()
// 	if len(*p) != len(*polygonType) {
// 		// TODO: print error
// 		return
// 	}

// 	sin := float32(math.Sin(float64(rot)))
// 	cos := float32(math.Cos(float64(rot)))

// 	for i, v := range *p {
// 		v.X = (*polygonType)[i].X*cos - (*polygonType)[i].Y*sin
// 		v.X += point.X
// 		v.Y = (*polygonType)[i].X*sin + (*polygonType)[i].Y*cos
// 		v.Y += point.Y
// 	}
// }

// Rotate2 rotates the Polygon in the xy-plane
func (p *Polygon) Rotate(rot float32, point *Vector3) {
	sin := float32(math.Sin(float64(rot)))
	cos := float32(math.Cos(float64(rot)))

	for _, v := range *p {
		v.X -= point.X
		v.Y -= point.Y

		xNew := v.X*cos - v.Y*sin
		yNew := v.X*sin + v.Y*cos

		v.X = xNew + point.X
		v.Y = yNew + point.Y
	}
}

func (p *Polygon) Translate(x float32, y float32, z float32) {
	for _, v := range *p {
		v.X += x
		v.Y += y
		v.Z += z
	}
}

func (p *Polygon) Print() {
	for i, v := range *p {
		fmt.Printf("idx: %d => x: %f, y: %f \n", i, v.X, v.Y)
	}
}

// projectVertices projects all vertices on a vector and then finds the max and min values
func (p Polygon) projectVertices(v *Vector3) (float32, float32) {
	var min, max, project float32
	for i := 0; i < len(p); i++ {
		project = v.X*p[i].X + v.Y*p[i].Y

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

// Collision checks if two object has collieded using the SAT (Separate Axis Theorem) algorithm
func (p Polygon) Collision(target *Polygon) (bool, *MTV) {
	wtf := -1

	var n, smallest Vector3
	overlap := math.MaxFloat32
	//overlap := float64(-1)

	lastIdx := len(p) - 1
	for i := 0; i <= lastIdx; i++ {
		if i == lastIdx {
			n = Vector3{
				X: p[lastIdx].Y - p[0].Y, // normals pointing outwards
				Y: p[0].X - p[lastIdx].X,
			}
		} else {
			n = Vector3{
				X: p[i].Y - p[i+1].Y,
				Y: p[i+1].X - p[i].X,
			}
		}
		n.Norm()

		minA, maxA := p.projectVertices(&n)
		minB, maxB := (*target).projectVertices(&n)

		if maxA < minB || maxB < minA {
			return false, nil
		}

		o := math.Abs(math.Min(float64(maxA), float64(maxB)) - math.Max(float64(minA), float64(minB)))
		//if o > 0 && o > overlap {
		if o > 0 && o < overlap {
			overlap = o
			smallest = n
			wtf = 1
		}
	}

	lastIdx = len(*target) - 1
	for i := 0; i <= lastIdx; i++ {
		if i == lastIdx {
			n = Vector3{
				X: (*target)[lastIdx].Y - (*target)[0].Y, // normals pointing outwards
				Y: (*target)[0].X - (*target)[lastIdx].X,
			}
		} else {
			n = Vector3{
				X: (*target)[i].Y - (*target)[i+1].Y,
				Y: (*target)[i+1].X - (*target)[i].X,
			}
		}
		n.Norm()

		minA, maxA := p.projectVertices(&n)
		minB, maxB := (*target).projectVertices(&n)

		if maxA < minB || maxB < minA {
			return false, nil
		}

		o := math.Abs(math.Min(float64(maxA), float64(maxB)) - math.Max(float64(minA), float64(minB)))
		//if o > 0 && o > overlap {
		if o > 0 && o < overlap {
			overlap = o
			smallest = n
			wtf = 2
		}
	}

	// lol := smallest.Clone()
	// lol.Norm()

	// //fmt.Println("SMALL:", wtf, smallest)
	// // fmt.Println("DOT:", lol.Dot(lol), target.GetCenter())

	// a := target.GetCenter()
	// b := p.GetCenter()

	// x := &Vector3{
	// 	X: a.X - b.X,
	// 	Y: a.Y - b.Y,
	// 	Z: 0,
	// }

	// // if x.Dot(lol) < 0 {
	// // 	overlap *= -1
	// // }
	// fmt.Println("DOT:", x.Dot(lol))

	return true, &MTV{
		// Vector:    smallest.Clone(),
		// Magnitude: 1.0,
		Vector:    smallest.Norm(),
		Magnitude: 1.01 * float32(overlap),
		Test:      wtf,
	}
}

// func (p *Polygon) GetCenter() *Vector3 {
// 	maxX := (*p)[0].X
// 	minX := (*p)[0].X
// 	maxY := (*p)[0].Y
// 	minY := (*p)[0].Y

// 	for _, v := range (*p)[1:] {
// 		//fmt.Println(i)
// 		if v.X > maxX {
// 			maxX = v.X
// 		}

// 		if v.X < minX {
// 			minX = v.X
// 		}

// 		if v.Y > maxY {
// 			maxY = v.Y
// 		}

// 		if v.Y < minY {
// 			minY = v.Y
// 		}
// 	}

// 	return &Vector3{
// 		X: minX + 0.5*float32(math.Abs(float64(minX-maxX))),
// 		Y: minY + 0.5*float32(math.Abs(float64(minY-maxY))),
// 		Z: 0,
// 	}
// }
