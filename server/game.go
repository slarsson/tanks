package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

type Player struct {
	ID       int
	Name     string
	Position *Vector3
	Velocity *Vector3
	Client   *Client
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
		ID:       playerID,
		Name:     "player" + strconv.Itoa(playerID),
		Position: &Vector3{X: 0, Y: 0, Z: 0},
		Velocity: &Vector3{X: 0, Y: 0, Z: 0},
		Client:   client,
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

// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"sync"
// )

// type Point struct {
// 	X float32
// 	Y float32
// 	Z float32
// }

// type Vector3 struct {
// 	X float32
// 	Y float32
// 	Z float32
// }

// type Player struct {
// 	ID       int
// 	Name     string
// 	Position *Vector3
// 	Client   *Client
// }

// type Game struct {
// 	Players map[int]*Player
// 	mutex   *sync.RWMutex
// }

// func newGame() *Game {
// 	return &Game{
// 		Players: make(map[int]*Player),
// 		mutex:   &sync.RWMutex{},
// 	}
// }

// func (g *Game) addPlayer(client *Client) {
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
// 		ID:       playerID,
// 		Name:     "player" + string(playerID),
// 		Position: &Vector3{X: 0, Y: 0, Z: 0},
// 		Client:   client,
// 	}

// 	fmt.Printf("\033[1;34m%s\033[0m", "Player added to server")
// }

// func (p *Point) print() {
// 	fmt.Printf("Point: {x: %f, y: %f, z: %f} \n", p.X, p.Y, p.Z)
// }

// func (p *Point) move() {
// 	p.X += 0.1
// }
