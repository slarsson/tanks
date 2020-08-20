package game

import (
	"encoding/binary"
	"math"
)

// SEND MESSAGE
// type 0: game state
// type 1: add player
// type 2: remove player

// type 0:
func (p Player) AppendPlayerState(buf *[]byte) {
	id := make([]byte, 4)
	px := make([]byte, 4)
	py := make([]byte, 4)
	pz := make([]byte, 4)
	vx := make([]byte, 4)
	vy := make([]byte, 4)
	vz := make([]byte, 4)
	rotation := make([]byte, 4)
	turretRotation := make([]byte, 4)
	shoot := make([]byte, 4)
	sq := make([]byte, 4)
	alive := make([]byte, 4)

	binary.LittleEndian.PutUint32(id, math.Float32bits(float32(p.ID)))
	binary.LittleEndian.PutUint32(px, math.Float32bits(p.Position.X))
	binary.LittleEndian.PutUint32(py, math.Float32bits(p.Position.Y))
	binary.LittleEndian.PutUint32(pz, math.Float32bits(p.Position.Z))
	binary.LittleEndian.PutUint32(vx, math.Float32bits(p.Velocity.X))
	binary.LittleEndian.PutUint32(vy, math.Float32bits(p.Velocity.Y))
	binary.LittleEndian.PutUint32(vz, math.Float32bits(p.Velocity.Z))
	binary.LittleEndian.PutUint32(rotation, math.Float32bits(p.Rotation))
	binary.LittleEndian.PutUint32(turretRotation, math.Float32bits(p.TurretRotation))

	if p.Controls.Shoot && p.WaitTime < 0.01 {
		binary.LittleEndian.PutUint32(shoot, math.Float32bits(float32(1)))
	} else {
		binary.LittleEndian.PutUint32(shoot, math.Float32bits(float32(0)))
	}

	binary.LittleEndian.PutUint32(sq, math.Float32bits(float32(p.SequenceNumber)))

	if p.IsAlive {
		binary.LittleEndian.PutUint32(alive, math.Float32bits(float32(1)))
	} else {
		binary.LittleEndian.PutUint32(alive, math.Float32bits(float32(0)))
	}

	*buf = append(*buf, id...)
	*buf = append(*buf, px...)
	*buf = append(*buf, py...)
	*buf = append(*buf, pz...)
	*buf = append(*buf, vx...)
	*buf = append(*buf, vy...)
	*buf = append(*buf, vz...)
	*buf = append(*buf, rotation...)
	*buf = append(*buf, turretRotation...)
	*buf = append(*buf, shoot...)
	*buf = append(*buf, sq...)
	*buf = append(*buf, alive...)
}

// func (p Player) AppendPlayerShoot(buf *[]byte) {
// 	id := make([]byte, 4)
// 	test := make([]byte, 4)

// 	binary.LittleEndian.PutUint32(id, math.Float32bits(float32(p.ID)))
// 	binary.LittleEndian.PutUint32(test, math.Float32bits(1337))

// 	*buf = append(*buf, id...)
// 	*buf = append(*buf, test...)
// }

func (p *Projectile) KillMessage(kk int) *[]byte {
	buf := make([]byte, 0, 12)
	mt := make([]byte, 4)
	killer := make([]byte, 4)
	killed := make([]byte, 4)
	binary.LittleEndian.PutUint32(mt, 99)
	binary.LittleEndian.PutUint32(killer, uint32(p.Owner.ID))
	binary.LittleEndian.PutUint32(killed, uint32(kk))
	buf = append(buf, mt...)
	buf = append(buf, killer...)
	buf = append(buf, killed...)
	return &buf
}

func TestMessage() *[]byte {
	buf := make([]byte, 0, 12)

	// mt := make([]byte, 4)
	// killer := make([]byte, 4)
	// killed := make([]byte, 4)

	// binary.LittleEndian.PutUint32(mt, 99)
	// binary.LittleEndian.PutUint32(killer, 99)
	// binary.LittleEndian.PutUint32(mt, 99)

	// buf = append(buf, mt...)
	return &buf
}

// mt = 98
func PlayerNameMessage(playerID int, data string) []byte {
	buf := make([]byte, 8)
	str := []byte(data)
	binary.LittleEndian.PutUint32(buf[0:4], 98)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(playerID))
	buf = append(buf, str...)
	return buf
}

// mt = 10
func SelfNameMessage(playerID int) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], 10)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(playerID))
	return buf
}

// mt = 9
func RemovePlayerMessage(playerID int) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], 9)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(playerID))
	return buf
}

func TestErrorMessage() *[]byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 33)
	return &buf
}
