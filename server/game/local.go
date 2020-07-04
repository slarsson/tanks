package game

import (
	"fmt"
	"math"
)

type LastState struct {
	Position       Vector3
	SequenceNumber uint32
	ShouldUpdate   bool
}

func NewLastState() *LastState {
	return &LastState{
		Position:       Vector3{X: 0, Y: 0, Z: 0},
		SequenceNumber: 0,
		ShouldUpdate:   true,
	}
}

func (l *LastState) Compare(x float32, y float32, z float32) bool {
	if math.Abs(float64(l.Position.X-x)) > 0.001 || math.Abs(float64(l.Position.Y-y)) > 0.001 || math.Abs(float64(l.Position.Z-z)) > 0.001 {
		fmt.Println("OFFSET ERROR", math.Abs(float64(l.Position.X-x)), math.Abs(float64(l.Position.Y-y)), math.Abs(float64(l.Position.Z-z)))
		return true
	}
	return false
}
