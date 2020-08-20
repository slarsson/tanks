package game

import (
	"fmt"
	"sync"

	"github.com/slarsson/game/network"
)

const (
	RespawnTime = 2000
)

type Game struct {
	Map      *Map
	Players  map[int]*Player
	PManager *ProjectileManager
	mutex    *sync.RWMutex
	Network  *network.Network
	counter  int
}

func NewGame(n *network.Network) *Game {
	return &Game{
		Map:      NewMap(),
		Players:  make(map[int]*Player),
		PManager: NewProjectileManager(),
		mutex:    &sync.RWMutex{},
		Network:  n,
		counter:  0,
	}
}

func (g *Game) AddPlayer(client *network.Client) {
	// broadcast the names of all the connected players to the new player
	for i, v := range g.Players {
		if v.Lobby {
			continue
		}
		client.NetworkOutput <- PlayerNameMessage(i, v.Name)
	}

	g.mutex.Lock()
	g.counter++
	ID := g.counter
	g.Players[ID] = NewPlayer(ID, client)
	g.mutex.Unlock()

	fmt.Printf("GAME: "+colorBlue+"NEW PLAYER ADDED (ID = %d)"+colorReset, ID)
	fmt.Println("")

	client.NetworkOutput <- SelfNameMessage(ID)
}

func (g *Game) RemovePlayer(playerID int) {
	g.mutex.Lock()
	delete(g.Players, playerID)
	g.mutex.Unlock()

	fmt.Printf("GAME: "+colorPurple+"REMOVED PLAYER (ID = %d)"+colorReset, playerID)
	fmt.Println("")

	g.Network.Broadcast <- RemovePlayerMessage(playerID)
}

func (g Game) PlayerNameExists(name *string) bool {
	for _, v := range g.Players {
		if v.Name == *name {
			return true
		}
	}
	return false
}

func (g *Game) SetPlayerName(playerID int, name string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, v := range g.Players {
		if v.Name == name {
			return false
		}
	}

	if p, ok := g.Players[playerID]; ok {
		p.Name = name
		return true
	}
	return false
}
