package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

type Player struct {
	ID             int
	Name           string
	Position       *Vector3
	Velocity       *Vector3
	Rotation       float32
	TurretRotation float32
	Client         *Client
}

type Game struct {
	Players map[int]*Player
	mutex   *sync.RWMutex
	Network *Network
}

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func newGame(n *Network) *Game {
	return &Game{
		Players: make(map[int]*Player),
		mutex:   &sync.RWMutex{},
		Network: n,
	}
}

func (g *Game) addPlayer(client *Client) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var playerID int
	for {
		playerID = rand.Intn(10000) // fejk random?
		_, ok := g.Players[playerID]
		if !ok {
			break
		}
	}

	g.Players[playerID] = &Player{
		ID:             playerID,
		Name:           "player" + strconv.Itoa(playerID),
		Position:       &Vector3{X: 0, Y: 0, Z: 0},
		Velocity:       &Vector3{X: 0, Y: 0, Z: 0},
		Rotation:       0,
		TurretRotation: 0,
		Client:         client,
	}

	fmt.Printf("\033[1;34m%s\033[0m", "new player added with id", playerID, "\n")

	buf := make([]byte, 0)
	mt := make([]byte, 4)
	id := make([]byte, 4)

	binary.LittleEndian.PutUint32(mt, 10)
	binary.LittleEndian.PutUint32(id, uint32(playerID))

	buf = append(buf, mt...)
	buf = append(buf, id...)

	fmt.Println(buf)

	g.Players[playerID].Client.NetworkOutput <- buf

}

func (g *Game) removePlayer(idx int) {
	g.mutex.Lock()
	delete(g.Players, idx)
	g.mutex.Unlock()

	buf := make([]byte, 0)
	mt := make([]byte, 4)
	id := make([]byte, 4)

	binary.LittleEndian.PutUint32(mt, 9)
	binary.LittleEndian.PutUint32(id, uint32(idx))

	buf = append(buf, mt...)
	buf = append(buf, id...)

	//fmt.Println(buf)

	g.Network.broadcast <- buf

	// broadcast new state to all clients
}

func (g *Game) crash(p1 *Player, p2 *Player) bool {

	// poly1 := newPoygon()
	// poly1.rotate(p1.Rotation, p1.Position)

	// poly2 := newPoygon()
	// poly2.rotate(p2.Rotation, p2.Position)

	// if poly1.test(poly2) {
	// 	fmt.Println("CRASH MY CRASH..")
	// 	return true
	// } else {
	// 	fmt.Println("not crash...")
	// 	return false
	// }

	// poly.print()

	// fmt.Println(poly.maxminX())
	// fmt.Println(poly.maxminY())

	// e1 := p1.corners()
	// fmt.Println(e1)
	// fmt.Println(e1.maxminX())
	// fmt.Println(e1.maxminY())
	// fmt.Println("======")

	//e2 := p2.corners()

	//fmt.Println(e1.A, e2.A)

	//normal := Vector3{X: 1, Y: 0, Z: 0}

	// A -> B
	// B -> C
	// C -> D
	// D -> A

	// slope := (e1.A.X - e1.B.X) / (e1.A.Y - e1.B.Y)
	// fmt.Println(slope)

	//for i := 0; i < 5; i++ {

	//}

	//return false

	// for i, v := range g.Players {
	// 	if i == p.ID {
	// 		continue
	// 	}

	// 	if math.Sqrt(math.Pow(float64(p.Position.X-v.Position.X), 2)+math.Pow(float64(p.Position.Y-v.Position.Y), 2)) < 5.7 {
	// 		//fmt.Println(math.Sqrt(math.Pow(float64(p.Position.X-v.Position.X), 2)))
	// 		return true
	// 	}
	// }

	return false
}

func (v *Vector3) rotate(rot float32) {
	v.X = v.X*float32(math.Cos(float64(rot))) - v.Y*float32(math.Sin(float64(rot)))
	v.Y = v.Y*float32(math.Cos(float64(rot))) + v.X*float32(math.Sin(float64(rot)))
}

func (v *Vector3) norm() {
	if v.X == 0 && v.Y == 0 && v.Z == 0 {
		return
	}

	x := 1 / math.Sqrt(float64(v.X*v.X+v.Y*v.Y+v.Z*v.Z))

	v.X = v.X * float32(x)
	v.Y = v.Y * float32(x)
	v.Z = v.Z * float32(x)
}

func (p *Player) corners() *Edges {
	//length := 6
	//width := 3

	v := Vector3{X: 1.5, Y: 3, Z: 0}
	// v.X -= p.Position.X
	// v.Y -= p.Position.Y
	v.rotate(p.Rotation)
	A := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
	//fmt.Println("A:", A)

	v = Vector3{X: -1.5, Y: 3, Z: 0}
	// v.X -= p.Position.X
	// v.Y -= p.Position.Y
	v.rotate(p.Rotation)
	B := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
	//fmt.Println("B:", B)

	v = Vector3{X: -1.5, Y: -3, Z: 0}
	// v.X -= p.Position.X
	// v.Y -= p.Position.Y
	v.rotate(p.Rotation)
	C := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
	//fmt.Println("C:", C)

	v = Vector3{X: 1.5, Y: -3, Z: 0}
	// v.X -= p.Position.X
	// v.Y -= p.Position.Y
	v.rotate(p.Rotation)
	D := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
	//fmt.Println("D:", D)

	fmt.Printf("A => x: %f, y: %f \n", A.X, A.Y)
	fmt.Printf("B => x: %f, y: %f \n", B.X, B.Y)
	fmt.Printf("C => x: %f, y: %f \n", C.X, C.Y)
	fmt.Printf("D => x: %f, y: %f \n", D.X, D.Y)
	fmt.Println("=====================")

	return &Edges{A: &A, B: &B, C: &C, D: &D}

	//A := 6/2*math.Cos(float64(p.Rotation)) + 3/2*math.Sin(float64(p.Rotation))

	//fmt.Println("A", A)

	//fmt.Println(p.Position.X + 2 + float32(math.Sin(float64(p.Rotation))))
}
