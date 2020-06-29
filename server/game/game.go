package game

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/slarsson/game/network"
)

type Player struct {
	ID             int
	Name           string
	Position       *Vector3
	Velocity       *Vector3
	Rotation       float32
	TurretRotation float32
	Client         *network.Client
}

type Game struct {
	Map     *Map
	Players map[int]*Player
	mutex   *sync.RWMutex
	Network *network.Network
}

func NewGame(n *network.Network) *Game {
	return &Game{
		Map:     NewMap(),
		Players: make(map[int]*Player),
		mutex:   &sync.RWMutex{},
		Network: n,
	}
}

func (g *Game) AddPlayer(client *network.Client) {
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

// func (g *Game) addBot() {
// 	g.mutex.Lock()
// 	defer g.mutex.Unlock()

// 	var playerID int
// 	for {
// 		playerID = rand.Intn(10000) // fejk random?
// 		_, ok := g.Players[playerID]
// 		if !ok {
// 			break
// 		}
// 	}

// 	g.Players[playerID] = &Player{
// 		ID:             playerID,
// 		Name:           "player" + strconv.Itoa(playerID),
// 		Position:       &Vector3{X: 0, Y: 0, Z: 0},
// 		Velocity:       &Vector3{X: 0, Y: 0, Z: 0},
// 		Rotation:       0,
// 		TurretRotation: 0,
// 		Client:         nil,
// 	}
// }

func (p *Player) moveBot(dt float32) {
	p.Position.Y += dt * 0.001
}

func (g *Game) RemovePlayer(idx int) {
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

	g.Network.Broadcast <- buf

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

// func (p *Player) corners() *Edges {
// 	// //length := 6
// 	// //width := 3

// 	// v := Vector3{X: 1.5, Y: 3, Z: 0}
// 	// // v.X -= p.Position.X
// 	// // v.Y -= p.Position.Y
// 	// v.rotate(p.Rotation)
// 	// A := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
// 	// //fmt.Println("A:", A)

// 	// v = Vector3{X: -1.5, Y: 3, Z: 0}
// 	// // v.X -= p.Position.X
// 	// // v.Y -= p.Position.Y
// 	// v.rotate(p.Rotation)
// 	// B := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
// 	// //fmt.Println("B:", B)

// 	// v = Vector3{X: -1.5, Y: -3, Z: 0}
// 	// // v.X -= p.Position.X
// 	// // v.Y -= p.Position.Y
// 	// v.rotate(p.Rotation)
// 	// C := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
// 	// //fmt.Println("C:", C)

// 	// v = Vector3{X: 1.5, Y: -3, Z: 0}
// 	// // v.X -= p.Position.X
// 	// // v.Y -= p.Position.Y
// 	// v.rotate(p.Rotation)
// 	// D := Vector3{X: p.Position.X + v.X, Y: p.Position.Y + v.Y, Z: 0}
// 	// //fmt.Println("D:", D)

// 	// fmt.Printf("A => x: %f, y: %f \n", A.X, A.Y)
// 	// fmt.Printf("B => x: %f, y: %f \n", B.X, B.Y)
// 	// fmt.Printf("C => x: %f, y: %f \n", C.X, C.Y)
// 	// fmt.Printf("D => x: %f, y: %f \n", D.X, D.Y)
// 	// fmt.Println("=====================")

// 	// return &Edges{A: &A, B: &B, C: &C, D: &D}

// 	//A := 6/2*math.Cos(float64(p.Rotation)) + 3/2*math.Sin(float64(p.Rotation))

// 	//fmt.Println("A", A)

// 	//fmt.Println(p.Position.X + 2 + float32(math.Sin(float64(p.Rotation))))
// }
