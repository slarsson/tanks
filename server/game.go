package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Point struct {
	X float32
	Y float32
	Z float32
}

type Player struct {
	ID       int
	Name     string
	Position *Point
	Client   *Client
}

type Game struct {
	Players map[int]*Player
	mutex   *sync.RWMutex
}

func newGame() *Game {
	return &Game{
		Players: make(map[int]*Player),
		mutex:   &sync.RWMutex{},
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
		Name:     "player" + string(playerID),
		Position: &Point{X: 0, Y: 0, Z: 0},
		Client:   client,
	}
}

func (p *Point) print() {
	fmt.Printf("Point: {x: %f, y: %f, z: %f} \n", p.X, p.Y, p.Z)
}

func (p *Point) move() {
	p.X += 0.1
}
