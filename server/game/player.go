package game

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/slarsson/game/network"
)

// Constant to handle speed of movment
const (
	VelocityConstant       = 0.0001
	MaxVelocityContant     = 0.5
	RotationConstant       = 0.001
	TurretRotationConstant = 0.002
	MaxSpeed               = 0.03
	ShootWaitTime          = 500
)

type Player struct {
	ID             int
	Name           string
	IsAlive        bool
	Position       *Vector3
	Velocity       *Vector3
	Direction      int8 // -1, 0, 1
	Rotation       float32
	TurretRotation float32
	Client         *network.Client
	Controls       *Controls
	WaitTime       float32
	OutOfMapTime   float32
	RespawnTime    float32
	SequenceNumber uint32
}

func NewLocalPlayer() *Player {
	return &Player{
		ID:             -1,
		Name:           "",
		IsAlive:        true,
		Position:       &Vector3{X: 0, Y: 0, Z: 0},
		Velocity:       &Vector3{X: 0, Y: 0, Z: 0},
		Direction:      0,
		Rotation:       0,
		TurretRotation: 0,
		Client:         nil,
		Controls:       NewControls(),
		WaitTime:       0,
		OutOfMapTime:   0,
		RespawnTime:    0,
	}
}

func (p *Player) SetSequenceNumber(payload *[]byte) {
	p.SequenceNumber = binary.LittleEndian.Uint32((*payload)[8:])
}

func (p *Player) Move(dt float32) {
	p.WaitTime += dt

	if p.Controls.Forward || p.Controls.Backward {
		if p.Controls.Forward {
			p.Velocity.X -= float32(math.Sin(float64(p.Rotation))) * VelocityConstant * dt
			p.Velocity.Y += float32(math.Cos(float64(p.Rotation))) * VelocityConstant * dt
			p.Direction = 1
		} else {
			p.Velocity.X += float32(math.Sin(float64(p.Rotation))) * VelocityConstant * dt
			p.Velocity.Y -= float32(math.Cos(float64(p.Rotation))) * VelocityConstant * dt
			p.Direction = -1
		}

		if p.Velocity.Length() > MaxSpeed {
			p.Velocity.Norm().Mult(MaxSpeed)
		}
	} else {
		p.Velocity.Y = 0
		p.Velocity.X = 0
		p.Direction = 0
	}

	if p.Controls.RotateLeft {
		p.Rotation += RotationConstant * dt
	}

	if p.Controls.RotateRight {
		p.Rotation -= RotationConstant * dt
	}

	if p.Controls.RotateTurretLeft {
		p.TurretRotation += TurretRotationConstant * dt
	}

	if p.Controls.RotateTurretRight {
		p.TurretRotation -= TurretRotationConstant * dt
	}

	p.Position.X += p.Velocity.X * dt
	p.Position.Y += p.Velocity.Y * dt

	//fmt.Println("DIR:", p.Direction)
	//fmt.Println("X:", p.Velocity.Length())
}

func (p *Player) Respawn() {
	p.IsAlive = true
	p.RespawnTime = 0
	p.Velocity.Set(0, 0, 0)
	p.Position.Set(0, 0, 0)
}

func (p *Player) Reset() {
	p.Velocity.Set(0, 0, 0)
	p.Position.Set(0, 0, 0)
	p.Rotation = 0
	p.TurretRotation = 0
}

func (p *Player) Kill() {
	p.IsAlive = false
}

func (p *Player) Shoot() (*Projectile, bool) {
	if p.Controls.Shoot && p.WaitTime > ShootWaitTime {
		p.WaitTime = 0
		return p.NewProjectile(), true
	}
	return nil, false
}

