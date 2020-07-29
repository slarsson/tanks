package game

import "math"

// Vector3 ...
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (v *Vector3) Set(x float32, y float32, z float32) {
	v.X = x
	v.Y = y
	v.Z = z
}

func (v *Vector3) SetFromVector(vec *Vector3) {
	v.X = vec.X
	v.Y = vec.Y
	v.Z = vec.Z
}

func (v *Vector3) Get() []float32 {
	return []float32{v.X, v.Y, v.Z}
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

func (v *Vector3) Mult(arg float32) *Vector3 {
	v.X *= arg
	v.Y *= arg
	v.Z *= arg
	return v
}

func (v *Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v *Vector3) Clone() *Vector3 {
	return &Vector3{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}
