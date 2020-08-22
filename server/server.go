package main

import (
	"fmt"
	"time"

	"github.com/slarsson/game/game"
	"github.com/slarsson/game/network"
)

// Server ...
type Server struct {
	Network    *network.Network
	Game       *game.Game
	Operations chan *network.Action
}

// NewServer creates the server and inits network connection
func NewServer() *Server {
	n := network.NewNetwork()
	s := &Server{
		Network:    n,
		Game:       game.NewGame(n),
		Operations: make(chan *network.Action),
	}

	go s.Network.Start(s.Operations)
	return s
}

// Manager handles operations coming from the client or from other threads/goroutines
func (s *Server) Manager() {
	for {
		select {
		case message, ok := <-s.Operations:
			if ok {
				fmt.Println("SERVER: new message, mt =", message.MessageType)

				switch message.MessageType {
				case 0:
					s.Game.AddPlayer(message.Client)
				case 1:
					for k, v := range s.Game.Players {
						if v.Client == message.Client {
							s.Game.RemovePlayer(k)
						}
					}
				case 99:
					for _, v := range s.Game.Players {
						if v.Client == message.Client {
							if game.ValidateName(message.Payload) {
								name := string(message.Payload)
								if s.Game.SetPlayerName(v.ID, name) {
									v.ExitLobby()
									s.Network.Broadcast <- game.PlayerNameMessage(v.ID, name)
								} else {
									message.Client.NetworkOutput <- game.ErrorMessage()
								}
							} else {
								message.Client.NetworkOutput <- game.ErrorMessage()
							}
						}
					}
				}
			}
		}
	}
}

// GameLoop runs the world updates and then broadcasts it to all clients
func (s *Server) GameLoop() {

	// new tick every 50ms, this should match the client broadcast rate
	step := float32(50)
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)

	for range ticker.C {
		//start := time.Now()

		// update projectiles
		s.Game.PManager.UpdateAll(step, &s.Game.Players, s.Game.Map, s.Game.Network.Broadcast)

		// create new binary buffer to store player state
		buf := game.NewPlayerState(len(s.Game.Players))

		// update all players and add to buffer
		for _, p := range s.Game.Players {
			if p.Client != nil {
				s.handleInputs(p.Client, p, step)
			}
			p.AppendPlayerState(&buf)
		}

		// broadcast the player state to all clients
		s.Network.Broadcast <- buf

		//fmt.Println("SERVER: executing time:", time.Since(start)*1000)
	}
}

// handleInputs applies player input to the game world
//
// only MessageType = 1 should be handle by this function
//
// The client is expected and should only send one input per tick, but this is not enforced by the server
func (s *Server) handleInputs(c *network.Client, p *game.Player, dt float32) {
	for len(c.NetworkInput) != 0 {
		message := <-c.NetworkInput
		if message[0] == 1 {
			p.SetSequenceNumber(&message)

			if !p.IsAlive {
				if p.Lobby {
					return
				}

				if p.RespawnTime > game.RespawnTime {
					p.Respawn(s.Game.Map.RandomSpawn())
				} else {
					p.RespawnTime += dt
				}
				return
			}

			p.Controls.Decode(&message)
			p.Move(dt)
			p.HandleCollsionWithObjects(&s.Game.Map.Obstacles, dt)
			p.HandleCollsionWithPlayers(&s.Game.Players, dt)

			if s.Game.Map.OutOfBounds(p.Position) {
				p.OutOfMapTime += dt
				if p.OutOfMapTime > 3000 {
					p.Kill()
				}
			} else {
				p.OutOfMapTime = 0
				if projectile, ok := p.Shoot(); ok {
					s.Game.PManager.Add(projectile)
				}
			}
		} else {
			fmt.Printf("SERVER @ handleInputs: UNKNOWN MESSAGE (mt = %d)", message[0])
			fmt.Println("")
		}
	}
}