func (p *Player) HandleCollsionWithObjects(objects *[]*Polygon, dt float32) {
	for _, v := range *objects {
		tank := NewTankHullPolygon()
		//tank.Rotate(p.Rotation, p.Position)
		tank.Translate(p.Position.X, p.Position.Y, 0)
		tank.Rotate(p.Rotation, p.Position)

		ok, mtv := tank.Collision(v)
		if ok {
			// // if p.Controls.RotateLeft || p.Controls.RotateRight {
			// // 	p.Velocity.X = 0
			// // 	p.Velocity.Y = 0
			// // 	continue
			// // }

			// if mtv.Test == 2 {
			// 	dx := -1 * mtv.Vector.X * mtv.Magnitude
			// 	dy := -1 * mtv.Vector.Y * mtv.Magnitude

			// 	// if p.Velocity.X < 0 && dx < 0 {
			// 	// 	dx = -dx
			// 	// }

			// 	// if p.Velocity.Y < 0 && dy < 0 {
			// 	// 	dy = -dy
			// 	// }

			// 	if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
			// 		dx = -dx
			// 	}

			// 	if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
			// 		dy = -dy
			// 	}

			// 	p.Position.X += dx
			// 	p.Position.Y += dy

			// 	p.Velocity.X = 0
			// 	p.Velocity.Y = 0

			// } else {
			// 	// x := p.Velocity.Clone()
			// 	// x.Norm()

			// 	// if p.Direction == 1 {
			// 	// 	p.Position.X -= mtv.Magnitude * x.X
			// 	// 	p.Position.Y -= mtv.Magnitude * x.Y
			// 	// }

			// 	// if p.Direction == -1 {
			// 	// 	p.Position.X += mtv.Magnitude * x.X
			// 	// 	p.Position.Y += mtv.Magnitude * x.Y
			// 	// }

			// 	dx := mtv.Vector.X * mtv.Magnitude
			// 	dy := mtv.Vector.Y * mtv.Magnitude

			// 	if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
			// 		dx = -dx
			// 	}

			// 	if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
			// 		dy = -dy
			// 	}

			// 	p.Position.X += dx
			// 	p.Position.Y += dy

			// 	// p.Position.X = 0
			// 	// p.Position.Y = 0

			// 	p.Velocity.X = 0
			// 	p.Velocity.Y = 0
			// }
			//fmt.Println("DOT:", p.Velocity.Dot(mtv.Vector))

			meh := p.Velocity.Dot(mtv.Vector)

			if meh > -0.001 && meh < 0.001 {
				fmt.Println("MAYBE PROBLEMZ")

				//fmt.Println("obj:", v)
				// p.Position.X = 0
				// p.Position.Y = 0

				// p.Velocity.X = 0
				// p.Velocity.Y = 0
				// continue
			}

			dx := mtv.Vector.X * mtv.Magnitude
			dy := mtv.Vector.Y * mtv.Magnitude

			if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
				dx = -dx
			}

			if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
				dy = -dy
			}

			p.Position.X += dx
			p.Position.Y += dy

			p.Velocity.X = 0
			p.Velocity.Y = 0
		}
	}
}

func (p *Player) HandleCollsionWithPlayers(players *map[int]*Player, dt float32) {
	for i, v := range *players {
		if i == p.ID || !v.IsAlive {
			continue
		}

		poly1 := NewTankHullPolygon()
		//poly1.Rotate(p.Rotation, p.Position)
		poly1.Translate(p.Position.X, p.Position.Y, 0)
		poly1.Rotate(p.Rotation, p.Position)

		poly2 := NewTankHullPolygon()
		//poly2.Rotate(v.Rotation, v.Position)
		poly2.Translate(v.Position.X, v.Position.Y, 0)
		poly2.Rotate(v.Rotation, v.Position)

		ok, mtv := poly1.Collision(poly2)
		if !ok {
			continue
		}

		dx := mtv.Vector.X * mtv.Magnitude
		dy := mtv.Vector.Y * mtv.Magnitude

		if (dx < 0 && p.Velocity.X < 0) || (dx > 0 && p.Velocity.X > 0) {
			dx = -dx
		}

		if (dy < 0 && p.Velocity.Y < 0) || (dy > 0 && p.Velocity.Y > 0) {
			dy = -dy
		}

		p.Position.X += dx
		p.Position.Y += dy

		p.Velocity.X = 0
		p.Velocity.Y = 0
	}
}

func (p *Player) SyncState(target *Player) {
	p.Position.SetFromVector(target.Position)
	p.Rotation = target.Rotation
	p.TurretRotation = target.TurretRotation
}
