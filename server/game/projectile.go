package game

type Projectile struct {
	Position  *Vector3
	Direction *Vector3
	Velocity  float32
	Owner     *Player
}

func (p *Projectile) Move(dt float32) {
	incr := dt * p.Velocity
	p.Position.X += incr * p.Direction.X
	p.Position.Y += incr * p.Direction.Y
	p.Position.Z += incr * p.Direction.Z
}
