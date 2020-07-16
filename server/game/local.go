package game

import (
	"fmt"
	"math"
	"time"
)

type LastState struct {
	Position       Vector3
	SequenceNumber uint32
	ShouldUpdate   bool
	Timestamp      time.Time
}

func NewLastState() *LastState {
	return &LastState{
		Position:       Vector3{X: 0, Y: 0, Z: 0},
		SequenceNumber: 0,
		ShouldUpdate:   true,
	}
}

// func (l *LastState) Compare(x float32, y float32, z float32) bool {
// 	if math.Abs(float64(l.Position.X-x)) > 0.001 || math.Abs(float64(l.Position.Y-y)) > 0.001 || math.Abs(float64(l.Position.Z-z)) > 0.001 {
// 		fmt.Println("OFFSET ERROR", math.Abs(float64(l.Position.X-x)), math.Abs(float64(l.Position.Y-y)), math.Abs(float64(l.Position.Z-z)))
// 		return true
// 	}
// 	return false
// }

func (l *LastState) Compare(vec *Vector3) bool {
	if math.Abs(float64(l.Position.X-vec.X)) > 0.001 || math.Abs(float64(l.Position.Y-vec.Y)) > 0.001 || math.Abs(float64(l.Position.Z-vec.Z)) > 0.001 {
		fmt.Println("OFFSET ERROR", math.Abs(float64(l.Position.X-vec.X)), math.Abs(float64(l.Position.Y-vec.Y)), math.Abs(float64(l.Position.Z-vec.Z)))
		return true
	}
	return false
}
