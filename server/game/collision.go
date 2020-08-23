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

// Rotate rotates the Polygon in the xy-plane
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

func (p *Polygon) FindRadius(v *Vector3) float32 {
	if len(*p) == 0 {
		return 0
	}

	radius := (*p)[0].Distance(v)
	for _, val := range (*p)[1:] {
		r := val.Distance(v)
		if r > radius {
			radius = r
		}
	}

	return radius
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
	var n, smallest Vector3
	overlap := math.MaxFloat32

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
		if o > 0 && o < overlap {
			overlap = o
			smallest = n
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
		if o > 0 && o < overlap {
			overlap = o
			smallest = n
		}
	}

	return true, &MTV{
		Vector:    smallest.Norm(),
		Magnitude: 1.01 * float32(overlap),
	}
}
