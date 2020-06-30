package game

import "fmt"

type Controls struct {
	Forward           bool
	Backward          bool
	RotateLeft        bool
	RotateRight       bool
	RotateTurretLeft  bool
	RotateTurretRight bool
	Shoot             bool
}

func DecodeControls(payload []byte) *Controls {
	fmt.Println(payload)
	return &Controls{
		Forward:           (payload[0] == 1),
		Backward:          (payload[2] == 1),
		RotateLeft:        (payload[1] == 1),
		RotateRight:       (payload[3] == 1),
		RotateTurretLeft:  (payload[4] == 1),
		RotateTurretRight: (payload[5] == 1),
		Shoot:             false,
	}
}
