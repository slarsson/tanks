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
)

type Player struct {
	ID             int
	Name           string
	IsAlive        bool
	Position       *Vector3
	Velocity       *Vector3
	Rotation       float32
	TurretRotation float32
	Client         *network.Client
	Controls       *Controls
	WaitTime       float32
	SequenceNumber uint32
}

func NewLocalPlayer() *Player {
	return &Player{
		ID:             -1,
		Name:           "uknownlocalplayer",
		IsAlive:        true,
		Position:       &Vector3{X: 0, Y: 0, Z: 0},
		Velocity:       &Vector3{X: 0, Y: 0, Z: 0},
		Rotation:       0,
		TurretRotation: 0,
		Client:         nil,
		Controls:       NewControls(),
		WaitTime:       0,
	}
}

func (p *Player) SetSequenceNumber(payload *[]byte) {
	p.SequenceNumber = binary.LittleEndian.Uint32((*payload)[8:])
}

func (p *Player) Move(dt float32) {
	p.WaitTime += dt

	if p.Controls.Forward || p.Controls.Backward {
		//if p.Velocity.Length() < 0.05 { // can still go over ..
		if p.Controls.Forward {
			p.Velocity.X -= float32(math.Sin(float64(p.Rotation))) * VelocityConstant * dt
			p.Velocity.Y += float32(math.Cos(float64(p.Rotation))) * VelocityConstant * dt
		} else {
			p.Velocity.X += float32(math.Sin(float64(p.Rotation))) * VelocityConstant * dt
			p.Velocity.Y -= float32(math.Cos(float64(p.Rotation))) * VelocityConstant * dt
		}

		if p.Velocity.Length() > MaxSpeed {
			fmt.Println("2fast")
			p.Velocity.Norm().Mult(MaxSpeed)
		}

		//}
	} else {
		p.Velocity.Y = 0
		p.Velocity.X = 0
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

}

func (p *Player) Shoot() (*Projectile, bool) {
	if p.Controls.Shoot && p.WaitTime > 200 {
		p.WaitTime = 0
		return p.NewProjectile(), true
	}
	//p.Controls.Shoot = false
	return nil, false
}

func (p *Player) HandleCollsionWithObjects(objects *[]*Polygon) {
	for _, v := range *objects {
		tank := NewTankHullPolygon()
		tank.Rotate(p.Rotation, p.Position)
		ok, mtv := tank.Collision(v)
		if ok {
			fmt.Println("CRASH WITH THE FKN WALL")
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
		if i == p.ID {
			continue
		}

		poly1 := NewTankHullPolygon()
		poly1.Rotate(p.Rotation, p.Position)

		poly2 := NewTankHullPolygon()
		poly2.Rotate(v.Rotation, v.Position)

		ok, mtv := poly1.Collision(poly2)
		if !ok {
			continue
		}

		if p.Controls.RotateLeft || p.Controls.RotateRight {
			fmt.Println("ROTATION WILL FUCK IT UP")
			if p.Controls.RotateLeft {
				p.Rotation -= 0.002 * dt
			}

			if p.Controls.RotateRight {
				p.Rotation += 0.002 * dt
			}

			poly1 = NewTankHullPolygon()
			poly1.Rotate(p.Rotation, p.Position)

			poly2 = NewTankHullPolygon()
			poly2.Rotate(v.Rotation, v.Position)

			ok, mtv = poly1.Collision(poly2)
			if !ok {
				continue
			}
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
