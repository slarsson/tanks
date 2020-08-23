package game

import (
	"encoding/binary"
	"math"
	"strconv"

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
	Lobby          bool
	Radius         float32
}

func NewPlayer(ID int, client *network.Client) *Player {
	return &Player{
		ID:             ID,
		Name:           "unknown_" + strconv.Itoa(ID),
		IsAlive:        false,
		Position:       &Vector3{X: 0, Y: 0, Z: 0},
		Velocity:       &Vector3{X: 0, Y: 0, Z: 0},
		Direction:      0,
		Rotation:       0,
		TurretRotation: 0,
		Client:         client,
		Controls:       NewControls(),
		WaitTime:       0,
		Lobby:          true,
		Radius:         5, // meh..
	}
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
		Radius:         5,
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

func (p *Player) Respawn(x float32, y float32) {
	p.IsAlive = true
	p.RespawnTime = 0
	p.Reset()
	p.Position.Set(x, y, 0)
}

func (p *Player) ExitLobby() {
	p.Lobby = false
	p.Respawn(0, 0)
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

func (p *Player) HandleCollsionWithObjects(objects *[]Obstacle, dt float32) {
	for _, v := range *objects {
		// broad-phase
		if p.Position.Distance(v.Centroid) > (p.Radius + v.Radius) {
			continue
		}

		// narrow-phase
		tank := NewTankHullPolygon()
		tank.Translate(p.Position.X, p.Position.Y, 0)
		tank.Rotate(p.Rotation, p.Position)

		ok, mtv := tank.Collision(v.Polygon)
		if ok {
			dx := mtv.Vector.X * mtv.Magnitude
			dy := mtv.Vector.Y * mtv.Magnitude

			after := p.Position.Clone()
			after.X += mtv.Vector.X
			after.Y += mtv.Vector.Y

			// if player is closer to the object then before => invert direction
			if v.Centroid.Distance(after) < v.Centroid.Distance(p.Position) {
				dx = -dx
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

		// broad-phase
		if p.Position.Distance(v.Position) > (p.Radius + v.Radius) {
			continue
		}

		// narrow-phase
		poly1 := NewTankHullPolygon()
		poly1.Translate(p.Position.X, p.Position.Y, 0)
		poly1.Rotate(p.Rotation, p.Position)

		poly2 := NewTankHullPolygon()
		poly2.Translate(v.Position.X, v.Position.Y, 0)
		poly2.Rotate(v.Rotation, v.Position)

		ok, mtv := poly1.Collision(poly2)
		if !ok {
			continue
		}

		dx := mtv.Vector.X * mtv.Magnitude
		dy := mtv.Vector.Y * mtv.Magnitude

		after := p.Position.Clone()
		after.X += mtv.Vector.X
		after.Y += mtv.Vector.Y

		// if player is closer to the object then before => invert direction
		if v.Position.Distance(after) < v.Position.Distance(p.Position) {
			dx = -dx
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
