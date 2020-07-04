package game

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/slarsson/game/network"
)

type Game struct {
	Map         *Map
	Players     map[int]*Player
	Projectiles map[int]*Projectile
	mutex       *sync.RWMutex
	Network     *network.Network
}

func NewGame(n *network.Network) *Game {
	return &Game{
		Map:         NewMap(),
		Players:     make(map[int]*Player),
		Projectiles: make(map[int]*Projectile),
		mutex:       &sync.RWMutex{},
		Network:     n,
	}
}

func (g *Game) AddProjectile(p *Player) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var id int
	for {
		id = rand.Intn(10000) // fejk random?
		_, ok := g.Projectiles[id]
		if !ok {
			break
		}
	}

	g.Projectiles[id] = p.NewProjectile()
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
		Controls:       NewControls(),
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

	g.Network.Broadcast <- buf
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

// func (p *Player) moveBot(dt float32) {
// 	p.Position.Y += dt * 0.001
// }
