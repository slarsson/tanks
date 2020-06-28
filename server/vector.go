package main

import "math"

// Vector3 ...
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// Rotate rotates a vector around origo
func (v *Vector3) Rotate(rot float32) *Vector3 {
	sin := float32(math.Sin(float64(rot)))
	cos := float32(math.Cos(float64(rot)))
	v.X = v.X*cos - v.Y*sin
	v.Y = v.Y*cos + v.X*sin
	return v
}

// Norm normalize a vector
func (v *Vector3) Norm() *Vector3 {
	if v.X == 0 && v.Y == 0 && v.Z == 0 {
		return v
	}

	val := float32(1 / math.Sqrt(float64(v.X*v.X+v.Y*v.Y+v.Z*v.Z)))
	v.X = v.X * val
	v.Y = v.Y * val
	v.Z = v.Z * val
	return v
}
