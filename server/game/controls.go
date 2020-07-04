package game

import (
	"fmt"
)

type Controls struct {
	Forward           bool
	Backward          bool
	RotateLeft        bool
	RotateRight       bool
	RotateTurretLeft  bool
	RotateTurretRight bool
	Shoot             bool
}

func NewControls() *Controls {
	return &Controls{
		Forward:           false,
		Backward:          false,
		RotateLeft:        false,
		RotateRight:       false,
		RotateTurretLeft:  false,
		RotateTurretRight: false,
		Shoot:             false,
	}
}

// OBS OBS: first index assumed to be messagetype
func (c *Controls) Decode(payload *[]byte) {
	// TODO: bad data input => PANIC!?
	c.Forward = ((*payload)[1] == 1)
	c.Backward = ((*payload)[3] == 1)
	c.RotateLeft = ((*payload)[2] == 1)
	c.RotateRight = ((*payload)[4] == 1)
	c.RotateTurretLeft = ((*payload)[5] == 1)
	c.RotateTurretRight = ((*payload)[6] == 1)
	c.Shoot = ((*payload)[7] == 1)

	//fmt.Println(binary.LittleEndian.Uint32((*payload)[8:]))
}

func (c Controls) Print() {
	fmt.Println(c)
}

func DecodeControls(payload []byte) *Controls {
	//fmt.Println(payload)
	return &Controls{
		Forward:           (payload[0] == 1),
		Backward:          (payload[2] == 1),
		RotateLeft:        (payload[1] == 1),
		RotateRight:       (payload[3] == 1),
		RotateTurretLeft:  (payload[4] == 1),
		RotateTurretRight: (payload[5] == 1),
		Shoot:             (payload[6] == 1),
	}
}
