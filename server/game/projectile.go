package game

import (
	"fmt"
	"math"
	"sync"
)

type Projectile struct {
	Position     *Vector3
	LastPosition *Vector3
	Direction    *Vector3
	Velocity     float32
	Owner        *Player
	IsAlive      bool
}

type ProjectileManager struct {
	Projectiles map[int]*Projectile
	mutex       *sync.RWMutex
	counter     int
}

func NewProjectileManager() *ProjectileManager {
	return &ProjectileManager{
		Projectiles: make(map[int]*Projectile),
		mutex:       &sync.RWMutex{},
		counter:     0,
	}
}

func (pm *ProjectileManager) Add(p *Projectile) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.counter++
	pm.Projectiles[pm.counter] = p
	fmt.Println("PManager: add:", pm.counter)
}

func (pm *ProjectileManager) UpdateAll(dt float32, players *map[int]*Player, m *Map, broadcast chan []byte) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for i, v := range pm.Projectiles {
		//fmt.Println(v)
		v.Move(dt)
		v.CollisionTest(m)

		id := v.CollisionTestPlayers(players)
		if id != -1 {
			broadcast <- *v.KillMessage(id)
		}

		if !v.IsAlive {
			fmt.Println("PManager: remove:", i)
			delete(pm.Projectiles, i)
		}
	}
}

func (pm *ProjectileManager) UpdateLocal(dt float32, players *map[int]*Player, m *Map) interface{} {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	buf := make([]interface{}, len(pm.Projectiles)*5)
	idx := 0

	for i, v := range pm.Projectiles {
		v.Move(dt)
		v.CollisionTest(m)

		id := v.CollisionTestPlayers(players)
		if id != -1 {
			fmt.Println("KILLED LOCAL!?")
		}

		if !v.IsAlive {
			fmt.Println("PManager: remove:", i)
			delete(pm.Projectiles, i)

			buf[idx] = i
			buf[idx+1] = 0
			buf[idx+2] = 0
			buf[idx+3] = 0
			buf[idx+4] = 0
		} else {
			buf[idx] = i
			buf[idx+1] = v.Position.X
			buf[idx+2] = v.Position.Y
			buf[idx+3] = v.Position.Z
			buf[idx+4] = 1
		}

		idx += 5
	}

	return buf
}

func (p *Player) NewProjectile() *Projectile {
	totRotation := p.TurretRotation + p.Rotation
	x := -float32(math.Sin(float64(totRotation)))
	y := float32(math.Cos(float64(totRotation)))

	dir := Vector3{X: x, Y: y, Z: 0}
	dir.Norm()

	return &Projectile{
		Position:     &Vector3{X: p.Position.X + 4*x, Y: p.Position.Y + 4*y, Z: 1.5},
		LastPosition: &Vector3{X: p.Position.X + 4*x, Y: p.Position.Y + 4*y, Z: 1.5},
		Direction:    &dir,
		Velocity:     0.05,
		Owner:        p,
		IsAlive:      true,
	}
}

func (p *Projectile) Move(dt float32) {
	// set LastPosition to current Position
	// NOTE: swapping the pointers might be better, but it generates fucked up behaviour with wasm
	p.LastPosition.X = p.Position.X
	p.LastPosition.Y = p.Position.Y
	p.LastPosition.Z = p.Position.Z

	// calulcate next position
	incr := dt * p.Velocity
	p.Position.X += incr * p.Direction.X
	p.Position.Y += incr * p.Direction.Y
	p.Position.Z += incr * p.Direction.Z
}

