package game

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

type Projectile struct {
	Position     *Vector3
	LastPosition *Vector3
	Direction    *Vector3
	Velocity     float32
	Owner        *Player
	IsAlive      bool
	Mutex        *sync.RWMutex
}

type ProjectileManager struct {
	Projectiles map[int]*Projectile
}

func NewProjectileManager() *ProjectileManager {
	return &ProjectileManager{
		Projectiles: make(map[int]*Projectile),
	}
}

func (pm *ProjectileManager) NewProjectile(p *Player) {
	var id int
	for {
		id = rand.Intn(10000) // fejk random?
		_, ok := pm.Projectiles[id]
		if !ok {
			break
		}
	}

	pm.Projectiles[id] = p.NewProjectile()

	fmt.Println("new projectile added:", id)
}

func (pm *ProjectileManager) UpdateAll(dt float32, players *map[int]*Player) {
	for i, v := range pm.Projectiles {
		//fmt.Println(v)
		v.Move(dt)
		v.CollisionTest()
		v.CollisionTestPlayers(players)

		if !v.IsAlive {
			pm.Remove(i)
		}
	}
}

func (pm *ProjectileManager) Remove(key int) {
	fmt.Println("remove projectile with ID:", key)
	delete(pm.Projectiles, key) // maybe check ? or panic?
}

func (p *Player) NewProjectile() *Projectile {
	totRotation := p.TurretRotation + p.Rotation
	x := -float32(math.Sin(float64(totRotation)))
	y := float32(math.Cos(float64(totRotation)))

	dir := Vector3{X: x, Y: y, Z: 0}
	dir.Norm()

	return &Projectile{
		Position:     &Vector3{X: p.Position.X + 4*x, Y: p.Position.Y + 4*y, Z: 1},
		LastPosition: &Vector3{X: p.Position.X + 4*x, Y: p.Position.Y + 4*y, Z: 1},
		Direction:    &dir,
		Velocity:     0.05,
		Owner:        p,
		IsAlive:      true,
		Mutex:        &sync.RWMutex{},
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

func (p *Projectile) CollisionTest() {
	//fmt.Println("my p:", p.Position)

	wall := Polygon{
		&Vector3{X: 0, Y: 16, Z: 0},
		&Vector3{X: 10, Y: 16, Z: 0},
		&Vector3{X: 10, Y: 15, Z: 0},
		&Vector3{X: 0, Y: 15, Z: 0},
	}

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
	lastIdx := len(wall) - 1
	var A *Vector3
	var B *Vector3
	C := p.LastPosition
	D := p.Position
	for i := 0; i <= lastIdx; i++ {
		if i == lastIdx {
			A = wall[i]
			B = wall[0]
		} else {
			A = wall[i]
			B = wall[i+1]
		}

		// golang can handle division with 0 (-/+Inf)
		uA := ((D.X-C.X)*(A.Y-C.Y) - (D.Y-C.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))
		uB := ((B.X-A.X)*(A.Y-C.Y) - (B.Y-A.Y)*(A.X-C.X)) / ((D.Y-C.Y)*(B.X-A.X) - (D.X-C.X)*(B.Y-A.Y))

		if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
			fmt.Println("HIITT")
			//p.Position.X = 1000

			p.IsAlive = false
			break
		}
	}
}

func (p *Projectile) CollisionTestPlayers(players *map[int]*Player) {
	for _, v := range *players {
		if v.ID == p.Owner.ID {
			continue
		}

		poly := NewTankHullPolygon()
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
				fmt.Println("TANK HIITT")
				fmt.Println(p.Owner.ID, "killed", v.ID)

				p.IsAlive = false
				break
			}
		}
	}

	//fmt.Println("chekc playerz")
}
