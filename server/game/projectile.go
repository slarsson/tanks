package game

import (
	"math"
)

type Projectile struct {
	Position  *Vector3
	Direction *Vector3
	Velocity  float32
	Owner     *Player
}

func (p *Player) NewProjectile() *Projectile {
	totRotation := p.TurretRotation + p.Rotation
	x := -float32(math.Sin(float64(totRotation)))
	y := float32(math.Cos(float64(totRotation)))

	dir := Vector3{X: x, Y: y, Z: 0}
	dir.Norm()

	return &Projectile{
		Position:  &Vector3{X: p.Position.X + 4*x, Y: p.Position.Y + 4*y, Z: 1},
		Direction: &dir,
		Velocity:  0.05,
		Owner:     p,
	}
}

func (p *Projectile) Move(dt float32) {
	incr := dt * p.Velocity
	p.Position.X += incr * p.Direction.X
	p.Position.Y += incr * p.Direction.Y
	p.Position.Z += incr * p.Direction.Z
}