func (p *Projectile) CollisionTest(m *Map) {
	if m.OffsetMap(p.Position, 50) {
		//if m.OutOfBounds(p.Position) {
		p.IsAlive = false
		return
	}

	//fmt.Println("my p:", p.Position)

	// wall := Polygon{
	// 	&Vector3{X: 0, Y: 16, Z: 0},
	// 	&Vector3{X: 10, Y: 16, Z: 0},
	// 	&Vector3{X: 10, Y: 15, Z: 0},
	// 	&Vector3{X: 0, Y: 15, Z: 0},
	// }

	// lastIdx := len(wall) - 1
	// var A *Vector3
	// var B *Vector3

	// wtf := 0
	// for i := 0; i <= lastIdx; i++ {
	// 	if i == lastIdx {
	// 		A = wall[i]
	// 		B = wall[0]
	// 	} else {
	// 		A = wall[i]
	// 		B = wall[i+1]
	// 	}

	// 	//fmt.Println(A, "=>", B)

	// 	ABx := B.X - A.X
	// 	ABy := B.Y - A.Y
	// 	APx := p.Position.X - A.X
	// 	APy := p.Position.Y - A.Y

	// 	val := ABx*APy - ABy*APx

	// 	if val < 0 {
	// 		wtf++
	// 	}
	// 	//fmt.Println(val, wtf)
	// }

	// if wtf == 4 {
	// 	fmt.Println("HIT AT:", p.Position)
	// }
	// // P = point to check
	// // A to B
	// // AB (dot) AP

	// //fmt.Println(wall)

	// TEST

	// basic fkn math: https://stackoverflow.com/a/565282/1671375
	// lastIdx := len(wall) - 1
	// var A *Vector3
	// var B *Vector3
	// C := p.LastPosition
	// D := p.Position
	// for i := 0; i <= lastIdx; i++ {
	// 	if i == lastIdx {
	// 		A = wall[i]
	// 		B = wall[0]
	// 	} else {
	// 		A = wall[i]
	// 		B = wall[i+1]
	// 	}

	// 	// golang can handle division with 0 (-/+Inf)
	// 	uA := ((D.X-C.X)*(A.Y-C.Y) - (D.Y-C.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))
	// 	uB := ((B.X-A.X)*(A.Y-C.Y) - (B.Y-A.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))

	// 	if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
	// 		//fmt.Println("HIITT")
	// 		//p.Position.X = 1000

	// 		p.IsAlive = false
	// 		break
	// 	}
	// }

	for _, obj := range m.Obstacles {
		lastIdx := len(*obj) - 1
		var A *Vector3
		var B *Vector3
		C := p.LastPosition
		D := p.Position
		for i := 0; i <= lastIdx; i++ {
			if i == lastIdx {
				A = (*obj)[i]
				B = (*obj)[0]
			} else {
				A = (*obj)[i]
				B = (*obj)[i+1]
			}

			// golang can handle division with 0 (-/+Inf)
			uA := ((D.X-C.X)*(A.Y-C.Y) - (D.Y-C.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))
			uB := ((B.X-A.X)*(A.Y-C.Y) - (B.Y-A.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))

			if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
				fmt.Println("HIITT WALL")
				//p.Position.X = 1000

				p.IsAlive = false
				break
			}
		}

	}
}

// TODO: this will fuck up when using localPlayer !?!?!
func (p *Projectile) CollisionTestPlayers(players *map[int]*Player) int {
	for _, v := range *players {
		if v.ID == p.Owner.ID || !v.IsAlive {
			continue
		}

		poly := NewTankHullPolygon()
		poly.Translate(v.Position.X, v.Position.Y, 0)
		poly.Rotate(v.Rotation, v.Position)

		lastIdx := len(*poly) - 1
		var A *Vector3
		var B *Vector3
		C := p.LastPosition
		D := p.Position
		for i := 0; i <= lastIdx; i++ {
			if i == lastIdx {
				A = (*poly)[i]
				B = (*poly)[0]
			} else {
				A = (*poly)[i]
				B = (*poly)[i+1]
			}

			// golang can handle division with 0 (-/+Inf)
			uA := ((D.X-C.X)*(A.Y-C.Y) - (D.Y-C.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))
			uB := ((B.X-A.X)*(A.Y-C.Y) - (B.Y-A.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))

			if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
				fmt.Println("TANK HIITT", v)
				//fmt.Println(p.Owner.ID, "killed", v.ID)
				//v.IsAlive = false
				//v.Position.Set(rand.Float32()*60-30, rand.Float32()*60-30, 0)
				v.IsAlive = false
				p.IsAlive = false
				//break
				return v.ID // might not be the correct player :()
			}
		}
	}
	return -1
	//fmt.Println("chekc playerz")
}
